package binance

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

/* ==================== Client =========================================== */

// NewClient returns a new Client.
func NewClient(apiKey, apiSecret string) *Client {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	c := &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		logger:     logger,
		recvWindow: 5000,
	}

	c.timeHandler = newTimeHandler(c)
	c.restClient = c.newRESTClient()
	c.wsClient = c.newWSClient()
	c.rlhs = map[BIEndpointType]*RateLimitHandler{
		endpointTypeAPI:  newRateLimitHandler(c, endpointTypeAPI),
		endpointTypeSAPI: newRateLimitHandler(c, endpointTypeSAPI),
		endpointTypeFAPI: newRateLimitHandler(c, endpointTypeFAPI),
	}
	c.logger.Debug("NewClient")
	return c
}

// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for managing time across all services.
// Client can make new REST sevices and WebSocket streams.
type Client struct {
	*timeHandler
	restClient *RESTClient
	wsClient   *WSClient
	rlhs       map[BIEndpointType]*RateLimitHandler // have one for each endpointType
	logger     *log.Logger
	apiKey     string
	apiSecret  string
	recvWindow int
}

// newTimeHandler returns a new TimeHandler.
func newTimeHandler(c *Client) *timeHandler {
	return &timeHandler{c: c, logger: c.logger.WithField("client", "TimeHandler")}
}

// newWSClient returns a new WSClient.
func (c *Client) newWSClient() *WSClient {
	return &WSClient{
		c:      c,
		logger: c.logger.WithField("client", "WSClient"),
	}
}

// newRESTClient returns a new RESTClient.
func (c *Client) newRESTClient() *RESTClient {
	return &RESTClient{
		c:          c,
		httpClient: http.DefaultClient,
		logger:     c.logger.WithField("client", "RESTClient"),
	}
}
