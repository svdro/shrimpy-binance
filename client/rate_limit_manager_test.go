package client

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/svdro/shrimpy-binance/common"
)

func TestAddAndGetRateLimitCounter(t *testing.T) {
	th := mockTimeHandler{}
	logger := log.NewEntry(log.StandardLogger())
	rlm := newRateLimitManager([]RateLimit{ /* ... */ }, &th, logger)

	// Define RateLimits to test
	rls := []RateLimit{
		{common.EndpointTypeAPI, common.RateLimitTypeIP, common.IntervalMinute, 1, 6000},
		{common.EndpointTypeAPI, common.RateLimitTypeIP, common.IntervalMinute, 60, 50000},
		{common.EndpointTypeAPI, common.RateLimitTypeUID, common.IntervalSecond, 10, 1200},
		{common.EndpointTypeSAPI, common.RateLimitTypeIP, common.IntervalMinute, 1, 1200},
		{common.EndpointTypeAPI, common.RateLimitTypeRAW, common.IntervalSecond, 1, 1200},
	}

	for i, rl := range rls {
		testName := fmt.Sprintf("i=%d, EndpointType=%s, RateLimitType=%s, IntervalType=%s, IntervalNum=%d, Limit=%d",
			i, rl.EndpointType, rl.RateLimitType, rl.RateLimitIntervalType, rl.RateLimitIntervalNum, rl.Limit)

		t.Run(testName, func(t *testing.T) {
			// add the rate limit counter
			rlm.addRLC(rl)

			// retrieve the rate limit counter, assert is not nil
			intervalSeconds := getSecondsInInterval(rl.RateLimitIntervalType, rl.RateLimitIntervalNum)
			rlc := rlm.getRLC(rl.EndpointType, rl.RateLimitType, intervalSeconds)
			assert.NotNil(t, rlc, "RateLimitCounter should not be nil")

			// assert that the rate limit counter has the correct values
			assert.Equal(t, rl.EndpointType, rlc.endpointType)
			assert.Equal(t, rl.RateLimitType, rlc.rateLimitType)
			assert.Equal(t, intervalSeconds, rlc.intervalSeconds)
			assert.Equal(t, rl.Limit, rlc.limit)
		})
	}
}

// I want to test updateing multiple rate limit counters at once
// I want to test updating a rate limit counter that does not exist
func TestUpdateUsed(t *testing.T) {
	// make mock time handler (not used in this test) & rate limit manager
	rls := []RateLimit{
		{common.EndpointTypeAPI, common.RateLimitTypeIP, common.IntervalMinute, 1, 6000},
		{common.EndpointTypeAPI, common.RateLimitTypeUID, common.IntervalSecond, 10, 1200},
	}
	th := &mockTimeHandler{tsl: 0, offset: 0}
	rlm := newRateLimitManager(rls, th, log.NewEntry(log.StandardLogger()))

	// Define RateLimitUpdates to test
	// 1. update a IP rate limit counter that exists
	// 2. update a IP rate limit counter with different interval that does not exist
	// 3. update a UID rate limit counter that exists
	updates := []common.RateLimitUpdate{
		{EndpointType: common.EndpointTypeAPI, RateLimitType: common.RateLimitTypeIP, IntervalSeconds: 60, Count: 300},
		{EndpointType: common.EndpointTypeAPI, RateLimitType: common.RateLimitTypeRAW, IntervalSeconds: 3600, Count: 300},
		{EndpointType: common.EndpointTypeAPI, RateLimitType: common.RateLimitTypeUID, IntervalSeconds: 10, Count: 100},
	}

	rlm.UpdateUsed(updates, 0)

	// retrieve rate limit counters, assert not nil, and very values
	for i, update := range updates {
		name := fmt.Sprintf("i=%d, update=%v", i, update)
		t.Run(name, func(t *testing.T) {
			rlc := rlm.getRLC(update.EndpointType, update.RateLimitType, update.IntervalSeconds)
			assert.NotNil(t, rlc, "RateLimitCounter should not be nil")
			assert.Equal(t, update.Count, rlc.countUsed)
		})
	}
}
