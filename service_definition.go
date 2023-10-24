package binance

import (
	"net/http"
)

/* ==================== serviceDefinition ================================ */

type serviceDefinition struct {
	scheme              string
	method              string
	endpoint            BIEndpoint
	path                string
	endpointType        BIEndpointType
	securityType        BISecurityType
	primaryDatasource   BIDataSource
	secondaryDatasource BIDataSource
	weightIP            int // WeightLimit
	weightUID           int
}

var (
	apiServices = map[string]serviceDefinition{
		"ping":       {"https", http.MethodGet, endpointAPI, "/api/v3/ping", endpointTypeAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
		"serverTime": {"https", http.MethodGet, endpointAPI, "/api/v3/time", endpointTypeAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
	}

	fapiServices = map[string]serviceDefinition{
		"ping":       {"https", http.MethodGet, endpointFAPI, "/fapi/v1/ping", endpointTypeFAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
		"serverTime": {"https", http.MethodGet, endpointFAPI, "/fapi/v1/time", endpointTypeFAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
	}
)

/* ==================== apiServices ====================================== */

// NewSpotMarginServerTimeService returns a new api ServerTimeService.
func (c *Client) NewSpotMarginServerTimeService() *ServerTimeService {
	return &ServerTimeService{
		rc:                c.restClient,
		serviceDefinition: apiServices["serverTime"],
	}
}

func (c *Client) NewFuturesServerTimeService() *ServerTimeService {
	return &ServerTimeService{
		rc:                c.restClient,
		serviceDefinition: fapiServices["serverTime"],
	}
}
