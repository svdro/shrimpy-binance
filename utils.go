package binance

import (
	"encoding/json"
	"fmt"
	"net/url"
)

/* ==================== params ============================================*/

// params is a map that holds the parameters for a request.
type params map[string]string

// Set Sets the value of key to value.
func (p params) Set(key, value string) {
	p[key] = value
}

func (p params) SetSlice(key string, value []string) {
	b, _ := json.Marshal(value)
	p.Set(key, fmt.Sprintf("%s", b))
}

// SetIfNotNil Sets the value of key to value if value is not nil.
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
