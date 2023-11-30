package client

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

// countPending is used to unregister pending requests from a rateLimitCounter and
// is also necessary for rolling back pending counts when a request fails.
type countPending struct {
	rlc   *rateLimitCounter
	count int
}

// utitlity for initializing a slice of countsPending with the same count.
func newPendingCounts(rlcs []*rateLimitCounter, count int) []countPending {
	countsPending := []countPending{}
	for _, rlc := range rlcs {
		countsPending = append(countsPending, countPending{rlc, count})
	}
	return countsPending
}

// rateLimitKey is used to access rateLimitCounters in the rateLimitManager.
type rateLimitKey struct {
	endpointType    common.BIEndpointType
	rateLimitType   common.BIRateLimitType
	intervalSeconds int
}

// newRateLimitManager creates a new rateLimitManager
func newRateLimitManager(
	rateLimits []RateLimit, th common.TimeHandler, logger *log.Entry) *rateLimitManager {
	rlm := &rateLimitManager{
		th:     th,
		rlcs:   make(map[rateLimitKey]*rateLimitCounter),
		logger: logger.WithField("_caller", "rateLimitManager"),
	}

	for _, rl := range rateLimits {
		if err := rlm.addRLC(rl); err != nil {
			logger := log.WithField("_caller", "newRateLimitManager")
			logger.WithError(err).Panic("failed to add rate limit counter")
		}

	}
	return rlm
}

// rateLimitManager is responsible for managing binance rate-limits across
// all endpoints, rate limit types, and intervals.
// when a request is made, the rateLimitManager is used to register the
// request as pending. If the request would exceed the rate limit, then
// the request is not registered, and a RateLimitError is returned.
// If the request is successful, then the rateLimitManager is used to
// update the used count for all rate limits that receive updates in the
// response header.
// In any case (request success or failure), the request must be unregistered
// as pending.
type rateLimitManager struct {
	mu     sync.Mutex
	th     common.TimeHandler
	rlcs   map[rateLimitKey]*rateLimitCounter
	logger *log.Entry
}

// addRateLimitCounter adds a new rateLimitCounter to rateLimitManager.rlcs map.
// Returns an error if a rateLimitCounter already exists for the given endpointType,
func (rlm *rateLimitManager) addRateLimitCounter(key rateLimitKey, rlc *rateLimitCounter) error {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	if _, ok := rlm.rlcs[key]; ok {
		errMsg := "rate limit already exists for endpointType: %s, rateLimitType: %s, intervalSeconds: %d"
		return fmt.Errorf(errMsg, key.endpointType, key.rateLimitType, key.intervalSeconds)
	}

	rlm.rlcs[key] = rlc
	return nil
}

// addRLC creates a new rateLimitCounter and adds it to the rateLimitManager.
// Returns an error if a rateLimitCounter already exists for the given endpointType,
// rateLimitType, and intervalSeconds.
func (rlm *rateLimitManager) addRLC(rl RateLimit) error {
	seconds := getSecondsInInterval(rl.RateLimitIntervalType, rl.RateLimitIntervalNum)
	key := rateLimitKey{rl.EndpointType, rl.RateLimitType, seconds}

	rlc := newRateLimitCounter(rlm.th, rl.EndpointType, rl.RateLimitType, seconds, rl.Limit, rlm.logger)
	return rlm.addRateLimitCounter(key, rlc)
}

// getRateLimitCounter returns the rateLimitCounter corresponding to the
// provided rateLimitKey. If no rateLimitCounter exists for the given
// rateLimitKey, then nil is returned.
func (rlm *rateLimitManager) getRateLimitCounter(key rateLimitKey) *rateLimitCounter {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	return rlm.rlcs[key]
}

