package services

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

type SystemStatusService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
}

func (s *SystemStatusService) toParams() params {
	return params{}
}

func (s *SystemStatusService) Do(ctx context.Context) (*StatusResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())
	if err != nil {
		return nil, err
	}

	resp, err := s.parseResponse(data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// parseResponse parses the request response into the PingResponse struct.
func (s *SystemStatusService) parseResponse(data []byte) (*StatusResponse, error) {
	resp := &StatusResponse{}
	if err := resp.ParseBaseResponse(&s.SM); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type StatusResponse struct {
	ServiceBaseResponse
	Status int    `json:"status"` // 0: normal, 1: system maintenance
	Msg    string `json:"msg"`    // "normal" or "system maintenance"
}
