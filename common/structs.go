package common

/* ==================== ServiceDefinition ================================ */

// ServiceDefinition holds all data needed to make a service call.
type ServiceDefinition struct {
	Scheme              string
	Method              string
	Endpoint            BIEndpoint
	Path                string
	EndpointType        BIEndpointType
	SecurityType        BISecurityType
	PrimaryDatasource   BIDataSource
	SecondaryDatasource BIDataSource
	WeightIP            int // WeightLimit
	WeightUID           int // UIDLimit
}

/* ==================== BIResponseHeader ================================= */

// RateLimitHeader contains used rate limit updates.
type RateLimitHeader struct {
	RateLimitType       BIRateLimitType
	IntervalType        BIRateLimitIntervalType
	IntervalNum         int
	IntervalNanoSeconds int64
	Count               int // weight count
}

// serviceResponseHeader contains all relevant information from a binance
// http response's header.
type ServiceResponseHeader struct {
	Server        string                    // (API, SAPI, FAPI, DAPI, EAPI)
	TSSRespHeader int64                     // (API, SAPI, FAPI, DAPI, EAPI)
	UIDLimits     map[int64]RateLimitHeader // (API, SAPI, FAPI, DAPI, EAPI)
	IPLimits      map[int64]RateLimitHeader // (API, SAPI, FAPI, DAPI, EAPI)
}

/* ==================== ServiceMeta ====================================== */

// NewServiceMeta
func NewServiceMeta(sd ServiceDefinition) *ServiceMeta {
	return &ServiceMeta{SD: sd}
}

// ServiceMeta holds all data collected during a service call, as well as
// the ServiceDefinition used to make the call.
type ServiceMeta struct {
	SD         ServiceDefinition
	SRH        *ServiceResponseHeader
	StatusCode int
	TSLSent    int64 // timestamp local sent in nanoseconds
	TSSSent    int64 // timestamp server sent in nanoseconds
	TSLRecv    int64 // timestamp local received in nanoseconds
	TSSRecv    int64 // timestamp server received in nanoseconds
}
