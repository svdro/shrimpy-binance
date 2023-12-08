package defaults

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	binance "github.com/svdro/shrimpy-binance"
	bsv "github.com/svdro/shrimpy-binance/services"
	bst "github.com/svdro/shrimpy-binance/streams"
)

/* ==================== OrderBookService ================================= */

// OrderBookService
// TODO: rethink how to publish snapshots to the rest of the application.
// a callback is probably a good idea, but there should be a way to
// communicate when the orderbook is out of sync, and when
// the connection is lost/ reconnecting, when the connection is restored.
// TODO: the user should be able to set snapshot depth (this is currently hardcoded)
// TODO: snapshots should contain information on the symbol and the endpoint type.
// TODO: make this generic so that it can also be used for the futures api.
type OrderBookService interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

// NewOrderbookService creates a new instance of OrderBookService.
func NewOrderBookService(
	client *binance.Client,
	symbol string,
	syncWithSnapshotInterval time.Duration,
	onSnaphsotCallback func(*OrderBookSnapshot),
	logger *log.Entry,
) OrderBookService {
	logger = logger.WithFields(log.Fields{"__package": "orderBook", "_caller": "NewOrderBookService", "_symbol": symbol})
	return &orderBookService{
		c:                        client,
		symbol:                   symbol,
		eventsBuffer:             newEventsBuffer(1000, logger),
		restErrorPolicy:          &defaultErrorPolicy{},
		wsErrorPolicy:            &defaultWSErrorPolicy{},
		syncWithSnapshotInterval: syncWithSnapshotInterval,
		onSnaphsotCallback:       onSnaphsotCallback,
		logger:                   logger,
	}
}

/* ==================== orderBookService ================================= */

// orderBookManager is responsible for maintaining the orderBook for a given symbol.
type orderBookService struct {
	c                        *binance.Client
	symbol                   string
	book                     OrderBook
	eventsBuffer             EventsBuffer
	restErrorPolicy          ErrorPolicy
	wsErrorPolicy            ErrorPolicy
	stream                   *bst.SpotMarginDiffDepthStream
	respChan                 chan *bsv.SpotMarginDepthResponse
	restErrChan              chan error
	syncWithSnapshotInterval time.Duration
	onSnaphsotCallback       func(*OrderBookSnapshot)
	inSync                   bool
	logger                   *log.Entry
}

// doDepthRequest blocks until all conditions for requesting a new snapshot
// are met and then proceeds to requests a new snapshot from the REST API.
func (s *orderBookService) doDepthRequest(
	ctx context.Context,
	restThrottleInterval time.Duration,
) {
	// wait for restThrottleInterval
	select {
	case <-ctx.Done():
		return
	case <-time.After(restThrottleInterval):
	}

	// wait for stream to be ready, put an error on restErrChan if
	// WaitForConnection() returns false
	select {
	case <-ctx.Done():
		return
	case isConnected := <-s.stream.WaitForConnection():
		if !isConnected {
			err := fmt.Errorf("doDepthRequest: error waiting for ws stream to connect (isConnected = %v)", isConnected)
			s.logger.WithError(err).Error("doDepthRequest: error connecting to ws stream")
			s.restErrChan <- err
			return
		}
	}

	// do the request, pass resp, err to the appropriate channel.
	resp, err := s.c.NewSpotMarginDepth5000Service().WithSymbol(s.symbol).Do(ctx)
	if err != nil {
		s.logger.WithError(err).Error("doDepthRequest: error doing depth request")
		s.restErrChan <- err
		return
	}

	s.respChan <- resp
}

// responsible for opening the eventsBuffer and starting a new request in a
// new goroutine. This will only allow for one pending request at a time.
// If one is already pending this will return an error.
func (s *orderBookService) startDepthRequest(
	ctx context.Context,
	restThrottleInterval time.Duration,
	logLevel log.Level,
) error {

	// this opens the eventsBuffer or returns an error if it is already open
	if err := s.eventsBuffer.Open(); err != nil {
		s.logger.WithError(err).Log(logLevel, "startDepthRequest: snapshot request already pending")
		return err
	}

	// start a new request in a new goroutine
	go s.doDepthRequest(ctx, restThrottleInterval)
	return nil
}

