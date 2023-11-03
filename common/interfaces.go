package common

import (
	"context"
	"net/url"
)

/* ==================== Interfaces ======================================= */

// Client
type Client interface{}

// RESTClient
type RESTClient interface {
	Do(ctx context.Context, sm *ServiceMeta, p url.Values) ([]byte, error)
	TimeHandler() TimeHandler
}

// WSClient
type WSClient interface{}

// TimeHandler
type TimeHandler interface {
	TSLNow() int64          // local time in nanoseconds
	TSSNow() int64          // server time in nanoseconds
	TSLToTSS(t int64) int64 // local time in nanoseconds to server time in nanoseconds
	TSSToTSL(t int64) int64 // server time in nanoseconds to local time in nanoseconds
	Offset() int64          // offset from server time in nanoseconds
}
