package services

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== CreateListenKeyService =========================== */

// CreateListenKeyResponse
type CreateListenKeyResponse struct {
	ServiceBaseResponse
	ListenKey string `json:"listenKey"`
}

// CreateListenKeyService
type CreateListenKeyService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
}

// toParams converts all parameter fields of the service to a params struct.
func (s *CreateListenKeyService) toParams() *params {
	return &params{}
}

// parseResponse parses the request response into the CreateListenKeyResponse struct.
func (s *CreateListenKeyService) parseResponse(data []byte) (*CreateListenKeyResponse, error) {
	resp := &CreateListenKeyResponse{}
	if err := resp.ParseBaseResponse(&s.SM, s.rc.TimeHandler()); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil

}

// Do does the CreateListenKey request.
func (s *CreateListenKeyService) Do(ctx context.Context) (*CreateListenKeyResponse, error) {
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

/* ==================== PingListenKeyService ============================= */

// PingListenKeyResponse
type PingListenKeyResponse struct {
	ServiceBaseResponse
}

// PingListenKeyService
type PingListenKeyService struct {
	SM        common.ServiceMeta
	rc        common.RESTClient
	logger    *log.Entry
	listenKey string
}

// WithListenKey returns a copy of the service with listenKey set to the
// provided value.
func (s PingListenKeyService) WithListenKey(listenKey string) *PingListenKeyService {
	s.listenKey = listenKey
	return &s
}

// toParams converts all parameter fields of the service to a params struct.
func (s *PingListenKeyService) toParams() *params {
	p := &params{}
	p.Set("listenKey", s.listenKey)
	return p
}

// parseResponse parses the request response into the PingListenKeyResponse struct.
func (s *PingListenKeyService) parseResponse(data []byte) (*PingListenKeyResponse, error) {
	resp := &PingListenKeyResponse{}
	if err := resp.ParseBaseResponse(&s.SM, s.rc.TimeHandler()); err != nil {
		return nil, err
	}
	return resp, nil
}

// Do does the PingListenKey request.
func (s *PingListenKeyService) Do(ctx context.Context) (*PingListenKeyResponse, error) {
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

/* ==================== CloseListenKeyService ============================ */

// CloseListenKeyResponse
type CloseListenKeyResponse struct {
	ServiceBaseResponse
}

// CloseListenKeyService
type CloseListenKeyService struct {
	SM        common.ServiceMeta
	rc        common.RESTClient
	logger    *log.Entry
	listenKey string
}

// WithListenKey returns a copy of the service with listenKey set to the
// provided value.
func (s CloseListenKeyService) WithListenKey(listenKey string) *CloseListenKeyService {
	s.listenKey = listenKey
	return &s
}

// toParams converts all parameter fields of the service to a params struct.
func (s *CloseListenKeyService) toParams() *params {
	p := &params{}
	p.Set("listenKey", s.listenKey)
	return p
}

// parseResponse parses the request response into the CloseListenKeyResponse struct.
func (s *CloseListenKeyService) parseResponse(data []byte) (*CloseListenKeyResponse, error) {
	resp := &CloseListenKeyResponse{}
	if err := resp.ParseBaseResponse(&s.SM, s.rc.TimeHandler()); err != nil {
		return nil, err
	}
	return resp, nil
}

// Do does the CloseListenKey request.
func (s *CloseListenKeyService) Do(ctx context.Context) (*CloseListenKeyResponse, error) {
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
