package client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ======================== Errors (shrimpy-binance/common) ============== */

// WSCloseError is an error returned when the websocket is closed
type WSCloseError struct {
	Code int
}

func (e *WSCloseError) Error() string {
	return fmt.Sprintf("websocket closed with code %d", e.Code)
}

/* ======================== wsClient ===================================== */

// NewWsClient creates a new wsClient
func newWsClient(
	th common.TimeHandler,
	connOpts WSConnOptions,
	defaultReconnectPolicy ReconnectPolicy,
	logger *log.Entry,
) *wsClient {
	logger = logger.WithField("_caller", "wsClient")
	return &wsClient{
		th:                     th,
		connOpts:               connOpts,
		defaultReconnectPolicy: defaultReconnectPolicy,
		logger:                 logger,
	}
}

// wsClient handles the creation of websocket streams
// wsClient is essentially a factory for websocket streams
type wsClient struct {
	th                     common.TimeHandler // needed for creating timestamps
	connOpts               WSConnOptions      // websocket connection options
	defaultReconnectPolicy ReconnectPolicy    // every stream has the same reconnect policy
	logger                 *log.Entry         // logger
}

// NewStream creates a common.Stream
func (wc *wsClient) NewStream(
	sm *common.StreamMeta, handler common.StreamHandler, logger *log.Entry) common.Stream {
	return &stream{
		handler:         handler,
		sm:              sm,
		connOpts:        wc.connOpts,
		reconnectPolicy: wc.defaultReconnectPolicy,
		pathFunc:        nil,
		logger:          logger,
	}
}
