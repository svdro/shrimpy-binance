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

/* ==================== restClient ======================================= */

func newRestClient(c *Client) *restClient {
	return &restClient{
		c:          c,
		httpClient: http.DefaultClient,
		logger:     c.logger.WithField("_caller", "restClient"),
	}
}

// restClient is responsible for making http requests to binance's REST APIs.
type restClient struct {
	c          *Client
	httpClient *http.Client
	logger     *log.Entry
}

// TimeHandler returns the TimeHandler associated with Client.
// This is needed when a service needs to know the current time.
func (rc *restClient) TimeHandler() common.TimeHandler {
	return rc.c.th
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

// doRequest makes the http request and records timestamps in MetaDataREST.
func (rc *restClient) doRequest(
	req *http.Request,
	sd *common.ServiceDefinition,
	sm *common.ServiceMeta,
) (*http.Response, error) {

	// get rlm
	rlm := rc.c.rlm

	// record timestamps
	sm.TSLSent = rc.c.th.TSLNow()
	sm.TSSSent = rc.c.th.TSSNow()
	defer func() {
		sm.TSLRecv = rc.c.th.TSLNow()
		sm.TSSRecv = rc.c.th.TSSNow()
	}()

	// register pending API call with rateLimitManager
	// return RetryAfterError if request would exceed the rate limit
	if err := rlm.RegisterPending(sd); err != nil {
		return nil, err
	}
	defer rlm.UnregisterPending(sd)

	// make request
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// parse response headers, update rateLimitManager
	if sm.SRH, err = parseServiceResponseHeader(resp.Header, sd.EndpointType); err != nil {
		return nil, err
	}

	rlm.UpdateUsed(sm.SRH.RateLimitUpdates, sm.SRH.TSSRespHeader)

	return resp, nil
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
		req.Header.Set("X-MBX-APIKEY", rc.c.apiKey)
	case common.SecurityTypeApiKey:
		req.Header.Set("X-MBX-APIKEY", rc.c.apiKey)
	case common.SecurityTypeNone:
	}

	return req, nil
}

// sign adds "timestamp", "recvWindow", and "signature" to urlValues.
func (rc *restClient) sign(urlValues *url.Values) {
	tsMilli := rc.c.th.TSSNow() / 1e6
	//tsMilli := rc.c.NanotoMilli(rc.c.TSNanoNowServer())
	urlValues.Add("timestamp", strconv.FormatInt(tsMilli, 10))
	urlValues.Add("recvWindow", strconv.Itoa(rc.c.recvWindow))
	queryString := urlValues.Encode()

	mac := hmac.New(sha256.New, []byte(rc.c.apiSecret))
	mac.Write([]byte(queryString))
	signature := mac.Sum(nil)

	urlValues.Add("signature", fmt.Sprintf("%x", signature))
}

// errResponse is a helper struct used to parse binance error responses.
type errResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// newRetryAfterError creates a new common.RetryAfterError.
func (rc *restClient) newRetryAfterError(
	tslRetryAt int64, statusCode int, errorCode int, errorMsg string,
) error {
	return &common.RetryAfterError{
		StatusCode:     statusCode,
		ErrorCode:      errorCode,
		Msg:            errorMsg,
		Producer:       "server",
		RetryTimeLocal: time.Unix(0, tslRetryAt),
		RetryAfter:     int(tslRetryAt-rc.c.th.TSLNow()) / 1e9,
	}
}

// getRetryTime calculates the the time at which the request should be retried
// based on the http response header.
func (rc *restClient) getRetryTime(srh *common.ServiceResponseHeader) (int64, error) {
	// make sure RetryAfter header has been parsed is not nil
	if srh.RetryAfter == nil {
		return 0, fmt.Errorf("calculateRetryTime: srh.RetryAfter is nil")
	}

	// better safe than sorry
	if srh.TSSRespHeader == 0 {
		return 0, fmt.Errorf("calculateRetryTime: srh.TSSRespHeader is 0")
	}

	// calculate and return retry time
	tssRetryAt := srh.TSSRespHeader + int64(*srh.RetryAfter)*1e9
	return rc.c.th.TSSToTSL(tssRetryAt), nil
}

// handleStatusCode handles the status code returned by the http response, as
// well as any binance error codes.
// TODO: implement other status codes and errors
func (rc *restClient) handleStatusCode(resp *http.Response, data []byte, sm *common.ServiceMeta) error {
	// update ServiceMeta status code
	sm.StatusCode = resp.StatusCode

	// don't do anything if status code is OK
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// parse error response. Don't do anything if error response can't be parsed.
	// if errResp is vital, then it should be handled downstream.
	errResp := &errResponse{}
	if err := json.Unmarshal(data, errResp); err != nil {
		rc.logger.WithField("data", string(data)).Error("handleStatusCode: error unmarshalling error response")
	}

	// handle status codes
	switch resp.StatusCode {

	// RetryAfterError (418, 429)
	case http.StatusTeapot, http.StatusTooManyRequests:
		tslRetryAt, err := rc.getRetryTime(sm.SRH)
		if err != nil {
			// should never happen, but if we can't get retry time, there's not point in continuing
			rc.logger.WithError(err).Fatal("handleStatusCode: error calculating retry time")
		}
		return rc.newRetryAfterError(tslRetryAt, resp.StatusCode, errResp.Code, errResp.Msg)

	// BadRequestError (400, 401)
	case http.StatusBadRequest, http.StatusUnauthorized:

	// UnexpectedStatusCodeError
	default:
		rc.logger.WithFields(log.Fields{
			"status_code": resp.StatusCode,
			"body":        fmt.Sprintf("%s", data),
		}).Fatal("Unexpected status code")
	}

	return nil
}
