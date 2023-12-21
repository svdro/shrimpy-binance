package services

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

/* ==================== Spot Margin Exchange Info Service ================ */

type SpotMarginExchangeInfoService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
}

func (s *SpotMarginExchangeInfoService) toParams() *params {
	return &params{}
}

func (s *SpotMarginExchangeInfoService) parseResponse(data []byte) (*SpotMarginExchangeInfoResponse, error) {
	resp := &SpotMarginExchangeInfoResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if err := resp.ParseBaseResponse(&s.SM); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SpotMarginExchangeInfoService) Do(ctx context.Context) (*SpotMarginExchangeInfoResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())

	if err != nil {
		s.logger.WithError(err).Error("Do")
		return nil, err
	}

	resp, err := s.parseResponse(data)
	return resp, nil
}

/* ==================== Futures Exchange Info Service ==================== */

type FuturesExchangeInfoService struct {
	SM     common.ServiceMeta
	rc     common.RESTClient
	logger *log.Entry
}

func (s *FuturesExchangeInfoService) toParams() *params {
	return &params{}
}

func (s *FuturesExchangeInfoService) parseResponse(data []byte) (*FuturesExchangeInfoResponse, error) {
	resp := &FuturesExchangeInfoResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if err := resp.ParseBaseResponse(&s.SM); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *FuturesExchangeInfoService) Do(ctx context.Context) (*FuturesExchangeInfoResponse, error) {
	params := s.toParams()
	data, err := s.rc.Do(ctx, &s.SM, params.UrlValues())

	if err != nil {
		s.logger.WithError(err).Error("Do")
		return nil, err
	}

	resp, err := s.parseResponse(data)
	return resp, nil
}

/* ==================== Shared Exchange Info Response ==================== */

type SharedExchangeInfoResponse struct {
	ServiceBaseResponse

	Timezone        string                   `json:"timezone"`
	TSSServerTime   common.TSNano            `json:"serverTime"`
	RateLimits      []RateLimitEventResponse `json:"rateLimits"`
	ExchangeFilters []interface{}            `json:"exchangeFilters"`
}

// RateLimitEvent (SPOT_MARGIN & FUTURES)
type RateLimitEventResponse struct {
	RateLimitType common.BIRateLimitType         `json:"rateLimitType"`
	Interval      common.BIRateLimitIntervalType `json:"interval"`
	IntervalNum   int                            `json:"intervalNum"`
	Limit         int                            `json:"limit"`
}

type SharedExchangeInfoSymbolResponse struct {
	SymbolREST         string                   `json:"symbol"`
	BaseAssetREST      string                   `json:"baseAsset"`
	BaseAssetPrecision int                      `json:"baseAssetPrecision"`
	QuoteAssetREST     string                   `json:"quoteAsset"`
	QuotePrecision     int                      `json:"quotePrecision"`
	OrderTypes         []common.BIOrderType     `json:"orderTypes"`
	Filters            []map[string]interface{} `json:"filters"`
}

// PriceFilter (SPOT_MARGIN & FUTURES)
type PriceFilter struct {
	MinPrice string `json:"minPrice"`
	MaxPrice string `json:"maxPrice"`
	TickSize string `json:"tickSize"`
}

// LotSizeFilter (SPOT_MARGIN & FUTURES)
type LotSizeFilter struct {
	MinQty   string `json:"minQty"`
	MaxQty   string `json:"maxQty"`
	StepSize string `json:"stepSize"`
}

// MarketLotSizeFilter (SPOT_MARGIN & FUTURES)
type MarketLotSizeFilter struct {
	MinQty   string `json:"minQty"`
	MaxQty   string `json:"maxQty"`
	StepSize string `json:"stepSize"`
}

// MaxNumOrdersFilter (SPOT_MARGIN & FUTURES)
type MaxNumOrdersFilter struct {
	MaxNumOrders int `json:"maxNumOrders"`
}

// MaxNumAlgoOrdersFilter (SPOT_MARGIN & FUTURES)
type MaxNumAlgoOrdersFilter struct {
	MaxNumAlgoOrders int `json:"maxNumAlgoOrders"`
}

