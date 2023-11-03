package common

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

	RateLimitTypeIP  BIRateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeUID BIRateLimitType = "ORDERS"
	RateLimitTypeRAW BIRateLimitType = "RAW_REQUESTSLNano"
)
