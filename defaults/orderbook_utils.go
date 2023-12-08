package defaults

import (
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
	bsv "github.com/svdro/shrimpy-binance/services"
	bst "github.com/svdro/shrimpy-binance/streams"
)

/* ==================== Level ============================================ */

// newLevelFromServiceLevel creates a new Level from a bsv.Level.
func newLevelFromServiceLevel(l *bsv.Level, tssEvent common.TSNano) *Level {
	return &Level{Price: l.Price, Qty: l.Qty, TSSEvent: tssEvent}
}

// newLevelFromStreamLevel creates a new Level from a bst.Level.
func newLevelFromStreamLevel(l *bst.Level, tssEvent common.TSNano) *Level {
	return &Level{Price: l.Price, Qty: l.Qty, TSSEvent: tssEvent}
}

// Level is a price level in the orderbook plus the timestamp the level was
// created/ last updated.
type Level struct {
	Price    string
	Qty      string
	TSSEvent common.TSNano
}

/* ==================== Snapshot ========================================= */

// OrderBookSnapshot is a snapshot of the state of the orderbook.
type OrderBookSnapshot struct {
	Asks         []Level
	Bids         []Level
	LastTSSEvent common.TSNano
}

/* ==================== OrderBookSide ==================================== */

type OrderBookSide interface {
	Len() int
	updateLevel(level *Level, logger *log.Entry)
	updateWithStreamLevels(streamLevels []bst.Level, tssEvent common.TSNano, logger *log.Entry)
	setFromServiceLevels(levels []bsv.Level, tssEvent common.TSNano)
	getSnapshot(depth int) ([]Level, bool)
	syncTimestampsFromSnapshot(snapshot []Level)
}

// neworderBookSide creates a new orderBookSide.
func newOrderBookSide(asc bool) OrderBookSide {
	return &orderBookSide{
		sideMap:     make(map[string]*Level),
		sortedSlice: newSortedSide(asc),
	}
}

type orderBookSide struct {
	sideMap     map[string]*Level
	sortedSlice SortedSide
}

func (s *orderBookSide) Len() int {
	return len(s.sideMap)
}

// syncTimestampsFromSnapshot iterates over the levels in the snapshot and
// updates the TSSEvent of the corresponding level in orderBookSide.
// the sortedSlice does not need to be updated since it is a slice of pointers
func (s *orderBookSide) syncTimestampsFromSnapshot(snapshot []Level) {
	for _, level := range snapshot {
		if lvl, ok := s.sideMap[level.Price]; ok {
			if lvl.Qty != level.Qty {
				log.WithFields(log.Fields{
					"price":       level.Price,
					"qtySide":     level.Qty,
					"qtySnapshot": lvl.Qty,
					"tssSide":     lvl.TSSEvent,
					"tssSnapshot": level.TSSEvent,
				}).Warn("syncTimestampsFromSnapshot: qty mismatch")
				continue
			}
			lvl.TSSEvent = level.TSSEvent
		}
	}
}

func (s *orderBookSide) getSnapshot(depth int) ([]Level, bool) {
	if depth == -1 {
		depth = s.sortedSlice.len()
	}
	return s.sortedSlice.snapshot(depth)
}

func (s *orderBookSide) updateLevel(level *Level, logger *log.Entry) {
	// make context logger
	l := logger.WithFields(log.Fields{
		"price": level.Price, "qty": level.Qty, "tssEvent": level.TSSEvent,
	})

	// parse qty
	qty, err := strconv.ParseFloat(level.Qty, 64)
	if err != nil {
		l.WithError(err).Panic("updateLevel: failed to parse qty")
	}

	// delete level
	if qty == 0.0 {
		delete(s.sideMap, level.Price)
		s.sortedSlice.delete(level)
		return
	}

	// update TSSEvent and qty
	// since we are using pointers to levels we don't need to update the level
	// in the map the sortedSlice explicitly.
	if lvl, ok := s.sideMap[level.Price]; ok {
		lvl.TSSEvent = level.TSSEvent
		lvl.Qty = level.Qty
		return
	}

	// insert level
	s.sideMap[level.Price] = level
	s.sortedSlice.insert(level)
}

func (s *orderBookSide) updateWithStreamLevels(
	streamLevels []bst.Level, tssEvent common.TSNano, logger *log.Entry) {

	for _, streamLevel := range streamLevels {
		level := newLevelFromStreamLevel(&streamLevel, tssEvent)
		s.updateLevel(level, logger)
	}
}

func (s *orderBookSide) setFromServiceLevels(levels []bsv.Level, tssEvent common.TSNano) {
	for _, level := range levels {
		s.sideMap[level.Price] = newLevelFromServiceLevel(&level, tssEvent)
		s.sortedSlice.insert(s.sideMap[level.Price])
	}
}

/* ==================== SortedSide ======================================= */

// SortedSide interface is designed for managing sorted orderbook sides.
// It works alongside a map-based side, which maintains pointers to the same
// level objects. The function of this interface is to provide an efficient
// method for sorted orderbook state retrieval without having to sort the
// map-based side on each snapshot.
// The map-based side updates existing levels. Hence, the sorted side
// does not need to update existing levels. This reduces the total number of
// operations on the sorted side by about 20%.
type SortedSide interface {
	insert(l *Level) bool
	delete(l *Level) bool
	snapshot(depth int) ([]Level, bool)
	len() int
}

