package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

// waitForInterval is a utility that waits for the given interval. If the
// context is cancelled, it returns the context error.
func waitForInterval(ctx context.Context, interval time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(interval):
		return nil
	}
}

// increaseInterval is a utility that increases the interval by the multiplier.
// If the interval is greater than the maxInterval, it returns the maxInterval.
func increaseInterval(current time.Duration, policy BackoffPolicy) time.Duration {
	next := time.Duration(float64(current) * policy.Multiplier)
	if next > policy.MaxInterval {
		next = policy.MaxInterval
	}
	return next
}

// incrConsecEarlyDisconnects is a utility that increments the counter if
// the connection was closed before the minConnDuration. Otherwise, it resets
// the counter.
func incrConsecEarlyDisconnects(current int, t0 time.Time, minConnDuration time.Duration) int {
	if time.Since(t0) < minConnDuration {
		return current + 1
	}
	return 0
}

// stream implements the common.Stream interface.
type stream struct {
	handler         common.StreamHandler
	sm              *common.StreamMeta
	th              common.TimeHandler
	connOpts        WSConnOptions
	reconnectPolicy ReconnectPolicy
	pathFunc        func() string
	pump            *wsPump
	isRunning       bool
	logger          *log.Entry
}

// SetPathFunc sets the pathFunc of the stream. This is necessary because
// stream paths can contain symbols or other dynamic data that may not be known
// when a stream is created.
func (s *stream) SetPathFunc(f func() string) {
	if s.isRunning {
		s.logger.Warn("cannot set pathFunc while stream is running")
		return
	}
	s.pathFunc = f
}

// getURI constructs the URI of the stream.
// e.g. "wss://stream.binance.com:9443/ws/bnbbtc@aggTrade"
func (s *stream) getURI() (url.URL, error) {
	if s.pathFunc == nil {
		return url.URL{}, fmt.Errorf("pathFunc is nil")
	}
	uri := url.URL{Scheme: s.sm.SD.Scheme, Host: string(s.sm.SD.Endpoint), Path: s.pathFunc()}
	return uri, nil
}

// connect creates the websocket connection.
func (s *stream) connect(uri string) (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(uri, nil)
	return conn, err
}

// reconnectWithPolicy attempts to reconnect to the websocket connection
// using the ReconnectPolicy (an exponential backoff strategy).
// It notifies the user of each FAILED attempt.
func (s *stream) reconnectWithPolicy(
	ctx context.Context, uri string, consecEarlyDisconnects int,
) (*websocket.Conn, error) {

	interval := s.reconnectPolicy.BackoffPolicy.InitialInterval

	for i := 0; i < s.reconnectPolicy.MaxAttempts; i++ {
		logger := s.logger.WithField("attempt", i)

		// always return when connection is established
		conn, err := s.connect(uri)
		if err == nil {
			logger.Trace("reconnected")
			return conn, nil
		}

		// on last attempt, return error
		if i == s.reconnectPolicy.MaxAttempts-1 {
			err = s.newWSConnError(err, "failed to reconnect", consecEarlyDisconnects, i+1, false)
			logger.WithError(err).Warn("failed to reconnect")
			s.handler.HandleError(err)
			return nil, err
		}

		// notify and wait for interval
		reason := fmt.Sprintf("failed to reconnect, trying again in %s", interval)
		err = s.newWSConnError(err, reason, consecEarlyDisconnects, i+1, true)
		logger.WithError(err).Debug(reason)
		s.handler.HandleError(err)

		// wait for interval
		logger.WithField("interval", interval).Debug("waiting for next attempt to reconnect")
		if err := waitForInterval(ctx, interval); err != nil {
			return nil, s.newWSConnError(err, "context cancelled", consecEarlyDisconnects, i+1, false)
		}

		// increase the interval for the next attemps
		interval = increaseInterval(interval, s.reconnectPolicy.BackoffPolicy)
	}

	return nil, fmt.Errorf("this is unreachable, but the compiler doesn't know that")
}

// newWSConnError is a utility for reating new WSConnErrors.
func (s *stream) newWSConnError(
	err error, reason string, consecEarlyDisconnects int, connAttempts int, isTransient bool,
) *common.WSConnError {
	return &common.WSConnError{
		Err:                       err,
		Reason:                    reason,
		ConsecEarlyDisconnects:    consecEarlyDisconnects,
		ReconnectionAttempts:      connAttempts,
		MaxConsecEarlyDisconnects: s.reconnectPolicy.MaxConsecEarlyDisconnects,
		MaxReconnectionAttempts:   s.reconnectPolicy.MaxAttempts,
		IsTransient:               isTransient,
	}
}

// isTransientConnError returns true if the error is in principle recoverable.
// it also returns the reason for the error.
func isTransientConnError(err error) (bool, string) {
	switch err := err.(type) {
	// websocket closed
	case *websocket.CloseError:
		return true, "websocket closed"
	// timeout
	case net.Error:
		if err.Timeout() {
			return true, "timeout"
		}
	}

	// websocket manually closed
	if strings.Contains(err.Error(), "websocket: close sent") {
		return false, "intentional close"
	}

	// unexpected error
	return false, "unknown error"
}

