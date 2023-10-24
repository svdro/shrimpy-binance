package binance

import (
	"context"
	"encoding/json"
)

/* ==================== ServerTimeService ================================ */

func (s *ServerTimeService) toParams() params {
	return params{}
}

type ServerTimeService struct {
	rc *RESTClient
	serviceDefinition
}

func (s *ServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	p := s.toParams()

	resp, meta, err := s.rc.Do(ctx, &s.serviceDefinition, p)

	if err != nil {
		return nil, err
	}

	return newServerTimeResponse(s.rc.c.timeHandler, meta, resp)
}

/* ==================== ServerTimeResponse ================================ */

func newServerTimeResponse(th *timeHandler, meta *MetaDataREST, data []byte) (*ServerTimeResponse, error) {
	r := &ServerTimeResponse{
		MetaDataREST: meta,
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	r.ServerTime = th.MillitoNano(r.ServerTime)
	return r, nil
}

type ServerTimeResponse struct {
	*MetaDataREST
	ServerTime int64 `json:"serverTime"`
}
