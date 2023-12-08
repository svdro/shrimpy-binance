package defaults

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
	bsv "github.com/svdro/shrimpy-binance/services"
	bst "github.com/svdro/shrimpy-binance/streams"
)

/* ==================== OrderBook ======================================== */

// OrderBook is an interface for orderbooks.
type OrderBook interface {
	// initFromDepthResponse initializes the orderbook from a depth response.
	initFromDepthResponse(resp *bsv.SpotMarginDepthResponse)

	// updateFromDepthEvent updates the orderbook from a depth event. It returns
	// true if the orderbook is in sync with the event, and false otherwise.
	updateFromDepthEvent(event *bst.SpotMarginDiffDepthEvent) bool

	// take snapshot returns the current state of the orderbook as a snapshot.
	// If depth is -1, the entire orderbook is returned.
	takeSnapshot(depth int) *OrderBookSnapshot

	// importTimestampsFromSnapshot uses the snapshot from an older version
	// of the orderbook to update the TSSEvent of the levels in the current
	// one. This is useful when initializing an orderbook from a depth event,
	// since depth events don't have timestamps for the levels.
	importTimestampsFromSnapshot(snapshot *OrderBookSnapshot)
}

// newOrderBook creates a new instance of OrderBook.
func newOrderBook(logger *log.Entry) OrderBook {
	logger.Info("newOrderBook")
	return &orderBook{
		asks:   newOrderBookSide(true),
		bids:   newOrderBookSide(false),
		logger: logger.WithFields(log.Fields{"_caller": "orderBook"}),
	}
}

/* ==================== orderBook ======================================== */

// orderBook
type orderBook struct {
	mu           sync.Mutex
	asks         OrderBookSide
	bids         OrderBookSide
	hasFirstID   bool
	lastUpdateID int64
	lastTSSEvent common.TSNano
	logger       *log.Entry
}

// importTimestampsFromSnapshot synchronizes the timestamps of the levels in
// each side of the orderbook with the timestamps of the levels in the snapshot.
func (b *orderBook) importTimestampsFromSnapshot(snapshot *OrderBookSnapshot) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.asks.syncTimestampsFromSnapshot(snapshot.Asks)
	b.bids.syncTimestampsFromSnapshot(snapshot.Bids)
}

// takeSnapshot returns the current state of the orderbook as a snapshot.
func (b *orderBook) takeSnapshot(depth int) *OrderBookSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	var asks, bids []Level
	var ok bool

	if asks, ok = b.asks.getSnapshot(depth); !ok {
		return nil
	}
	if bids, ok = b.bids.getSnapshot(depth); !ok {
		return nil
	}

	return &OrderBookSnapshot{Asks: asks, Bids: bids, LastTSSEvent: b.lastTSSEvent}
}

// initFromDepthResponse initializes the orderbook from a depth response.
func (b *orderBook) initFromDepthResponse(resp *bsv.SpotMarginDepthResponse) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// simulate a tss event for snapshots (annoyingly this is not included in the response)
	tssEvent := common.TSNano((resp.TSSRecv.Int64() + resp.TSSSent.Int64()) / 2)

	// set bids, asks and lastUpdateID, and lastTSSEvent
	b.asks.setFromServiceLevels(resp.Asks, tssEvent)
	b.bids.setFromServiceLevels(resp.Bids, tssEvent)
	b.lastTSSEvent = common.TSNano(resp.TSSRecv.Int64() - resp.TSSSent.Int64())
	b.lastUpdateID = resp.LastUpdateID

	// log
	b.logger.WithFields(log.Fields{
		"asks": b.asks.Len(), "bids": b.bids.Len(), "lastUpdateID": b.lastUpdateID},
	).Info("initFromDepthResponse")
}

// updateFromDepthEvent updates the orderbook from a depth event. It returns
// true if the orderbook is in sync with the event, and false otherwise.
func (b *orderBook) updateFromDepthEvent(event *bst.SpotMarginDiffDepthEvent) bool {
	// make context logger
	logger := b.logger.WithFields(log.Fields{
		"asks":                len(event.Asks),
		"bids":                len(event.Bids),
		"lastUpdateID (book)": b.lastUpdateID,
		"firstID (event)":     event.FirstID,
		"finalID (event)":     event.FinalID,
	})

	// lock
	b.mu.Lock()
	defer b.mu.Unlock()

	// check if event is in order and update bids and asks
	switch {

	// drop event if it is too old (only if first event has not been identified)
	case !b.hasFirstID && event.FinalID <= b.lastUpdateID:
		logger.Trace("dropping event")
		return true

	// identify first event to be processed (only if first event has not been identified)
	case !b.hasFirstID && event.FirstID <= b.lastUpdateID+1:
		logger.Info("first event")

		b.asks.updateWithStreamLevels(event.Asks, event.TSSEvent, logger)
		b.bids.updateWithStreamLevels(event.Bids, event.TSSEvent, logger)

		b.lastTSSEvent = event.TSSEvent
		b.lastUpdateID = event.FinalID
		b.hasFirstID = true
		return true

	// handle event
	case event.FirstID == b.lastUpdateID+1:
		logger.Trace("handling event")

		b.asks.updateWithStreamLevels(event.Asks, event.TSSEvent, logger)
		b.bids.updateWithStreamLevels(event.Bids, event.TSSEvent, logger)

		b.lastTSSEvent = event.TSSEvent
		b.lastUpdateID = event.FinalID

		return true

	// any event that does not match the above criteria is out of order
	default:
		logger.Warn("event out of order")
		return false
	}
}
