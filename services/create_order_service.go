package services

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== CreateMarginOrderResponse =========================*/

type CreateMarginOrderFillResponse struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

type CreateMarginOrderAckResponse struct {
	ServiceBaseResponse
	TSSTransact   common.TSNano `json:"transactTime"`
	Symbol        string        `json:"symbol"`
	OrderID       int64         `json:"orderId"`
	ClientOrderID string        `json:"clientOrderId"`
	IsIsolated    bool          `json:"isIsolated"`
}

type CreateMarginOrderResultResponse struct {
	CreateMarginOrderAckResponse
	Price                   string                           `json:"price"`
	OrigQty                 string                           `json:"origQty"`
	ExecutedQty             string                           `json:"executedQty"`
	CumQuoteQty             string                           `json:"cummulativeQuoteQty"`
	Status                  common.BIOrderStatus             `json:"status"`
	TimeInForce             common.BIOrderTimeInForce        `json:"timeInForce"`
	OrderType               common.BIOrderType               `json:"type"`
	Side                    common.BIOrderSide               `json:"side"`
	SelfTradePreventionMode common.BISelfTradePreventionMode `json:"selfTradePreventionMode"`
}

type CreateMarginOrderFullResponse struct {
	CreateMarginOrderResultResponse
	MarginBuyBorrowAmount string                          `json:"marginBuyBorrowAmount"` // (MARGIN) will not return if no margin trade happens
	MarginBuyBorrowAsset  string                          `json:"marginBuyBorrowAsset"`  // (MARGIN) will not return if no margin trade happens
	Fills                 []CreateMarginOrderFillResponse `json:"fills"`
}

/* ==================== CreateMarginOrderService ==========================*/

// CreateMarginOrderService
type CreateMarginOrderService struct {
	SM                      common.ServiceMeta
	rc                      common.RESTClient
	logger                  *log.Entry
	symbolREST              string  // (ALL ORDERS)
	side                    string  // (ALL ORDERS)
	orderType               string  // (ALL ORDERS)
	price                   *string // (LIMIT, STOP_LOSS_LIMIT, TAKE_PROFIT_LIMIT, LIMIT_MAKER)
	stopPrice               *string // (STOP_LOSS, TAKE_PROFIT, STOP_LOSS_LIMIT, TAKE_PROFIT_LIMIT)
	quantity                *string // (ALL ORDERS)
	quoteOrderQty           *string // (MARKET)
	icebergQty              *string // (LIMIT, STOP_LOSS_LIMIT, TAKE_PROFIT_LIMIT)
	sideEffectType          *string // (???) (NO_SIDE_EFFECT, MARGIN_BUY, AUTO_REPAY) (default: NO_SIDE_EFFECT)
	autoRepayAtCancel       *string // (???) (true, false) (default: true)
	selfTradePreventionMode *string // (???) (EXPIRE_TAKER, EXPIRE_MAKER, EXPIRE_BOTH, NONE)
	timeInForce             *string // (LIMIT, STOP_LOSS_LIMIT, TAKE_PROFIT_LIMIT)
	isIsolated              *string // (ALL ORDERS)
	newClientOrderID        *string // (ALL ORDERS)
}

// Do sends the request and returns a CreateMarginOrderFullEvent.
func (s *CreateMarginOrderService) Do(ctx context.Context) (*CreateMarginOrderFullResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())
	if err != nil {
		s.logger.WithError(err).Error("Do")
		return nil, err
	}

	resp, err := s.parseResponse(data)
	if err != nil {
		s.logger.WithError(err).Error("Do")
		return nil, err
	}
	return resp, nil
}

// toParams converts all parameter fields of the service to a params struct.
func (s *CreateMarginOrderService) toParams() params {
	p := params{}
	p.Set("symbol", s.symbolREST)
	p.Set("side", s.side)
	p.Set("type", s.orderType)

	// always set newOrderRespType to FULL
	p.Set("newOrderRespType", string(common.OrderResponseTypeFull))

	p.SetIfNotNil("price", s.price)
	p.SetIfNotNil("stopPrice", s.stopPrice)

	p.SetIfNotNil("quantity", s.quantity)
	p.SetIfNotNil("quoteOrderQty", s.quoteOrderQty)
	p.SetIfNotNil("icebergQty", s.icebergQty)

	p.SetIfNotNil("sideEffectType", s.sideEffectType)
	p.SetIfNotNil("autoRepayAtCancel", s.autoRepayAtCancel)
	p.SetIfNotNil("selfTradePreventionMode", s.selfTradePreventionMode)

	p.SetIfNotNil("timeInForce", s.timeInForce)
	p.SetIfNotNil("isIsolated", s.isIsolated)
	p.SetIfNotNil("newClientOrderId", s.newClientOrderID)

	return p
}

