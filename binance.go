package binance

import "github.com/svdro/shrimpy-binance/client"

/* ==================== Client =========================================== */

// Client is a convenience wrapper around client.Client.
// Client is the main entry-point for interacting with shrimpy-binance.
// It is responsible for creating new REST services and Websocket streams.
type Client = client.Client

// BinanceClient returns a pointer to a new client.Client.
func BinanceClient(apiKey string, secretKey string) *Client {
	return client.NewClient(apiKey, secretKey)
}
