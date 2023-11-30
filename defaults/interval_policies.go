package defaults

import (
	"time"

	log "github.com/sirupsen/logrus"
)

/* ==================== Warmup Policy ==================================== */

// WarmupPolicy is the base interface for all warmup policies. A warmup policy
// is used to to ramp up interval wait times between consecutive requests.
type WarmupPolicy interface {
	CalcNextInterval() time.Duration
}

/* ==================== Exponential Warmup Policy ========================= */

func newExponentialWarmupPolicy(
	initialInterval int64, interval int64, multiplier float64, logger *log.Entry,
) WarmupPolicy {
	return &exponentialWarmupPolicy{
		currInterval:    0,
		InitialInterval: initialInterval,
		Interval:        interval,
		Multiplier:      multiplier,
		logger:          logger.WithField("_caller", "exponentialWarmupPolicy"),
	}
}

type exponentialWarmupPolicy struct {
	currInterval    int64
	InitialInterval int64 // milliseconds
	Interval        int64
	Multiplier      float64
	logger          *log.Entry
}

func (p *exponentialWarmupPolicy) CalcNextInterval() time.Duration {
	// granularity
	g := time.Millisecond

	// always return initial interval on first call
	if p.currInterval == 0 {
		p.currInterval = p.InitialInterval
		return time.Duration(p.currInterval) * g
	}

	// return interval if it is greater than or equal to max interval
	if p.currInterval >= p.Interval {
		p.currInterval = p.Interval
		return time.Duration(p.currInterval) * g
	}

	// calculate and return next interval
	next := int64(float64(p.currInterval) * p.Multiplier)
	if next > p.Interval {
		next = p.Interval
	}
	p.currInterval = next
	return time.Duration(next) * g
}
