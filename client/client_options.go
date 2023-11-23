package client

import (
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== ClientOptions ==================================== */
// BackoffPolicy
type BackoffPolicy struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// ReconnectPolicy
// Fields:
//   - Enabled: whether or not the reconnect policy is enabled.
//   - MaxAttempts: the max number of attempts to reconnect.
//   - BackoffPolicy: the backoff policy to use.
//   - MinConnDuration: the minimum amount of time a connection must be open
//     to not be considered an early disconnect.
//   - MaxConsecEarlyDisconnects: the max number of consecutive early
//     disconnects before the reconnect policy is disabled.
type ReconnectPolicy struct {
	Enabled                   bool
	MaxAttempts               int
	BackoffPolicy             BackoffPolicy
	MinConnDuration           time.Duration
	MaxConsecEarlyDisconnects int
}

// RateLimit
type RateLimit struct {
	EndpointType          common.BIEndpointType
	RateLimitType         common.BIRateLimitType
	RateLimitIntervalType common.BIRateLimitIntervalType
	RateLimitIntervalNum  int
	Limit                 int
}

// ClientOptions
type ClientOptions struct {
	LogReportCaller          bool          // default: false
	LogFormatter             log.Formatter // default: &log.TextFormatter{}
	LogLevel                 log.Level     // default: log.PanicLevel
	LogOutput                io.Writer     // default: os.Stderr
	RateLimits               []RateLimit   // default: []RateLimit{}
	WSConnOpts               WSConnOptions
	WSDefaultReconnectPolicy ReconnectPolicy
}

type WSConnOptions struct {
	WSWriteWait  time.Duration
	WSPongWait   time.Duration
	WSPingPeriod time.Duration
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
		RateLimits: []RateLimit{
			{
				EndpointType:          common.EndpointTypeAPI,
				RateLimitType:         common.RateLimitTypeIP,
				RateLimitIntervalType: common.IntervalMinute,
				RateLimitIntervalNum:  1,
				Limit:                 6000,
			},
		},
		WSConnOpts: WSConnOptions{
			WSWriteWait:  3 * time.Second,
			WSPongWait:   5 * time.Second,
			WSPingPeriod: (5 * time.Second * 8) / 10,
		},
		WSDefaultReconnectPolicy: ReconnectPolicy{
			Enabled:                   false,
			MaxAttempts:               0,
			BackoffPolicy:             BackoffPolicy{},
			MinConnDuration:           0,
			MaxConsecEarlyDisconnects: 0,
		},
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
