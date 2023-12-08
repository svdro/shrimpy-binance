package client

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

// newRateLimitCounter returns a new rateLimitCounter.
func newRateLimitCounter(
	th common.TimeHandler,
	endpointType common.BIEndpointType,
	rateLimitType common.BIRateLimitType,
	intervalSeconds int,
	limit int,
	logger *log.Entry,
) *rateLimitCounter {

	return &rateLimitCounter{
		th:              th,
		limit:           limit,
		endpointType:    endpointType,
		rateLimitType:   rateLimitType,
		intervalSeconds: intervalSeconds,
		logger: logger.WithField("_caller", fmt.Sprintf(
			"rateLimitCounter(%s|%s|%d sec|limit=%d)",
			endpointType, rateLimitType, intervalSeconds, limit),
		),
	}
}

// rateLimitCounter manages the rate limiting for Binance http requests.
// Binance servers enforce rate limits on various endpoints (e.g., 'api',
// 'sapi'), types (e.g., IP, UID, RAW), and across different time intervals
// (e.g., 1 second, 1 minute, 1 day). An instance of rateLimitCounter keeps
// track of these limits for one endpoint, type, and interval combination.
//
// Binance servers accumulate rate limit counts over the interval and reset
// at the beginning of each new interval. Exceeding the rate limit within an
// interval will result in the server responding with a status 429 or 418
// status code, which also indicates the number of seconds to wait before
// issuing another request. This wait time does not necessarily align with
// the start of the next interval.
//
// rateLimitCounter maintains two main counts:
// 1. countUsed: The number of requests utilized in the current interval.
// 2. countPending: The number of requests initiated but not yet completed.
//
// Usage:
// Upon initiating a request, countPending is incremented by the weight of
// the request. Once the request completes, countPending is decremented by
// the same amount, and countUsed is updated based on the value returned in
// the response header.
// If doing the request would exceed the rate limit, then countPending is not
// incremented, and a RateLimitError is returned, signaling to the caller
// that the request should be retried at the time specified in the error.
//
// Mutex:
// all public facing functions of rateLimitCounter have a mutex to ensure
// that there are no race conditions. The use of mutex is not scoped more
// narrowly because the order of operations is important, and we want to
// avoid simultaneous calls from messing with the logic.
type rateLimitCounter struct {
	mu              sync.Mutex
	th              common.TimeHandler
	endpointType    common.BIEndpointType
	rateLimitType   common.BIRateLimitType
	intervalSeconds int
	limit           int // -1 means no limit
	countUsed       int
	countPending    int
	currInterval    int64
	logger          *log.Entry
}

// tssToInterval returns the interval that the given nano timestamp is in.
// The interval is calculated as the nth interval since the unix epoch.
func (rlc *rateLimitCounter) tssToInterval(tss common.TSNano) int64 {
	return tss.Int64() / int64(rlc.intervalSeconds*1e9)
}

// intervalToTSS returns the first nano timestamp that occurs in the given
// interval. The interval is calculated as the nth interval since the unix
func (rlc *rateLimitCounter) intervalToTSS(interval int64) common.TSNano {
	return common.TSNano(interval * int64(rlc.intervalSeconds*1e9))
}

// newRateLimitError returns a new common.RateLimitError.
func (rlc *rateLimitCounter) newRateLimitError(tslRetryAt common.TSNano, countProjected int) error {
	reason := fmt.Sprintf("request would exceed limit. (%d/ %d)", countProjected, rlc.limit)
	return &common.RateLimitError{
		StatusCode:     0,
		ErrorCode:      0,
		Msg:            reason,
		Producer:       "shrimpy-binance",
		RetryTimeLocal: time.Unix(0, tslRetryAt.Int64()),
		RetryAfter:     int(tslRetryAt-rlc.th.TSLNow()) / 1e9,
	}
}

