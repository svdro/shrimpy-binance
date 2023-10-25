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

/* ==================== Examples of error Codes ========================== */
/*
Bad Status Code Error Messages Examples:
	* e.g. Filter failure
		* 400 -> {"code":-1013,"msg":"Filter failure: PERCENT_PRICE_BY_SIDE"}
		* 400 -> {"code":-1013,"msg":"Filter failure: NOTIONAL"}
	* e.g. Mandatory parameter not sent/ Too many parameters
		* 400 -> {"code":-1102,"msg":"Mandatory parameter 'quantity' was not sent, was empty/null, or malformed."}
		* 400 -> {"code":-1101,"msg":"Too many parameters; expected '0' and received '2'."}
	* e.g. Failed signature
		* 400 -> {"code":-1022,"msg":"Signature for this request is not valid."}
	* e.g. insufficient balance
		* 400 -> {"code":-2010,"msg":"Account has insufficient balance for requested action."}
	* e.g. Invalid symbol / Invalid characters in symbol
		* 400 -> {"code":-1121,"msg":"Invalid symbol."}
		* 400 -> {"code":-1100,"msg":"Illegal characters found in parameter 'symbols'; legal range is '^\\[(\"[A-Z0-9-_.]{1,20}\"(,\"[A-Z0-9-_.]{1,20}\"){0,}){0,1}\\]$'."}
	* e.g. API Key Problems
		* 400 -> {"code":-2015,"msg":"Invalid API-key, IP, or permissions for action."}
	* e.g. CancelOrder
	  * 400 -> {"code":-2011,"msg":"Unknown order sent."} (trying to cancel an order that doesnot exist)
	* e.g. Duplicate order sent (2 orders with same client order id)
	  * 400 -> {"code":-2010,"msg":"Duplicate order sent."}
	* e.g. ListenKey does not exist
    * 400 -> {"code":-1105,"msg":"This listenKey does not exist."}
*/

//// BadRequestError is an error returned when the server returns a 400 status code.
//// It contains the error code and message returned by the server.
//// The error code should match either a BISpotMarginErrorCode or BIFuturesErrorCode,
//// depending on the API used.
//type BadRequestError struct {
//StatusCode HTTPResponseCode
//ErrorCode  int    `json:"code"`
//Msg        string `json:"msg"`
//}

//func (e *BadRequestError) Error() string {
//return fmt.Sprintf("Bad Request (code: %d, msg: %s)", e.ErrorCode, e.Msg)
//}
