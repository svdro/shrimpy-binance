package client

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== parseServiceResponseHeader ======================= */

var (
	reIPLimit        = regexp.MustCompile(`x-mbx-used-weight-\d+[a-z]`)
	reUIDLimit       = regexp.MustCompile(`x-mbx-order-count-\d+[a-z]`)
	reSapiIPLimit    = regexp.MustCompile(`x-sapi-used-ip-weight-\d+[a-z]`)
	reSapiUIDLimit   = regexp.MustCompile(`x-sapi-used-uid-weight-\d+[a-z]`)
	reIntervalNumber = regexp.MustCompile(`\d`)
	reIntervalLetter = regexp.MustCompile(`[a-z]$`)
)

// intervalSecondsMap is used for calclulating nanoseconds from header.
var intervalSecondsMap = map[common.BIRateLimitIntervalType]int{
	common.IntervalSecond: 1,
	common.IntervalMinute: 60,
	common.IntervalDay:    60 * 60 * 24,
}

// intervalTypeMap is used for parsing intervalLetter from header.
var intervalTypeMap = map[string]common.BIRateLimitIntervalType{
	"s": common.IntervalSecond,
	"m": common.IntervalMinute,
	"d": common.IntervalDay,
}

// parseRateLimitHeader parses the response header key-value pair, and
// returns a common.RateLimitHeader of BIRateLimitType rateLimitType.
func parseRateLimitHeader(
	k, v string, rateLimitType common.BIRateLimitType,
) (*common.RateLimitHeader, error) {

	// parse limit from header value
	limit, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}

	// parseIntervalType from header key
	intervalLetter := reIntervalLetter.FindString(k)
	if intervalLetter == "" {
		err := fmt.Errorf("error parsing intervalLetter from key %s", k)
		return nil, err
	}

	// parse intervalNumber from header key
	intervalNumberStr := reIntervalNumber.FindString(k)
	intervalNumber, err := strconv.Atoi(intervalNumberStr)
	if err != nil {
		return nil, err
	}

	// calculate interval nanoseconds
	intervalTypeSeconds := intervalSecondsMap[intervalTypeMap[intervalLetter]]
	intervalNanoSeconds := int64(intervalTypeSeconds * intervalNumber * 1e9)

	return &common.RateLimitHeader{
		RateLimitType:       rateLimitType,
		IntervalType:        intervalTypeMap[intervalLetter],
		IntervalNum:         intervalNumber,
		IntervalNanoSeconds: intervalNanoSeconds,
		Count:               limit,
	}, nil

}

// parseServiceResponseHeader parses the http response header and returns a
// common.ServiceResponseHeader.
func parseServiceResponseHeader(h http.Header) (*common.ServiceResponseHeader, error) {
	logger := log.WithField("caller", "parseServiceResponseHeader")
	IPLimits := map[int64]common.RateLimitHeader{}
	UIDLimits := map[int64]common.RateLimitHeader{}

	// loop over key value pairs to parse IP & UID limits
	// this is necessary, because the header keys are not always known.
	for k := range h {
		value := h.Get(k)
		key := strings.ToLower(k)
		logger := logger.WithFields(log.Fields{"key": key, "value": value})

		switch {
		case reIPLimit.MatchString(key) || reSapiIPLimit.MatchString(key):
			rlHeader, err := parseRateLimitHeader(key, value, common.RateLimitTypeIP)
			if err != nil {
				logger.WithError(err).Fatal("error parsing IPLimit")
				continue
			}
			IPLimits[rlHeader.IntervalNanoSeconds] = *rlHeader
		case reUIDLimit.MatchString(key) || reSapiUIDLimit.MatchString(key):
			rlHeader, err := parseRateLimitHeader(key, value, common.RateLimitTypeUID)
			if err != nil {
				logger.WithError(err).Fatal("error parsing UIDLimit")
				continue
			}
			UIDLimits[rlHeader.IntervalNanoSeconds] = *rlHeader
		}
	}

	// get server header
	server := h.Get("Server")
	if server == "" {
		err := fmt.Errorf("error parsing server header")
		return nil, err
	}

	// get date header and parse time
	dateValue := h.Get("Date")
	date, err := time.Parse(time.RFC1123, dateValue)
	if err != nil {
		return nil, err
	}

	serverResponseHeader := &common.ServiceResponseHeader{
		Server:        server,
		TSSRespHeader: date.UnixNano(),
		IPLimits:      IPLimits,
		UIDLimits:     UIDLimits,
	}

	return serverResponseHeader, nil
}