// SetUsed does two things. It keeps track of what interval we are in, and it
// updates the countUsed based on the response header from the last request.
// To deal with concurrent requests we only update countUsed if tssResp is
// greater than the current interval, or if tssResp is in the currentInteval
// and countUsed is greater than the current countUsed.
// tssResp is the timestamp of the response header. If the response header
// is not available (e.g. a WebsocketAPI response), then tssResp should be
// set to the timestamp (in server time) the request was made, that way we
// don't accidentally carry over countUsed from a previous interval.
func (rlc *rateLimitCounter) SetUsed(countUsed int, tssResp common.TSNano) {

	// lock the mutex (ensure order of operations)
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	// calculate tssResp Interval
	currInterval := rlc.tssToInterval(tssResp)
	lastInterval := rlc.currInterval

	logger := rlc.logger.WithFields(log.Fields{
		"currInterval":  currInterval,
		"countUsed":     countUsed,
		"lastInterval":  lastInterval,
		"rlc.countUsed": rlc.countUsed,
	})

	// don't do anything if countUsed refers to an interval that is in the past.
	if currInterval < lastInterval {
		logger.Debug("SetUsed: currInterval < rlc.currInterval")
		return
	}

	// always log if countUsed exceeds 75% of the limit
	if countUsed > rlc.limit*3/4 {
		logger.Debug("SetUsed: countUsed approaches limit")
	}

	// always set countUsed if currInterval is greater than rlc.currInterval,
	// also set rlc.currInterval to the new currInterval.
	if currInterval > lastInterval {
		logger.Debug("SetUsed: currInterval > rlc.currInterval")
		rlc.currInterval = currInterval
		rlc.countUsed = countUsed
		return
	}

	// if we're in the same interval, only set countUsed if countUsed is
	// greater than the current countUsed.
	if countUsed > rlc.countUsed {
		rlc.countUsed = countUsed
	}
}

// IncrementPending increments the countPending by incrPending. If the
// projected count exceeds the limit, a RateLimitError is returned.
//
// When calculating the projected count, we need to first check if we're still
// in the same interval. If we're already in a new interval, then we can
// ignore rlc.countUsed from the previous interval, and assume it is 0.
// NOTE: We do not update the value of rlc.currInterval here. SetUsed is
// responsible for updating rlc.currInterval.
func (rlc *rateLimitCounter) IncrementPending(incrPending int) error {

	// lock the mutex (to ensure order of operations)
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	// don't do anything if incrPending is 0.
	if incrPending == 0 {
		return nil
	}

	// always update count pending if limit is -1.
	if rlc.limit == -1 {
		rlc.countPending += incrPending
		return nil
	}

	// get current interval and countUsed
	tss := rlc.th.TSSNow()
	currInterval := rlc.tssToInterval(tss)
	logger := rlc.logger.WithFields(log.Fields{
		"currInterval":     currInterval,
		"incrPending":      incrPending,
		"rlc.countPending": rlc.countPending,
		"rlc.countUsed":    rlc.countUsed,
		"rlc.currInterval": rlc.currInterval,
		"rlc.limit":        rlc.limit,
	})

	countUsed := rlc.countUsed
	if currInterval > rlc.currInterval {
		logger.Debug("IncrementPending: currInterval > rlc.currInterval")
		countUsed = 0
	}

	// calculate projected count. Return a RateLimitError if the projected
	// count exceeds the limit.
	countProjected := countUsed + rlc.countPending + incrPending
	logger = logger.WithField("countProjected", countProjected)
	if countProjected > rlc.limit {
		tssRetryAt := rlc.intervalToTSS(currInterval + 1)
		tslRetryAt := rlc.th.TSSToTSL(tssRetryAt)
		err := rlc.newRateLimitError(tslRetryAt, countProjected)
		logger.WithError(err).Warn("IncrementPending: countProjected > rlc.limit")
		return err
	}

	// always log when countProjected exceeds 75% of the limit.
	if countProjected > rlc.limit*3/4 {
		logger.Debug("IncrementPending: countProjected approaches limit")
	}

	// if the projected count does not exceed the limit, then increment pending
	// and return nil (greenlight the request).
	rlc.countPending += incrPending
	return nil
}

// DecrementPending decrements the countPending by decrPending.
func (rlc *rateLimitCounter) DecrementPending(decrPending int) {

	// lock the mutex (ensure order of operations)
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	// don't do anything if decrPending is 0.
	if decrPending == 0 {
		return
	}

	// decrement
	newCountPending := rlc.countPending - decrPending
	rlc.countPending = newCountPending

	// Something is seriously wrong if rlc.countPending falls below 0.
	// Correct it, but issue a serious warning.
	// NOTE: if rateLimitManager receives a rateLimitUpdate for an uninitialized
	// rateLimitCounter, it will create a new rateLimitCounter. In this case
	// countPending will be 0, and if we decrement it, it will always fall below 0.
	// This only happens on the first request, so it's not that big of a deal.
	if newCountPending < 0 {
		rlc.countPending = 0

		tss := rlc.th.TSSNow()
		currInterval := rlc.tssToInterval(tss)
		err := fmt.Errorf("countPending < 0. Correcting for now, but fix this as soon as possible!")
		rlc.logger.WithError(err).WithFields(log.Fields{
			"currInterval":     currInterval,
			"decrPending":      decrPending,
			"rlc.countPending": rlc.countPending,
			"rlc.countUsed":    rlc.countUsed,
			"rlc.currInterval": rlc.currInterval,
			"rlc.limit":        rlc.limit,
		}).Error("DecrementPending")
	}
}
