package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ======================== wsClient ===================================== */

// NewWsClient creates a new wsClient
func newWSClient(
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
// this is essentially a factory for websocket streams
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
		th:              wc.th,
		connOpts:        wc.connOpts,
		reconnectPolicy: wc.defaultReconnectPolicy,
		pathFunc:        nil,
		logger:          logger,
	}
}
