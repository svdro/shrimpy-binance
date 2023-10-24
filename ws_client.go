package binance

import log "github.com/sirupsen/logrus"

/* ==================== WSClient ========================================= */

// WSClient is responsible for interacting with ALL Binance WebSocket APIs.
type WSClient struct {
	c      *Client
	logger *log.Entry
}
