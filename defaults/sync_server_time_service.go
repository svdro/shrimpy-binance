package defaults

import (
	"context"
	"time"

	"sync"

	log "github.com/sirupsen/logrus"
	binance "github.com/svdro/shrimpy-binance"
	bc "github.com/svdro/shrimpy-binance/common"
	bsv "github.com/svdro/shrimpy-binance/services"
)

/* ==================== helpers ========================================== */

// calcOffAndRTT calculates the offset and round trip time from the 4
// timestamps (TSLSent, TSSServer, TSSServer, TSLRecv) in a ServerTimeResponse.
func calcOffAndRTT(t0, t1, t2, t3 bc.TSNano) (int64, int64) {
	off := ((t1.Int64() - t0.Int64()) + (t2.Int64() - t3.Int64())) / 2
	rtt := t3.Int64() - t0.Int64()
	return off, rtt
}

/* ==================== Offset Smoothing Strategty ======================= */

type ServerTimeOffsetSmoothingStrategy interface {
	Update(off, rtt int64) (int64, bool)
}

// newOffsetSmoothing creates a new OffsetSmoothingStrategy with the given
// windowSize and rttOutlierFactor.
func newOffsetSmoothingStrategy(windowSize int, rttOutlierFactor float64, logger *log.Entry) *OffsetSmoothingStrategy {
	return &OffsetSmoothingStrategy{
		windowSize:       windowSize,
		rttOutlierFactor: rttOutlierFactor,
		offs:             make([]int64, windowSize),
		rtts:             make([]int64, windowSize),
		logger:           logger.WithField("_caller", "OffsetSmoothingStrategy"),
	}
}

// OffsetSmoothingStrategy is a helper struct that implements an offset
// smoothing strategy that uses a round trip time filter to avoid outliers,
// and a rolling mean to calculate the mean offset for a given window size.
// TODO: what happens when the system clock is adjusted?
type OffsetSmoothingStrategy struct {
	windowSize       int
	rttOutlierFactor float64
	rtts             []int64
	offs             []int64
	count            int
	logger           *log.Entry
}

// add adds a new offset and round trip time to the OffsetSmoothingStrategy.
func (s *OffsetSmoothingStrategy) add(off, rtt int64) {

	// s.rtts is initialized with zeros, so we need to filter them out before
	// calculating the median if we don't have enough data yet.
	var medianRTT int64
	switch {
	case s.count >= s.windowSize:
		medianRTT, _ = median(s.rtts)
	case s.count >= 1 && s.count < s.windowSize:
		rtts := filterZeros(s.rtts)
		medianRTT, _ = median(rtts)
	}

	// if the round trip time is an outlier, we don't add the offset and round
	// trip time to the strategy.
	if medianRTT > 0 && float64(rtt) > float64(medianRTT)*s.rttOutlierFactor {
		s.logger.WithFields(log.Fields{
			"rtt":          rtt / 1e6,
			"median_rtt":   medianRTT / 1e6,
			"outlier_fact": s.rttOutlierFactor,
		}).Debug("rejected outlier round trip time")
		return
	}

	// add offset and round trip time to the strategy
	s.offs[s.count%s.windowSize] = off
	s.rtts[s.count%s.windowSize] = rtt
	s.count++
}

// meanOffset calculates the mean offset for the current window size.
func (s *OffsetSmoothingStrategy) meanOffset() (int64, bool) {
	// return 0 if we don't have enough data
	if s.count < s.windowSize {
		return 0, false
	}

	// calculate and return the mean offset
	meanOff, _ := mean(s.offs)
	return meanOff, true
}

// Update updates the OffsetSmoothingStrategy with a new offset and round trip
// time, and returns the mean offset for the current window size.
// If the window size is not reached yet, it returns 0 and false as the second
// return value. If the round trip time is an outlier, it returns the previous
// mean offset and true as the second return value.
func (s *OffsetSmoothingStrategy) Update(off, rtt int64) (int64, bool) {
	s.add(off, rtt)
	return s.meanOffset()
}

/* ==================== SyncServerTimeService ============================ */

// NewSyncServerTimeService creates a new SyncServerTimeService.
func NewSyncServerTimeService(client *binance.Client, logger *log.Entry) *SyncServerTimeService {
	errorPolicy := &defaultErrorPolicy{}
	warmupPolicy := newExponentialWarmupPolicy(25, 5000, 1.2, logger)
	offsetStrategy := newOffsetSmoothingStrategy(10, 1.2, logger)

	return &SyncServerTimeService{
		client:         client,
		offsetStrategy: offsetStrategy,
		warmupPolicy:   warmupPolicy,
		errorPolicy:    errorPolicy,
		IsSyncedChan:   make(chan struct{}),
		ErrChan:        make(chan error),
		logger:         logger.WithField("_caller", "SyncServerTimeService"),
	}
}

