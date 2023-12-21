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

		"exchangeInfo": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/api/v3/exchangeInfo",
			EndpointType:        common.EndpointTypeAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceMemory,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            10,
			WeightUID:           0,
		},
	}

	SAPIServices = map[string]common.ServiceDefinition{
		"systemStatus": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointAPI,
			Path:                "/sapi/v1/system/status",
			EndpointType:        common.EndpointTypeSAPI,
			SecurityType:        common.SecurityTypeNone,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            1,
			WeightUID:           0,
		},

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

		"createMarginOrder": {
			Scheme:              "https",
			Method:              http.MethodPost,
			Endpoint:            common.EndpointAPI,
			Path:                "/sapi/v1/margin/order",
			EndpointType:        common.EndpointTypeSAPI,
			SecurityType:        common.SecurityTypeSigned,
			PrimaryDatasource:   common.DataSourceNone,
			SecondaryDatasource: common.DataSourceNone,
			WeightIP:            0,
			WeightUID:           6,
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

		"serverTime": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointFAPI,
			Path:                "/fapi/v1/time",
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

		"exchangeInfo": {
			Scheme:              "https",
			Method:              http.MethodGet,
			Endpoint:            common.EndpointFAPI,
			Path:                "/fapi/v1/exchangeInfo",
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

func NewSpotMarginExchangeInfoService(rc common.RESTClient, logger *log.Entry) *SpotMarginExchangeInfoService {
	return &SpotMarginExchangeInfoService{
		SM:     *common.NewServiceMeta(APIServices["exchangeInfo"]),
		rc:     rc,
		logger: logger.WithField("_caller", "SpotMarginExchangeInfoService"),
	}
}

/* ==================== SAPIServices ===================================== */

func NewMarginSystemStatusService(rc common.RESTClient, logger *log.Entry) *SystemStatusService {
	return &SystemStatusService{
		SM:     *common.NewServiceMeta(SAPIServices["systemStatus"]),
		rc:     rc,
		logger: logger.WithField("_caller", "MarginSystemStatusService"),
	}
}

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

func NewCreateMarginOrderService(rc common.RESTClient, logger *log.Entry) *CreateMarginOrderService {
	return &CreateMarginOrderService{
		SM:     *common.NewServiceMeta(SAPIServices["createMarginOrder"]),
		rc:     rc,
		logger: logger.WithField("_caller", "CreateMarginOrderService"),
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

func NewFuturesServerTimeService(rc common.RESTClient, logger *log.Entry) *ServerTimeService {
	return &ServerTimeService{
		SM:     *common.NewServiceMeta(FAPIServices["serverTime"]),
		rc:     rc,
		logger: logger.WithField("_caller", "FuturesServerTimeService"),
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

func NewFuturesExchangeInfoService(rc common.RESTClient, logger *log.Entry) *FuturesExchangeInfoService {
	return &FuturesExchangeInfoService{
		SM:     *common.NewServiceMeta(FAPIServices["exchangeInfo"]),
		rc:     rc,
		logger: logger.WithField("_caller", "FuturesExchangeInfoService"),
	}
}
