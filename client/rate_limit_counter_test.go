package client

import (
	"fmt"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/svdro/shrimpy-binance/common"
)

// mockTimeHandler is a mock implementation of common.TimeHandler for testing.
// it allows for the time to be set and changed manually as needed.
type mockTimeHandler struct {
	offset int64
	tsl    common.TSNano
}

func (th *mockTimeHandler) TSLNow() common.TSNano {
	return th.tsl
}

func (th *mockTimeHandler) TSSNow() common.TSNano {
	return th.TSLToTSS(th.TSLNow())
}

func (th *mockTimeHandler) TSLToTSS(tsl common.TSNano) common.TSNano {
	return common.TSNano(tsl.Int64() - th.offset)
}

func (th *mockTimeHandler) TSSToTSL(tss common.TSNano) common.TSNano {
	return common.TSNano(tss.Int64() + th.offset)
}

func (th *mockTimeHandler) Offset() int64 {
	return th.offset
}

func (th *mockTimeHandler) SetTSL(tsl common.TSNano) {
	th.tsl = tsl
}

func TestIncrementAndDecrementPending(t *testing.T) {
	tsl := 1700080339 * 1e9 // 19 seconds into the minute
	th := &mockTimeHandler{tsl: common.TSNano(tsl), offset: 0}
	rlc := newRateLimitCounter(th, common.EndpointTypeAPI, common.RateLimitTypeIP, 60, 6000, log.NewEntry(log.StandardLogger()))

	// Test incrementing countPending
	err := rlc.IncrementPending(3)
	assert.Nil(t, err)
	assert.Equal(t, 3, rlc.countPending)

	// Test decrementing countPending
	rlc.DecrementPending(2)
	assert.Equal(t, 1, rlc.countPending)
}

type testRetryErrorGenerateionTestCase struct {
	limit     int           // rate limit, -1 for no limit
	pending   int           // pending weight to be added
	tsl0      common.TSNano // current time
	tsl1      common.TSNano // first timestamp in next interval
	shouldErr bool          // whether or not a RetryAfterError is expected
}

func TestRetryErrorGeneration(t *testing.T) {
	testCases := []testRetryErrorGenerateionTestCase{
		{limit: 5, pending: 6, tsl0: 1700080339 * 1e9, tsl1: 1700080380 * 1e9, shouldErr: true},
		{limit: 2, pending: 3, tsl0: 1700080380 * 1e9, tsl1: 1700080440 * 1e9, shouldErr: true},
		{limit: 6000, pending: 30, tsl0: 0, tsl1: 0, shouldErr: false}, // tsl not used if !shouldErr
		{limit: -1, pending: 2000, tsl0: 0, tsl1: 0, shouldErr: false}, // tsl not used if !shouldErr
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("i: %d, limit: %d, pending: %d, tsl0: %d, tsl1: %d", i, tc.limit, tc.pending, tc.tsl0, tc.tsl1)
		t.Run(name, func(t *testing.T) {
			th := &mockTimeHandler{tsl: common.TSNano(tc.tsl0), offset: 0}
			rlc := newRateLimitCounter(th, common.EndpointTypeAPI, common.RateLimitTypeIP, 60, tc.limit, log.NewEntry(log.StandardLogger()))

			// increment pending
			err := rlc.IncrementPending(tc.pending)

			// Assert the error, lack of error
			if !tc.shouldErr {
				assert.Nil(t, err)
				return
			}
			assert.NotNil(t, err)

			// Assert the type and the content of the error
			deltaSeconds := int((tc.tsl1 - tc.tsl0) / 1e9)
			if retryErr, ok := err.(*common.RetryAfterError); ok {
				assert.Equal(t, "shrimpy-binance", retryErr.Producer)
				assert.Equal(t, retryErr.RetryAfter, deltaSeconds)
				assert.Equal(t, retryErr.RetryTimeLocal.UnixNano(), int64(tc.tsl1))
			} else {
				t.Errorf("Expected error type *common.RetryAfterError, got %T", err)
			}
		})
	}
}

type testIntervalTransitionTestCase struct {
	tsl             common.TSNano
	countUsed       int
	targetInterval  int64
	targetCountUsed int
	mockSleepTime   int64
}

