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

// intervalSecondsMap is used for calclulating seconds from header.
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

// getSecondsInInterval returns the number of seconds in the given interval.
func getSecondsInInterval(intervalType common.BIRateLimitIntervalType, intervalNumber int) int {
	intervalTypeSeconds := intervalSecondsMap[intervalType]
	intervalSeconds := intervalTypeSeconds * intervalNumber
	return intervalSeconds
}

var (
	reIPLimit        = regexp.MustCompile(`x-mbx-used-weight-(\d+)[a-z]`)
	reUIDLimit       = regexp.MustCompile(`x-mbx-order-count-(\d+)[a-z]`)
	reSapiIPLimit    = regexp.MustCompile(`x-sapi-used-ip-weight-\d+[a-z]`)
	reSapiUIDLimit   = regexp.MustCompile(`x-sapi-used-uid-weight-\d+[a-z]`)
	reIntervalNumber = regexp.MustCompile(`(\d+)`)
	reIntervalLetter = regexp.MustCompile(`[a-z]$`)
)

// parseRateLimitHeader parses the response header key-value pair, and
// returns a common.RateLimitUpdate of BIRateLimitType rateLimitType.
func parseRateLimitHeader(
	k, v string, rateLimitType common.BIRateLimitType, endpointType common.BIEndpointType,
) (*common.RateLimitUpdate, error) {

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

	// parse interval type and calculate interval nanoseconds
	intervalType := intervalTypeMap[intervalLetter]
	intervalSeconds := getSecondsInInterval(intervalType, intervalNumber)

	return &common.RateLimitUpdate{
		EndpointType:    endpointType,
		RateLimitType:   rateLimitType,
		IntervalSeconds: intervalSeconds,
		Count:           limit,
	}, nil

}

// parseServiceResponseHeader parses the http response header and returns a
// common.ServiceResponseHeader.
func parseServiceResponseHeader(
	h http.Header, endpointType common.BIEndpointType,
) (*common.ServiceResponseHeader, error) {
	logger := log.WithField("_caller", "parseServiceResponseHeader")
	rateLimitUpdates := []common.RateLimitUpdate{}

	// loop over key value pairs to parse IP & UID limits
	// this is necessary, because the header keys are not always known.
	for k := range h {
		value := h.Get(k)
		key := strings.ToLower(k)
		logger := logger.WithFields(log.Fields{"key": key, "value": value})

		var err error
		var rlHeader *common.RateLimitUpdate

		switch {
		case reIPLimit.MatchString(key) || reSapiIPLimit.MatchString(key):
			rlHeader, err = parseRateLimitHeader(key, value, common.RateLimitTypeIP, endpointType)
		case reUIDLimit.MatchString(key) || reSapiUIDLimit.MatchString(key):
			rlHeader, err = parseRateLimitHeader(key, value, common.RateLimitTypeUID, endpointType)
		}

		if err != nil {
			logger.WithError(err).Error("error parsing rate limit header")
			return nil, err
		}

		if rlHeader != nil {
			rateLimitUpdates = append(rateLimitUpdates, *rlHeader)
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

	// make serviceResponseHeader
	serviceResponseHeader := &common.ServiceResponseHeader{
		Server:           server,
		TSSRespHeader:    date.UnixNano(),
		RateLimitUpdates: rateLimitUpdates,
	}

	// get retry after header, if included in header
	if retryAfter, err := strconv.Atoi(h.Get("Retry-After")); err == nil {
		serviceResponseHeader.RetryAfter = &retryAfter
	}

	return serviceResponseHeader, nil
}
