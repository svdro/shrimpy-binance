package streams

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== constants ========================================= */

const (
	// most  common.StreamHandlers will not implement HandleSend. Use this.
	handleSendWarning = "HandleSend method is not be implemented for this StreamHandler. This handler is designed for receiving and processing incoming data only. Any attempts to send data using HandleSend will not be processed and should be avoided."
)

/* ==================== StreamDefinitions ================================ */

var APIStreams = map[string]common.StreamDefinition{
	"aggTrades": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointAPI,
		EndpointType: common.EndpointTypeAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
	},
	"depth100ms": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointAPI,
		EndpointType: common.EndpointTypeAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  100, // 100ms
	},
	"depth1000ms": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointAPI,
		EndpointType: common.EndpointTypeAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  1000, // 1000ms
	},
	"userDataStream": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointAPI,
		EndpointType: common.EndpointTypeAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
	},
}

var SAPIStreams = map[string]common.StreamDefinition{
	"userDataStream": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointAPI,
		EndpointType: common.EndpointTypeSAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
	},
}

var FAPIStreams = map[string]common.StreamDefinition{
	"aggTrades": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointFAPI,
		EndpointType: common.EndpointTypeFAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
	},
	"depth100ms": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointFAPI,
		EndpointType: common.EndpointTypeFAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  100, // 100ms
	},
}

var WSAPIStreams = map[string]common.StreamDefinition{
	"wsAPIStream": {
		Scheme:       "wss",
		Endpoint:     common.WSAPIEndpointAPI,
		EndpointType: common.EndpointTypeAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
	},
}

/* ==================== APIStreams Factory =============================== */

func NewSpotUserDataStream(wc common.WSClient, logger *log.Entry) *SpotUserDataStream {
	sm := common.NewStreamMeta(APIStreams["userDataStream"])
	handler := newSpotUserDataStreamHandler(logger.WithField("_caller", "SpotUserDataHandler"))
	return &SpotUserDataStream{
		Handler: handler,
		Stream:  wc.NewStream(sm, handler, logger.WithField("_caller", "SpotUserDataStream")),
	}
}

func NewSpotMarginDiffDepth100Stream(wc common.WSClient, logger *log.Entry) *SpotMarginDiffDepthStream {
	sm := common.NewStreamMeta(APIStreams["depth100ms"])
	handler := newSpotMarginDiffDepthHandler(logger.WithField("_caller", "SpotMarginDiffDepthHandler"))

	return &SpotMarginDiffDepthStream{
		Handler:     handler,
		Stream:      wc.NewStream(sm, handler, logger.WithField("_caller", "SpotMarginDiffDepthStream")),
		UpdateSpeed: &sm.SD.UpdateSpeed, // hardcoded to 100ms
	}
}

func NewSpotMarginAggTradesStream(wc common.WSClient, logger *log.Entry) *SpotMarginAggTradesStream {
	sm := common.NewStreamMeta(APIStreams["aggTrades"])
	handler := newSpotMarginAggTradesHandler(logger.WithField("_caller", "SpotMarginAggTradesHandler"))

	return &SpotMarginAggTradesStream{
		Handler: handler,
		Stream:  wc.NewStream(sm, handler, logger.WithField("_caller", "SpotMarginAggTradesStream")),
	}
}

/* ==================== SAPIStreams Factory ============================== */

func NewMarginUserDataStream(wc common.WSClient, logger *log.Entry) *MarginUserDataStream {
	sm := common.NewStreamMeta(SAPIStreams["userDataStream"])
	handler := newMarginUserDataStreamHandler(logger.WithField("_caller", "MarginUserDataHandler"))
	return &MarginUserDataStream{
		Handler: handler,
		Stream:  wc.NewStream(sm, handler, logger.WithField("_caller", "MarginUserDataStream")),
	}
}

/* ==================== FAPIStreams Factory ============================== */

func NewFuturesDiffDepth100Stream(wc common.WSClient, logger *log.Entry) *FuturesDiffDepthStream {
	sm := common.NewStreamMeta(FAPIStreams["depth100ms"])
	handler := newFuturesDiffDepthHandler(logger.WithField("_caller", "FuturesDiffDepthHandler"))

	return &FuturesDiffDepthStream{
		Handler:     handler,
		Stream:      wc.NewStream(sm, handler, logger.WithField("_caller", "FuturesDiffDepthStream")),
		UpdateSpeed: &sm.SD.UpdateSpeed,
	}
}

func NewFuturesAggTradesStream(wc common.WSClient, logger *log.Entry) *FuturesAggTradesStream {
	sm := common.NewStreamMeta(FAPIStreams["aggTrades"])
	handler := newFuturesAggTradesHandler(logger.WithField("_caller", "FuturesAggTradesHandler"))

	return &FuturesAggTradesStream{
		Handler: handler,
		Stream:  wc.NewStream(sm, handler, logger.WithField("_caller", "FuturesAggTradesStream")),
	}
}
