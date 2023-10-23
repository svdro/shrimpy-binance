package binance_api

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

type Client struct {
	apiKey    string
	apiSecret string
}
