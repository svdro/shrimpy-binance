package common

import (
	"net/http"
)

/* ==================== Symbol & Asset & TSNano, etc ... ================= */
/*
I want to have custom types for these.
I want to always unmarshal timestamps to TSNano
Ideally I would like to have an exchange symbol mapping for symbols
Ideally a similar mapping for assets
*/

/* ==================== Client & Common ================================== */

type BIDataSource int
type BIEndpoint string
type BIEndpointType string
type BIRateLimitIntervalType string
type BISecurityType int
type BIRateLimitType string
type BIHttpResponseCode int
type BIWSEndpoint string
type BIWSSecurityType int

const (
	DataSourceNone BIDataSource = iota
	DataSourceMatchingEngine
	DataSourceMemory
	DataSourceDatabase

	EndpointAPI  BIEndpoint = "api.binance.com"
	EndpointFAPI BIEndpoint = "fapi.binance.com"

	EndpointTypeAPI  BIEndpointType = "api"
	EndpointTypeSAPI BIEndpointType = "sapi"
	EndpointTypeFAPI BIEndpointType = "fapi"

	HTTPStatusOK              BIHttpResponseCode = http.StatusOK              // (200)
	HTTPStatusTeapot          BIHttpResponseCode = http.StatusTeapot          // (418) IP ban
	HTTPStatusTooManyRequests BIHttpResponseCode = http.StatusTooManyRequests // (429) Backoff
	HTTPStatusBadRequest      BIHttpResponseCode = http.StatusBadRequest      // (400) Invalid request
	HTTPStatusUnauthorized    BIHttpResponseCode = http.StatusUnauthorized    // (401) e.g invalid format

	IntervalSecond BIRateLimitIntervalType = "SECOND"
	IntervalMinute BIRateLimitIntervalType = "MINUTE"
	IntervalDay    BIRateLimitIntervalType = "DAY"

	SecurityTypeNone   BISecurityType = iota // NONE
	SecurityTypeApiKey                       // USER_STREAM, MARKET_DATA
	SecurityTypeSigned                       // TRADE, MARGIN, USER_DATA

	RateLimitTypeIP  BIRateLimitType = "IP"
	RateLimitTypeUID BIRateLimitType = "UID"
	RateLimitTypeRAW BIRateLimitType = "RAW"

	WSEndpointAPI    BIWSEndpoint = "stream.binance.com:9443"
	WSEndpointFAPI   BIWSEndpoint = "fstream.binance.com"
	WSAPIEndpointAPI BIWSEndpoint = "ws-api.binance.com:443/ws-api/v3"

	WSSecurityTypeNone BIWSSecurityType = iota
	WSSecurityTypeListenKey
)

/* ==================== Order ============================================ */

type BIOrderSide string               // (SPOT & MARGIN & FUTURES)
type BIOrderType string               // (SPOT & MARGIN & FUTURES)
type BIOrderResponseType string       // (SPOT & MARGIN & FUTURES)
type BIOrderTimeInForce string        // (SPOT & MARGIN & FUTURES)
type BIOrderStatus string             // (SPOT & MARGIN & FUTURES)
type BISelfTradePreventionMode string // (SPOT & MARGIN & FUTURES)
type BIOrderSideEffect string         //  (MARGIN)
type BIExecutionType string           // (SPOT & MARGIN & FUTURES)

