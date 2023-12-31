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

	// add rateLimitManager, timeHandler, restClient, and wsClient to client
	c.th = newTimeHandler(c)
	c.rlm = newRateLimitManager(opts.RateLimits, c.th, c.logger)
	c.rc = newRestClient(c.th, c.rlm, apiConfig, c.logger)
	c.wc = newWSClient(c.th, opts.WSConnOpts, opts.WSDefaultReconnectPolicy, c.logger)

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

// SetServerTimeOffset sets the server time offset on client.TimeHandler.
func (c *Client) SetServerTimeOffset(offset int64) {
	c.th.setServerTimeOffset(offset)
}

// GetServerTimeOffset gets the server time offset from client.TimeHandler.
func (c *Client) GetServerTimeOffset() int64 {
	return c.th.getServerTimeOffset()
}

/* ==================== API-Streams Factory ============================== */

func (c *Client) NewSpotUserDataStream() *streams.SpotUserDataStream {
	return streams.NewSpotUserDataStream(c.wc, c.logger)
}

func (c *Client) NewSpotMarginDiffDepth100Stream() *streams.SpotMarginDiffDepthStream {
	return streams.NewSpotMarginDiffDepth100Stream(c.wc, c.logger)
}

func (c *Client) NewSpotMarginAggTradesStream() *streams.SpotMarginAggTradesStream {
	return streams.NewSpotMarginAggTradesStream(c.wc, c.logger)
}

/* ==================== SAPI-Streams Factory ============================= */

func (c *Client) NewMarginUserDataStream() *streams.MarginUserDataStream {
	return streams.NewMarginUserDataStream(c.wc, c.logger)
}

/* ==================== FAPI-Streams Factory ============================= */

func (c *Client) NewFuturesDiffDepth100Stream() *streams.FuturesDiffDepthStream {
	return streams.NewFuturesDiffDepth100Stream(c.wc, c.logger)
}

func (c *Client) NewFuturesAggTradesStream() *streams.FuturesAggTradesStream {
	return streams.NewFuturesAggTradesStream(c.wc, c.logger)
}

/* ==================== API-Services Factory ============================= */

func (c *Client) NewSpotMarginPingService() *services.PingService {
	return services.NewSpotMarginPingService(c.rc, c.logger)
}

func (c *Client) NewSpotMarginServerTimeService() *services.ServerTimeService {
	return services.NewSpotMarginServerTimeService(c.rc, c.logger)
}

func (c *Client) NewSpotMarginDepth100Service() *services.SpotMarginDepthService {
	return services.NewSpotMarginDepth100Service(c.rc, c.logger)
}

func (c *Client) NewSpotMarginDepth5000Service() *services.SpotMarginDepthService {
	return services.NewSpotMarginDepth5000Service(c.rc, c.logger)
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

func (c *Client) NewSpotMarginExchangeInfoService() *services.SpotMarginExchangeInfoService {
	return services.NewSpotMarginExchangeInfoService(c.rc, c.logger)
}

/* ==================== SAPI-Services Factory ============================ */

func (c *Client) NewMarginSystemStatusService() *services.SystemStatusService {
	return services.NewMarginSystemStatusService(c.rc, c.logger)
}

func (c *Client) NewMarginCreateListenKeyService() *services.CreateListenKeyService {
	return services.NewMarginCreateListenKeyService(c.rc, c.logger)
}

func (c *Client) NewMarginPingListenKeyService() *services.PingListenKeyService {
	return services.NewMarginPingListenKeyService(c.rc, c.logger)
}

func (c *Client) NewMarginCloseListenKeyService() *services.CloseListenKeyService {
	return services.NewMarginCloseListenKeyService(c.rc, c.logger)
}

func (c *Client) NewCreateMarginOrderService() *services.CreateMarginOrderService {
	return services.NewCreateMarginOrderService(c.rc, c.logger)
}

/* ==================== FAPI-Services Factory ============================ */

func (c *Client) NewFuturesPingService() *services.PingService {
	return services.NewFuturesPingService(c.rc, c.logger)
}

func (c *Client) NewFuturesServerTimeService() *services.ServerTimeService {
	return services.NewFuturesServerTimeService(c.rc, c.logger)
}

func (c *Client) NewFuturesDepth1000Service() *services.FuturesDepthService {
	return services.NewFuturesDepth1000Service(c.rc, c.logger)
}

func (c *Client) NewFuturesExchangeInfoService() *services.FuturesExchangeInfoService {
	return services.NewFuturesExchangeInfoService(c.rc, c.logger)
}
