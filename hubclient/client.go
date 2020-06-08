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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"time"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client will need to support CSRF tokens for session-based auth for Hub 4.1.x (or was it 4.0?)
type Client struct {
	httpClient      *http.Client
	baseURL         string
	authToken       string
	useAuthToken    bool
	haveCsrfToken   bool
	csrfToken       string
	debugFlags      HubClientDebug
	headerOverrides http.Header
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
		httpClient:      httpClient,
		baseURL:         baseURL,
		useAuthToken:    false,
		debugFlags:      debugFlags,
		headerOverrides: http.Header{},
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
		httpClient:      httpClient,
		baseURL:         baseURL,
		authToken:       authToken,
		useAuthToken:    true,
		debugFlags:      debugFlags,
		headerOverrides: http.Header{},
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

func validateHTTPResponse(resp *http.Response, expectedStatusCode int, debugFlags HubClientDebug) error {
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	if statusCode != expectedStatusCode { // Should this be a list at some point?
		body := readResponseBody(resp, debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("got a %d response instead of a %d", statusCode, expectedStatusCode))
	}

	return nil
}

func (c *Client) processResponse(resp *http.Response, result interface{}, expectedStatusCode int) error {

	if err := validateHTTPResponse(resp, expectedStatusCode, c.debugFlags); err != nil {
		return AnnotateHubClientError(err, "error validating HTTP Response")
	}

	if result == nil {
		// Don't have a result to deserialize to, skip it
		return nil
	}

	bodyBytes, err := readBytes(resp.Body, resp.ContentLength)
	if err != nil {
		return newHubClientError(bodyBytes, resp, fmt.Sprintf("error reading HTTP Response: %+v", err))
	}

	debugReportBytes(bodyBytes, c.debugFlags)

	if err := json.Unmarshal(bodyBytes, result); err != nil {
		return newHubClientError(bodyBytes, resp, fmt.Sprintf("error parsing HTTP Response: %+v", err))
	}

	return nil
}

func (c *Client) HttpGetJSON(url string, result interface{}, expectedStatusCode int, mimetypes ...string) error {

	var resp *http.Response

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING GET REQUEST: %s", url)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http get request for %s: %+v", url, err))
	}

	c.applyHeaderValues(req)

	if len(mimetypes) > 0 {
		for _, mimetype := range mimetypes {
			if mimetype != "" {
				req.Header.Add(HeaderNameAccept, mimetype)
			}
		}
	} else if bdJsonType := hubapi.GetMimeType(result); bdJsonType != "" {
		req.Header.Add(HeaderNameAccept, bdJsonType)
	}

	httpStart := time.Now()
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from %s: %+v", url, err))
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP GET ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	err = c.processResponse(resp, result, expectedStatusCode)
	return AnnotateHubClientErrorf(err, "unable to process response from GET to %s", url)
}

func (c *Client) HttpPutJSON(url string, data interface{}, contentType string, expectedStatusCode int) error {
	var resp *http.Response

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s REQUEST: %s", http.MethodPut, url)
	}

	// Encode json
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(&data); err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err))
	}

	req, err := http.NewRequest(http.MethodPut, url, &buf)
	if err != nil {
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http put request for %s: %+v", url, err))
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req)
	log.Debugf("PUT Request: %+v.", req)

	httpStart := time.Now()
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from PUT to %s: %+v", url, err))
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP PUT ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	return AnnotateHubClientErrorf(c.processResponse(resp, nil, expectedStatusCode), "unable to process response from PUT to %s", url) // TODO: Maybe need a response too?
}

func (c *Client) HttpPostJSON(url string, data interface{}, contentType string, expectedStatusCode int) (string, error) {

	var resp *http.Response

	if c.debugFlags&HubClientDebugTimings != 0 {
		log.Debugf("DEBUG HTTP STARTING %s REQUEST: %s", http.MethodPost, url)
	}

	// Encode json
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(&data); err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err))
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error creating http post request for %s: %+v", url, err))
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req)

	log.Debugf("POST Request: %+v.", req)

	httpStart := time.Now()

	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return "", newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from POST to %s: %+v", url, err))
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
		return "", newHubClientError(nil, nil, fmt.Sprintf("error encoding json: %+v", err))
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", newHubClientError(nil, nil, fmt.Sprintf("error creating http post request for %s: %+v", url, err))
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req)

	log.Debugf("POST Request: %+v.", req)

	httpStart := time.Now()
	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return "", newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from POST to %s: %+v", url, err))
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
		return newHubClientError(nil, nil, fmt.Sprintf("error creating http delete request for %s: %+v", url, err))
	}

	req.Header.Set(HeaderNameContentType, contentType)

	c.applyHeaderValues(req)

	log.Debugf("DELETE Request: %+v.", req)

	httpStart := time.Now()

	if resp, err = c.httpClient.Do(req); err != nil {
		body := readResponseBody(resp, c.debugFlags)
		return newHubClientError(body, resp, fmt.Sprintf("error getting HTTP Response from DELETE to %s: %+v", url, err))
	}

	if c.debugFlags&HubClientDebugTimings != 0 {
		httpElapsed := time.Since(httpStart)
		log.Debugf("DEBUG HTTP DELETE ELAPSED TIME: %d ms.   -- Request: %s", (httpElapsed / 1000 / 1000), url)
	}

	return AnnotateHubClientErrorf(c.processResponse(resp, nil, expectedStatusCode), "unable to process response from DELETE to %s", url)
}

// Applies authentication headers (see setAuthHeaders) and also applies any custom header values to the provided request
func (c *Client) applyHeaderValues(request *http.Request) {
	for key, values := range c.headerOverrides {
		// remove any old values in the provided request header
		request.Header.Del(key)
		for _, value := range values {
			// add the new values
			request.Header.Add(key, value)
		}

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

// Sets header override values for the provided key for any subsequent requests by this client
func (c *Client) SetHeaderValue(key string, value string) {
	c.headerOverrides.Set(key, value)
}

// Clears a previously set header override value for this client
func (c *Client) DeleteHeaderValue(key string) {
	c.headerOverrides.Del(key)
}

// Adds header override values for the provided key for any subsequent requests by this client
func (c *Client) AddHeaderValue(key string, value string) {
	c.headerOverrides.Add(key, value)
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

func newHubClientError(respBody []byte, resp *http.Response, message string) *HubClientError {
	var hre HubResponseError

	// Do not try to read the body of the response or response itself
	var statusCode int
	if resp != nil {
		statusCode = resp.StatusCode
	}

	hce := &HubClientError{errors.New(message), statusCode, hre}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &hre); err != nil {
			hce = AnnotateHubClientError(hce, fmt.Sprintf("error unmarshaling HTTP response body: %+v", err)).(*HubClientError)
		}
		hce.HubError = hre
	}

	return hce
}
