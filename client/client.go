package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
	"github.com/svdro/shrimpy-binance/services"
)

/* ==================== Client =========================================== */

func NewClient(apiKey string, secretKey string) *Client {
	c := &Client{
		apiKey:    apiKey,
		apiSecret: secretKey,
		logger:    log.New(),
	}

	c.rc = newRestClient(c)
	c.th = &timeHandler{c: c}
	c.rlhs = map[common.BIEndpointType]*rateLimitHandler{
		common.EndpointTypeAPI:  newRateLimitHandler(c),
		common.EndpointTypeSAPI: newRateLimitHandler(c),
		common.EndpointTypeFAPI: newRateLimitHandler(c),
	}
	return c
}

// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for creating new REST services and Websocket streams.
type Client struct {
	th *timeHandler
	//rlh    *rateLimitHandler
	rlhs       map[common.BIEndpointType]*rateLimitHandler // have one for each endpointType
	rc         *restClient
	apiKey     string
	apiSecret  string
	recvWindow int
	logger     *log.Logger
}

// getRlh returns the rateLimitHandler associated with endpointType.
// when rlh with endpointType is not found, panic.
func (c *Client) getRlh(endpointType common.BIEndpointType) *rateLimitHandler {
	rlh, ok := c.rlhs[endpointType]

	if !ok {
		c.logger.Panicf("rateLimitHander with endpointType %s not found", endpointType)
	}
	return rlh
}

/* ==================== API-Services ===================================== */

func (c *Client) NewSpotMarginPingService() *services.PingService {
	return services.NewSpotMarginPingService(c.rc)
}

func (c *Client) NewSpotMarginServerTimeService() *services.ServerTimeService {
	return services.NewSpotMarginServerTimeService(c.rc)
}

func (c *Client) NewSpotCreateListenKeyService() *services.CreateListenKeyService {
	return services.NewSpotCreateListenKeyService(c.rc)
}

func (c *Client) NewSpotPingListenKeyService() *services.PingListenKeyService {
	return services.NewSpotPingListenKeyService(c.rc)
}

func (c *Client) NewSpotCloseListenKeyService() *services.CloseListenKeyService {
	return services.NewSpotCloseListenKeyService(c.rc)
}

/* ==================== SAPI-Services ==================================== */

func (c *Client) NewMarginCreateListenKeyService() *services.CreateListenKeyService {
	return services.NewMarginCreateListenKeyService(c.rc)
}

func (c *Client) NewMarginPingListenKeyService() *services.PingListenKeyService {
	return services.NewMarginPingListenKeyService(c.rc)
}

func (c *Client) NewMarginCloseListenKeyService() *services.CloseListenKeyService {
	return services.NewMarginCloseListenKeyService(c.rc)
}

/* ==================== FAPI-Services ==================================== */

func (c *Client) NewFuturesPingService() *services.PingService {
	return services.NewFuturesPingService(c.rc)
}