func TestIntervalTransition(t *testing.T) {
	// setup the test cases
	tsl0 := common.TSNano(1700080339 * 1e9) // 19 seconds into the minute (current time)
	tsl1 := common.TSNano(1700080380 * 1e9) // 0 seconds into the next minute (start time of next interval)
	testCases := []testIntervalTransitionTestCase{
		{tsl: tsl0, countUsed: 300, targetInterval: int64(tsl0 / (60 * 1e9)), targetCountUsed: 300}, // setUsed (t0)-> should work
		{tsl: tsl0, countUsed: 5, targetInterval: int64(tsl0 / (60 * 1e9)), targetCountUsed: 300},   // setUsed (t0) -> should fail
		{tsl: tsl1, countUsed: 5, targetInterval: int64(tsl1 / (60 * 1e9)), targetCountUsed: 5},     // setUsed (t1) -> should work
		{tsl: tsl0, countUsed: 200, targetInterval: int64(tsl1 / (60 * 1e9)), targetCountUsed: 5},   // setUsed (t0) -> should fail
		{tsl: tsl1, countUsed: 300, targetInterval: int64(tsl1 / (60 * 1e9)), targetCountUsed: 300}, // setUsed (t1) -> should work
	}

	// initialize the rate limit counter.
	// assert that the value of currInterval is 0 after initialization.
	// NOTE: th is not used in this test, as we are manually setting the tsl in setUsed.
	th := &mockTimeHandler{tsl: common.TSNano(0), offset: int64(0)}
	rlc := newRateLimitCounter(th, common.EndpointTypeAPI, common.RateLimitTypeIP, 60, 6000, log.NewEntry(log.StandardLogger()))
	assert.Equal(t, int64(0), rlc.currInterval)

	// run the test cases (test loop)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("i: %d, Interval: %d, CountUsed: %d", i, tc.targetInterval, tc.countUsed), func(t *testing.T) {
			rlc.SetUsed(tc.countUsed, tc.tsl)
			assert.Equal(t, tc.targetCountUsed, rlc.countUsed)
			assert.Equal(t, tc.targetInterval, rlc.currInterval)
		})
	}
}

type testIntervalTransitionConcurrentTestCase struct {
	countUsed     int
	mockSleepTime int64
}

func TestIntervalTransitionConcurrently(t *testing.T) {
	// setup the test cases
	tsl0 := common.TSNano(1700080339 * 1e9) // 19 seconds into the minute (current time)
	testCases := []testIntervalTransitionConcurrentTestCase{
		{countUsed: 45, mockSleepTime: 103}, // setUsed (t4, i2)
		{countUsed: 17, mockSleepTime: 163}, // setUsed (t5, i3)
		{countUsed: 300, mockSleepTime: 0},  // setUsed (t0, i0)
		{countUsed: 2, mockSleepTime: 41},   // setUsed (t2, i1)
		{countUsed: 1, mockSleepTime: 45},   // setUsed (t3, i1)
		{countUsed: 1000, mockSleepTime: 7}, // setUsed (t1, i0)
		{countUsed: 14, mockSleepTime: 101}, // setUsed (t4, i2)
	}
	maxMockSleepTime := 163
	targetInterval := int64((tsl0.Int64() + int64(maxMockSleepTime*1e9)) / (60 * 1e9))
	targetCountUsed := 17

	// initialize the rate limit counter.
	// assert that the value of currInterval is 0 after initialization.
	// NOTE: th is not used in this test, as we are manually setting the tsl in setUsed.
	th := &mockTimeHandler{tsl: common.TSNano(0), offset: int64(0)}
	rlc := newRateLimitCounter(th, common.EndpointTypeAPI, common.RateLimitTypeIP, 60, 6000, log.NewEntry(log.StandardLogger()))
	assert.Equal(t, int64(0), rlc.currInterval)

	wg := &sync.WaitGroup{}
	for i, tc := range testCases {
		// time.Sleep(1 * time.Millisecond)
		wg.Add(1)
		go func(i int, tc testIntervalTransitionConcurrentTestCase) {
			defer wg.Done()
			tsl := common.TSNano(tsl0.Int64() + int64(tc.mockSleepTime*1e9))
			rlc.SetUsed(tc.countUsed, tsl)
		}(i, tc)
	}

	wg.Wait()
	assert.Equal(t, targetCountUsed, rlc.countUsed)
	assert.Equal(t, targetInterval, rlc.currInterval)
}

func TestRateLimitCounter(t *testing.T) {
	tsl0 := common.TSNano(1700080339 * 1e9)             // 19 seconds into the minute (current time)
	tsl1 := common.TSNano(1700080340 * 1e9)             // 20 seconds into the next minute (start time of next interval)
	th := &mockTimeHandler{tsl: tsl0, offset: int64(0)} // current time is tsl0

	// initialize the rate limit counter.
	// assert that the value of currInterval is 0 after initialization.
	rlc := newRateLimitCounter(th, common.EndpointTypeAPI, common.RateLimitTypeIP, 60, 6000, log.NewEntry(log.StandardLogger()))
	assert.Equal(t, int64(0), rlc.currInterval)

	// register pending requests
	err := rlc.IncrementPending(3)
	assert.Nil(t, err)
	assert.Equal(t, 3, rlc.countPending)

	// time passes until response is received, assume TSSResp is (tsl1- tsl0) / 2
	tssResp := (tsl1 - tsl0) / 2
	th.SetTSL(tsl1) // current time is tsl1

	// update countUsed
	rlc.SetUsed(3, tssResp)
	assert.Equal(t, 3, rlc.countUsed)

	// unregister pending requests
	rlc.DecrementPending(3)
	assert.Equal(t, 0, rlc.countPending)
}
