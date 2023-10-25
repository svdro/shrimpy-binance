package binance

import (
	"context"
	"encoding/json"
)

/* ==================== CreateListenKeyService =========================== */

// CreateListenKeyService
type CreateListenKeyService struct {
	rc *RESTClient
	serviceDefinition
}

// toParams converts all parameter fields of the service to a params struct.
func (s *CreateListenKeyService) toParams() params {
	return params{}
}

// Do sends the request and returns a CreateListenKeyResponse and error, if any.
func (s *CreateListenKeyService) Do(ctx context.Context) (*CreateListenKeyResponse, error) {
	p := s.toParams()

	resp, meta, err := s.rc.Do(ctx, &s.serviceDefinition, p)

	if err != nil {
		return nil, err
	}

	return newCreateListenKeyResponse(meta, resp)
}

// newCreateListenKeyResponse
func newCreateListenKeyResponse(meta *MetaDataREST, data []byte) (*CreateListenKeyResponse, error) {
	r := &CreateListenKeyResponse{
		MetaDataREST: meta,
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// CreateListenKeyResponse
type CreateListenKeyResponse struct {
	*MetaDataREST
	ListenKey string `json:"listenKey"`
}

/* ==================== PingListenKeyService ============================ */

// PingListenKeyService
type PingListenKeyService struct {
	rc *RESTClient
	serviceDefinition
	listenKey string
}

// WithListenKey returns a copy of the service with listenKey set to the provided value.
func (s PingListenKeyService) WithListenKey(listenKey string) *PingListenKeyService {
	s.listenKey = listenKey
	return &s
}

// toParams converts all parameter fields of the service to a params struct.
func (s *PingListenKeyService) toParams() params {
	p := params{}
	p.Set("listenKey", s.listenKey)
	return p
}

// Do sends the request and returns a PingListenKeyResponse and error, if any.
func (s *PingListenKeyService) Do(ctx context.Context) (*PingListenKeyResponse, error) {
	p := s.toParams()

	resp, meta, err := s.rc.Do(ctx, &s.serviceDefinition, p)

	if err != nil {
		return nil, err
	}

	return newPingListenKeyResponse(meta, resp)
}

// newPingListenKeyResponse
func newPingListenKeyResponse(meta *MetaDataREST, data []byte) (*PingListenKeyResponse, error) {
	r := &PingListenKeyResponse{
		MetaDataREST: meta,
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// PingListenKeyResponse
type PingListenKeyResponse struct {
	*MetaDataREST
}

/* ==================== CloseListenKeyService =========================== */
// TODO: implement
