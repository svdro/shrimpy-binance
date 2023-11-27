package streams

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== UserDataStream =================================== */

type SpotMarginUserDataStream[A, B, O Event] struct {
	common.Stream
	Handler   *SpotMarginUserDataStreamHandler[A, B, O]
	ListenKey *string
}

func (s *SpotMarginUserDataStream[A, B, O]) SetListenKey(listenKey string) *SpotMarginUserDataStream[A, B, O] {
	s.ListenKey = &listenKey
	s.SetPathFunc(s.path)
	return s
}

func (s *SpotMarginUserDataStream[A, B, O]) path() string {
	path := "/ws/%s"
	return fmt.Sprintf(path, *s.ListenKey)
}

/* ==================== UserDataEvents =================================== */

type Balance struct {
	Asset  string `json:"a"`
	Free   string `json:"f"`
	Locked string `json:"l"`
}

type AccountUpdateEvent struct {
	StreamBaseEvent
	TSSLUpdate common.TSNano `json:"u"` // time of last account update
	Balances   []Balance     `json:"B"`
}

type BalanceUpdateEvent struct {
	StreamBaseEvent
	Asset    string        `json:"a"`
	Delta    string        `json:"d"`
	TSSClear common.TSNano `json:"T"`
}

type OrderUpdateEvent struct {
	StreamBaseEvent
	Symbol      string                    `json:"s"`
	Side        common.BIOrderSide        `json:"S"`
	OrderType   common.BIOrderType        `json:"o"`
	TimeInForce common.BIOrderTimeInForce `json:"f"`

	ExecutionStatus common.BIExecutionType `json:"x"` // current execution type
	OrderStatus     common.BIOrderStatus   `json:"X"`

	ClientOrderID     string `json:"c"`
	OrderListID       int64  `json:"g"`
	OrigClientOrderID string `json:"C"` // this is the ID of the order being canceled
	OrderID           int64  `json:"i"`
	TradeID           int64  `json:"t"`

	Price             string `json:"p"`
	StopPrice         string `json:"P"`
	LastExecutedPrice string `json:"L"`

	Qty               string `json:"q"`
	IcebergQty        string `json:"F"`
	LastExecutedQty   string `json:"l"`
	CumFilledQty      string `json:"z"`
	CumQuoteAssetQty  string `json:"Z"` // cumulative quote asset transacted quantity
	LastQuoteAssetQty string `json:"Y"` // last quote asset transacted quantity (i.e. lastPrice * lastQty)
	QuoteOrderQty     string `json:"Q"` // quote order quantity

	CommissionAmount string `json:"n"`
	CommissionAsset  string `json:"N"`

	TSSTransact     common.TSNano `json:"T"` // transaction time
	TSSOrderCreated common.TSNano `json:"O"` // order creation time
	TSSWorking      common.TSNano `json:"W"` // working time; this is only visible if the order has been placed on the book

	IsOnBook bool `json:"w"`
	IsMaker  bool `json:"m"`

	OrderRejectReason       string                           `json:"r"` // this will be an error code
	SelfTradePreventionMode common.BISelfTradePreventionMode `json:"V"`

	Ignore1 interface{} `json:"I"`
	Ignore2 interface{} `json:"M"`
}

/* ==================== Spot ============================================= */

// SpotAccountUpdateEvent is an account update event for spot markets.
type SpotAccountUpdateEvent = AccountUpdateEvent

// SpotBalanceUpdateEvent is a balance update event for spot markets.
type SpotBalanceUpdateEvent = BalanceUpdateEvent

// SpotOrderUpdateEvent is an order update event for spot markets.
type SpotOrderUpdateEvent = OrderUpdateEvent

// SpotUserDataStreamHandler is a handler for spot user data streams.
type SpotUserDataStreamHandler = SpotMarginUserDataStreamHandler[
	*SpotAccountUpdateEvent, *SpotBalanceUpdateEvent, *SpotOrderUpdateEvent]

// newSpotUserDataStreamHandler creates a new SpotUserDataStreamHandler.
func newSpotUserDataStreamHandler(logger *log.Entry) *SpotUserDataStreamHandler {
	return newSpotMarginUserDataStreamHandler[
		*SpotAccountUpdateEvent, *SpotBalanceUpdateEvent, *SpotOrderUpdateEvent,
	](logger)
}

// SpotUserDataStream is a user data stream for spot markets.
type SpotUserDataStream = SpotMarginUserDataStream[
	*SpotAccountUpdateEvent, *SpotBalanceUpdateEvent, *SpotOrderUpdateEvent]

/* ==================== Margin =========================================== */

// MarginAccountUpdateEvent is an account update event for margin markets.
type MarginAccountUpdateEvent = AccountUpdateEvent

// MarginBalanceUpdateEvent is a balance update event for margin markets.
type MarginBalanceUpdateEvent = BalanceUpdateEvent

// MarginOrderUpdateEvent is an order update event for margin markets.
type MarginOrderUpdateEvent = OrderUpdateEvent

// MarginUserDataStreamHandler is a handler for margin user data streams.
type MarginUserDataStreamHandler = SpotMarginUserDataStreamHandler[
	*MarginAccountUpdateEvent, *MarginBalanceUpdateEvent, *MarginOrderUpdateEvent]

// newMarginUserDataStreamHandler creates a new MarginUserDataStreamHandler.
func newMarginUserDataStreamHandler(logger *log.Entry) *MarginUserDataStreamHandler {
	return newSpotMarginUserDataStreamHandler[
		*MarginAccountUpdateEvent, *MarginBalanceUpdateEvent, *MarginOrderUpdateEvent,
	](logger)
}

// MarginUserDataStream is a user data stream for margin markets.
type MarginUserDataStream = SpotMarginUserDataStream[
	*MarginAccountUpdateEvent, *MarginBalanceUpdateEvent, *MarginOrderUpdateEvent]
