package streams

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== AggTradesStream ================================== */

// AggTradesStream is a websocket stream for aggregate trades.
type AggTradesStream struct {
	common.Stream
	SM       *common.StreamMeta
	Handler  *AggTradesHandler
	WSSymbol *string
}

// SetSymbol sets the symbol of the stream to the provided value,
// calls SetPathFunc on common.Stream, and returns the stream.
func (s *AggTradesStream) SetSymbol(restSymbol string) *AggTradesStream {
	wsSymbol := strings.ToLower(restSymbol)
	s.WSSymbol = &wsSymbol
	s.SetPathFunc(s.path)
	return s
}

// path constructs the path of the stream. Used as a callback in common.Stream.
func (s *AggTradesStream) path() string {
	path := "/ws/%s@aggTrade"
	return fmt.Sprintf(path, *s.WSSymbol)
}

/* ==================== AggTradesEvent (shrimpy-binance/streams) ========= */

func NewAggTradesHandler() *AggTradesHandler {
	return &AggTradesHandler{
		EventChan: make(chan AggTradesEvent, 256),
		ErrChan:   make(chan error, 1),
	}
}

type AggTradesHandler struct {
	EventChan chan AggTradesEvent
	ErrChan   chan error
}

func (h *AggTradesHandler) HandleSend(req common.WSRequest) {
	log.Warn(handleSendWarning)
}

func (h *AggTradesHandler) HandleError(err error) {
	//log.WithError(err).Warn("error in streamHandler")
	h.ErrChan <- err
}

func (h *AggTradesHandler) HandleRecv(msg []byte) {
	event := AggTradesEvent{}
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Fatalf("error unmarshaling response: %s", err)
		h.ErrChan <- err
	}
	//log.WithField("event", event).Info("handling recv")
	h.EventChan <- event
}

/* ==================== AggTradesEvent (shrimpy-binance/streams) ========= */

type AggTradesEvent struct {
	EventType        string `json:"e"`
	EventTime        int64  `json:"E"`
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
	IsBestMatch      bool   `json:"M"`
}
