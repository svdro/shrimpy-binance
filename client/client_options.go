package client

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== ClientOptions ==================================== */

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
	LogReportCaller bool          // default: false
	LogFormatter    log.Formatter // default: &log.TextFormatter{}
	LogLevel        log.Level     // default: log.PanicLevel
	LogOutput       io.Writer     // default: os.Stderr
	RateLimits      []RateLimit   // default: []RateLimit{}
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
