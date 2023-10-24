package binance

import "net/http"

/* ==================== Constants ======================================== */

type BIDataSource int
type BIEndpoint string
type BIEndpointType string
type BIRateLimitIntervalType string
type BISecurityType int
type BIRateLimitType string
type BIHttpResponseCode int

const (
	DataSourceNone BIDataSource = iota
	DataSourceMatchingEngine
	DataSourceMemory
	DataSourceDatabase

	endpointAPI  BIEndpoint = "api.binance.com"
	endpointFAPI BIEndpoint = "fapi.binance.com"

	endpointTypeAPI  BIEndpointType = "api"
	endpointTypeSAPI BIEndpointType = "sapi"
	endpointTypeFAPI BIEndpointType = "fapi"

	HTTPStatusOK              BIHttpResponseCode = http.StatusOK              // (200)
	HTTPStatusTeapot          BIHttpResponseCode = http.StatusTeapot          // (418) IP ban
	HTTPStatusTooManyRequests BIHttpResponseCode = http.StatusTooManyRequests // (429) Backoff
	HTTPStatusBadRequest      BIHttpResponseCode = http.StatusBadRequest      // (400) Invalid request

	intervalSecond BIRateLimitIntervalType = "SECOND"
	intervalMinute BIRateLimitIntervalType = "MINUTE"
	intervalDay    BIRateLimitIntervalType = "DAY"

	SecurityTypeNone   BISecurityType = iota // NONE
	SecurityTypeApiKey                       // USER_STREAM, MARKET_DATA
	SecurityTypeSigned                       // TRADE, MARGIN, USER_DATA

	rateLimitTypeIP  BIRateLimitType = "REQUEST_WEIGHT"
	rateLimitTypeUID BIRateLimitType = "ORDERS"
	rateLimitTypeRAW BIRateLimitType = "RAW_REQUESTS"
)