// SyncServerTimeService is responsible for synchronizing shrimpy-binance
// client (client.TimeHandler).
// It periodically makes api calls to the server to determine the difference
// between the local time and the server time, and updates the time offset in
// the client by calling Client.SetServerTimeOffset.
// This is intended to be run only once. If SyncServerTimeService.Run crashes,
// make a new SyncServerTimeService and run it again.
type SyncServerTimeService struct {
	client         *binance.Client
	offsetStrategy ServerTimeOffsetSmoothingStrategy
	warmupPolicy   WarmupPolicy
	errorPolicy    ErrorPolicy
	interval       time.Duration
	ticker         *time.Ticker
	IsSyncedChan   chan struct{}
	ErrChan        chan error
	logger         *log.Entry
}

// resetTicker updates the current interval and resets the ticker with the
// new value.
func (s *SyncServerTimeService) resetTicker(nextInterval time.Duration, level log.Level) {
	s.logger.WithFields(log.Fields{
		"nextInterval": nextInterval,
		"currInterval": s.interval,
	}).Log(level, "updating interval")

	s.interval = nextInterval
	s.ticker.Reset(s.interval)
}

// requestServerTime makes an api call to the server to get the server time
// and puts the response on the serverTimeRespChan, errors are put on the errChan.
func (s *SyncServerTimeService) requestServerTime(ctx context.Context, c chan<- *bsv.ServerTimeResponse, errC chan<- error) {
	service := s.client.NewSpotMarginServerTimeService()
	resp, err := service.Do(ctx)
	if err != nil {
		errC <- err
		return
	}
	c <- resp
}

// syncWithServerTime calculates and keeps track off the server time offset
// and round trip time. It starts setting the time offset in the client on every call
// once it has achieved sync. The first time it achieves sync, it will close the
// IsSyncedChan, so as to notify the caller that the client is now synchronized
// with the server time.
func (s *SyncServerTimeService) syncWithServerTime(resp *bsv.ServerTimeResponse) {

	// calculate offset and round trip time
	off, rtt := calcOffAndRTT(resp.TSLSent, resp.TSSServerTime, resp.TSSServerTime, resp.TSLRecv)

	// update offset smoothing strategy
	avgOff, ok := s.offsetStrategy.Update(off, rtt)
	if !ok {
		s.logger.Trace("not synced with server time yet")
		return
	}

	// update time offset in client
	s.client.SetServerTimeOffset(-avgOff)

	// close IsSyncedChan if it is not nil
	if s.IsSyncedChan != nil {
		close(s.IsSyncedChan)
		s.IsSyncedChan = nil
	}

	// log
	s.logger.WithFields(
		log.Fields{
			"offAvg": avgOff / 1e6,
			"off":    off / 1e6,
			"rtt":    rtt / 1e6,
			"t0":     resp.TSSSent,
			"t1":     resp.TSSServerTime,
			"t3":     resp.TSSRecv,
		}).Debug("synced with server time")

}

// Run starts the SyncServerTimeService. This can only be run once!
// If this crashes, make a new SyncServerTimeService and run it again.
func (s SyncServerTimeService) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(s.ErrChan)
		if s.IsSyncedChan != nil {
			close(s.IsSyncedChan)
		}
	}()

	// garbage collector will clean these up
	serverTimeRespChan := make(chan *bsv.ServerTimeResponse)
	errChan := make(chan error)

	//interval := s.warmupPolicy.CalcNextInterval()
	//ticker := time.NewTicker(interval)
	s.interval = s.warmupPolicy.CalcNextInterval()
	s.ticker = time.NewTicker(s.interval)
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case err := <-errChan:
			// consult error policy to determin if error is transient and how it
			// should be handled
			isTransientErr, nextInterval, reason := s.errorPolicy.Handle(err)
			s.logger.WithError(err).WithField("reason", reason).Error("error in sync server time service")

			// if error is not transient, return
			if !isTransientErr {
				return
			}

			// if the error policy says to update the interval, update the interval
			if nextInterval > 0 {
				s.resetTicker(nextInterval, log.InfoLevel)
			}

		case <-s.ticker.C:
			// do request
			go s.requestServerTime(ctx, serverTimeRespChan, errChan)

			// update interval
			nextInterval := s.warmupPolicy.CalcNextInterval()
			if nextInterval != s.interval {
				s.resetTicker(nextInterval, log.TraceLevel)
			}

		case resp := <-serverTimeRespChan:
			// attempt to update time offset in client & close IsSyncedChan
			s.syncWithServerTime(resp)

		}
	}
}
