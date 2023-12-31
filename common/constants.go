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

//func (t *BIRateLimitType) UnmarshalJSON(b []byte) error {
//*t = BIRateLimitType(b)
//return nil
//}

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

	//RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
	//RateLimitTypeOrders        RateLimitType = "ORDERS"
	//RateLimitTypeRawRequests   RateLimitType = "RAW_REQUESTS"
	RateLimitTypeIP  BIRateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeUID BIRateLimitType = "ORDERS"
	RateLimitTypeRAW BIRateLimitType = "RAW_REQUESTS"
	//RateLimitTypeIP  BIRateLimitType = "IP"
	//RateLimitTypeUID BIRateLimitType = "UID"
	//RateLimitTypeRAW BIRateLimitType = "RAW"

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

	OrderTypeLimit              BIOrderType = "LIMIT"                // (SPOT & MARGIN & FUTURES)
	OrderTypeMarket             BIOrderType = "MARKET"               // (SPOT & MARGIN & FUTURES)
	OrderTypeStopLoss           BIOrderType = "STOP_LOSS"            // (SPOT & MARGIN)
	OrderTypeStopLossLimit      BIOrderType = "STOP_LOSS_LIMIT"      // (SPOT & MARGIN)
	OrderTypeTakeProfit         BIOrderType = "TAKE_PROFIT"          // (SPOT & MARGIN & FUTURES)
	OrderTypeTakeProfitLimit    BIOrderType = "TAKE_PROFIT_LIMIT"    // (SPOT & MARGIN)
	OrderTypeLimitMaker         BIOrderType = "LIMIT_MAKER"          // (SPOT & MARGIN)
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

	OrderStatusNew             BIOrderStatus = "NEW"              // (SPOT & MARGIN & FUTURES) Order has been accepted by the engine
	OrderStatusPartiallyFilled BIOrderStatus = "PARTIALLY_FILLED" // (SPOT & MARGIN & FUTURES) Order has been partially filled
	OrderStatusFilled          BIOrderStatus = "FILLED"           // (SPOT & MARGIN & FUTURES) Order has been filled
	OrderStatusCanceled        BIOrderStatus = "CANCELED"         // (SPOT & MARGIN & FUTURES) Order has been canceled
	OrderStatusPendingCancel   BIOrderStatus = "PENDING_CANCEL"   // (SPOT & MARGIN) Currently unused
	OrderStatusRejected        BIOrderStatus = "REJECTED"         // (SPOT & MARGIN & FUTURES) Order not accepted by the engine
	OrderStatusExpired         BIOrderStatus = "EXPIRED"          // (SPOT & MARGIN & FUTURES) Order expired
	OrderStatusExpiredInMatch  BIOrderStatus = "EXPIRED_IN_MATCH" // (SPOT & MARGIN & FUTURES) Order expired in match engine

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

/* ==================== ExchangeInfo ===================================== */

type BISymbolStatusType string           // (SPOT & MARGIN)
type BIAccountAndSymbolPermission string // (SPOT & MARGIN)
type BIContractType string               //(FUTURES)
type BIContractStatus string             // (FUTURES)

