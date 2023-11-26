package common

import "encoding/json"

/* ==================== Custom Types ===================================== */

// newTSNano creates a new TSNano from any int64 that represents a timestamp.
func NewTSNano(ts int64) TSNano {
	n := countDigitsInInt64(ts)
	pow := 19 - n

	// loop pow times to multiply ts by 10^pow
	for i := 0; i < pow; i++ {
		ts *= 10
	}
	return TSNano(ts)
}

// TSNano is an alias for int64 that represents a timestamp in nanoseconds.
// The main purpose for this type is to ensure that timestamps are always
// represented in nanoseconds, and to provide a consistent way to parse
// milli-second timestamps from json to nanoseconds.
type TSNano int64

func (ts TSNano) Int64() int64 {
	return int64(ts)
}

// UnmarshalJSON always unmarshals timestamps to nanoseconds.
func (ts *TSNano) UnmarshalJSON(data []byte) error {
	var tsTmp int64
	if err := json.Unmarshal(data, &tsTmp); err != nil {
		return err
	}

	*ts = NewTSNano(tsTmp)
	return nil
}

/* ==================== Structs ========================================== */

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
	TSSRespHeader    TSNano // (API, SAPI, FAPI, DAPI, EAPI)
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
	TSLSent    TSNano // timestamp local sent in nanoseconds
	TSSSent    TSNano // timestamp server sent in nanoseconds
	TSLRecv    TSNano // timestamp local received in nanoseconds
	TSSRecv    TSNano // timestamp server received in nanoseconds
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
	TSLRecv TSNano // timestamp local received in nanoseconds
	TSSRecv TSNano // timestamp server received in nanoseconds
	TSLSent TSNano // only applicable to requests (stream.DO)
	TSSent  TSNano // only applicable to requests (stream.DO)
}
