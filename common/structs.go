package common

// ServiceDefinition holds all hardcoded data needed to make a service call.
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

// RateLimitUpdate is used by ServiceResponseHeader
// ServiceResponseHeader is used by ServiceMeta,
// so this should be in common.
type RateLimitUpdate struct {
	EndpointType    BIEndpointType
	RateLimitType   BIRateLimitType
	IntervalSeconds int
	Count           int
}

// serviceResponseHeader contains all relevant information from a binance
// http response's header. Optional headers are included as pointers.
type ServiceResponseHeader struct {
	Server           string // (API, SAPI, FAPI, DAPI, EAPI)
	TSSRespHeader    int64  // (API, SAPI, FAPI, DAPI, EAPI)
	RateLimitUpdates []RateLimitUpdate
	RetryAfter       *int // (seconds) (API, SAPI, FAPI, DAPI, EAPI)
}

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

// StreamDefinition holds all hardcoded data needed to create a stream.
type StreamDefinition struct {
	Scheme       string
	Endpoint     BIWSEndpoint
	EndpointType BIEndpointType   // api, fapi, etc
	SecurityType BIWSSecurityType // WSSecurityTypeListenKey, WSSecurityTypeNone
	UpdateSpeed  int              // (milliseconds)
}

func NewStreamMeta(sd StreamDefinition) *StreamMeta {
	return &StreamMeta{SD: sd}
}

// StreamMeta
// Stream Meta persists over the lifetime of a stream.
// As such, it cannot be used to store data on individual events.
// Possibly use this later to collect metadata on a stream (e.g. disconnects,
// reconnects, etc), and/or use this to configure the (dynamic) stream path.
type StreamMeta struct {
	SD StreamDefinition
}

// StreamEventMeta holds metadata on a single websocket event,
// such as timestamps.
type StreamEventMeta struct {
	TSLRecv int64 // timestamp local received in nanoseconds
	TSSRecv int64 // timestamp server received in nanoseconds
	TSLSent int64 // only applicable to requests (stream.DO)
	TSSent  int64 // only applicable to requests (stream.DO)
}
