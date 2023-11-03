package services

import (
	"context"

	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== PingService ====================================== */

// PingService
type PingService struct {
	SM common.ServiceMeta
	rc common.RESTClient
}

// toParams converts all parameter fields of the service to a params struct.
func (s *PingService) toParams() *params {
	return &params{}
}

// Do does the Ping request.
func (s *PingService) Do(ctx context.Context) (*PingResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())

	if err != nil {
		return nil, err
	}
	return s.parseResponse(data)
}

// parseResponse parses the request response into the PingResponse struct.
func (s *PingService) parseResponse(data []byte) (*PingResponse, error) {
	resp := &PingResponse{}
	if err := resp.ParseBaseResponse(&s.SM, s.rc.TimeHandler()); err != nil {
		return nil, err
	}
	return resp, nil
}

/* ==================== PingResponse ===================================== */

type PingResponse struct {
	ServiceBaseResponse
}
