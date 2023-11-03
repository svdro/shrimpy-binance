package services

import (
	"net/http"

	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== ServicesDefinitions ============================== */

var (
	APIServices = map[string]common.ServiceDefinition{
		"ping": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/ping",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
		"serverTime": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/time",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
		"createListenKey": {
			Scheme:              "https",
			Method:              http.MethodPost,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/userDataStream",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            2,
			WeightUID:           0,
		},
		"pingListenKey": {
			Scheme:              "https",
			Method:              http.MethodPut,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/userDataStream",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            2,
			WeightUID:           0,
		},
		"closeListenKey": {
			Scheme:              "https",
			Method:              http.MethodDelete,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/userDataStream",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            2,
			WeightUID:           0,
		},
	}

	SAPIServices = map[string]common.ServiceDefinition{
		"createListenKey": {
			Scheme:              "https",
			Method:              http.MethodPost,
			Endpoint:            common.EndpointAPI,
			Path:                "/sapi/v1/userDataStream",
			EndpointType:        common.EndpointTypeSAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
		"pingListenKey": {
			Scheme:              "https",
			Method:              http.MethodPut,
			Endpoint:            common.EndpointAPI,
			Path:                "/sapi/v1/userDataStream",
			EndpointType:        common.EndpointTypeSAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
		"closeListenKey": {
			Scheme:              "https",
			Method:              http.MethodDelete,
			Endpoint:            common.EndpointAPI,
			Path:                "/sapi/v1/userDataStream",
			EndpointType:        common.EndpointTypeSAPI,
			SecurityType:        common.SecurityTypeApiKey,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
	}

	FAPIServices = map[string]common.ServiceDefinition{
		"ping": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointFAPI,
			Path:                "/fapi/v1/ping",
			EndpointType:        common.EndpointTypeFAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},
	}
)

/* ==================== APIServices ====================================== */

func NewSpotMarginPingService(rc common.RESTClient) *PingService {
	return &PingService{
		SM: *common.NewServiceMeta(APIServices["ping"]),
		rc: rc,
	}
}

func NewSpotMarginServerTimeService(rc common.RESTClient) *ServerTimeService {
	return &ServerTimeService{
		SM: *common.NewServiceMeta(APIServices["serverTime"]),
		rc: rc,
	}
}

func NewSpotCreateListenKeyService(rc common.RESTClient) *CreateListenKeyService {
	return &CreateListenKeyService{
		SM: *common.NewServiceMeta(APIServices["createListenKey"]),
		rc: rc,
	}
}

func NewSpotPingListenKeyService(rc common.RESTClient) *PingListenKeyService {
	return &PingListenKeyService{
		SM: *common.NewServiceMeta(APIServices["pingListenKey"]),
		rc: rc,
	}
}

func NewSpotCloseListenKeyService(rc common.RESTClient) *CloseListenKeyService {
	return &CloseListenKeyService{
		SM: *common.NewServiceMeta(APIServices["closeListenKey"]),
		rc: rc,
	}
}

/* ==================== SAPIServices ===================================== */
func NewMarginCreateListenKeyService(rc common.RESTClient) *CreateListenKeyService {
	return &CreateListenKeyService{
		SM: *common.NewServiceMeta(SAPIServices["createListenKey"]),
		rc: rc,
	}
}

func NewMarginPingListenKeyService(rc common.RESTClient) *PingListenKeyService {
	return &PingListenKeyService{
		SM: *common.NewServiceMeta(SAPIServices["pingListenKey"]),
		rc: rc,
	}
}

func NewMarginCloseListenKeyService(rc common.RESTClient) *CloseListenKeyService {
	return &CloseListenKeyService{
		SM: *common.NewServiceMeta(SAPIServices["closeListenKey"]),
		rc: rc,
	}
}

/* ==================== FAPIServices ===================================== */

func NewFuturesPingService(rc common.RESTClient) *PingService {
	return &PingService{
		SM: *common.NewServiceMeta(FAPIServices["ping"]),
		rc: rc,
	}
}
