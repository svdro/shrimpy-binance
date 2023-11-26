package streams

import (
	"github.com/svdro/shrimpy-binance/common"
)

// Event is an interface for all websocket events.
// All events must implement addEventMeta.
type Event interface {
	addEventMeta(TSLRecv, TSSRecv common.TSNano)
}

// StreamBaseEvent holds meta data that is common to all events.
// TSS is the timestamp in server time.
// TSL is the timestamp in local time.
type StreamBaseEvent struct {
	EventType string        `json:"e"`
	TSSEvent  common.TSNano `json:"E"`
	TSLRecv   common.TSNano
	TSSRecv   common.TSNano
}

// addEventMeta adds meta data that was generated when the event was received.
func (s *StreamBaseEvent) addEventMeta(TSLRecv, TSSRecv common.TSNano) {
	s.TSLRecv = TSLRecv
	s.TSSRecv = TSSRecv
}
