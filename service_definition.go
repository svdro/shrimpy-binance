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
		"ping":                {"https", http.MethodGet, endpointAPI, "/api/v3/ping", endpointTypeAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
		"serverTime":          {"https", http.MethodGet, endpointAPI, "/api/v3/time", endpointTypeAPI, SecurityTypeNone, DataSourceMemory, DataSourceNone, 1, 0},
		"createListenKeySpot": {"https", http.MethodPost, endpointAPI, "/api/v3/userDataStream", endpointTypeAPI, SecurityTypeApiKey, DataSourceMemory, DataSourceNone, 2, 0},
		"pingListenKeySpot":   {"https", http.MethodPut, endpointAPI, "/api/v3/userDataStream", endpointTypeAPI, SecurityTypeApiKey, DataSourceMemory, DataSourceNone, 2, 0},
		//"closeListenKeySpot":  {"https", http.MethodDelete, endpointAPI, "/api/v3/userDataStream", endpointTypeAPI, SecurityTypeApiKey, DataSourceMemory, DataSourceNone, 2, 0},
		"createListenKeyMargin": {"https", http.MethodPost, endpointAPI, "/sapi/v1/userDataStream", endpointTypeSAPI, SecurityTypeApiKey, DataSourceNone, DataSourceNone, 1, 0},
		"pingListenKeyMargin":   {"https", http.MethodPut, endpointAPI, "/sapi/v1/userDataStream", endpointTypeSAPI, SecurityTypeApiKey, DataSourceNone, DataSourceNone, 1, 0},
		//"closeListenKeyMargin":  {"https", http.MethodDelete, endpointAPI, "/sapi/v1/userDataStream", endpointTypeSAPI, SecurityTypeApiKey, DataSourceNone, DataSourceNone, 1, 0},
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

// NewSpotMarginListenKeySpotService returns a new api ListenKeyService.
func (c *Client) NewSpotMarginCreateListenKeySpotService() *CreateListenKeyService {
	return &CreateListenKeyService{
		rc:                c.restClient,
		serviceDefinition: apiServices["createListenKeySpot"],
	}
}

func (c *Client) NewSpotMarginCreateListenKeyMarginService() *CreateListenKeyService {
	return &CreateListenKeyService{
		rc:                c.restClient,
		serviceDefinition: apiServices["createListenKeyMargin"],
	}
}

func (c *Client) NewSpotMarginPingListenKeySpotService() *PingListenKeyService {
	return &PingListenKeyService{
		rc:                c.restClient,
		serviceDefinition: apiServices["pingListenKeySpot"],
	}
}

func (c *Client) NewSpotMarginPingListenKeyMarginService() *PingListenKeyService {
	return &PingListenKeyService{
		rc:                c.restClient,
		serviceDefinition: apiServices["pingListenKeyMargin"],
	}
}

/* ==================== fapiServices ===================================== */

func (c *Client) NewFuturesServerTimeService() *ServerTimeService {
	return &ServerTimeService{
		rc:                c.restClient,
		serviceDefinition: fapiServices["serverTime"],
	}
}
