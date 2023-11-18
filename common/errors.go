package common

import (
	"fmt"
	"time"
)

/* ==================== Errors =========================================== */

// UnexpectedStatusCodeError is an error returned when the server returns a
// status code that is not expected.
// ErrorCode and Msg may be empty.
type UnexpectedStatusCodeError struct {
	StatusCode int
	ErrorCode  int
	Msg        string
}

func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected Status Code (code: %d, msg: %s)", e.StatusCode, e.Msg)
}

// RetryAfterError is an error returned when the server returns either a
// 418 (IP Ban) or 429 (Backoff) status code.
// A rateLimitManager can also return this error, if executing a request
// would exceed the rate limit.
//
// Fields:
// StatusCode:     the status code returned by the server
// RetryTimeLocal: the local time at which the request can be retried
// RetryAfter:     the number of seconds to wait before retrying
// Producer:       who produced the error (e.g. "server", "shrimpy-binance")
// Reason:         the reason for the error (e.g. "IP ban", "backoff")
type RetryAfterError struct {
	StatusCode     int
	ErrorCode      int
	Msg            string
	Producer       string
	RetryTimeLocal time.Time // local time
	RetryAfter     int
}

func (e *RetryAfterError) Error() string {
	return fmt.Sprintf("Retry after %d seconds @ %s local time", e.RetryAfter, e.RetryTimeLocal)
}

// BadRequestError is an error returned when the server returns a 400 status code.
// It contains the error code and message returned by the server.
// The error code should match either a BISpotMarginErrorCode or BIFuturesErrorCode,
// depending on the API used.
// Bad Status Code Error Messages Examples:
// * e.g. Filter failure
// * 400 -> {"code":-1013,"msg":"Filter failure: PERCENT_PRICE_BY_SIDE"}
// * 400 -> {"code":-1013,"msg":"Filter failure: NOTIONAL"}
// * e.g. Mandatory parameter not sent/ Too many parameters
// * 400 -> {"code":-1102,"msg":"Mandatory parameter 'quantity' was not sent, was empty/null, or malformed."}
// * 400 -> {"code":-1101,"msg":"Too many parameters; expected '0' and received '2'."}
// * e.g. Failed signature
// * 400 -> {"code":-1022,"msg":"Signature for this request is not valid."}
// * e.g. insufficient balance
// * 400 -> {"code":-2010,"msg":"Account has insufficient balance for requested action."}
// * e.g. Invalid symbol / Invalid characters in symbol
// * 400 -> {"code":-1121,"msg":"Invalid symbol."}
// * 400 -> {"code":-1100,"msg":"Illegal characters found in parameter 'symbols'; legal range is '^\\[(\"[A-Z0-9-_.]{1,20}\"(,\"[A-Z0-9-_.]{1,20}\"){0,}){0,1}\\]$'."}
// * e.g. API Key Problems
// * 400 -> {"code":-2015,"msg":"Invalid API-key, IP, or permissions for action."}
// * 401 -> {"code":-2014,"msg":"API-key format invalid."}
// * e.g. CancelOrder
// * 400 -> {"code":-2011,"msg":"Unknown order sent."} (trying to cancel an order that doesnot exist)
// * e.g. Duplicate order sent (2 orders with same client order id)
// * 400 -> {"code":-2010,"msg":"Duplicate order sent."}
// * e.g. ListenKey does not exist
// * 400 -> {"code":-1105,"msg":"Parameter 'listenKey' was empty."}
// * 400 -> {"code":-1105,"msg":"This listenKey does not exist."}
// * 400 -> {"code":-1125,"msg":"This listenKey does not exist."}
// * 400 -> {"code":-1100,"msg":"Illegal characters found in parameter 'listenKey'; legal range is '^[a-zA-Z0-9]{1,60}$'."}
type BadRequestError struct {
	StatusCode int
	ErrorCode  int    `json:"code"`
	Msg        string `json:"msg"`
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request (code: %d, msg: %s)", e.ErrorCode, e.Msg)
}
