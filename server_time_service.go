package binance

import (
	"context"
	"encoding/json"
)

/* ==================== ServerTimeService ================================ */

// ServerTimeService
type ServerTimeService struct {
	rc *RESTClient
	serviceDefinition
}

// toParams converts all parameter fields of the service to a params struct.
func (s *ServerTimeService) toParams() params {
	return params{}
}

// Do sends the request and returns a ServerTimeResponse and error, if any.
func (s *ServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	p := s.toParams()

	resp, meta, err := s.rc.Do(ctx, &s.serviceDefinition, p)

	if err != nil {
		return nil, err
	}

	return newServerTimeResponse(s.rc.c.timeHandler, meta, resp)
}

// newServerTimeResponse
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

// ServerTimeResponse
type ServerTimeResponse struct {
	*MetaDataREST
	ServerTime int64 `json:"serverTime"`
}
