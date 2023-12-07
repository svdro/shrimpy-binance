package services

import (
	"net/http"

	log "github.com/sirupsen/logrus"
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
		"depth100": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/depth",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            5,
			WeightUID:           0,
		},
		"depth5000": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/depth",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            250,
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
		"depth1000": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointFAPI,
			Path:                "/fapi/v1/depth",
			EndpointType:        common.EndpointTypeFAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            20,
			WeightUID:           0,
		},
	}
)

/* ==================== APIServices ====================================== */

func NewSpotMarginPingService(rc common.RESTClient, logger *log.Entry) *PingService {
	return &PingService{
		SM:     *common.NewServiceMeta(APIServices["ping"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotMarginPingService"),
	}
}

func NewSpotMarginServerTimeService(rc common.RESTClient, logger *log.Entry) *ServerTimeService {
	return &ServerTimeService{
		SM:     *common.NewServiceMeta(APIServices["serverTime"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotMarginServerTimeService"),
	}
}

func NewSpotMarginDepth100Service(rc common.RESTClient, logger *log.Entry) *SpotMarginDepthService {
	return &SpotMarginDepthService{
		SM:     *common.NewServiceMeta(APIServices["depth100"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotMarginDepth5000Service"),
		depth:  100,
	}
}

func NewSpotMarginDepth5000Service(rc common.RESTClient, logger *log.Entry) *SpotMarginDepthService {
	return &SpotMarginDepthService{
		SM:     *common.NewServiceMeta(APIServices["depth5000"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotMarginDepth5000Service"),
		depth:  5000,
	}
}

func NewSpotCreateListenKeyService(rc common.RESTClient, logger *log.Entry) *CreateListenKeyService {
	return &CreateListenKeyService{
		SM:     *common.NewServiceMeta(APIServices["createListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotCreateListenKeyService"),
	}
}

func NewSpotPingListenKeyService(rc common.RESTClient, logger *log.Entry) *PingListenKeyService {
	return &PingListenKeyService{
		SM:     *common.NewServiceMeta(APIServices["pingListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotPingListenKeyService"),
	}
}

func NewSpotCloseListenKeyService(rc common.RESTClient, logger *log.Entry) *CloseListenKeyService {
	return &CloseListenKeyService{
		SM:     *common.NewServiceMeta(APIServices["closeListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotCloseListenKeyService"),
	}
}

/* ==================== SAPIServices ===================================== */
func NewMarginCreateListenKeyService(rc common.RESTClient, logger *log.Entry) *CreateListenKeyService {
	return &CreateListenKeyService{
		SM:     *common.NewServiceMeta(SAPIServices["createListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "MarginCreateListenKeyService"),
	}
}

func NewMarginPingListenKeyService(rc common.RESTClient, logger *log.Entry) *PingListenKeyService {
	return &PingListenKeyService{
		SM:     *common.NewServiceMeta(SAPIServices["pingListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "MarginPingListenKeyService"),
	}
}

func NewMarginCloseListenKeyService(rc common.RESTClient, logger *log.Entry) *CloseListenKeyService {
	return &CloseListenKeyService{
		SM:     *common.NewServiceMeta(SAPIServices["closeListenKey"]),
		rc:     rc,
		logger: logger.WithField("_caller", "MarginCloseListenKeyService"),
	}
}

/* ==================== FAPIServices ===================================== */

func NewFuturesPingService(rc common.RESTClient, logger *log.Entry) *PingService {
	return &PingService{
		SM:     *common.NewServiceMeta(FAPIServices["ping"]),
		rc:     rc,
		logger: logger.WithField("_caller", "FuturesPingService"),
	}
}

func NewFuturesDepth1000Service(rc common.RESTClient, logger *log.Entry) *FuturesDepthService {
	return &FuturesDepthService{
		SM:     *common.NewServiceMeta(FAPIServices["depth1000"]),
		rc:     rc,
		logger: logger.WithField("_caller", "FuturesDepth1000Service"),
		depth:  1000,
	}
}
