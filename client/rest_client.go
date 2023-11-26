package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/svdro/shrimpy-binance/common"
)

// newRestClient creates a new restClient.
func newRestClient(th common.TimeHandler, rlm *rateLimitManager, apiConfig *APIConfig, logger *log.Entry) *restClient {
	return &restClient{
		th:         th,
		rlm:        rlm,
		apiConfig:  apiConfig,
		httpClient: http.DefaultClient,
		logger:     logger.WithField("_caller", "restClient"),
	}
}

// restClient is responsible for making http requests to binance's REST APIs.
type restClient struct {
	th         common.TimeHandler
	rlm        *rateLimitManager
	apiConfig  *APIConfig
	httpClient *http.Client
	logger     *log.Entry
}

// TimeHandler returns the TimeHandler associated with Client.
// This is needed when a service that uses the restClient needs to
// access the TimeHandler.
func (rc *restClient) TimeHandler() common.TimeHandler {
	return rc.th
}

// Do makes an http request to a binance REST API.
// All data needed to make the request is contained in ServiceMeta (SD).
// Any meta data that is created during the request is stored in ServiceMeta.
func (rc *restClient) Do(ctx context.Context, sm *common.ServiceMeta, p url.Values) ([]byte, error) {

	// create request/ handle security
	req, err := rc.createRequest(ctx, &sm.SD, p)
	if err != nil {
		return nil, err
	}

	// do request
	resp, err := rc.doRequest(req, &sm.SD, sm)
	if err != nil {
		rc.logger.WithError(err).Error("Error doing request")
		return nil, err
	}
	defer resp.Body.Close()

	// read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle status code and binance error codes
	if err := rc.handleStatusCode(resp, data, sm); err != nil {
		return nil, err
	}

	return data, nil
}

// doRequest makes the http request.
// It records timestamps in ServiceMeta.
// It updates the rateLimitManager.
func (rc *restClient) doRequest(
	req *http.Request,
	sd *common.ServiceDefinition,
	sm *common.ServiceMeta,
) (*http.Response, error) {

	// record timestamps
	sm.TSLSent = rc.th.TSLNow()
	sm.TSSSent = rc.th.TSSNow()
	defer func() {
		sm.TSLRecv = rc.th.TSLNow()
		sm.TSSRecv = rc.th.TSSNow()
	}()

	// register pending API call with rateLimitManager
	// return RetryAfterError if request would exceed the rate limit
	if err := rc.rlm.RegisterPending(sd); err != nil {
		return nil, err
	}
	defer rc.rlm.UnregisterPending(sd)

	// make request
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// parse response headers, update rateLimitManager
	if sm.SRH, err = parseServiceResponseHeader(resp.Header, sd.EndpointType); err != nil {
		return nil, err
	}

	rc.rlm.UpdateUsed(sm.SRH.RateLimitUpdates, sm.SRH.TSSRespHeader)

	return resp, nil
}

// sign adds "timestamp", "recvWindow", and "signature" to urlValues.
// NOTE: the "X-MBX-APIKEY" header is not added here.
func (rc *restClient) sign(urlValues *url.Values) {
	// add timestamp and recvWindow to urlValues, encode query string
	tsMilli := rc.th.TSSNow().Int64() / 1e6
	urlValues.Add("timestamp", strconv.FormatInt(tsMilli, 10))
	urlValues.Add("recvWindow", strconv.Itoa(rc.apiConfig.recvWindow))
	queryString := urlValues.Encode()

	// calculate signature and add to urlValues
	mac := hmac.New(sha256.New, []byte(rc.apiConfig.apiSecret))
	mac.Write([]byte(queryString))
	signature := mac.Sum(nil)
	urlValues.Add("signature", fmt.Sprintf("%x", signature))
}

// createRequest creates a new http request and implements binance's endpoint
// security protocol.
func (rc *restClient) createRequest(
	ctx context.Context,
	sd *common.ServiceDefinition,
	urlValues url.Values,
) (*http.Request, error) {

	uri := url.URL{Scheme: "https", Host: string(sd.Endpoint), Path: sd.Path, RawQuery: urlValues.Encode()}
	req, err := http.NewRequestWithContext(ctx, sd.Method, uri.String(), nil)
	if err != nil {
		rc.logger.WithError(err).Error("Error creating request")
		return nil, err
	}

	// handle Security
	switch sd.SecurityType {
	case common.SecurityTypeSigned:
		rc.sign(&urlValues)
		req.URL.RawQuery = urlValues.Encode()
		req.Header.Set("X-MBX-APIKEY", rc.apiConfig.apiKey)
	case common.SecurityTypeApiKey:
		req.Header.Set("X-MBX-APIKEY", rc.apiConfig.apiKey)
	case common.SecurityTypeNone:
	}

	return req, nil
}

