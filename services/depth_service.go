package services

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== Depth Service ==================================== */

// DepthService is a generic depth service that can be used for spot/margin
// and futures.
type DepthService[R Response] struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
	symbol string // mandatory
	depth  int64  // mandatory (this needs to be initialized when the service is created)
}

// WithSymbol returns a copy of the service with symbol set to the provided
func (s DepthService[R]) WithSymbol(symbol string) *DepthService[R] {
	s.symbol = symbol
	return &s
}

// toParams converts all parameter fields of the service to a params struct.
func (s *DepthService[R]) toParams() *params {
	p := &params{}
	p.Set("symbol", s.symbol)
	p.Set("limit", strconv.FormatInt(s.depth, 10))
	return p
}

// parseResponse parses the response data into a either a SpotMarginDepthResponse
// or a FuturesDepthResponse.
// NOTE: resp = new(R)
// this does not correctly initialize the embedded struct.
// it creates a pointer to a pointer of a nil value (or sth like that).
// it works because json fixes the initialization, but only when
// calling json.Unmarshal before calling methods on the embedded struct.
// maybe fix this cause it's ugly, but also maybe not cause it works.
func (s *DepthService[R]) parseResponse(data []byte) (R, error) {
	resp := new(R)

	//log.WithField("resp", resp).Info("parseResponse")

	if err := json.Unmarshal(data, &resp); err != nil {
		return *resp, err
	}

	if err := (*resp).ParseBaseResponse(&s.SM); err != nil {
		return *resp, err
	}

	return *resp, nil
}

// Do does the depth request.
func (s DepthService[R]) Do(ctx context.Context) (R, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())
	if err != nil {
		s.logger.WithError(err).Error("Do: rc.Do")
		return *new(R), err
	}

	resp, err := s.parseResponse(data)
	if err != nil {
		s.logger.WithError(err).Error("Do: parseResponse")
	}
	return resp, err
}

/* ==================== DepthResponse ==================================== */

// Level represents a price and quantity pair.
// NOTE: this struct is defined both in streams and services. Consider
// refactoring if more shared structs pop up.
type Level struct {
	Price string `json:"p"`
	Qty   string `json:"q"`
}

// UnmarshalJSON unmarshals a price qty pair from a JSON array to a Level.
func (l *Level) UnmarshalJSON(data []byte) error {
	var tmp [2]string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	l.Price = tmp[0]
	l.Qty = tmp[1]
	return nil
}

// SharedDepthResponse is a has the fields that are shared between
// SpotMarginDepthResponse and FuturesDepthResponse.
type SharedDepthResponse struct {
	LastUpdateID int64   `json:"lastUpdateId"`
	Bids         []Level `json:"bids"`
	Asks         []Level `json:"asks"`
}

/* ==================== SpotMargin ======================================= */

// SpotMarginDepthResponse is the response from the spot/margin depth service.
type SpotMarginDepthResponse struct {
	ServiceBaseResponse
	SharedDepthResponse
}

// SpotMarginDepthService is a depth service for spot/margin markets.
type SpotMarginDepthService = DepthService[*SpotMarginDepthResponse]

/* ==================== Futures ========================================== */

// FuturesDepthResponse is the response from the futures depth service.
type FuturesDepthResponse struct {
	ServiceBaseResponse
	SharedDepthResponse
	TSSEvent    common.TSNano `json:"E"` // Message output time
	TSSTransact common.TSNano `json:"T"` // Transaction time
}

// FuturesDepthService is a depth service for futures markets.
type FuturesDepthService = DepthService[*FuturesDepthResponse]
