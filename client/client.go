package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/services"
	"github.com/svdro/shrimpy-binance/streams"
)

/* ==================== Client =========================================== */

type APIConfig struct {
	apiKey     string
	apiSecret  string
	recvWindow int
}

// NewClient
func NewClient(apiKey string, secretKey string, opts *ClientOptions) *Client {
	apiConfig := &APIConfig{
		apiKey:     apiKey,
		apiSecret:  secretKey,
		recvWindow: 5000,
	}

	c := &Client{
		logger: newLogger(opts),
	}

	// add restClient and timeHandler to client
	c.th = newTimeHandler(c)
	c.rlm = newRateLimitManager(opts.RateLimits, c.th, c.logger)
	c.rc = newRestClient(c.th, c.rlm, apiConfig, c.logger)
	c.wc = newWsClient(c.th, opts.WSConnOpts, opts.WSDefaultReconnectPolicy, c.logger)

	return c
}

// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for creating new REST services and Websocket streams.
type Client struct {
	th     *timeHandler
	rlm    *rateLimitManager
	rc     *restClient
	wc     *wsClient
	logger *log.Entry
}

/* ==================== API-Streams ====================================== */

func (c *Client) NewSpotMarginAggTradesStream() *streams.AggTradesStream {
	return streams.NewSpotMarginAggTradesStream(c.wc, c.logger)
}

/* ==================== FAPI-Streams ====================================== */

func (c *Client) NewFuturesAggTradesStream() *streams.AggTradesStream {
	return streams.NewFuturesAggTradesStream(c.wc, c.logger)
}

/* ==================== API-Services ===================================== */

func (c *Client) NewSpotMarginPingService() *services.PingService {
	return services.NewSpotMarginPingService(c.rc, c.logger)
}

func (c *Client) NewSpotMarginServerTimeService() *services.ServerTimeService {
	return services.NewSpotMarginServerTimeService(c.rc, c.logger)
}

func (c *Client) NewSpotCreateListenKeyService() *services.CreateListenKeyService {
	return services.NewSpotCreateListenKeyService(c.rc, c.logger)
}

func (c *Client) NewSpotPingListenKeyService() *services.PingListenKeyService {
	return services.NewSpotPingListenKeyService(c.rc, c.logger)
}

func (c *Client) NewSpotCloseListenKeyService() *services.CloseListenKeyService {
	return services.NewSpotCloseListenKeyService(c.rc, c.logger)
}

/* ==================== SAPI-Services ==================================== */

func (c *Client) NewMarginCreateListenKeyService() *services.CreateListenKeyService {
	return services.NewMarginCreateListenKeyService(c.rc, c.logger)
}

func (c *Client) NewMarginPingListenKeyService() *services.PingListenKeyService {
	return services.NewMarginPingListenKeyService(c.rc, c.logger)
}

func (c *Client) NewMarginCloseListenKeyService() *services.CloseListenKeyService {
	return services.NewMarginCloseListenKeyService(c.rc, c.logger)
}

/* ==================== FAPI-Services ==================================== */

func (c *Client) NewFuturesPingService() *services.PingService {
	return services.NewFuturesPingService(c.rc, c.logger)
}