// getRLC returns the rateLimitCounter corresponding to the provided
// endpointType, rateLimitType, and intervalSeconds.
// If no rateLimitCounter exists for the given endpointType, rateLimitType,
// and intervalSeconds, then a new rateLimitCounter is created with no limit.
// If that fails, then nil is returned.
func (rlm *rateLimitManager) getRLC(
	endpointType common.BIEndpointType, rateLimitType common.BIRateLimitType, intervalSeconds int,
) *rateLimitCounter {

	key := rateLimitKey{endpointType, rateLimitType, intervalSeconds}
	rlc := rlm.getRateLimitCounter(key)
	if rlc == nil {
		fields := log.Fields{"endpointType": endpointType, "rateLimitType": rateLimitType, "intervalSeconds": intervalSeconds}
		rlm.logger.WithFields(fields).Warn("received rate limit update for unknown key. Making new rate limit counter.")

		rateLimit := RateLimit{
			EndpointType:          endpointType,
			RateLimitType:         rateLimitType,
			RateLimitIntervalType: common.IntervalSecond,
			RateLimitIntervalNum:  intervalSeconds,
			Limit:                 -1,
		}

		err := rlm.addRLC(rateLimit)
		if err != nil {
			rlm.logger.WithFields(fields).Error("failed to add rate limit counter. Returning nil.")
			return nil
		}
		rlc = rlm.rlcs[key]
	}
	return rlc
}

// getRLCs returns a slice of rateLimitCounters for the given endpointType
// and rateLimitType. If no rateLimitCounters exist for the given
// endpointType and rateLimitType, then an empty slice is returned.
func (rlm *rateLimitManager) getRLCs(
	endpointType common.BIEndpointType, rateLimitType common.BIRateLimitType,
) []*rateLimitCounter {
	var rlcs []*rateLimitCounter

	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	for key, rlc := range rlm.rlcs {
		if key.endpointType == endpointType && key.rateLimitType == rateLimitType {
			rlcs = append(rlcs, rlc)
		}
	}
	return rlcs
}

// RegisterPending registers pending requests for all rateLimitCounters that
// are associated with the provided ServiceDefinition.
// If registering a request would exceed the limit for any rateLimitCounter,
// all pending counts are rolled back and a RateLimitError is returned.
// TODO: add RAW limits
func (rlm *rateLimitManager) RegisterPending(sd *common.ServiceDefinition) error {
	rlm.logger.WithFields(log.Fields{
		"endpointType": sd.EndpointType,
		"weightIP":     sd.WeightIP,
		"weightUID":    sd.WeightUID,
	}).Debug("RegisterPending")
	countsPending := []countPending{}

	// IP Limits
	for _, rlc := range rlm.getRLCs(sd.EndpointType, common.RateLimitTypeIP) {
		if err := rlc.IncrementPending(sd.WeightIP); err != nil {
			rlm.unregisterPending(countsPending)
			return err
		}
		countsPending = append(countsPending, countPending{rlc, sd.WeightIP})
	}

	// UID Limits
	for _, rlc := range rlm.getRLCs(sd.EndpointType, common.RateLimitTypeUID) {
		if err := rlc.IncrementPending(sd.WeightUID); err != nil {
			rlm.unregisterPending(countsPending)
			return err
		}
		countsPending = append(countsPending, countPending{rlc, sd.WeightUID})
	}

	return nil
}

// unregisterPending unregisters pending requests for rlcsPending.
func (rlm *rateLimitManager) unregisterPending(countsPending []countPending) {
	for _, cp := range countsPending {
		cp.rlc.DecrementPending(cp.count)
	}
}

// UnregisterPending unregisters pending requests for all rateLimitCounters
// that are associated with the provided ServiceDefinition.
// TODO: add RAW limits
func (rlm *rateLimitManager) UnregisterPending(sd *common.ServiceDefinition) {
	countsPendingIP := newPendingCounts(rlm.getRLCs(sd.EndpointType, common.RateLimitTypeIP), sd.WeightIP)
	countsPendingUID := newPendingCounts(rlm.getRLCs(sd.EndpointType, common.RateLimitTypeUID), sd.WeightUID)
	rlm.unregisterPending(countsPendingIP)
	rlm.unregisterPending(countsPendingUID)
}

// updateUsed updates the used count for all rateLimits that receive updates
// in a binance response header. tssResp is the timestamp in nanoseconds
// when the response was generated by the server (e.g. response Date Header).
// If not available, then tssResp should be the timestamp in nanoseconds when
// the request was made.
// TODO: think this through more
func (rlm *rateLimitManager) UpdateUsed(rateLimitUpdates []common.RateLimitUpdate, tssResp common.TSNano) {
	for _, rlu := range rateLimitUpdates {
		rlc := rlm.getRLC(rlu.EndpointType, rlu.RateLimitType, rlu.IntervalSeconds)
		if rlc != nil {
			rlc.SetUsed(rlu.Count, tssResp)
		}
	}
}