// errResponse is a helper struct used to parse binance error responses.
type errResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// newBadRequestError creates a new common.BadRequestError.
// It is used when the http response status code is 400 or 401.
func (rc *restClient) newBadRequestError(statusCode int, errResp *errResponse) error {
	return &common.BadRequestError{
		StatusCode: statusCode,
		ErrorCode:  errResp.Code,
		Msg:        errResp.Msg,
	}
}

// newUnexpectedStatusCodeError creates a new common.UnexpectedStatusCodeError.
// It is used for all other status codes.
func (rc *restClient) newUnexpectedStatusCodeError(statusCode int, errResp *errResponse) error {
	return &common.UnexpectedStatusCodeError{
		StatusCode: statusCode,
		ErrorCode:  errResp.Code,
		Msg:        errResp.Msg,
	}
}

// newRetryAfterError creates a new common.RetryAfterError.
// It is used when the http response status code is 418 or 429.
func (rc *restClient) newRetryAfterError(
	tslRetryAt common.TSNano, statusCode int, errorCode int, errorMsg string,
) error {
	return &common.RetryAfterError{
		StatusCode:     statusCode,
		ErrorCode:      errorCode,
		Msg:            errorMsg,
		Producer:       "server",
		RetryTimeLocal: time.Unix(0, tslRetryAt.Int64()),
		RetryAfter:     int(tslRetryAt-rc.th.TSLNow()) / 1e9,
	}
}

// getRetryTime calculates the the time at which the request should be retried
// based on the http response header.
func (rc *restClient) getRetryTime(srh *common.ServiceResponseHeader) (common.TSNano, error) {
	// make sure RetryAfter header has been parsed is not nil
	if srh.RetryAfter == nil {
		return 0, fmt.Errorf("calculateRetryTime: srh.RetryAfter is nil")
	}

	// better safe than sorry
	if srh.TSSRespHeader == 0 {
		return 0, fmt.Errorf("calculateRetryTime: srh.TSSRespHeader is 0")
	}

	// calculate and return retry time
	tssRetryAt := common.TSNano(srh.TSSRespHeader.Int64() + int64(*srh.RetryAfter)*1e9)
	return rc.th.TSSToTSL(tssRetryAt), nil
}

// handleRetryAfterError generates and returns a new common.RetryAfterError.
// It is used when the http response status code is 418 or 429.
// Panic if the RetryAfter header has not been parsed. In practice this should
// never happen, but if it does, there's no point in continuing.
func (rc *restClient) handleRetryAfterError(
	statusCode int, srh *common.ServiceResponseHeader, errResp *errResponse) error {
	tslRetryAt, err := rc.getRetryTime(srh)
	if err != nil {
		rc.logger.WithError(err).Panic("handleStatusCode: error calculating retry time")
	}
	return rc.newRetryAfterError(tslRetryAt, statusCode, errResp.Code, errResp.Msg)
}

// handleStatusCode handles the status code returned by the http response, as
// well as any binance error codes.
// If the status code is not OK, it returns one of three errors:
// * 418, 428: common.RetryAfterError
// * 400, 401: common.BadRequestError
// * other: common.UnexpectedStatusCodeError
func (rc *restClient) handleStatusCode(resp *http.Response, data []byte, sm *common.ServiceMeta) error {
	// update ServiceMeta status code
	sm.StatusCode = resp.StatusCode

	// don't do anything if status code is OK
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// parse error response. Don't do anything if error response can't be parsed.
	// if errResp fields are vital, then it should be handled downstream.
	errResp := &errResponse{}
	if err := json.Unmarshal(data, errResp); err != nil {
		rc.logger.WithField("data", string(data)).Error("handleStatusCode: error unmarshalling error response")
	}

	// handle status codes
	switch resp.StatusCode {

	case http.StatusTeapot, http.StatusTooManyRequests: // 418, 429
		return rc.handleRetryAfterError(resp.StatusCode, sm.SRH, errResp)

	case http.StatusBadRequest, http.StatusUnauthorized: // 400, 401
		return rc.newBadRequestError(resp.StatusCode, errResp)

	default:
		return rc.newUnexpectedStatusCodeError(resp.StatusCode, errResp) // other
	}
}
