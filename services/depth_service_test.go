package services

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svdro/shrimpy-binance/common"
)

type MockRestClient struct{}

func (m *MockRestClient) Do(ctx context.Context, sm *common.ServiceMeta, p url.Values) ([]byte, error) {
	return apiDepthData, nil
}

var (
	apiDepthData = []byte(`{
  "lastUpdateId": 1027024,
  "bids": [
    [
      "4.00000000",
      "431.00000000"
    ]
  ],
  "asks": [
    [
      "4.00000200",
      "12.00000000"
    ]
  ]
}`)

	apiDepthTarget = &SpotMarginDepthResponse{
		SharedDepthResponse: SharedDepthResponse{
			LastUpdateID: 1027024,
			Bids:         []Level{{Price: "4.00000000", Qty: "431.00000000"}},
			Asks:         []Level{{Price: "4.00000200", Qty: "12.00000000"}},
		},
		ServiceBaseResponse: ServiceBaseResponse{}}
)

func TestUnmarshalDepthResponse(t *testing.T) {
	resp := &SpotMarginDepthResponse{}

	err := json.Unmarshal(apiDepthData, resp)
	assert.Nil(t, err)
	assert.Equal(t, apiDepthTarget, resp)
}

func TestSpotMarginDepthService(t *testing.T) {
	service := &SpotMarginDepthService{
		SM:     *common.NewServiceMeta(APIServices["depth5000"]),
		rc:     &MockRestClient{},
		logger: nil,
		depth:  5000,
	}

	resp, err := service.parseResponse(apiDepthData)
	assert.Nil(t, err)
	assert.IsType(t, &SpotMarginDepthResponse{}, resp)
}
