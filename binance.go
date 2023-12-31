package binance

import (
	"github.com/svdro/shrimpy-binance/client"
	"github.com/svdro/shrimpy-binance/services"
)

/* ==================== Client =========================================== */

// Client is a convenience wrapper around client.Client.
// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for creating new REST services and Websocket streams.
type Client = client.Client

// ClientOptions is a convenience wrapper around client.ClientOptions.
type ClientOptions = client.ClientOptions

// RateLimit is a convenience wrapper around client.RateLimit.
type RateLimit = client.RateLimit

// WSConnOptions is a convenience wrapper around client.WSConnOptions.
type WSConnOptions = client.WSConnOptions

// ReconnectPolicy is a convenience wrapper around client.ReconnectPolicy.
type ReconnectPolicy = client.ReconnectPolicy

// BackoffPolicy is a convenience wrapper around client.BackoffPolicy.
type BackoffPolicy = client.BackoffPolicy

// ServiceBaseResponse is a convenience wrapper around services.ServiceBaseResponse.
type ServiceBaseResponse = services.ServiceBaseResponse

// BinanceClient returns a pointer to a new client.Client.
func BinanceClient(apiKey string, secretKey string) *Client {
	return client.NewClient(apiKey, secretKey, client.DefaultClientOptions())
}

// BinanceClientWithOptions returns a pointer to a new client.Client.
func BinanceClientWithOptions(apiKey string, secretKey string, opts *ClientOptions) *Client {
	return client.NewClient(apiKey, secretKey, opts)
}
