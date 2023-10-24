package binance

import "net/url"

/* ==================== params ============================================*/

// params is a map that holds the parameters for a request.
type params map[string]string

func (p params) urlValues() url.Values {
	values := url.Values{}
	for k, v := range p {
		values.Set(k, v)
	}
	return values
}
