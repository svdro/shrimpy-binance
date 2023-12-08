package streams

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== StreamHandler utils ============================== */

// parseEventType parses the event type from the websocket message.
// if the message is not valid json, or the event type is empty, an error is returned.
func parseEventTypeOld(msg []byte) (string, error) {
	var event StreamBaseEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return "", err
	}
	if event.EventType == "" {
		return "", fmt.Errorf("event type is empty")
	}
	return event.EventType, nil
}

// parseEventType parses the event type from the websocket message.
func parseEventType(msg []byte, logger *log.Entry) (string, *common.WSHandlerError) {

	// unmarshal StreamBaseEvent
	var event StreamBaseEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		logger.WithField("msg", string(msg)).WithError(err).Error("failed to parse event type")
		return "", &common.WSHandlerError{Err: err, Reason: "failed to parse event type", IsFatal: true}
	}

	// assert event type is not empty
	if event.EventType == "" {
		err := fmt.Errorf("event type is empty")
		logger.WithField("msg", string(msg)).WithError(err).Error("failed to parse event type")
		return "", &common.WSHandlerError{Err: err, Reason: "failed to parse event type", IsFatal: true}
	}

	return event.EventType, nil
}

// unmarshalAndSendEvent unmarshals the message into the event and sends it to
// the eventChan. If an error occurs, it is logged and returned to the caller.
func unmarshalAndSendEvent[E Event](
	msg []byte, event *E, TSLRecv, TSSRecv common.TSNano, eventChan chan<- E, logger *log.Entry,
) *common.WSHandlerError {

	// unmarshal event
	if err := json.Unmarshal(msg, event); err != nil {
		logger.WithField("msg", string(msg)).WithError(err).Error("failed to unmarshal event")
		return &common.WSHandlerError{Err: err, Reason: "failed to unmarshal event", IsFatal: true}
	}

	// add event meta
	(*event).addEventMeta(TSLRecv, TSSRecv)
	eventChan <- (*event)

	return nil
}

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
// NOTE: when passing new(E) to unmarshalAndSendEvent,
// this does not correctly initialize the struct with its embedded structs.
// it creates a pointer to a pointer of a nil value (or sth like that).
// it works because json fixes the initialization, but only when
// calling json.Unmarshal before calling methods on the embedded struct.
// maybe fix this cause it's ugly, but also maybe not cause it works.
func (h *MarketStreamHandler[E]) HandleRecv(msg []byte, TSLRecv, TSSRecv common.TSNano) *common.WSHandlerError {
	wshErr := unmarshalAndSendEvent(msg, new(E), TSLRecv, TSSRecv, h.EventChan, h.logger)
	return wshErr
}

/* ==================== SpotMarginUserDataStreamHandler ===================*/

// newSpotMarginUserDataStreamHandler creates a new SpotMarginUserDataStreamHandler.
func newSpotMarginUserDataStreamHandler[A, B, O Event](
	logger *log.Entry) *SpotMarginUserDataStreamHandler[A, B, O] {
	return &SpotMarginUserDataStreamHandler[A, B, O]{
		AccountUpdateEventChan: make(chan A, 256),
		BalanceUpdateEventChan: make(chan B, 256),
		OrderUpdateEventChan:   make(chan O, 256),
		ErrChan:                make(chan error, 1),
		logger:                 logger,
	}
}

// SpotMarginUserDataStreamHandler implements the common.StreamHandler
// interface. It is a generic handler for spot/margin user data streams.
type SpotMarginUserDataStreamHandler[A, B, O Event] struct {
	AccountUpdateEventChan chan A
	BalanceUpdateEventChan chan B
	OrderUpdateEventChan   chan O
	ErrChan                chan error
	logger                 *log.Entry
}

// HandleError puts the error on the ErrChan and expects the caller to handle
// the error. In this case stream.Run will is the default caller and will
// pass either a common.WSConnError or a common.WSHandlerError to the caller.
func (h *SpotMarginUserDataStreamHandler[A, B, O]) HandleError(err error) {
	h.ErrChan <- err
}

// HandleSend is not implemented. It is not used spot/margin user data streams.
func (h *SpotMarginUserDataStreamHandler[A, B, O]) HandleSend(req common.WSRequest) *common.WSHandlerError {
	log.Warn(handleSendWarning)
	return nil
}

// HandleRecv parses the message and sends it to the corresponding eventChan.
// If an error occurs, it is logged and returned to the caller. It is then
// the caller's responsibility to handle the error. In this case stream.Run
// will put the error on the ErrChan as a non-transient error, and shutdown.
func (h *SpotMarginUserDataStreamHandler[A, B, O]) HandleRecv(msg []byte, TSLRecv, TSSRecv common.TSNano) *common.WSHandlerError {

	// parse event type
	eventType, wshErr := parseEventType(msg, h.logger)
	if wshErr != nil {
		return wshErr
	}

	switch eventType {
	case "outboundAccountPosition":
		wshErr = unmarshalAndSendEvent(msg, new(A), TSLRecv, TSSRecv, h.AccountUpdateEventChan, h.logger)
	case "balanceUpdate":
		wshErr = unmarshalAndSendEvent(msg, new(B), TSLRecv, TSSRecv, h.BalanceUpdateEventChan, h.logger)
	case "executionReport":
		wshErr = unmarshalAndSendEvent(msg, new(O), TSLRecv, TSSRecv, h.OrderUpdateEventChan, h.logger)
	default:
		err := fmt.Errorf("unknown event type: %s", eventType)
		return &common.WSHandlerError{Err: err, Reason: "unknown event type", IsFatal: true}
	}

	return wshErr
}
