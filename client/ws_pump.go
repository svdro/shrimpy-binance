package client

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// newWsPump creates a new wsPump
func newWsPump(
	conn *websocket.Conn, wsWriteWait, wsPongWait, wsPingPeriod time.Duration, logger *log.Entry,
) *wsPump {
	caller := "wsPump"
	if c, ok := logger.Data["_caller"]; ok {
		caller = c.(string) + " > " + caller
	}

	return &wsPump{
		wsWriteWait:  wsWriteWait,
		wsPongWait:   wsPongWait,
		wsPingPeriod: wsPingPeriod,
		errChan:      make(chan error, 1),
		writeChan:    make(chan []byte, 256),
		readChan:     make(chan []byte, 256), // better safe than sorry
		conn:         conn,
		logger:       logger.WithField("_caller", caller),
	}
}

// wsPump is responsible for reading and writing to the websocket
// It is also responsible for keeping track of the websocket connection,
// closing it when necessary and propagating websocket errors to the
// corresponding WsStream.
type wsPump struct {
	connOpts     WSConnOptions // websocket connection options
	wsWriteWait  time.Duration // time to wait for a write
	wsPongWait   time.Duration // time to wait for a pong
	wsPingPeriod time.Duration // time between pings
	errChan      chan error    // channel for writing errors
	writeChan    chan []byte   // channel for writing to websocket
	readChan     chan []byte   // channel for reading from websocket
	conn         *websocket.Conn
	logger       *log.Entry
}

// readPump reads messages from the websocket and sends them to the readChan.
// It also handles pongs and closes the websocket when necessary.
// It exits when a close message is received or when an error occurs.
func (p *wsPump) readPump() {
	var closeErr error

	defer func() {
		p.logger.Debug("exiting readPump")
		p.conn.Close()
		p.errChan <- closeErr
		close(p.readChan)
		close(p.errChan)
	}()

	p.conn.SetReadDeadline(time.Now().Add(p.wsPongWait))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(p.wsPongWait))
		p.logger.Trace("pong")
		return nil
	})

	for {
		_, msg, err := p.conn.ReadMessage()

		if err != nil {
			closeErr = err
			return
		}

		p.readChan <- msg
	}
}

// writeWithDeadline writes a message to the websocket with a deadline
func (p *wsPump) writeWithDeadline(msgType int, msg []byte) error {
	p.conn.SetWriteDeadline(time.Now().Add(p.wsWriteWait))
	err := p.conn.WriteMessage(msgType, msg)
	if err != nil {
		p.logger.WithError(err).WithField("msg", string(msg)).Debug("error writing to websocket")
	}
	return err
}

// writePump writes messages from the writeChan to the websocket
// It also handles pings and closes the websocket when necessary
// It exits when the context is done or when an error occurs.
// NOTE: the error handling logic is that when an error occurs, writePump logs
// the error and returns, no more ping messages are sent, and readPump will
// throw an error when it doesn't receive a pong message in time and return.
func (p *wsPump) writePump(ctx context.Context) {
	pingTicker := time.NewTicker(p.wsPingPeriod)
	defer func() {
		pingTicker.Stop()
	}()

	for {
		select {

		// write messages to websocket
		case msg, ok := <-p.writeChan:
			if !ok {
				p.logger.Trace("writeChan closed. exiting writePump")
				return
			}
			if err := p.writeWithDeadline(websocket.TextMessage, msg); err != nil {
				p.logger.WithError(err).Warn("error writing to websocket. exiting writePump")
				return
			}

		// send ping messages
		case <-pingTicker.C:
			if err := p.writeWithDeadline(websocket.PingMessage, nil); err != nil {
				p.logger.WithError(err).Warn("error writing ping message to websocket. exiting writePump")
				return
			}
			p.logger.Trace("ping")

		// close websocket when context is done
		case <-ctx.Done():
			p.logger.Debug("context done, writing close message. exiting writePump")
			closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			if err := p.writeWithDeadline(websocket.CloseMessage, closeMsg); err != nil {
				p.logger.WithError(err).Warn("error writing close message to websocket")
			}
			return
		}
	}
}
