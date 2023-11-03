package services

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== ServerTimeService ================================ */

// ServerTimeResponse
type ServerTimeResponse struct {
	ServiceBaseResponse
	TSSServerTime int64 `json:"serverTime"`
}

// ServerTimeService
type ServerTimeService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
}

// toParams converts all parameter fields of the service to a params struct.
func (s *ServerTimeService) toParams() *params {
	return &params{}
}

// parseResponse parses the request response into the ServerTimeResponse struct.
func (s *ServerTimeService) parseResponse(data []byte) (*ServerTimeResponse, error) {
	resp := &ServerTimeResponse{}
	if err := resp.ParseBaseResponse(&s.SM, s.rc.TimeHandler()); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	resp.TSSServerTime = resp.TSSServerTime * 1e6 // convert to nanoseconds
	return resp, nil
}

// Do does the ServerTime request.
func (s *ServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())
	if err != nil {
		s.logger.WithError(err).Error("Do")
		return nil, err
	}

	resp, err := s.parseResponse(data)
	if err != nil {
		s.logger.WithError(err).Debug("Do")
	}
	return resp, err
}
