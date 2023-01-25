// Copyright 2018 Synopsys, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hubclient

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client will need to support CSRF tokens for session-based auth for Hub 4.1.x (or was it 4.0?)
type Client struct {
	httpClient    *http.Client
	baseURL       string
	authToken     string
	useAuthToken  bool
	haveCsrfToken bool
	csrfToken     string
	debugFlags    HubClientDebug
	// Unix time in seconds at which the authToken expires
	authTokenExpiryInUnixSec int64
	userAgent                string
}

func NewWithSession(baseURL string, debugFlags HubClientDebug, timeout time.Duration) (*Client, error) {
	client := createHttpClient(timeout)
	return NewWithClient(baseURL, debugFlags, client)
}

func NewWithClient(baseURL string, debugFlags HubClientDebug, httpClient *http.Client) (*Client, error) {

	if httpClient == nil {
		httpClient = createHttpClient(time.Minute)
	}

	if httpClient.Jar == nil {
		jar, err := cookiejar.New(nil) // Look more at this function

		if err != nil {
			return nil, AnnotateHubClientError(err, "unable to instantiate cookie jar")
		}

		httpClient.Jar = jar
	}

	return &Client{
		httpClient:   httpClient,
		baseURL:      baseURL,
		useAuthToken: false,
		debugFlags:   debugFlags,
	}, nil
}

func NewWithToken(baseURL string, authToken string, debugFlags HubClientDebug, timeout time.Duration) (*Client, error) {
	client := createHttpClient(timeout)
	return NewWithTokenAndClient(baseURL, authToken, debugFlags, client)
}

func NewWithTokenAndClient(baseURL string, authToken string, debugFlags HubClientDebug, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = createHttpClient(time.Minute)
	}

	return &Client{
		httpClient:   httpClient,
		baseURL:      baseURL,
		authToken:    authToken,
		useAuthToken: true,
		debugFlags:   debugFlags,
	}, nil
}

func createHttpClient(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	return client
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func readBytes(readCloser io.ReadCloser, expected int64) ([]byte, error) {

	defer readCloser.Close()

	buf := bytes.Buffer{}
	if expected > 0 && int(expected) > 0 {
		buf.Grow(int(expected))
	}

	if _, err := buf.ReadFrom(readCloser); err != nil {
		return nil, TraceHubClientError(err)
	}

	return buf.Bytes(), nil
}

func validateHTTPResponse(resp *http.Response, debugFlags HubClientDebug, expectedStatusCodes ...int) error {
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	statusCodeExpected := false

	for _, expectedCode := range expectedStatusCodes {
		if statusCode == expectedCode {
			statusCodeExpected = true
			break
		}
	}

	if !statusCodeExpected {
		body := readResponseBody(resp, debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("got a %d response instead of %d", statusCode, expectedStatusCodes), nil)
	}

	return nil
}

func (c *Client) processResponse(resp *http.Response, result interface{}, expectedStatusCodes ...int) error {

	if err := validateHTTPResponse(resp, c.debugFlags, expectedStatusCodes...); err != nil {
		return AnnotateHubClientError(err, "error validating HTTP Response")
	}

	if result == nil {
		// Don't have a result to deserialize to, skip it
		return nil
	}

	bodyBytes, err := readBytes(resp.Body, resp.ContentLength)
	if err != nil {
		return newHubClientError(bodyBytes, resp, fmt.Sprintf("error reading HTTP Response: %+v", err), err)
	}

	debugReportBytes(bodyBytes, c.debugFlags)

	if err := json.Unmarshal(bodyBytes, result); err != nil {
		return newHubClientError(bodyBytes, resp, fmt.Sprintf("error parsing HTTP Response: %+v", err), err)
	}

	return nil
}

func (c *Client) HttpGetString(url string, result *string, expectedStatusCode []int, mimetypes ...string) (error, int) {
	err, response := c.httpGet(url, nil, expectedStatusCode, mimetypes...)

	statusCode := getResponseStatus(response)
	body := readResponseBody(response, c.debugFlags)
	*result = string(body)

	return err, statusCode
}

func (c *Client) HttpGetJSON(url string, result interface{}, expectedStatusCode int, mimetypes ...string) error {
	err, _ := c.httpGet(url, result, []int{expectedStatusCode}, mimetypes...)
	return err
}

func (c *Client) httpGet(url string, result interface{}, expectedStatusCodes []int, mimetypes ...string) (error, *http.Response) {

	var resp *http.Response

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING GET REQUEST: %s", url)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http get request for %s: %+v", url, err), err), resp
	}

	c.applyHeaderValues(req, nil)

	if len(mimetypes) > 0 {
		for _, mimetype := range mimetypes {
			if mimetype != "" {
				req.Header.Add(HeaderNameAccept, mimetype)
			}
		}
	} else if result != nil {
		if bdJsonType := hubapi.GetMimeType(result); bdJsonType != "" {
			req.Header.Add(HeaderNameAccept, bdJsonType)
		}
	}

	httpStart := time.Now()
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from %s: %+v", url, err), err), resp
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP GET ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	err = c.processResponse(resp, result, expectedStatusCodes...)
	return AnnotateHubClientErrorf(err, "unable to process response from GET to %s", url), resp
}