// syncBookWithEventsBuffer is used to synchronize an orderbook with buffered
// events. It returns true if the orderbook is in sync and false otherwise.
func (s *orderBookService) syncBookWithEventsBuffer(book OrderBook, bufferChan <-chan *bst.SpotMarginDiffDepthEvent) bool {
	for event := range bufferChan {
		if inSync := book.updateFromDepthEvent(event); !inSync {
			return false
		}
	}
	return true
}

// handleDepthResponse does two things:
// It closes and flushes the eventsBuffer.
// It creates a new orderBook from the response and syncs it with the events
// buffer and the old ordeBook's timestamps.
// If the orderbook is out of sync after the update, it attempts to request a
// snapshot and returns.
func (s *orderBookService) handleDepthResponse(ctx context.Context, resp *bsv.SpotMarginDepthResponse) {
	// close and flush the eventsBuffer (to signal that the depth request is complete)
	b := s.eventsBuffer.CloseAndFlush()

	// make a new orderBook, populate it with the snapshot
	book := newOrderBook(s.logger)
	book.initFromDepthResponse(resp)

	// sync book with the eventsBuffer. If this fails, and the orderbook is out
	// of sync, attempt to request a snapshot, and return.
	if inSync := s.syncBookWithEventsBuffer(book, b); !inSync {
		s.logger.Warn("handleDepthResponse: orderbook is out of sync")
		if !s.eventsBuffer.IsOpen() {
			s.logger.Warn("handleDepthResponse: requesting new snapshot")
			s.startDepthRequest(ctx, 0, log.DebugLevel)
		}
		return
	}

	// if the orderBookService has an existing orderbook, transfer timestamps to
	// the new one where appropriate.
	if s.book != nil {
		if snap := s.book.takeSnapshot(-1); snap != nil {
			book.importTimestampsFromSnapshot(snap)
		} else {
			s.logger.Warn("handleDepthResponse: error taking snapshot")
		}
	}

	// replace the current orderbook with the new one
	s.logger.Info("setting new orderbook")
	s.book = book
}

// handleDepthEvent does two things:
// It adds the event to the eventsBuffer if it is open.
// It synchronizes the orderbook with the event if an orderbook exists.
// If the orderbook loses sync after the update, it attempts to request a
// snapshot and returns.
func (s *orderBookService) handleDepthEvent(ctx context.Context, event *bst.SpotMarginDiffDepthEvent) {

	// if the eventsBuffer is open, add the event to it
	s.eventsBuffer.AddIfOpen(event)

	// if no orderbook exists, return. don't request a snapshot here.
	if s.book == nil {
		s.logger.Trace("handleDepthEvent: orderbook is nil")
		return
	}

	// try to update orderbook with event. If this fails and the orderbook is
	// out of sync after the update, attempt to request a snapshot and return.
	if inSync := s.book.updateFromDepthEvent(event); !inSync {
		s.logger.Debug("handleDepthEvent: orderbook is out of sync")

		if !s.eventsBuffer.IsOpen() {
			s.logger.Warn("handleDepthEvent: requesting new snapshot")
			s.startDepthRequest(ctx, 0, log.DebugLevel)
		}
		return
	}
}

// handleRestError determines if the rest error is transient. If yes, it
// starts a new depth request with the appropriate delay.
func (s *orderBookService) handleRestError(ctx context.Context, err error) error {
	// always close the eventsBuffer to signal that the snapshot request is complete
	_ = s.eventsBuffer.CloseAndFlush()

	// consult error policy to determine if the rest error is transient and
	// whether the next request should be delayed and by how much (d).
	t, d, r := s.restErrorPolicy.Handle(err)
	s.logger.WithFields(log.Fields{"reason": r, "isTransient": t, "waitDuration": d}).WithError(err).Error("handleRestError")

	// if the error is transient, start a new depth request. None should be
	// pending at this point so this should never fail.
	if t {
		s.logger.WithField("restThrottleInterval", d).Warn("handleRestError: retrying")
		s.startDepthRequest(ctx, d, log.FatalLevel)
		return nil
	}
	return err
}

