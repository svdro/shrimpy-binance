package streams

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== MarketStreamHandler ============================== */

// newMarketStreamHandler creates a new MarketStreamHandler.
func newMarketStreamHandler[E Event](logger *log.Entry) *MarketStreamHandler[E] {
	return &MarketStreamHandler[E]{
		EventChan: make(chan E, 256),
		ErrChan:   make(chan error, 1),
		logger:    logger,
	}
}

// MarketStreamHandler implements the common.StreamHandler interface.
// it is a generic handler for websocket market streams.
type MarketStreamHandler[E Event] struct {
	EventChan chan E
	ErrChan   chan error
	logger    *log.Entry
}

// HandleError puts the error on the ErrChan and expects the caller to handle
// the error. In this case stream.Run will is the default caller and will
// pass either a common.WSConnError or a common.WSHandlerError to the caller.
func (h *MarketStreamHandler[E]) HandleError(err error) {
	h.ErrChan <- err
}

// HandleSend is not implemented. It is not used for market streams.
func (h *MarketStreamHandler[E]) HandleSend(req common.WSRequest) *common.WSHandlerError {
	log.Warn(handleSendWarning)
	return nil
}

// HandleRecv parses the message and sends it to the EventChan.
// If an error occurs, it is logged and returned to the caller. It is then
// the caller's responsibility to handle the error. In this case stream.Run
// will put the error on the ErrChan as a non-transient error, and shutdown.
func (h *MarketStreamHandler[E]) HandleRecv(msg []byte, TSLRecv, TSSRecv int64) *common.WSHandlerError {
	var event E
	if err := json.Unmarshal(msg, &event); err != nil {
		h.logger.WithField("msg", string(msg)).WithError(err).Error("failed to unmarshal event")
		return &common.WSHandlerError{Err: err, Reason: "failed to unmarshal event", IsFatal: true}
	}
	event.addEventMeta(TSLRecv, TSSRecv)
	h.EventChan <- event
	return nil
}

/* ==================== UserDataStreamHandler =============================*/
//.. implement this