// handleConnError checks if the error is transient (recoverable) and forwards
// the correct common.WSConnError to the handler. Retruns true if the error is
// transient, reconnectPolicy is enabled, and the maxConsecEarlyDisconnects is
// not reached, indicating that the stream should attempt to reconnect.
// Always notify the user with a common.WSConnError.
func (s *stream) handleConnError(err error, consecEarlyDisconnects int) bool {

	// Check if the error is Transient
	isTransient, reason := isTransientConnError(err)
	if !isTransient {
		s.handler.HandleError(s.newWSConnError(err, reason, consecEarlyDisconnects, 0, false))
		return false
	}

	// Check if the reconnectPolicy is enabled.
	if !s.reconnectPolicy.Enabled {
		err = s.newWSConnError(err, "reconnectPolicy disabled", consecEarlyDisconnects, 0, false)
		s.handler.HandleError(err)
		return false
	}

	// Check if the maxConsecEarlyDisconnects is reached.
	if consecEarlyDisconnects >= s.reconnectPolicy.MaxConsecEarlyDisconnects {
		err = s.newWSConnError(err, "maxConsecEarlyDisconnects reached", consecEarlyDisconnects, 0, false)
		s.handler.HandleError(err)
		return false
	}

	// Always notify user that a reconnect attempt will be made
	err = s.newWSConnError(err, reason, consecEarlyDisconnects, 0, true)
	s.handler.HandleError(err)
	return true
}

// listen listens for websocket messages and errors. It routes messages to
// StreamHandler.HandleRecv(). On error, it returns the error, expecting the
// caller to handle the error and decide whether or not to reconnect.
func (s *stream) listen(p *wsPump, ctx context.Context) error {
	go p.readPump()
	go p.writePump(ctx)

	ctxCancelled := false

	for {
		select {
		case <-ctx.Done():
			ctxCancelled = true

		case err := <-p.errChan:
			// err will be handled by the caller
			return err

		case msg, ok := <-p.readChan:
			// readChan is closed, continue to wait for an error on errChan
			if !ok {
				continue
			}

			// if context is cancelled, don't handle recv. continue to wait for
			// an error on errChan
			if ctxCancelled {
				continue
			}

			// handle recv
			// if HandleRecv returns an error, it's always a common.WSHandlerError.
			// If the error is not fatal, notify the user that it occured, but don't
			// take any action. If the error is fatal, return it to the caller (Run()).
			// Run() should shut down the stream and notify the user.
			s.logger.WithField("msg", string(msg)).Trace("trying to handle recv")
			if err := s.handler.HandleRecv(msg, s.th.TSLNow(), s.th.TSSNow()); err != nil {
				if err.IsFatal {
					return err
				}
				s.handler.HandleError(err)
			}
		}
	}
}

// cleanupPump cleans up the wsPump, and sets stream.pump to nil.
// close wsPump.writeChan here. stream.Do() is the only function that
// writes to writeChan. First disable stream.Do() from writing to the
// writeChan by setting s.pump to nil. Then close the writeChan.
func (s *stream) cleanupPump() {
	pump := s.pump
	s.pump = nil
	close(pump.writeChan)
}

// Run starts the websocket stream and handles reconnections.
func (s *stream) Run(ctx context.Context) {

	// set isRunning flag to true, and defer setting it to false
	s.isRunning = true
	defer func() {
		s.isRunning = false
	}()

	// if the setPathFunc is not set, don't run the stream
	uri, err := s.getURI()
	if err != nil {
		s.handler.HandleError(s.newWSConnError(err, "failed to get URI", 0, 0, false))
		return
	}

	// create websocket connection, init consecEarlyDisconnects counter
	conn, err := s.connect(uri.String())
	var consecEarlyDisconnects int = 0

	//conn, err := s.connect()
	if err != nil {
		s.handler.HandleError(s.newWSConnError(err, "failed to connect", 0, 0, false))
		return
	}

	for {
		// make a new wsPump, take the startTime  and listen()
		s.pump = newWsPump(conn, s.connOpts.WSWriteWait, s.connOpts.WSPongWait, s.connOpts.WSPingPeriod, s.logger)
		t0 := time.Now()
		err = s.listen(s.pump, ctx)

		// incr or reset consecEarlyDisconnects counter based on MinConnDuration
		consecEarlyDisconnects = incrConsecEarlyDisconnects(consecEarlyDisconnects, t0, s.reconnectPolicy.MinConnDuration)

		// cleanup s.pump and set it to nil
		s.cleanupPump()

		// if a common.WSHanlderError is returned, it's always fatal.
		// shutdown the stream and notify the user.
		if err, ok := err.(*common.WSHandlerError); ok {
			s.handler.HandleError(err)
			return
		}

		// if it's not a common.WSHandlerError, it's an error that occured in
		// wsPump, in which case it should be parsed to a common.WSConnError
		// and handled accordingly.
		if !s.handleConnError(err, consecEarlyDisconnects) {
			return
		}

		// attempt to reconnect, Notify user!
		if conn, err = s.reconnectWithPolicy(ctx, uri.String(), consecEarlyDisconnects); err != nil {
			return
		}
	}
}

func (s *stream) Do(req common.WSRequest) {
	// don't send if stream is not running
	if !s.isRunning {
		s.logger.Warn("stream is not running. cannot send request")
		return
	}

	// don't send if pump is nil (e.g. stream is currently reconnecting)
	if s.pump == nil {
		s.logger.Warn("pump is nil. cannot send request")
		return
	}

	// TODO: this is not necessarily fatal. Make a transient error and put it
	// on the errChan, but don't shut down run.
	if err := s.handler.HandleSend(req); err != nil {
		s.logger.WithError(err).Warn("error handling send")
		return
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("error marshaling request: %s", err)
	}

	s.pump.writeChan <- reqBytes
}

// SetReconnectPolicy allows the user to set a reconnect policy other than
// the default policy. This may be useful for crucial streams that should
// never be disconnected.
func (s *stream) SetReconnectPolicy(policy ReconnectPolicy) {
	if s.isRunning {
		s.logger.Warn("cannot set reconnect policy while stream is running")
		return
	}
	s.reconnectPolicy = policy
}
