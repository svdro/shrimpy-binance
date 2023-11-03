package client

import (
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== rateLimitHandler ================================= */

// newRateLimitHandler returns a new RateLimitHandler.
func newRateLimitHandler(c *Client) *rateLimitHandler {
	return &rateLimitHandler{
		c: c,
	}
}

// rateLimitHandler is responsible for keeping track of rate limits and
// throttling requests across services and streams.
type rateLimitHandler struct {
	c *Client
}

func (h *rateLimitHandler) RegisterPending(sd *common.ServiceDefinition) error {
	// get WeightIP, WeightUID from ServiceMeta.ServiceDefinition
	return nil
}

func (h *rateLimitHandler) UnregisterPending(sd *common.ServiceDefinition) {}

func (h *rateLimitHandler) UpdateUsed(sm *common.ServiceMeta) {}
