package common

import (
	"context"
	"net/url"

	log "github.com/sirupsen/logrus"
)

/* ==================== Interfaces (shrimpy-binance/client) ============== */

// Client
type Client interface {
	SetServerTimeOffset(offset int64)
	GetServerTimeOffset() int64
}

// RESTClient
type RESTClient interface {
	Do(ctx context.Context, sm *ServiceMeta, p url.Values) ([]byte, error)
}

// WSClient
type WSClient interface {
	NewStream(sm *StreamMeta, handler StreamHandler, logger *log.Entry) Stream
}

// TimeHandler
type TimeHandler interface {
	TSLNow() TSNano           // local time in nanoseconds
	TSSNow() TSNano           // server time in nanoseconds
	TSLToTSS(t TSNano) TSNano // local time in nanoseconds to server time in nanoseconds
	TSSToTSL(t TSNano) TSNano // server time in nanoseconds to local time in nanoseconds
	Offset() int64            // offset from server time in nanoseconds
}

/* ==================== Interfaces (shrimpy-binance/common) ============== */

// WSRequest is a request to the websocket pump
type WSRequest interface {
	GetID() string
}

// StreamHandler is a handler for websocket events
type StreamHandler interface {
	HandleSend(req WSRequest) *WSHandlerError
	HandleRecv(msg []byte, TSLRecv, TSSRecv TSNano) *WSHandlerError
	HandleError(err error)
}

// Stream is responsible for handling the underlying websocket connection.
// It is responsible for reading and writing to the websocket connection.
type Stream interface {
	Run(ctx context.Context)
	SetPathFunc(f func() string)
}