func (c *Client) HttpPutString(url string, data string, contentType string, expectedStatusCode int) error {
	return c.HttpPutStringWithHeader(url, data, contentType, expectedStatusCode, nil)
}

func (c *Client) HttpPutStringWithHeader(url string, data string, contentType string, expectedStatusCode int, header http.Header) error {
	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s STRING REQUEST: %s", http.MethodPut, url)
	}

	reader := strings.NewReader(data)
	return c.putRequestWithHeader(url, reader, contentType, expectedStatusCode, header)
}

func (c *Client) HttpPutJSON(url string, data interface{}, contentType string, expectedStatusCode int) error {
	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s JSON REQUEST: %s", http.MethodPut, url)
	}

	// Encode json
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&data); err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err), err)
	}
	reader := &buf

	return c.putRequest(url, reader, contentType, expectedStatusCode)
}

func (c *Client) putRequest(url string, reader io.Reader, contentType string, expectedStatusCode int) error {
	return c.putRequestWithHeader(url, reader, contentType, expectedStatusCode, nil)
}

func (c *Client) putRequestWithHeader(url string, reader io.Reader, contentType string, expectedStatusCode int, header http.Header) error {
	req, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http put request for %s: %+v", url, err), err)
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req, header)
	log.Debugf("PUT Request: %+v.", maskAuthHeader(req))

	httpStart := time.Now()
	var resp *http.Response
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from PUT to %s: %+v", url, err), err)
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP PUT ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	return AnnotateHubClientErrorf(c.processResponse(resp, nil, expectedStatusCode), "unable to process response from PUT to %s", url) // TODO: Maybe need a response too?
}

func (c *Client) HttpPostString(url string, data string, contentType string, expectedStatusCode int) (string, error) {
	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s STRING REQUEST: %s", http.MethodPost, url)
	}

	reader := strings.NewReader(data)
	return c.postRequest(url, reader, contentType, expectedStatusCode)
}

func (c *Client) HttpPostJSON(url string, data interface{}, contentType string, expectedStatusCode int) (string, error) {
	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s JSON REQUEST: %s", http.MethodPost, url)
	}

	// Encode json
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&data); err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err), err)
	}
	reader := &buf

	return c.postRequest(url, reader, contentType, expectedStatusCode)
}

func (c *Client) HttpPostFile(url string, filePath string, contentType string) (string, int, error) {

	var resp *http.Response
	var err error

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING POST REQUEST: %s", url)
	}

	file, err := os.Open(filePath)

	if err != nil {
		log.Errorf("Error opening file: %+v", err)
		return "", -1, err
	}

	httpStart := time.Now()
	req, err := http.NewRequest(http.MethodPost, url, file)
	req.Header.Set(HeaderNameContentType, contentType)

	c.setAuthHeaders(req)

	if err != nil {
		log.Errorf("Error making http post request: %+v.", err)
		return "", -1, err
	}

	log.Debugf("POST Request: %+v.", maskAuthHeader(req))

	if resp, err = c.httpClient.Do(req); err != nil {
		log.Errorf("Error getting HTTP Response: %+v.", err)
		readResponseBody(resp, c.debugFlags)
		return "", resp.StatusCode, err
	}

	httpElapsed := time.Since(httpStart)

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP POST ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}
	return resp.Header.Get("Location"), resp.StatusCode, nil
}

func (c *Client) postRequest(url string, reader io.Reader, contentType string, expectedStatusCode int) (string, error) {
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error creating http post request for %s: %+v", url, err), err)
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req, nil)

	log.Debugf("POST Request: %+v.", maskAuthHeader(req))

	httpStart := time.Now()

	var resp *http.Response
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return "", newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from POST to %s: %+v", url, err), err)
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP POST ELAPSED TIME: %d ms. -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	if err := c.processResponse(resp, nil, expectedStatusCode); err != nil {
		return "", AnnotateHubClientErrorf(err, "unable to process response from POST to %s", url)
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) HttpPostJSONExpectResult(url string, data interface{}, result interface{}, contentType string, expectedStatusCode int) (string, error) {

	var resp *http.Response

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING POST REQUEST: %s", url)
	}

	// Encode json
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(&data); err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err), err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error creating http post request for %s: %+v", url, err), err)
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req, nil)

	log.Debugf("POST Request: %+v.", maskAuthHeader(req))

	httpStart := time.Now()
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return "", newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from POST to %s: %+v", url, err), err)
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP POST ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	if err := c.processResponse(resp, result, expectedStatusCode); err != nil {
		return "", AnnotateHubClientErrorf(err, "unable to process response from POST to %s", url)
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) HttpDelete(url string, contentType string, expectedStatusCode int) error {

	var resp *http.Response
	var err error

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING DELETE REQUEST: %s", url)
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http delete request for %s: %+v", url, err), err)
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req, nil)

	log.Debugf("DELETE Request: %+v.", maskAuthHeader(req))

	httpStart := time.Now()

	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from DELETE to %s: %+v", url, err), err)
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP DELETE ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	return AnnotateHubClientErrorf(c.processResponse(resp, nil, expectedStatusCode), "unable to process response from DELETE to %s", url)
}

