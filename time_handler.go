package binance

import (
	"time"

	log "github.com/sirupsen/logrus"
)

/* ==================== TimeHandler ====================================== */

// timeHandler is responsible for keeping the Client's time in sync with
// binance servers. In practice, time across different binance servers is
// consistent, so we only need to sync time with one server.
// The most consistent way to sync time is to binance's WebSocketAPI.
type timeHandler struct {
	c      *Client
	logger *log.Entry
	offset int64
}

func (th *timeHandler) SyncTime() error {
	return nil
}

func (th *timeHandler) TSNanoNowLocal() int64 {
	return time.Now().UnixNano()
}

func (th *timeHandler) TSNanoNowServer() int64 {
	return th.TSNanoNowLocal() + th.offset
}

func (th *timeHandler) TSNanoLocalToServer(tsNanoLocal int64) int64 {
	return tsNanoLocal + th.offset
}

func (th *timeHandler) TSNanoServerToLocal(tsNanoServer int64) int64 {
	return tsNanoServer - th.offset
}

func (th *timeHandler) MillitoNano(ts int64) int64 {
	return ts * 1e6
}

func (th *timeHandler) NanotoMilli(ts int64) int64 {
	return ts / 1e6
}