/* ==================== SpotMargin Exchange Info Response ================ */

type SpotMarginExchangeInfoResponse struct {
	SharedExchangeInfoResponse
	Symbols []SpotMarginExchangeInfoSymbolResponse `json:"symbols"`
}

type SpotMarginAllowedTradeTypesResponse struct {
	IcebergAllowed             bool `json:"icebergAllowed"`
	OCOAllowed                 bool `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop          bool `json:"allowTrailingStop"`
	CancelReplaceAllowed       bool `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed       bool `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool `json:"isMarginTradingAllowed"`
}

type SpotMarginExchangeInfoSymbolResponse struct {
	SharedExchangeInfoSymbolResponse
	SpotMarginAllowedTradeTypesResponse

	Status                          common.BISymbolStatusType             `json:"status"`
	QuoteAssetPrecision             int                                   `json:"quoteAssetPrecision"`
	BaseCommissionPrecision         int                                   `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int                                   `json:"quoteCommissionPrecision"`
	Permissions                     []common.BIAccountAndSymbolPermission `json:"permissions"`
	DefaultSelfTradePreventionMode  string                                `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string                              `json:"allowedSelfTradePreventionModes"`
}

// create aliases for shared filters for naming consistency
type SpotMarginPriceFilter = PriceFilter
type SpotMarginLotSizeFilter = LotSizeFilter
type SpotMarginMarketLotSizeFilter = MarketLotSizeFilter
type SpotMarginMaxNumOrdersFilter = MaxNumOrdersFilter
type SpotMarginMaxNumAlgoOrdersFilter = MaxNumAlgoOrdersFilter

// SpotMarginPercentPriceFilter (SPOT_MARGIN)
// NOTE: this is structurally different from FuturesPercentPriceFilter
type SpotMarginPercentPriceFilter struct {
	MultiplierUp   string `json:"multiplierUp"`
	MultiplierDown string `json:"multiplierDown"`
	AvgPriceMins   int    `json:"avgPriceMins"`
}

// SpotMarginMinNotionalFilter (SPOT_MARGIN)
// NOTE: this is structurally different from FuturesMinNotionalFilter
type SpotMarginMinNotionalFilter struct {
	MinNotional   string `json:"minNotional"`
	ApplyToMarket bool   `json:"applyToMarket"`
	AvgPriceMins  int    `json:"avgPriceMins"`
}

// SpotMarginPercentPriceBySideFilter (SPOT_MARGIN)
type SpotMarginPercentPriceBySideFilter struct {
	BidMultiplierUp   string `json:"bidMultiplierUp"`
	BidMultiplierDown string `json:"bidMultiplierDown"`
	AskMultiplierUp   string `json:"askMultiplierUp"`
	AskMultiplierDown string `json:"askMultiplierDown"`
	AvgPriceMins      int    `json:"avgPriceMins"`
}

// SpotMarginIcebergPartsFilter (SPOT_MARGIN)
type SpotMarginIcebergPartsFilter struct {
	Limit int `json:"limit"`
}

// SpotMarginMaxNumIcebergOrdersFilter (SPOT_MAGRIN)
type SpotMarginMaxNumIcebergOrdersFilter struct {
	MaxNumIcebergOrders int `json:"maxNumIcebergOrders"`
}

// SpotMarginMaxPositionFilter (SPOT_MARGIN)
type SpotMarginMaxPositionFilter struct {
	MaxPosition string `json:"maxPosition"`
}

// SpotMarginTrailingDeltaFilter (SPOT_MARGIN)
type SpotMarginTrailingDeltaFilter struct {
	MinTrailingAboveDelta int `json:"minTrailingAboveDelta"`
	MaxTrailingAboveDelta int `json:"maxTrailingAboveDelta"`
	MinTrailingBelowDelta int `json:"minTrailingBelowDelta"`
	MaxTrailingBelowDelta int `json:"maxTrailingBelowDelta"`
}

// SpotMarginNotionalFilter (SPOT_MARGIN)
type SpotMarginNotionalFilter struct {
	MinNotional      string `json:"minNotional"`
	ApplyMinToMarket bool   `json:"applyMinToMarket"`
	MaxNotional      string `json:"maxNotional"`
	ApplyMaxToMarket bool   `json:"applyMaxToMarket"`
	AvgPriceMins     int    `json:"avgPriceMins"`
}

// GetPriceFilter (SPOT_MARGIN & FUTURES)
func (r *SpotMarginExchangeInfoSymbolResponse) GetPriceFilter() (*SpotMarginPriceFilter, error) {
	return getFilter[SpotMarginPriceFilter](r.Filters, symbolFilterTypePriceFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetLotSizeFilter() (*SpotMarginLotSizeFilter, error) {
	return getFilter[SpotMarginLotSizeFilter](r.Filters, symbolFilterTypeLotSizeFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetMarketLotSizeFilter() (*SpotMarginMarketLotSizeFilter, error) {
	return getFilter[SpotMarginMarketLotSizeFilter](r.Filters, symbolFilterTypeMarketLotSizeFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetMaxNumOrdersFilter() (*SpotMarginMaxNumOrdersFilter, error) {
	return getFilter[SpotMarginMaxNumOrdersFilter](r.Filters, symbolFilterTypeMaxNumOrdersFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetMaxNumAlgoOrdersFilter() (*SpotMarginMaxNumAlgoOrdersFilter, error) {
	return getFilter[SpotMarginMaxNumAlgoOrdersFilter](r.Filters, symbolFilterTypeMaxNumAlgoOrdersFilter)
}

// GetPercentPriceFilter returns the PercentPriceFilter for binance's api/sapi
// API. PercentPriceFilter is listed in the api documentation, but in practice
// it's not in the response.
func (r *SpotMarginExchangeInfoSymbolResponse) GetPercentPriceFilter() (*SpotMarginPercentPriceFilter, error) {
	return getFilter[SpotMarginPercentPriceFilter](r.Filters, symbolFilterTypePercentPriceFilterSpotMargin)
}

// GetMinNotionalFilter returns the MinNotionalFilter for binance's api/sapi
// API. MinNotionalFilter is listed in the api documentation, but in practice
// it's not in the response.
func (r *SpotMarginExchangeInfoSymbolResponse) GetMinNotionalFilter() (*SpotMarginMinNotionalFilter, error) {
	return getFilter[SpotMarginMinNotionalFilter](r.Filters, symbolFilterTypeMinNotionalFilterSpotMargin)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetPercentPriceBySideFilter() (*SpotMarginPercentPriceBySideFilter, error) {
	return getFilter[SpotMarginPercentPriceBySideFilter](r.Filters, symbolFilterTypePercentPriceBySideFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetIcebergPartsFilter() (*SpotMarginIcebergPartsFilter, error) {
	return getFilter[SpotMarginIcebergPartsFilter](r.Filters, symbolFilterTypeIcebergPartsFilter)
}

// GetMaxNumIcebergOrdersFilter returns the MaxNumIcebergOrdersFilter for
// binance's api/sapi API. MaxNumIcebergOrdersFilter is listed in the api
// documentation, but in practice it's not in the response.
func (r *SpotMarginExchangeInfoSymbolResponse) GetMaxNumIcebergOrdersFilter() (*SpotMarginMaxNumIcebergOrdersFilter, error) {
	return getFilter[SpotMarginMaxNumIcebergOrdersFilter](r.Filters, symbolFilterTypeMaxNumIcebergOrdersFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetMaxPositionFilter() (*SpotMarginMaxPositionFilter, error) {
	return getFilter[SpotMarginMaxPositionFilter](r.Filters, symbolFilterTypeMaxPositionFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetTrailingDeltaFilter() (*SpotMarginTrailingDeltaFilter, error) {
	return getFilter[SpotMarginTrailingDeltaFilter](r.Filters, symbolFilterTypeTrailingDeltaFilter)
}

func (r *SpotMarginExchangeInfoSymbolResponse) GetNotionalFilter() (*SpotMarginNotionalFilter, error) {
	return getFilter[SpotMarginNotionalFilter](r.Filters, symbolFilterTypeNotionalFilter)
}

/* ==================== Futures Exchange Info Response =================== */

type FuturesExchangeInfoResponse struct {
	SharedExchangeInfoResponse
	FuturesMarginResponse

	Symbols []FuturesExchangeInfoSymbolResponse `json:"symbols"`
}

type FuturesExchangeInfoSymbolResponse struct {
	SharedExchangeInfoSymbolResponse
	FuturesContractResponse
	FuturesMarginResponse

	PairName          string                      `json:"pair"`
	PricePrecision    int                         `json:"pricePrecision"`
	QuantityPrecision int                         `json:"quantityPrecision"`
	UnderlyingType    string                      `json:"underlyingType"`
	UnderlyingSubType []string                    `json:"underlyingSubType"`
	SettlePlan        int                         `json:"settlePlan"`
	TriggerProtect    string                      `json:"triggerProtect"`
	TimeInForce       []common.BIOrderTimeInForce `json:"timeInForce"`
	LiquidationFee    string                      `json:"liquidationFee"`
	MarketTakeBound   string                      `json:"marketTakeBound"`
}

type FuturesMarginResponse struct {
	MaintMarginPercent    string `json:"maintMarginPercent"`    // Ignore
	RequiredMarginPercent string `json:"requiredMarginPercent"` // Ignore
	MarginAssetREST       string `json:"marginAsset"`
}

type FuturesContractResponse struct {
	ContractStatus common.BIContractStatus `json:"status"`
	ContractType   common.BIContractType   `json:"contractType"`
	DeliveryDate   int64                   `json:"deliveryDate"`
	OnboardDate    int64                   `json:"onboardDate"`
}

// create aliases for shared filters for naming consistency
type FuturesPriceFilter = PriceFilter
type FuturesLotSizeFilter = LotSizeFilter
type FuturesMarketLotSizeFilter = MarketLotSizeFilter
type FuturesMaxNumOrdersFilter = MaxNumOrdersFilter
type FuturesMaxNumAlgoOrdersFilter = MaxNumAlgoOrdersFilter

// FuturesPercentPriceFilter (FUTURES)
// NOTE: this is structurally different from SpotMarginPercentPriceFilter
type FuturesPercentPriceFilter struct {
	MultiplierUp      string `json:"multiplierUp"`
	MultiplierDown    string `json:"multiplierDown"`
	MultiplierDecimal string `json:"multiplierDecimal"`
}

// FuturesMinNotionalFilter (FUTURES)
// NOTE: this is structurally different from SpotMarginMinNotionalFilter
type FuturesMinNotionalFilter struct {
	MinNotional string `json:"notional"`
}

func (r *FuturesExchangeInfoSymbolResponse) GetPriceFilter() (*FuturesPriceFilter, error) {
	return getFilter[FuturesPriceFilter](r.Filters, symbolFilterTypePriceFilter)
}

func (r *FuturesExchangeInfoSymbolResponse) GetLotSizeFilter() (*FuturesLotSizeFilter, error) {
	return getFilter[FuturesLotSizeFilter](r.Filters, symbolFilterTypeLotSizeFilter)
}

func (r *FuturesExchangeInfoSymbolResponse) GetMarketLotSizeFilter() (*FuturesMarketLotSizeFilter, error) {
	return getFilter[FuturesMarketLotSizeFilter](r.Filters, symbolFilterTypeMarketLotSizeFilter)
}

func (r *FuturesExchangeInfoSymbolResponse) GetMaxNumOrdersFilter() (*FuturesMaxNumOrdersFilter, error) {
	return getFilter[FuturesMaxNumOrdersFilter](r.Filters, symbolFilterTypeMaxNumOrdersFilter)
}

func (r *FuturesExchangeInfoSymbolResponse) GetMaxNumAlgoOrdersFilter() (*FuturesMaxNumAlgoOrdersFilter, error) {
	return getFilter[FuturesMaxNumAlgoOrdersFilter](r.Filters, symbolFilterTypeMaxNumAlgoOrdersFilter)
}

func (r *FuturesExchangeInfoSymbolResponse) GetPercentPriceFilter() (*FuturesPercentPriceFilter, error) {
	return getFilter[FuturesPercentPriceFilter](r.Filters, symbolFilterTypePercentPriceFilterFutures)
}

func (r *FuturesExchangeInfoSymbolResponse) GetMinNotionalFilter() (*FuturesMinNotionalFilter, error) {
	return getFilter[FuturesMinNotionalFilter](r.Filters, symbolFilterTypeMinNotionalFilterFutures)
}

/* ==================== response utils =================================== */

// constants for identifying filters.
// these are never used outside of this file, so should not be in shrimpy-binance/common
const (
	symbolFilterTypePriceFilter                  string = "PRICE_FILTER"           // (SPOT & MARGIN & FUTURES)
	symbolFilterTypeLotSizeFilter                string = "LOT_SIZE"               // (SPOT & MARGIN & FUTURES)
	symbolFilterTypeMarketLotSizeFilter          string = "MARKET_LOT_SIZE"        // (SPOT & MARGIN & FUTURES)
	symbolFilterTypeMaxNumOrdersFilter           string = "MAX_NUM_ORDERS"         // (SPOT & MARGIN & FUTURES)
	symbolFilterTypeMaxNumAlgoOrdersFilter       string = "MAX_NUM_ALGO_ORDERS"    // (SPOT & MARGIN & FUTURES)
	symbolFilterTypePercentPriceFilterSpotMargin string = "PERCENT_PRICE"          // (SPOT & MARGIN)
	symbolFilterTypePercentPriceBySideFilter     string = "PERCENT_PRICE_BY_SIDE"  // (SPOT & MARGIN)
	symbolFilterTypeMinNotionalFilterSpotMargin  string = "MIN_NOTIONAL"           // (SPOT & MARGIN)
	symbolFilterTypeNotionalFilter               string = "NOTIONAL"               // (SPOT & MARGIN)
	symbolFilterTypeIcebergPartsFilter           string = "ICEBERG_PARTS"          // (SPOT & MARGIN)
	symbolFilterTypeMaxNumIcebergOrdersFilter    string = "MAX_NUM_ICEBERG_ORDERS" // (SPOT & MARGIN)
	symbolFilterTypeMaxPositionFilter            string = "MAX_POSITION"           // (SPOT & MARGIN)
	symbolFilterTypeTrailingDeltaFilter          string = "TRAILING_DELTA"         // (SPOT & MARGIN)
	symbolFilterTypeMinNotionalFilterFutures     string = "MIN_NOTIONAL"           // (FUTURES)
	symbolFilterTypePercentPriceFilterFutures    string = "PERCENT_PRICE"          // (FUTURES)
)

// findFilter returns the first filter that matches filterType from the filters.
func findFilter(filters []map[string]interface{}, filterType string) (map[string]interface{}, error) {
	for _, filter := range filters {
		filterTypeStr, ok := filter["filterType"].(string)
		if !ok {
			return nil, fmt.Errorf("findFilter: filterType is not a string")
		}

		if filterTypeStr == string(filterType) {
			return filter, nil
		}
	}
	return nil, nil
}

// getFilter is a helper function that searches the filters for the first
// filter that matches filterType parses it into a struct of type T, and
// returns it.
func getFilter[T any](filters []map[string]interface{}, filterType string) (*T, error) {
	result := new(T)

	filter, err := findFilter(filters, filterType)
	if err != nil {
		return result, err
	}

	if filter == nil {
		return result, fmt.Errorf("no filter found for filterType %s", filterType)
	}

	rawFilter, err := json.Marshal(filter)
	if err != nil {
		return result, err
	}

	//var result T
	if err = json.Unmarshal(rawFilter, &result); err != nil {
		return result, err
	}
	return result, nil
}