// Applies authentication headers (see setAuthHeaders) and also applies any custom header values to the provided request
func (c *Client) applyHeaderValues(request *http.Request, header http.Header) {
	for key, values := range header {
		// remove any old values in the provided request header
		request.Header.Del(key)
		for _, value := range values {
			// add the new values
			request.Header.Add(key, value)
		}

	}

	if c.userAgent != "" {
		request.Header.Add(HeaderNameUserAgent, c.userAgent)
	}

	c.setAuthHeaders(request)
}

func (c *Client) setAuthHeaders(request *http.Request) {

	if c.useAuthToken {
		request.Header.Set(HeaderNameAuthorization, fmt.Sprintf("Bearer %s", c.authToken))
	}

	if c.haveCsrfToken {
		request.Header.Set(HeaderNameCsrfToken, c.csrfToken)
	}
}

func (c *Client) Count(link string) (int, error) {
	var list hubapi.ItemsListBase
	err := c.HttpGetJSON(link+"?offset=0&limit=1", &list, 200)

	if err != nil {
		return 0, AnnotateHubClientError(err, "Error trying to retrieve count")
	}

	return list.TotalCount, nil
}

func (c *Client) ForEachPage(link string, listOptions *hubapi.GetListOptions, list interface{}, pageFunc func() error) (err error) {
	listOptions = hubapi.EnsureLimits(listOptions)

	for totalCount := 1; err == nil && *listOptions.Offset < totalCount; listOptions.NextPage() {
		resetList(list)
		err = c.GetPage(link, listOptions, list)
		if err == nil {
			err = pageFunc()
		}

		if t, ok := list.(hubapi.TotalCountable); ok {
			totalCount = t.Total()
		}
	}

	return err
}

func resetList(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (c *Client) GetPage(link string, options *hubapi.GetListOptions, list interface{}) error {
	listUrl := link + hubapi.ParameterString(options)

	err := c.HttpGetJSON(listUrl, list, 200)

	if err != nil {
		return AnnotateHubClientError(err, fmt.Sprintf("Error trying to retrieve list %T", list))
	}

	return nil
}

func getResponseStatus(resp *http.Response) int {
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	return statusCode
}

func readResponseBody(resp *http.Response, debugFlags HubClientDebug) (bodyBytes []byte) {

	if resp == nil {
		log.Errorf("Empty HTTP Response")
		return nil
	}

	var err error
	if bodyBytes, err = readBytes(resp.Body, resp.ContentLength); err != nil {
		log.Errorf("Error reading HTTP Response: %+v.", err)
	}

	debugReportBytes(bodyBytes, debugFlags)

	return bodyBytes
}

func newHubClientError(respBody []byte, resp *http.Response, message string, underlyingErr error) *HubClientError {
	var hre HubResponseError

	// Do not try to read the body of the response or response itself
	var statusCode int
	if resp != nil {
		statusCode = resp.StatusCode
	}

	var err error
	if underlyingErr == nil {
		err = errors.New(message)
	} else {
		err = &errorWithMessage{
			err:     underlyingErr,
			message: message,
		}
	}

	hce := &HubClientError{err, statusCode, hre}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &hre); err != nil {
			hce = AnnotateHubClientError(hce, fmt.Sprintf("error unmarshaling HTTP response body: %+v", err)).(*HubClientError)
		}
		hce.HubError = hre
	}

	return hce
}

func (c *Client) SetBearerToken(token string) {
	c.authToken = token
	c.useAuthToken = true
}

// GetAuthTokenExpiryTime returns the unix time in seconds at which the cached auth token is set to expire
// This func returns -1 if the Client is not configured to use auth token
func (c *Client) GetAuthTokenExpiryTime() int64 {
	if c == nil || !c.useAuthToken || c.authToken == "" {
		return -1
	}
	return c.authTokenExpiryInUnixSec
}

// Sets the User-Agent header value to be used in all http/https requests made by the client
func (c *Client) SetUserAgent(agent string) {
	c.userAgent = agent
}

func maskAuthHeader(req *http.Request) *http.Request {
	// Deep clone the request
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Debugf("error cloning request: %v", err)
		return nil
	}
	var reqClone *http.Request
	if reqClone, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(reqBytes))); err != nil {
		log.Debugf("error cloning request: %v", err)
		return nil
	}
	if val := reqClone.Header.Get(HeaderNameAuthorization); val != "" {
		reqClone.Header.Set(HeaderNameAuthorization, fmt.Sprintf("Bearer %s", "<NOT_SHOWN>"))
	}
	return reqClone
}
