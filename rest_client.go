package binance

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
)

/* ==================== BIResponseHeader ================================= */

// RateLimitHeader contains used rate limit updates.
type rateLimitHeader struct {
	RateLimitType BIRateLimitType
	IntervalType  BIRateLimitIntervalType
	IntervalNum   int
	Count         int // weight count
}

func newResponseHeader(headers http.Header) *responseHeader {
	return &responseHeader{}
}

// BIResponseHeader contains all relevant information from binance http
// response headers.
type responseHeader struct {
	Server    string                    // (API, SAPI, FAPI, DAPI, EAPI)
	TSNano    int64                     // (API, SAPI, FAPI, DAPI, EAPI)
	UIDLimits map[int64]rateLimitHeader // (API, SAPI, FAPI, DAPI, EAPI)
	IPLimits  map[int64]rateLimitHeader // (API, SAPI, FAPI, DAPI, EAPI)
}

/* ==================== MetaData ========================================= */

// MetaDataREST holds the meta data that RESTClient produces when making a
// request. It is included in all REST responses.
// Optionally, RESTClient can produce a channel of MetaDataREST for analytics
// and logging purposes.
// TODO: Figure out what this should actually be.
// TODO: I want this to have most fields from serviceDefinition for analytics,
// TODO: but I don't want to include it in every response.
type MetaDataREST struct {
	Server      string // server that handled the request
	StatusCode  int    // status code
	TRespHeader int64  // ts (server time) in response header
	T0_Local    int64  // ts nano (local time) request dispatched
	T0_Server   int64  // ts nano (server time) request dispatched
	T3_Local    int64  // ts nano (local time) response received
	T3_Server   int64  // ts nano (server time) response received
}

/* ==================== RESTClient ======================================= */

// RESTClient is responsible for interacting with All Binance REST APIs.
type RESTClient struct {
	c          *Client
	httpClient *http.Client
	logger     *log.Entry
}

// Do makes the http request. It returns the response or error, if any.
func (rc *RESTClient) Do(
	ctx context.Context,
	sd *serviceDefinition,
	p params,
) ([]byte, *MetaDataREST, error) {
	meta := &MetaDataREST{}

	// fetch rate limit handler, panic if not found
	rlh, ok := rc.c.rlhs[sd.endpointType]
	if !ok {
		err := fmt.Errorf("rateLimitHander with endpointType %s not found", sd.endpointType)
		rc.logger.WithError(err).Panic("Do")
	}

	// create request
	req, err := rc.createRequest(ctx, sd, p.UrlValues())
	if err != nil {
		return nil, nil, err
	}

	// make request
	resp, err := rc.doRequest(req, sd, rlh, meta)
	if err != nil {
		return nil, nil, err
	}

	// read response body
	data, _ := io.ReadAll(resp.Body)

	// handle status code
	switch resp.StatusCode {
	case int(HTTPStatusOK):
	default:
		fmt.Printf("%s", data)
		rc.logger.WithFields(log.Fields{
			"status_code": resp.StatusCode,
			"body":        fmt.Sprintf("%s", data),
		}).Fatal("Unexpected status code")
	}

	return data, meta, nil
}

// createRequest creates a new http request and implements binance's
// endpoint security protocol.
func (rc *RESTClient) createRequest(
	ctx context.Context,
	sd *serviceDefinition,
	urlValues url.Values,
) (*http.Request, error) {

	uri := url.URL{Scheme: "https", Host: string(sd.endpoint), Path: sd.path, RawQuery: urlValues.Encode()}
	req, err := http.NewRequestWithContext(ctx, sd.method, uri.String(), nil)

	if err != nil {
		rc.logger.WithError(err).Fatal("Error creating request")
		return nil, err
	}

	// handle Security
	switch sd.securityType {
	case SecurityTypeSigned:
		rc.sign(&urlValues)
		req.URL.RawQuery = urlValues.Encode()
		req.Header.Set("X-MBX-APIKEY", rc.c.apiKey)
	case SecurityTypeApiKey:
		req.Header.Set("X-MBX-APIKEY", rc.c.apiKey)
	case SecurityTypeNone:
	}

	return req, nil
}

// sign adds "timestamp", "recvWindow", and "signature" to urlValues.
func (rc *RESTClient) sign(urlValues *url.Values) {
	tsMilli := rc.c.NanotoMilli(rc.c.TSNanoNowServer())
	urlValues.Add("timestamp", strconv.FormatInt(tsMilli, 10))
	urlValues.Add("recvWindow", strconv.Itoa(rc.c.recvWindow))
	queryString := urlValues.Encode()

	mac := hmac.New(sha256.New, []byte(rc.c.apiSecret))
	mac.Write([]byte(queryString))
	signature := mac.Sum(nil)

	urlValues.Add("signature", fmt.Sprintf("%x", signature))
}

// do makes the http request and records timestamps in MetaDataREST.
func (rc *RESTClient) doRequest(
	req *http.Request,
	sd *serviceDefinition,
	rlh *RateLimitHandler,
	meta *MetaDataREST,
) (*http.Response, error) {

	// record timestamps
	meta.T0_Local = rc.c.TSNanoNowLocal()
	meta.T0_Server = rc.c.TSNanoNowServer()
	defer func() {
		meta.T3_Local = rc.c.TSNanoNowLocal()
		meta.T3_Server = rc.c.TSNanoNowServer()
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
	rh := newResponseHeader(resp.Header)
	meta.TRespHeader = rh.TSNano
	rlh.UpdateUsed(rh)

	return resp, nil
}
