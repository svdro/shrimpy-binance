package binance

import (
	"fmt"
	"time"
)

/* ==================== Errors =========================================== */

// RetryAfterError is an error returned when the server returns either a
// 418 (IP Ban) or 429 (Backoff) status code.
// rateLimitHandlers can also return this error, if the additional of a new
// request would exceed the rate limit.
type RetryAfterError struct {
	RetryAfter     int
	RetryTimeLocal time.Time // local time
	//RetryTimeServer time.Time // server time
	//StatusCode      BIHttpResponseCode
	//Reason          string
}

func (e *RetryAfterError) Error() string {
	return fmt.Sprintf("Retry after %d seconds @ %s local time", e.RetryAfter, e.RetryTimeLocal)
}
