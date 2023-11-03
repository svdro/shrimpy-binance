package client

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
	"github.com/svdro/shrimpy-binance/services"
)

/* ==================== ClientOptions ==================================== */

// ClientOptions
type ClientOptions struct {
	LogReportCaller bool          // default: false
	LogFormatter    log.Formatter // default: &log.TextFormatter{}
	LogLevel        log.Level     // default: log.PanicLevel
	LogOutput       io.Writer     // default: os.Stderr
}

// DefaultClientOptions returns a new ClientOptions with default values.
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		LogReportCaller: false,
		LogFormatter: &log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "15:04:05.000000",
		},
		LogLevel:  log.PanicLevel,
		LogOutput: os.Stderr,
	}
}

/* ==================== ClientUtils ====================================== */

// newLogger returns a new *log.Entry with the given options.
func newLogger(opts *ClientOptions) *log.Entry {

	logger := log.New()

	if opts.LogFormatter != nil {
		logger.SetFormatter(opts.LogFormatter)
	}

	if opts.LogOutput != nil {
		logger.SetOutput(opts.LogOutput)
	}

	if opts.LogReportCaller {
		logger.SetReportCaller(true)
	}

	logger.SetLevel(opts.LogLevel)

	//logger.SetFormatter(options.formatter)
	return logger.WithFields(log.Fields{
		"__package": "shrimpy-binance",
		"_caller":   "client",
	})
}

/* ==================== Client =========================================== */

// a synonym for package that starts on "" is
func NewClient(apiKey string, secretKey string, opts *ClientOptions) *Client {
	c := &Client{
		apiKey:    apiKey,
		apiSecret: secretKey,
		logger:    newLogger(opts),
	}

	c.rc = newRestClient(c)
	c.th = newTimeHandler(c)
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
	logger     *log.Entry
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
