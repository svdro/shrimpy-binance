package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/services"
)

/* ==================== Client =========================================== */

// NewClient
func NewClient(apiKey string, secretKey string, opts *ClientOptions) *Client {
	c := &Client{
		apiKey:    apiKey,
		apiSecret: secretKey,
		logger:    newLogger(opts),
	}

	// add restClient and timeHandler to client
	c.rc = newRestClient(c)
	c.th = newTimeHandler(c)
	c.rlm = newRateLimitManager(opts.RateLimits, c.th, c.logger)

	return c
}

// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for creating new REST services and Websocket streams.
type Client struct {
	th         *timeHandler
	rlm        *rateLimitManager
	rc         *restClient
	apiKey     string
	apiSecret  string
	recvWindow int
	logger     *log.Entry
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