const (
	OrderSideBuy  BIOrderSide = "BUY"  // (SPOT & MARGIN & FUTURES)
	OrderSideSell BIOrderSide = "SELL" // (SPOT & MARGIN & FUTURES)

	OrderTypeLimit              BIOrderType = "LIMIT"                // (SPOT_MARGIN & FUTURES)
	OrderTypeMarket             BIOrderType = "MARKET"               // (SPOT_MARGIN & FUTURES)
	OrderTypeStopLoss           BIOrderType = "STOP_LOSS"            // (SPOT_MARGIN)
	OrderTypeStopLossLimit      BIOrderType = "STOP_LOSS_LIMIT"      // (SPOT_MARGIN)
	OrderTypeTakeProfit         BIOrderType = "TAKE_PROFIT"          // (SPOT_MARGIN & FUTURES)
	OrderTypeTakeProfitLimit    BIOrderType = "TAKE_PROFIT_LIMIT"    // (SPOT_MARGIN)
	OrderTypeLimitMaker         BIOrderType = "LIMIT_MAKER"          // (SPOT_MARGIN)
	OrderTypeStop               BIOrderType = "STOP"                 // (FUTURES)
	OrderTypeStopMarket         BIOrderType = "STOP_MARKET"          // (FUTURES)
	OrderTypeTakeProfitMarket   BIOrderType = "TAKE_PROFIT_MARKET"   // (FUTURES)
	OrderTypeTrailingStopMarket BIOrderType = "TRAILING_STOP_MARKET" // (FUTURES)

	OrderResponseTypeAcknowledge BIOrderResponseType = "ACK"    // (SPOT & MARGIN & FUTURES)
	OrderResponseTypeResult      BIOrderResponseType = "RESULT" // (SPOT & MARGIN & FUTURES)
	OrderResponseTypeFull        BIOrderResponseType = "FULL"   // (SPOT & MARGIN)

	OrderTimeInForceGTC BIOrderTimeInForce = "GTC" // (SPOT & FUTURES) (Good-Till-Canceled)
	OrderTimeInForceIOC BIOrderTimeInForce = "IOC" // (SPOT & FUTURES) (Immediate-or-Cancel)
	OrderTimeInForceFOK BIOrderTimeInForce = "FOK" // (SPOT & FUTURES) (Fill-or-Kill)
	OrderTimeInForceGTX BIOrderTimeInForce = "GTX" // (FUTURES)        (Good-Till-Executed)
	OrderTimeInForceGTD BIOrderTimeInForce = "GTD" // (FUTURES)        (Good-Till-Date)

	OrderStatusNew             BIOrderStatus = "NEW"              // (SPOT_MARGIN & FUTURES) Order has been accepted by the engine
	OrderStatusPartiallyFilled BIOrderStatus = "PARTIALLY_FILLED" // (SPOT_MARGIN & FUTURES) Order has been partially filled
	OrderStatusFilled          BIOrderStatus = "FILLED"           // (SPOT_MARGIN & FUTURES) Order has been filled
	OrderStatusCanceled        BIOrderStatus = "CANCELED"         // (SPOT_MARGIN & FUTURES) Order has been canceled
	OrderStatusPendingCancel   BIOrderStatus = "PENDING_CANCEL"   // (SPOT_MARGIN) Currently unused
	OrderStatusRejected        BIOrderStatus = "REJECTED"         // (SPOT_MARGIN & FUTURES) Order not accepted by the engine
	OrderStatusExpired         BIOrderStatus = "EXPIRED"          // (SPOT_MARGIN & FUTURES) Order expired
	OrderStatusExpiredInMatch  BIOrderStatus = "EXPIRED_IN_MATCH" // (SPOT_MARGIN & FUTURES) Order expired in match engine

	SelfTradePreventionModeExpireMaker BISelfTradePreventionMode = "EXPIRE_MAKER" // (SPOT & MARGIN & FUTURES)
	SelfTradePreventionModeExpireTaker BISelfTradePreventionMode = "EXPIRE_TAKER" // (SPOT & MARGIN & FUTURES)
	SelfTradePreventionModeExpireBoth  BISelfTradePreventionMode = "EXPIRE_BOTH"  // (SPOT & MARGIN & FUTURES)
	SelfTradePreventionModeNone        BISelfTradePreventionMode = "NONE"         // (SPOT & MARGIN & FUTURES)

	SideEffectNoSideEffect BIOrderSideEffect = "NO_SIDE_EFFECT" // (MARGIN)
	SideEffectMarginBuy    BIOrderSideEffect = "MARGIN_BUY"     // (MARGIN)
	SideEffectAutoRepay    BIOrderSideEffect = "AUTO_REPAY"     // (MARGIN)

	ExecutionTypeNew             BIExecutionType = "NEW"              // (SPOT & MARGIN & FUTURES)
	ExecutionTypeCanceled        BIExecutionType = "CANCELED"         // (SPOT & MARGIN & FUTURES)
	ExecutionTypeReplaced        BIExecutionType = "REPLACED"         // (SPOT & MARGIN) (not currently used)
	ExecutionTypeTrade           BIExecutionType = "TRADE"            // (SPOT & MARGIN & FUTURES)
	ExecutionTypeExpired         BIExecutionType = "EXPIRED"          // (SPOT & MARGIN & FUTURES)
	ExecutionTypeRejected        BIExecutionType = "REJECTED"         // (SPOT & MARGIN)
	ExecutionTypeTradePrevention BIExecutionType = "TRADE_PREVENTION" // (SPOT & MARGIN)
	ExecutionTypeAmendment       BIExecutionType = "AMENDMENT"        // (FUTURES)
	ExecutionTypeCalculated      BIExecutionType = "CALCULATED"       // (FUTURES)
)
