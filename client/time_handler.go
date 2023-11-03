package client

import (
	"time"

	log "github.com/sirupsen/logrus"
)

/* ==================== timeHandler ====================================== */

// newTimeHandler returns a new timeHandler.
func newTimeHandler(c *Client) *timeHandler {
	return &timeHandler{
		c:      c,
		logger: c.logger.WithField("_caller", "timeHandler"),
	}
}

// timeHandler is responsible for keeping in sync with server time and
// converting between local time and server time.
type timeHandler struct {
	c      *Client
	offset int64 // offset from server time in nanoseconds
	logger *log.Entry
}

// sync
func (th *timeHandler) sync() error {
	_ = th.c // access client and do stuff to set offset
	return nil
}

// TSLNow timestamp loscal now (nanoseconds)
func (th *timeHandler) TSLNow() int64 {
	return time.Now().UnixNano()
}

// TSSNow timestamp server now (nanoseconds)
func (th *timeHandler) TSSNow() int64 {
	return th.TSLToTSS(th.TSLNow())
}

// TSLToTSS timestamp local (nanoseconds) to timestamp server (nanoseconds)
func (th *timeHandler) TSLToTSS(tsl int64) int64 {
	return tsl - th.offset
}

// TSSToTSL timestamp server (nanoseconds) to timestamp local (nanoseconds)
func (th *timeHandler) TSSToTSL(tss int64) int64 {
	return tss + th.offset
}

// Offset offset from server time (nanoseconds)
func (th *timeHandler) Offset() int64 {
	return th.offset
}