// newSortedSide creates a new sortedSide.
func newSortedSide(asc bool) SortedSide {
	return &sortedSide{
		slice: []*Level{},
		asc:   asc,
	}
}

/* ==================== sortedSide ======================================= */

// sortedSide is a sorted slice of levels that can be set up to be either
// ascending (asks) or descending (bids).
type sortedSide struct {
	slice []*Level
	asc   bool
}

// len returns the length of the sorted side.
func (s *sortedSide) len() int {
	return len(s.slice)
}

// snapshot returns a copy of the first depth levels in the sorted side.
func (s *sortedSide) snapshot(depth int) ([]Level, bool) {
	snapshot := make([]Level, depth)
	if len(s.slice) < depth {
		return snapshot, false
	}

	for i, l := range s.slice[:depth] {
		snapshot[i] = *l
	}

	return snapshot, true
}

// getSearchFunc returns the search function to use in sort.Search.
// if asc is true, the search function will be ascending, otherwise descending.
func (s *sortedSide) getSearchFunc(l *Level) func(int) bool {
	// convert price string to float64
	price, _ := strconv.ParseFloat(l.Price, 64)

	// ascending
	if s.asc {
		return func(i int) bool {
			priceI, _ := strconv.ParseFloat(s.slice[i].Price, 64)
			return priceI >= price
		}
	}

	// descending
	return func(i int) bool {
		priceI, _ := strconv.ParseFloat(s.slice[i].Price, 64)
		return priceI <= price
	}
}

// insert inserts a level into the sorted side. If an insert operation would
// result in a duplicate, the operation is not performed and false is returned.
func (s *sortedSide) insert(l *Level) bool {
	// find the index where to insert
	f := s.getSearchFunc(l)
	i := sort.Search(len(s.slice), f)

	// don't allow duplicates
	if i != len(s.slice) && s.slice[i].Price == l.Price {
		fmt.Println("already exists")
		return false
	}

	// insert
	s.slice = append(s.slice[:i], append([]*Level{l}, s.slice[i:]...)...)
	return true
}

// delete deletes a level from the sorted side. If the level does not exist,
// the operation is not performed and false is returned.
func (s *sortedSide) delete(l *Level) bool {
	// find the index to delete
	f := s.getSearchFunc(l)
	i := sort.Search(len(s.slice), f)

	// attempt to delete
	if i < len(s.slice) && s.slice[i].Price == l.Price {
		s.slice = append(s.slice[:i], s.slice[i+1:]...)
		return true
	}
	return false
}

/* ==================== EventsBuffer ===================================== */

type EventsBuffer interface {
	// Open opens the buffer channel.
	Open() error

	// IsOpen returns true if the buffer channel is open.
	IsOpen() bool

	// AddIfOpen adds the event to the buffer channel if it is open.
	AddIfOpen(event *bst.SpotMarginDiffDepthEvent) bool

	// CloseAndFlush closes the buffer channel and returns the channel.
	CloseAndFlush() <-chan *bst.SpotMarginDiffDepthEvent
}

/* ==================== eventsBuffer ====================================== */

// newEventsBuffer creates a new event buffer.
func newEventsBuffer(maxSize int, logger *log.Entry) EventsBuffer {
	return &eventsBuffer{
		maxSize: maxSize,
		logger:  logger.WithFields(log.Fields{"__package": "orderBook", "_caller": "eventsBuffer"}),
	}
}

// eventsBuffer is used to store events for a certain time period, and then
// flush them to the orderbook.
type eventsBuffer struct {
	bufferChan chan *bst.SpotMarginDiffDepthEvent
	maxSize    int
	count      int
	logger     *log.Entry
}

// Open creates a new buffer channel. If the buffer channel already exists,
// this function does nothing.
func (b *eventsBuffer) Open() error {
	if b.bufferChan != nil {
		return fmt.Errorf("buffer channel already exists")
	}

	b.logger.Info("opening events buffer")
	b.bufferChan = make(chan *bst.SpotMarginDiffDepthEvent, b.maxSize)
	b.count = 0
	return nil
}

// IsOpen returns true if the buffer channel is open.
func (b *eventsBuffer) IsOpen() bool {
	return b.bufferChan != nil
}

// AddIfOpen adds the event to the buffer channel if it is open.
// Returns true if the event was added to the buffer channel.
func (b *eventsBuffer) AddIfOpen(event *bst.SpotMarginDiffDepthEvent) bool {
	if !b.IsOpen() {
		return false
	}

	if b.count >= b.maxSize {
		<-b.bufferChan
		b.logger.Debug("buffer is full, removing oldest event")
	}

	b.logger.Trace("adding event to buffer")
	b.bufferChan <- event
	b.count++
	return true
}

// CloseAndFlush closes the buffer channel and returns the channel.
func (b *eventsBuffer) CloseAndFlush() <-chan *bst.SpotMarginDiffDepthEvent {
	if b.bufferChan == nil {
		return nil
	}

	b.logger.Info("flushing events buffer")
	close(b.bufferChan)
	defer func() {
		b.bufferChan = nil
	}()
	return b.bufferChan
}
