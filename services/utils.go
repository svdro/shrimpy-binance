package services

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== params ============================================*/

// params
type params map[string]string

// Set sets the value of key to value.
func (p params) Set(key, value string) {
	p[key] = value
}

func (p params) SetSlice(key string, value []string) {
	b, _ := json.Marshal(value)
	p.Set(key, fmt.Sprintf("%s", b))
}

// setIfNotNil sets the value of key to value if value is not nil.
func (p params) SetIfNotNil(key string, value *string) bool {
	if value == nil {
		return false
	}
	p.Set(key, *value)
	return true
}

func (p params) SetSliceIfNotNil(key string, value *[]string) bool {
	if value == nil {
		return false
	}

	p.SetSlice(key, *value)
	return true

}

func (p params) UrlValues() url.Values {
	values := url.Values{}
	for k, v := range p {
		values.Set(k, v)
	}
	return values
}

/* ==================== ServiceBaseResponse ============================== */

// ServiceBaseResponse holds meta data that is common to all responses.
type ServiceBaseResponse struct {
	TSLSent int64
	TSLRecv int64
	TSSSent int64
	TSSRecv int64
}

// ParseBaseResponse
func (s *ServiceBaseResponse) ParseBaseResponse(sm *common.ServiceMeta, th common.TimeHandler) error {
	s.TSLSent = sm.TSLSent
	s.TSLRecv = sm.TSLRecv
	s.TSSSent = sm.TSSSent
	s.TSSRecv = sm.TSSRecv
	return nil
}
