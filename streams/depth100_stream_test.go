package streams

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	apiDepthMsg = []byte(`{
  "e": "depthUpdate",
  "E": 123456789,
  "s": "BNBBTC",
  "U": 157,
  "u": 160,
  "b": [ [ "0.0024", "10" ], [ "0.0023", "50" ] ],
  "a": [ [ "0.0026", "100" ] ]
}`)

	apiDepthTarget = &SpotMarginDiffDepthEvent{
		StreamBaseEvent: StreamBaseEvent{EventType: "depthUpdate", TSSEvent: 123456789},
		sharedDiffDepthEvent: sharedDiffDepthEvent{
			Symbol:  "BNBBTC",
			FirstID: 157,
			FinalID: 160,
			Bids:    []Level{{Price: "0.0024", Qty: "10"}, {Price: "0.0023", Qty: "50"}},
			Asks:    []Level{{Price: "0.0026", Qty: "100"}},
		},
	}

	fapiDepthMsg = []byte(`{
  "e": "depthUpdate",
  "E": 123456789,
  "T": 123456788,
  "s": "BTCUSDT",
  "U": 157,
  "u": 160,
  "pu": 149,
  "b": [ [ "0.0024", "10" ] ], "a": [ [ "0.0026", "100" ] ]
}`)

	fapiDepthTarget = &FuturesDiffDepthEvent{
		StreamBaseEvent: StreamBaseEvent{EventType: "depthUpdate", TSSEvent: 123456789},
		sharedDiffDepthEvent: sharedDiffDepthEvent{
			Symbol:  "BTCUSDT",
			FirstID: 157,
			FinalID: 160,
			Bids:    []Level{{Price: "0.0024", Qty: "10"}},
			Asks:    []Level{{Price: "0.0026", Qty: "100"}},
		},
		TransTime:   123456788,
		LastFinalID: 149,
	}
)

func TestUnmarshalDepthUpdateEvent(t *testing.T) {
	apiEvent := &SpotMarginDiffDepthEvent{}
	err := json.Unmarshal(apiDepthMsg, apiEvent)
	assert.Nil(t, err)
	assert.Equal(t, apiDepthTarget, apiEvent)

	fapiEvent := &FuturesDiffDepthEvent{}
	err = json.Unmarshal(fapiDepthMsg, fapiEvent)
	assert.Nil(t, err)
	assert.Equal(t, fapiDepthTarget, fapiEvent)
}

func TestDepthHandler(t *testing.T) {
	handler := newFuturesDiffDepthHandler(nil)
	assert.NotNil(t, handler)

	handler.HandleRecv(fapiDepthMsg, 0, 0)
	event := <-handler.EventChan
	assert.NotNil(t, event)
	assert.IsType(t, &FuturesDiffDepthEvent{}, event)
	assert.Equal(t, fapiDepthTarget, event)
}