const (
	SymbolStatusPreTrading   BISymbolStatusType = "PRE_TRADING"   // (SPOT & MARGIN)
	SymbolStatusTrading      BISymbolStatusType = "TRADING"       // (SPOT & MARGIN)
	SymbolStatusPostTrading  BISymbolStatusType = "POST_TRADING"  // (SPOT & MARGIN)
	SymbolStatusEndOfDay     BISymbolStatusType = "END_OF_DAY"    // (SPOT & MARGIN)
	SymbolStatusHalt         BISymbolStatusType = "HALT"          // (SPOT & MARGIN)
	SymbolStatusAuctionMatch BISymbolStatusType = "AUCTION_MATCH" // (SPOT & MARGIN)
	SymbolStatusBreak        BISymbolStatusType = "BREAK"         // (SPOT & MARGIN)

	BIAccountAndSymbolPermissionSPOT        BIAccountAndSymbolPermission = "SPOT"        // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionMARGIN      BIAccountAndSymbolPermission = "MARGIN"      // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionLEVERAGED   BIAccountAndSymbolPermission = "LEVERAGED"   // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_002 BIAccountAndSymbolPermission = "TRD_GRP_002" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_003 BIAccountAndSymbolPermission = "TRD_GRP_003" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_004 BIAccountAndSymbolPermission = "TRD_GRP_004" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_005 BIAccountAndSymbolPermission = "TRD_GRP_005" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_006 BIAccountAndSymbolPermission = "TRD_GRP_006" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_007 BIAccountAndSymbolPermission = "TRD_GRP_007" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_008 BIAccountAndSymbolPermission = "TRD_GRP_008" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_009 BIAccountAndSymbolPermission = "TRD_GRP_009" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_010 BIAccountAndSymbolPermission = "TRD_GRP_010" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_011 BIAccountAndSymbolPermission = "TRD_GRP_011" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_012 BIAccountAndSymbolPermission = "TRD_GRP_012" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_013 BIAccountAndSymbolPermission = "TRD_GRP_013" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_014 BIAccountAndSymbolPermission = "TRD_GRP_014" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_015 BIAccountAndSymbolPermission = "TRD_GRP_015" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_016 BIAccountAndSymbolPermission = "TRD_GRP_016" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_017 BIAccountAndSymbolPermission = "TRD_GRP_017" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_018 BIAccountAndSymbolPermission = "TRD_GRP_018" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_019 BIAccountAndSymbolPermission = "TRD_GRP_019" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_020 BIAccountAndSymbolPermission = "TRD_GRP_020" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_021 BIAccountAndSymbolPermission = "TRD_GRP_021" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_022 BIAccountAndSymbolPermission = "TRD_GRP_022" // (SPOT & MARGIN)
	BIAccountAndSymbolPermissionTRD_GRP_023 BIAccountAndSymbolPermission = "TRD_GRP_023" // (SPOT & MARGIN)

	ContractTypePerpetual      BIContractType = "PERPETUAL"       // (FUTURES)
	ContractTypeCurrent        BIContractType = "CURRENT_MONTH"   // (FUTURES)
	ContractTypeNext           BIContractType = "NEXT_MONTH"      // (FUTURES)
	ContractTypeCurrentQuarter BIContractType = "CURRENT_QUARTER" // (FUTURES)
	ContractTypeNextQuarter    BIContractType = "NEXT_QUARTER"    // (FUTURES)

	ContractStatusPendingTrading BIContractStatus = "PENDING_TRADING" // (FUTURES)
	ContractStatusTrading        BIContractStatus = "TRADING"         // (FUTURES)
	ContractStatusPreDelivering  BIContractStatus = "PRE_DELIVERING"  // (FUTURES)
	ContractStatusDelivering     BIContractStatus = "DELIVERING"      // (FUTURES)
	ContractStatusDelivered      BIContractStatus = "DELIVERED"       // (FUTURES)
	ContractStatusPreSettled     BIContractStatus = "PRE_SETTLED"     // (FUTURES)
	ContractStatusSettled        BIContractStatus = "SETTLED"         // (FUTURES)
	ContractStatusClosed         BIContractStatus = "CLOSED"          // (FUTURES)
)

// this is only used in services for parsing, so it should not be here!
//type BISymbolFilterType string // (SPOT & MARGIN & FUTURES)
//const (
//SymbolFilterTypePriceFilter                  BISymbolFilterType = "PRICE_FILTER"           // (SPOT & MARGIN & FUTURES)
//SymbolFilterTypeLotSizeFilter                BISymbolFilterType = "LOT_SIZE"               // (SPOT & MARGIN & FUTURES)
//SymbolFilterTypeMarketLotSizeFilter          BISymbolFilterType = "MARKET_LOT_SIZE"        // (SPOT & MARGIN & FUTURES)
//SymbolFilterTypeMaxNumOrdersFilter           BISymbolFilterType = "MAX_NUM_ORDERS"         // (SPOT & MARGIN & FUTURES)
//SymbolFilterTypeMaxNumAlgoOrdersFilter       BISymbolFilterType = "MAX_NUM_ALGO_ORDERS"    // (SPOT & MARGIN & FUTURES)
//SymbolFilterTypePercentPriceFilterSpotMargin BISymbolFilterType = "PERCENT_PRICE"          // (SPOT & MARGIN)
//SymbolFilterTypePercentPriceBySideFilter     BISymbolFilterType = "PERCENT_PRICE_BY_SIDE"  // (SPOT & MARGIN)
//SymbolFilterTypeMinNotionalFilterSpotMargin  BISymbolFilterType = "MIN_NOTIONAL"           // (SPOT & MARGIN)
//SymbolFilterTypeNotionalFilter               BISymbolFilterType = "NOTIONAL"               // (SPOT & MARGIN)
//SymbolFilterTypeIcebergPartsFilter           BISymbolFilterType = "ICEBERG_PARTS"          // (SPOT & MARGIN)
//SymbolFilterTypeMaxNumIcebergOrdersFilter    BISymbolFilterType = "MAX_NUM_ICEBERG_ORDERS" // (SPOT & MARGIN)
//SymbolFilterTypeMaxPositionFilter            BISymbolFilterType = "MAX_POSITION"           // (SPOT & MARGIN)
//SymbolFilterTypeTrailingDeltaFilter          BISymbolFilterType = "TRAILING_DELTA"         // (SPOT & MARGIN)
//SymbolFilterTypeMinNotionalFilterFutures     BISymbolFilterType = "MIN_NOTIONAL"           // (FUTURES)
//SymbolFilterTypePercentPriceFilterFutures    BISymbolFilterType = "PERCENT_PRICE"          // (FUTURES)
//)
