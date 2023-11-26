package client

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
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
func (th *timeHandler) TSLNow() common.TSNano {
	return common.TSNano(time.Now().UnixNano())
}

// TSSNow timestamp server now (nanoseconds)
func (th *timeHandler) TSSNow() common.TSNano {
	return th.TSLToTSS(th.TSLNow())
}

// TSLToTSS timestamp local (nanoseconds) to timestamp server (nanoseconds)
func (th *timeHandler) TSLToTSS(tsl common.TSNano) common.TSNano {
	return common.TSNano(tsl.Int64() - th.offset)
}

// TSSToTSL timestamp server (nanoseconds) to timestamp local (nanoseconds)
func (th *timeHandler) TSSToTSL(tss common.TSNano) common.TSNano {
	return common.TSNano(tss.Int64() + th.offset)
}

// Offset offset from server time (nanoseconds)
func (th *timeHandler) Offset() int64 {
	return th.offset
}
