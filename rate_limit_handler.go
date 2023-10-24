package binance

import (
	"sync"
)

/* ==================== RateLimitHandler ================================= */

type rateLimitCounter struct {
	used                int
	pending             int
	max                 int
	intervalNumber      int
	intervalType        BIRateLimitIntervalType
	tsNanoLastReqHeader int64 // value from binance "Date" header; has second granularity
	tsNanoCurrInterval  int64 // first timesamp in current interval
}

func newRateLimitHandler(c *Client, endpointType BIEndpointType) *RateLimitHandler {
	return &RateLimitHandler{endpointType: endpointType}
}

// RateLimitHandler
type RateLimitHandler struct {
	mu           sync.Mutex
	endpointType BIEndpointType
	rateLimitIP  rateLimitCounter
	rateLimitUID rateLimitCounter
}

func (h *RateLimitHandler) RegisterPending(sd *serviceDefinition) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return nil
}

func (h *RateLimitHandler) UnregisterPending(sd *serviceDefinition) {
	h.mu.Lock()
	defer h.mu.Unlock()
}

func (h *RateLimitHandler) UpdateUsed(rh *responseHeader) {
	h.mu.Lock()
	defer h.mu.Unlock()
}