func (s *CreateMarginOrderService) parseResponse(data []byte) (*CreateMarginOrderFullResponse, error) {
	resp := &CreateMarginOrderFullResponse{}

	if err := resp.ParseBaseResponse(&s.SM); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil

}

// WithMarketOrderParams returns a copy of the service with all parameters that
// are mandatory for a market order set.
// Any optional params must be set before or after calling this method.
func (s CreateMarginOrderService) WithMarketOrderParams(
	symbolREST string, side common.BIOrderSide, quantity string,
) *CreateMarginOrderService {
	return s.
		WithBaseOrderParams(symbolREST, side, common.OrderTypeMarket).
		WithQuantity(quantity)
}

// WithBaseOrderParams returns a copy of the service with all parameters that
// are mandatory and shared between all order types set.
func (s CreateMarginOrderService) WithBaseOrderParams(
	symbolREST string, side common.BIOrderSide, orderType common.BIOrderType,
) *CreateMarginOrderService {
	return s.WithSymbolREST(symbolREST).WithSide(side).WithOrderType(orderType)
}

// WithLimitOrderParams returns a copy of the service with all parameters that
// are mandatory for a limit order set.
// Any optional params must be set before or after calling this method.
func (s CreateMarginOrderService) WithLimitOrderParams(
	symbolREST string,
	side common.BIOrderSide,
	quantity,
	price string,
	timeInForce common.BIOrderTimeInForce,
) *CreateMarginOrderService {
	return s.
		WithBaseOrderParams(symbolREST, side, common.OrderTypeLimit).
		WithQuantity(quantity).
		WithPrice(price).
		WithTimeInForce(timeInForce)
}

// WithSymbolREST returns a copy of the service with symbolREST set to the given value.
func (s CreateMarginOrderService) WithSymbolREST(symbolREST string) *CreateMarginOrderService {
	s.symbolREST = symbolREST
	return &s
}

// WithSide returns a copy of the service with side set to the given value.
func (s CreateMarginOrderService) WithSide(side common.BIOrderSide) *CreateMarginOrderService {
	s.side = string(side)
	return &s
}

// WithOrderType returns a copy of the service with orderType set to the given value.
func (s CreateMarginOrderService) WithOrderType(orderType common.BIOrderType) *CreateMarginOrderService {
	s.orderType = string(orderType)
	return &s
}

// WithPrice returns a copy of the service with price set to the given value.
func (s CreateMarginOrderService) WithPrice(price string) *CreateMarginOrderService {
	s.price = &price
	return &s
}

// WithStopPrice returns a copy of the service with stopPrice set to the given value.
func (s CreateMarginOrderService) WithStopPrice(stopPrice string) *CreateMarginOrderService {
	s.stopPrice = &stopPrice
	return &s
}

// WithQuantity returns a copy of the service with quantity set to the given value.
func (s CreateMarginOrderService) WithQuantity(quantity string) *CreateMarginOrderService {
	s.quantity = &quantity
	return &s
}

// WithQuoteOrderQty returns a copy of the service with quoteOrderQty set to the given value.
func (s CreateMarginOrderService) WithQuoteOrderQty(quoteOrderQty string) *CreateMarginOrderService {
	s.quoteOrderQty = &quoteOrderQty
	return &s
}

// WithIcebergQty returns a copy of the service with icebergQty set to the given value.
func (s CreateMarginOrderService) WithIcebergQty(icebergQty string) *CreateMarginOrderService {
	s.icebergQty = &icebergQty
	return &s
}

// WithSideEffectType returns a copy of the service with sideEffectType set to the given value.
func (s CreateMarginOrderService) WithSideEffectType(sideEffectType common.BIOrderSideEffect) *CreateMarginOrderService {
	sideEffectTypeStr := string(sideEffectType)
	s.sideEffectType = &sideEffectTypeStr
	return &s
}

// WithAutoRepayAtCancel returns a copy of the service with autoRepayAtCancel set to the given value.
// TODO: this is almost certainly a bool rather than an upper bool.
func (s CreateMarginOrderService) WithAutoRepayAtCancel(autoRepayAtCancel bool) *CreateMarginOrderService {
	autoRepayAtCancelStr := boolToUpperStr(autoRepayAtCancel)
	s.autoRepayAtCancel = &autoRepayAtCancelStr
	return &s
}

