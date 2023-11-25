package streams

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== Shared AggTradesStream =========================== */

// AggTradesStream is a shared Stream implementation for agg trades streams.
// It includes a Stream, a handler for market streams, and a symbol for
// WebSocket communication.
type AggTradesStream[E Event] struct {
	common.Stream
	Handler  *MarketStreamHandler[E]
	WSSymbol *string
}

// SetSymbol sets the wsSymbol that is used in generating the path for the
// stream. It also sets the path function for the stream, if all path params
// are set (in this case wsSymbol is the only param).
func (s *AggTradesStream[E]) SetSymbol(restSymbol string) *AggTradesStream[E] {
	wsSymbol := strings.ToLower(restSymbol)
	s.WSSymbol = &wsSymbol
	s.SetPathFunc(s.path)
	return s
}

// path returns the path for the agg trades stream. It is used as the path
// function for the stream.
func (s *AggTradesStream[E]) path() string {
	path := "/ws/%s@aggTrade"
	return fmt.Sprintf(path, *s.WSSymbol)
}

/* ==================== sharedAggTradesEvent ============================= */

// sharedAggTradesEvent is a shared event for agg trades streams.
type sharedAggTradesEvent struct {
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
}

/* ==================== SpotMargin ======================================= */

// SpotMarginAggTradesEvent is an agg trades event for spot/margin streams.
type SpotMarginAggTradesEvent struct {
	StreamBaseEvent
	sharedAggTradesEvent
	Ignore interface{} `json:"M"`
}

// SpotMarginAggTradesHandler is a handler for spot/margin agg trades streams.
type SpotMarginAggTradesHandler = MarketStreamHandler[*SpotMarginAggTradesEvent]

// newSpotMarginAggTradesHandler creates a new SpotMarginAggTradesHandler.
func newSpotMarginAggTradesHandler(logger *log.Entry) *SpotMarginAggTradesHandler {
	return newMarketStreamHandler[*SpotMarginAggTradesEvent](logger)
}

// SpotMarginAggTradesStream is a stream for spot/margin agg trades streams.
type SpotMarginAggTradesStream = AggTradesStream[*SpotMarginAggTradesEvent]

/* ==================== Futures ======================================= */

// FuturesAggTradesEvent is an agg trades event for futures streams.
type FuturesAggTradesEvent struct {
	sharedAggTradesEvent
	StreamBaseEvent
}

// FuturesAggTradesHandler is a handler for futures agg trades streams.
type FuturesAggTradesHandler = MarketStreamHandler[*FuturesAggTradesEvent]

func newFuturesAggTradesHandler(logger *log.Entry) *FuturesAggTradesHandler {
	return newMarketStreamHandler[*FuturesAggTradesEvent](logger)
}

// FuturesAggTradesStream is a stream for futures agg trades streams.
type FuturesAggTradesStream = AggTradesStream[*FuturesAggTradesEvent]
