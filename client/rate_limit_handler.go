package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== rateLimitHandler ================================= */

// newRateLimitHandler returns a new rateLimitHandler.
func newRateLimitHandler(c *Client) *rateLimitHandler {
	return &rateLimitHandler{
		c:      c,
		logger: c.logger.WithField("_caller", "rateLimitHandler"),
	}
}

// rateLimitHandler is responsible for keeping track of rate limits and
// throttling requests across services and streams.
type rateLimitHandler struct {
	c      *Client
	logger *log.Entry
}

// RegisterPending
func (h *rateLimitHandler) RegisterPending(sd *common.ServiceDefinition) error {
	// get WeightIP, WeightUID from ServiceMeta.ServiceDefinition
	h.logger.Info("RegisterPending")
	return nil
}

// UnregisterPending
func (h *rateLimitHandler) UnregisterPending(sd *common.ServiceDefinition) {
	h.logger.Info("UnregisterPending")
}

// UpdateUsed
func (h *rateLimitHandler) UpdateUsed(sm *common.ServiceMeta) {
	h.logger.Info("UpdateUsed")
}
