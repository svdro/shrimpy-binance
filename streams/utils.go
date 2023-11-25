package streams

// Event is an interface for all websocket events.
// All events must implement addEventMeta.
type Event interface {
	addEventMeta(TSLRecv, TSSRecv int64)
}

// StreamBaseEvent holds meta data that is common to all events.
// TSS is the timestamp in server time.
// TSL is the timestamp in local time.
type StreamBaseEvent struct {
	EventType string `json:"e"`
	TSSEvent  int64  `json:"E"`
	TSLRecv   int64
	TSSRecv   int64
}

// addEventMeta adds meta data that was generated when the event was received.
func (s *StreamBaseEvent) addEventMeta(TSLRecv, TSSRecv int64) {
	s.TSLRecv = TSLRecv
	s.TSSRecv = TSSRecv
}
