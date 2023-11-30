package services

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== PingService ====================================== */

// PingResponse
type PingResponse struct {
	ServiceBaseResponse
}

// PingService
type PingService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
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
		s.logger.WithError(err).Error("Do")
		return nil, err
	}
	resp, err := s.parseResponse(data)
	if err != nil {
		s.logger.WithError(err).Error("Do")
	}
	return resp, err
}

// parseResponse parses the request response into the PingResponse struct.
func (s *PingService) parseResponse(data []byte) (*PingResponse, error) {
	resp := &PingResponse{}
	if err := resp.ParseBaseResponse(&s.SM); err != nil {
		return nil, err
	}
	return resp, nil
}