// WithSelfTradePreventionMode returns a copy of the service with selfTradePreventionMode set to the given value.
func (s CreateMarginOrderService) WithSelfTradePreventionMode(selfTradePreventionMode common.BISelfTradePreventionMode) *CreateMarginOrderService {
	selfTradePreventionModeStr := string(selfTradePreventionMode)
	s.selfTradePreventionMode = &selfTradePreventionModeStr
	return &s
}

// WithTimeInForce returns a copy of the service with timeInForce set to the given value.
func (s CreateMarginOrderService) WithTimeInForce(timeInForce common.BIOrderTimeInForce) *CreateMarginOrderService {
	timeInForceStr := string(timeInForce)
	s.timeInForce = &timeInForceStr
	return &s
}

// WithIsIsolated returns a copy of the service with isIsolated set to the given value.
func (s CreateMarginOrderService) WithIsIsolated(isIsolated bool) *CreateMarginOrderService {
	isIsolatedStr := boolToUpperStr(isIsolated)
	s.isIsolated = &isIsolatedStr
	return &s
}

// WithNewClientOrderId returns a copy of the service with newClientOrderId
// set to the given value.
func (s CreateMarginOrderService) WithNewClientOrderId(newClientOrderId string) *CreateMarginOrderService {
	s.newClientOrderID = &newClientOrderId
	return &s
}

//[> ==================== CancelMarginOrderService ==========================<]

//// CancelMarginOrderService
//// (symbolREST and either orderId or origClientOrderId must be sent)
//type CancelMarginOrderService struct {
//api *API
//ServiceDefinition
//symbolREST        string
//isIsolated        *string
//orderID           *string
//origClientOrderID *string
//newClientOrderID  *string
//}

//// Do sends the request and returns a CancelMarginOrderEvent.
//func (s *CancelMarginOrderService) Do(ctx context.Context) (*CancelMarginOrderEvent, error) {
//p := s.toParams()
//resp, reqMetaData, err := s.api.Do(ctx, &s.ServiceDefinition, p)
//if err != nil {
//return nil, err
//}

//e := &CancelMarginOrderEvent{}
//e.RequestMetaData = reqMetaData
//err = json.Unmarshal(resp, e)
//return e, err
//}

//// toParams converts all parameter fields of the service to a params struct.
//func (s *CancelMarginOrderService) toParams() params {
//p := params{}
//p.set("symbol", s.symbolREST)
//p.setIfNotNil("isIsolated", s.isIsolated)
//p.setIfNotNil("orderId", s.orderID)
//p.setIfNotNil("origClientOrderId", s.origClientOrderID)
//p.setIfNotNil("newClientOrderId", s.newClientOrderID)
//return p
//}

//func (s *CancelMarginOrderService) GetParams() params {
//return s.toParams()
//}

//// WithSymbolREST returns a copy of the service with symbolREST set to the given value.
//func (s CancelMarginOrderService) WithSymbol(symbolREST string) *CancelMarginOrderService {
//s.symbolREST = symbolREST
//return &s
//}

//// WithIsIsolated returns a copy of the service with isIsolated set to the given value.
//func (s CancelMarginOrderService) WithIsIsolated(isIsolated bool) *CancelMarginOrderService {
//isIsolatedStr := boolToUpperStr(isIsolated)
//s.isIsolated = &isIsolatedStr
//return &s
//}

//// WithOrderID returns a copy of the service with orderID set to the given value.
//func (s CancelMarginOrderService) WithOrderID(orderID int64) *CancelMarginOrderService {
//orderIDStr := strconv.FormatInt(orderID, 10)
//s.orderID = &orderIDStr
//return &s
//}

//// WithOrigClientOrderID returns a copy of the service with origClientOrderID
//// set to the given value.
//func (s CancelMarginOrderService) WithOrigClientOrderID(origClientOrderID string) *CancelMarginOrderService {
//s.origClientOrderID = &origClientOrderID
//return &s
//}

//// WithNewClientOrderID returns a copy of the service with newClientOrderID
//// set to the given value.
//func (s CancelMarginOrderService) WithNewClientOrderID(newClientOrderID string) *CancelMarginOrderService {
//s.newClientOrderID = &newClientOrderID
//return &s
//}
