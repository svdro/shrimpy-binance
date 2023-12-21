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

// Response is the interface that wraps the ParseBaseResponse method.
type Response interface {
	ParseBaseResponse(sm *common.ServiceMeta) error
}

// ServiceBaseResponse holds meta data that is common to all responses.
type ServiceBaseResponse struct {
	TSLSent common.TSNano
	TSLRecv common.TSNano
	TSSSent common.TSNano
	TSSRecv common.TSNano
}

// ParseBaseResponse
func (s *ServiceBaseResponse) ParseBaseResponse(sm *common.ServiceMeta) error {
	s.TSLSent = sm.TSLSent
	s.TSLRecv = sm.TSLRecv
	s.TSSSent = sm.TSSSent
	s.TSSRecv = sm.TSSRecv
	return nil
}

// boolToUpperStr converts a bool to an upper case string.
// some binance API endpoint expect "TRUE" or "FALSE" instead of regular
// true or false (e.g. MarginOrder IsIsolated).
func boolToUpperStr(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}
