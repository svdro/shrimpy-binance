package streams

import (
	"encoding/json"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/svdro/shrimpy-binance/common"
)

var (
	unexpectedMsg       = []byte(`{"e":"unexpected", "E": 123}`)
	missingEventTypeMsg = []byte(`{"E": 123}`)

	accountUpdateMsg = []byte(` {
  "e": "outboundAccountPosition", 
  "E": 1564034571105,
  "u": 1564034571073,
  "B": [ { "a": "ETH", "f": "10000.000000", "l": "0.000000" } ]
}`)

	accountUpdateTarget = &AccountUpdateEvent{
		StreamBaseEvent: StreamBaseEvent{EventType: "outboundAccountPosition", TSSEvent: common.NewTSNano(1564034571105)},
		TSSLUpdate:      common.NewTSNano(1564034571073),
		Balances:        []Balance{{Asset: "ETH", Free: "10000.000000", Locked: "0.000000"}},
	}

	balanceUpdateMsg = []byte(`{
  "e": "balanceUpdate",
  "E": 1573200697110,
  "a": "BTC",
  "d": "100.00000000",
  "T": 1573200697068
}`)

	balanceUpdateTarget = &BalanceUpdateEvent{
		StreamBaseEvent: StreamBaseEvent{EventType: "balanceUpdate", TSSEvent: common.NewTSNano(1573200697110)},
		Asset:           "BTC",
		Delta:           "100.00000000",
		TSSClear:        common.NewTSNano(1573200697068),
	}

	orderUpdateMsg = []byte(`{
  "e": "executionReport",
  "E": 1499405658658,
  "s": "ETHBTC",
  "c": "mUvoqJxFIILMdfAW5iGSOW",
  "S": "BUY",
  "o": "LIMIT",
  "f": "GTC",
  "q": "1.00000000",
  "p": "0.10264410",
  "P": "0.00000000",
  "F": "0.00000000",
  "g": -1,
  "C": "",
  "x": "NEW",
  "X": "NEW",
  "r": "NONE",
  "i": 4293153,
  "l": "0.00000000",
  "z": "0.00000000",
  "L": "0.00000000",
  "n": "0",
  "N": null,
  "T": 1499405658657,
  "t": -1,
  "I": 8641984,
  "w": true,
  "m": false,
  "M": false,
  "O": 1499405658657,
  "Z": "0.00000000",
  "Y": "0.00000000",
  "Q": "0.00000000",
  "W": 1499405658657,
  "V": "NONE"
}`)

	orderUpdateEventTarget = &OrderUpdateEvent{
		StreamBaseEvent:         StreamBaseEvent{EventType: "executionReport", TSSEvent: common.NewTSNano(1499405658658)},
		Symbol:                  "ETHBTC",
		ClientOrderID:           "mUvoqJxFIILMdfAW5iGSOW",
		Side:                    common.BIOrderSide("BUY"),
		OrderType:               common.OrderTypeLimit,
		TimeInForce:             common.OrderTimeInForceGTC,
		Qty:                     "1.00000000",
		Price:                   "0.10264410",
		StopPrice:               "0.00000000",
		IcebergQty:              "0.00000000",
		OrderListID:             -1,
		OrigClientOrderID:       "",
		ExecutionStatus:         common.ExecutionTypeNew,
		OrderStatus:             common.OrderStatusNew,
		OrderRejectReason:       "NONE",
		OrderID:                 4293153,
		LastExecutedQty:         "0.00000000",
		CumFilledQty:            "0.00000000",
		LastExecutedPrice:       "0.00000000",
		CommissionAmount:        "0",
		CommissionAsset:         "",
		TSSTransact:             common.NewTSNano(1499405658657),
		TradeID:                 -1,
		IsOnBook:                true,
		IsMaker:                 false,
		TSSOrderCreated:         common.NewTSNano(1499405658657),
		CumQuoteAssetQty:        "0.00000000",
		LastQuoteAssetQty:       "0.00000000",
		QuoteOrderQty:           "0.00000000",
		TSSWorking:              common.NewTSNano(1499405658657),
		SelfTradePreventionMode: common.SelfTradePreventionModeNone,
		Ignore1:                 float64(8641984), // does not really matter
		Ignore2:                 false,
	}
)

func TestUnmarshalAccountUpdateEvent(t *testing.T) {
	event := &AccountUpdateEvent{}
	err := json.Unmarshal(accountUpdateMsg, event)
	assert.Nil(t, err)
	assert.Equal(t, accountUpdateTarget, event)
}

func TestUnmarshalBalanceUpdateEvent(t *testing.T) {
	event := &BalanceUpdateEvent{}
	err := json.Unmarshal(balanceUpdateMsg, event)
	assert.Nil(t, err)
	assert.Equal(t, balanceUpdateTarget, event)
}

func TestUnmarshalOrderUpdateEvent(t *testing.T) {
	event := &OrderUpdateEvent{}
	err := json.Unmarshal(orderUpdateMsg, event)
	assert.Nil(t, err)
	assert.Equal(t, orderUpdateEventTarget, event)
}

func TestUnmarshalBaseEvent(t *testing.T) {
	event := &StreamBaseEvent{}
	err := json.Unmarshal(accountUpdateMsg, event)
	assert.Nil(t, err)
	target := &StreamBaseEvent{
		EventType: "outboundAccountPosition",
		TSSEvent:  common.NewTSNano(1564034571105),
	}
	assert.Equal(t, target, event)
}

func TestSpotUserDataStreamHadler(t *testing.T) {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	handler := newSpotUserDataStreamHandler(log.NewEntry(logger))
	assert.NotNil(t, handler)

	// accountUpdateEvent
	err := handler.HandleRecv(accountUpdateMsg, 0, 0)
	assert.Nil(t, err)

	accountEvent := <-handler.AccountUpdateEventChan
	assert.IsType(t, &SpotAccountUpdateEvent{}, accountEvent)
	assert.Equal(t, accountUpdateTarget, accountEvent)

	// balanceUpdateEvent
	err = handler.HandleRecv(balanceUpdateMsg, 0, 0)
	assert.Nil(t, err)

	balanceEvent := <-handler.BalanceUpdateEventChan
	assert.IsType(t, &SpotBalanceUpdateEvent{}, balanceEvent)
	assert.Equal(t, balanceUpdateTarget, balanceEvent)

	// orderUpdateEvent
	err = handler.HandleRecv(orderUpdateMsg, 0, 0)
	assert.Nil(t, err)

	orderEvent := <-handler.OrderUpdateEventChan
	assert.IsType(t, &SpotOrderUpdateEvent{}, orderEvent)
	assert.Equal(t, orderUpdateEventTarget, orderEvent)

	// unknown event type error
	err = handler.HandleRecv(unexpectedMsg, 0, 0)
	assert.NotNil(t, err)
	assert.IsType(t, &common.WSHandlerError{}, err)
	assert.Equal(t, "unknown event type: unexpected", err.Err.Error())

	// missing event type error
	err = handler.HandleRecv(missingEventTypeMsg, 0, 0)
	assert.NotNil(t, err)
	assert.IsType(t, &common.WSHandlerError{}, err)
	assert.Equal(t, "event type is empty", err.Err.Error())
}
