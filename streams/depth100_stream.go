package streams

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== Shared DiffDepthStream ============================*/

// DiffDepth100Stream is a shared Stream implementation for diff depth streams.
// It includes a Stream, a handler for market streams, and a symbol and the
// update speed for WebSocket communication.
type DiffDepthStream[E Event] struct {
	common.Stream
	Handler     *MarketStreamHandler[E]
	WSSymbol    *string
	UpdateSpeed *int
}

// SetSymbol sets the wsSymbol that is used in generating the path for the
// stream. It also sets the path function for the stream, if all path params
// are set.
func (s *DiffDepthStream[E]) SetSymbol(restSymbol string) *DiffDepthStream[E] {
	wsSymbol := strings.ToLower(restSymbol)
	s.WSSymbol = &wsSymbol
	s.SetPathFunc(s.path)
	return s
}

// path returns the path for the diff depth stream. It is used as the path
// function for the stream.
func (s *DiffDepthStream[E]) path() string {
	//path := "/ws/%s@depth@100ms"
	path := "/ws/%s@depth@%dms"
	log.Info("path: ", fmt.Sprintf(path, *s.WSSymbol, *s.UpdateSpeed))
	return fmt.Sprintf(path, *s.WSSymbol, *s.UpdateSpeed)
}

/* ==================== sharedDiffDepthEvent ============================= */

// Level represents a price and quantity pair.
type Level struct {
	Price string `json:"p"`
	Qty   string `json:"q"`
}

// UnmarshalJSON unmarshals a price qty pair from a JSON array to a Level.
func (l *Level) UnmarshalJSON(data []byte) error {
	var tmp [2]string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	l.Price = tmp[0]
	l.Qty = tmp[1]
	return nil
}

// sharedDiffDepthEvent holds all fields that are common to diff depth events.
type sharedDiffDepthEvent struct {
	Symbol  string  `json:"s"`
	FirstID int64   `json:"U"`
	FinalID int64   `json:"u"`
	Bids    []Level `json:"b"`
	Asks    []Level `json:"a"`
}

/* ==================== SpotMargin ======================================= */

// SpotMarginDiffDepthEvent is a diff depth event for spot and margin markets.
// e.g. https://api.binance.com/api/v3/depth?symbol=BTCUSDT&limit=100
type SpotMarginDiffDepthEvent struct {
	StreamBaseEvent
	sharedDiffDepthEvent
}

// SpotMarginDiffDepthHandler is a handler for spot and margin diff depth streams.
type SpotMarginDiffDepthHandler = MarketStreamHandler[*SpotMarginDiffDepthEvent]

// newSpotMarginDiffDepthHandler creates a new SpotMarginDiffDepthHandler.
func newSpotMarginDiffDepthHandler(logger *log.Entry) *SpotMarginDiffDepthHandler {
	return newMarketStreamHandler[*SpotMarginDiffDepthEvent](logger)
}

// SpotMarginDiffDepthStream is a diff depth stream for spot and margin markets.
type SpotMarginDiffDepthStream = DiffDepthStream[*SpotMarginDiffDepthEvent]

/* ==================== Futures ========================================== */

// FuturesDiffDepthEvent is a diff depth event for futures markets.
// TODO: convert millisecond timestamps to ns timestamps.
// This is not trivial since this would logically be done in the
// MarketStreamHandler when it unmarshals the event. However, the
// MarketStreamHandler is generic and does not know about the
// specific event type. Maybe make a TimespampNano type that
// implements json.Unmarshaler and json.Marshaler, and use that
// in the event?
// That would break TimeHandler, which expects int64, so a custom TimespampNano
// type would require changing TimeHandler, which is a little bit of work.
type FuturesDiffDepthEvent struct {
	StreamBaseEvent
	sharedDiffDepthEvent
	TSSTransact common.TSNano `json:"T"`  // FAPI
	LastFinalID int64         `json:"pu"` // FAPI
}

// FuturesDiffDepthHandler is a handler for futures diff depth streams.
type FuturesDiffDepthHandler = MarketStreamHandler[*FuturesDiffDepthEvent]

// newFuturesDiffDepthHandler creates a new FuturesDiffDepthHandler.
func newFuturesDiffDepthHandler(logger *log.Entry) *FuturesDiffDepthHandler {
	return newMarketStreamHandler[*FuturesDiffDepthEvent](logger)
}

// FuturesDiffDepthStream is a diff depth stream for futures markets.
type FuturesDiffDepthStream = DiffDepthStream[*FuturesDiffDepthEvent]
