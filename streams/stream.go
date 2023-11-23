package streams

import (
	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== AggTradesHandler (shrimpy-binance/streams/utils) = */

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
}

var FAPIStreams = map[string]common.StreamDefinition{
	"aggTrades": {
		Scheme:       "wss",
		Endpoint:     common.WSEndpointFAPI,
		EndpointType: common.EndpointTypeFAPI,
		SecurityType: common.WSSecurityTypeNone,
		UpdateSpeed:  0, // Real-time
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

/* ==================== APIStreams ======================================= */

func NewSpotMarginAggTradesStream(wc common.WSClient, logger *log.Entry) *AggTradesStream {
	sm := common.NewStreamMeta(APIStreams["aggTrades"])
	tradesHandler := NewAggTradesHandler()
	//logger = logger.WithField("_caller", "SpotMarginAggTradesStream")

	return &AggTradesStream{
		SM:      sm,
		Handler: tradesHandler,
		Stream:  wc.NewStream(sm, tradesHandler, logger.WithField("_caller", "SpotMarginAggTradesStream")),
	}
}

/* ==================== FAPIStreams ====================================== */

func NewFuturesAggTradesStream(wc common.WSClient, logger *log.Entry) *AggTradesStream {
	sm := common.NewStreamMeta(FAPIStreams["aggTrades"])
	tradesHandler := NewAggTradesHandler()

	return &AggTradesStream{
		SM:      sm,
		Handler: tradesHandler,
		Stream:  wc.NewStream(sm, tradesHandler, logger.WithField("_caller", "FuturesAggTradesStream")),
	}
}
