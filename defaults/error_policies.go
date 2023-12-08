package defaults

import (
	"net"
	"time"

	bc "github.com/svdro/shrimpy-binance/common"
)

/* ==================== Error Policy ===================================== */

// ErrorPolicy is the base interface for all error policies.
// an implementation of errror policy should check service errors and do 3 things:
// determine whether the error is transient or not,
// calculate the next interval to wait before retrying,
// give the reason for the error.
type ErrorPolicy interface {
	Handle(err error) (bool, time.Duration, string)
}

/* ==================== Default Error Policy ============================== */

// defaultErrorPolicy is the default implementation of ErrorPolicy.
type defaultErrorPolicy struct {
	currRetryTime time.Time
}

func (p *defaultErrorPolicy) CalcNextInterval(retryAfterTime time.Time) time.Duration {
	if retryAfterTime.Before(p.currRetryTime) {
		return 0
	}

	p.currRetryTime = retryAfterTime
	timeToWait := retryAfterTime.Sub(time.Now())
	return timeToWait
}

// Handle checks the error and returns a bool signaling whether the error is
// transient and the loop should continue, a time.Duration that is the time
// to wait before retrying, and a string message that can be logged (a reason).
// TODO: count consec errors before within a certain duration, and exit if
//
//	the count exceeds a certain threshold.
func (p *defaultErrorPolicy) Handle(err error) (bool, time.Duration, string) {
	switch err.(type) {

	// RateLimitError
	case *bc.RateLimitError:
		retryTime := err.(*bc.RateLimitError).RetryTimeLocal
		interval := p.CalcNextInterval(retryTime)
		return true, interval, "rate limit exceeded: throttling"

	// BadRequestError
	case *bc.BadRequestError:
		return false, 0, "bad request: exiting"

	// UnexpectedStatusCodeError
	case *bc.UnexpectedStatusCodeError:
		return false, 0, "unexpected status code: exiting"

	// TODO: what do I do with these?
	case net.Error:
		if err.(net.Error).Timeout() {
			return false, time.Duration(0), "timeout: exiting"
		}
		return false, time.Duration(0), "network error: exiting"
	}

	return false, time.Duration(0), "unknown error: exiting"
}

/* ==================== Default WS Error Policy =========================== */

type defaultWSErrorPolicy struct {
}

func (p *defaultWSErrorPolicy) Handle(err error) (bool, time.Duration, string) {
	switch err.(type) {
	case *bc.WSConnError:
		if err.(*bc.WSConnError).IsTransient {
			return true, time.Duration(0), "websocket connection error: retrying"
		}
		return false, time.Duration(0), "websocket connection error: exiting"

	case *bc.WSHandlerError:
		return false, time.Duration(0), "websocket handler error: exiting"
	}
	return false, time.Duration(0), "websocket error: exiting"
}
