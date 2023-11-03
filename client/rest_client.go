package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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
	rc.logger.Info("Do")

	// fetch rate limit handler, panic if not found
	sd := sm.SD
	rlh := rc.c.getRlh(sd.EndpointType)

	// create request/ handle security
	req, err := rc.createRequest(ctx, &sd, p)
	if err != nil {
		return nil, err
	}

	// do request
	resp, err := rc.doRequest(req, &sd, rlh, sm)
	if err != nil {
		return nil, err
	}

	// read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle status code and binance error codes
	rc.handleStatusCode(resp, data, sm)

	return data, nil
}

// doRequest makes the http request and records timestamps in MetaDataREST.
func (rc *restClient) doRequest(
	req *http.Request,
	sd *common.ServiceDefinition,
	rlh *rateLimitHandler,
	sm *common.ServiceMeta,
) (*http.Response, error) {

	// record timestamps
	sm.TSLSent = rc.c.th.TSLNow()
	sm.TSSSent = rc.c.th.TSSNow()
	defer func() {
		sm.TSLRecv = rc.c.th.TSLNow()
		sm.TSSRecv = rc.c.th.TSSNow()
	}()

	// register pending API call with rate limit handler
	// return RetryAfterError if request would exceed rate limit
	err := rlh.RegisterPending(sd)
	if err != nil {
		return nil, err
	}
	defer rlh.UnregisterPending(sd)

	// make request
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// parse response headers, update rate limit handler
	if sm.SRH, err = parseServiceResponseHeader(resp.Header); err != nil {
		return nil, err
	}
	rlh.UpdateUsed(sm)

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
	tsMilli := rc.c.th.TSLNow() / 1e6
	//tsMilli := rc.c.NanotoMilli(rc.c.TSNanoNowServer())
	urlValues.Add("timestamp", strconv.FormatInt(tsMilli, 10))
	urlValues.Add("recvWindow", strconv.Itoa(rc.c.recvWindow))
	queryString := urlValues.Encode()

	mac := hmac.New(sha256.New, []byte(rc.c.apiSecret))
	mac.Write([]byte(queryString))
	signature := mac.Sum(nil)

	urlValues.Add("signature", fmt.Sprintf("%x", signature))
}

// handleStatusCode handles the status code returned by the http response, as
// well as any binance error codes.
// TODO: implement other status codes and errors
func (rc *restClient) handleStatusCode(resp *http.Response, data []byte, sm *common.ServiceMeta) error {
	// update ServiceMeta status code
	sm.StatusCode = resp.StatusCode

	// handle status code
	switch resp.StatusCode {

	// ok
	case http.StatusOK:

	// RetryAfterError
	case http.StatusTeapot:

	// RetryAfterError
	case http.StatusTooManyRequests:

	// BadRequestError
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