// handleWSError determines if the ws error is transient. If yes, it attempts
// to start a new depth request.
// NOTE: when a stream error occurs while a snapshot request is pending (and
// already dispatched) it is possible that the orderbook created from the
// pending request cannot sync until the stream reconnects. This leads to the
// new orderbook mistakenly believing it's in sync until the first update after
// reconnection. At that point, it will detect the sync discrepancy and initiate
// a new snapshot request. This scenario is quite likely, especially in situations
// where a connection issue triggers both a REST error (leading to a new
// snapshot request) and a subsequent stream disconnection.
func (s *orderBookService) handleWSError(ctx context.Context, err error) error {
	// consult error policy to determine if the ws error is transient.
	t, _, r := s.wsErrorPolicy.Handle(err)
	s.logger.WithFields(log.Fields{"reason": r, "isTransient": t}).WithError(err).Error("handleWSError")

	// Try to start a depth request if none is pending.
	if t {
		if !s.eventsBuffer.IsOpen() {
			s.logger.WithField("case", "streamErrChan").Warn("stream has disconnected. Requesting new snapshot.")
			s.startDepthRequest(ctx, 0, log.DebugLevel)
		}
		return nil
	}

	return err
}

// handleTick tries to start a new depth request if none is pending.
func (s *orderBookService) handleTick(ctx context.Context) {
	if !s.eventsBuffer.IsOpen() {
		s.logger.WithField("case", "ticker").Debug("requesting new snapshot")
		s.startDepthRequest(ctx, 0, log.DebugLevel)
	}
}

// Run starts the orderBookService.
func (s *orderBookService) Run(ctx context.Context, wg *sync.WaitGroup) {
	// runCtx is needed for closing the stream when an unexpected error occurs
	runCtx, cancel := context.WithCancel(ctx)

	// initialize channels
	s.restErrChan = make(chan error)
	s.respChan = make(chan *bsv.SpotMarginDepthResponse)

	// cleanup
	defer func() {
		cancel()
		close(s.restErrChan)
		close(s.respChan)
		wg.Done()
	}()

	// create a ticker that requests a snapshot from the REST API periodically
	ticker := time.NewTicker(s.syncWithSnapshotInterval)
	defer ticker.Stop()

	// start the stream
	s.stream = s.c.NewSpotMarginDiffDepth100Stream().SetSymbol(s.symbol)
	go s.stream.Run(runCtx)

	// do the initial snapshot request
	if s.startDepthRequest(runCtx, 0, log.InfoLevel) != nil {
		s.logger.Error("Run: error starting initial snapshot request")
	}

	// this exits when ctx or runCtx is canceled. either one will shut down the
	// stream, causing a non-transient error in handleWSError, which will then return.
	for {
		select {
		// request snapshot from REST API periodically
		case <-ticker.C:
			s.handleTick(runCtx)

		// handle rest errors
		// If the error is not transient, cancel the runCtx. This will cause the
		// stream to close, which in turn will cause Run() to return.
		case err := <-s.restErrChan:
			if err := s.handleRestError(runCtx, err); err != nil {
				s.logger.WithError(err).Error("Run: terminal error handling rest error")
				cancel()
			}

		// handle ws stream errors
		case err := <-s.stream.Handler.ErrChan:
			if err := s.handleWSError(runCtx, err); err != nil {
				s.logger.WithError(err).Error("Run: terminal error handling ws error")
				return
			}

		// handle depth response
		case resp := <-s.respChan:
			s.handleDepthResponse(runCtx, resp)

		// handle depthEvent and publish snapshots
		case event := <-s.stream.Handler.EventChan:
			s.handleDepthEvent(runCtx, event)

			if s.book == nil {
				s.logger.Warn("Run: orderbook is nil. Cannot take snapshot.")
				continue
			}

			snap := s.book.takeSnapshot(5)
			s.onSnaphsotCallback(snap)
		}
	}
}
