package hubclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type HubClientDebug uint16

const (
	HubClientDebugTimings HubClientDebug = 1 << iota
	HubClientDebugContent
)

type Client struct {
	httpClient   *http.Client
	baseURL      string
	authToken    string
	useAuthToken bool
	debugFlags   HubClientDebug
}

func NewWithSession(baseURL string, debugFlags HubClientDebug) (*Client, error) {

	jar, err := cookiejar.New(nil) // Look more at this function

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: tr,
	}

	return &Client{
		httpClient:   client,
		baseURL:      baseURL,
		useAuthToken: false,
		debugFlags:   debugFlags,
	}, nil
}

func NewWithToken(baseURL string, authToken string, debugFlags HubClientDebug) (*Client, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	return &Client{
		httpClient:   client,
		baseURL:      baseURL,
		authToken:    authToken,
		useAuthToken: true,
		debugFlags:   debugFlags,
	}, nil
}

func readBytes(readCloser io.ReadCloser) ([]byte, error) {

	defer readCloser.Close()
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(readCloser); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func validateHTTPResponse(resp *http.Response) error {

	// TODO: Check for general success and maybe redirect?
	if resp.StatusCode != 200 {
		return fmt.Errorf("received http status %d when expected 200", resp.StatusCode)
	}

	return nil
}

func (c *Client) processResponse(resp *http.Response, result interface{}) error {

	var bodyBytes []byte
	var err error

	if err := validateHTTPResponse(resp); err != nil {
		fmt.Println("Error validating HTTP Response:")
		fmt.Println(err)
		return err
	}

	if result == nil {
		// Don't have a result to deserialize to, skip it
		return nil
	}

	if bodyBytes, err = readBytes(resp.Body); err != nil {
		fmt.Println("Error reading HTTP Response:")
		fmt.Println(err)
		return err
	}

	if c.debugFlags&HubClientDebugContent != 0 {
		fmt.Printf("START DEBUG: --------------------------------------------------------------------------- \n\n")
		fmt.Printf("TEXT OF RESPONSE: \n %s\n", string(bodyBytes[:]))
		fmt.Printf("END DEBUG: --------------------------------------------------------------------------- \n\n\n\n")
	}

	if err := json.Unmarshal(bodyBytes, result); err != nil {
		fmt.Println("Error parsing HTTP Response:")
		fmt.Println(err)
		fmt.Println("\n\n--------------------")
		fmt.Println("Response:")
		fmt.Println("--------------------\n\n")
		fmt.Println(resp)
		return err
	}

	return nil
}

func (c *Client) httpGetJSON(url string, result interface{}, expectedStatusCode int) error {

	var resp *http.Response
	var err error

	if c.debugFlags&HubClientDebugTimings != 0 {
		fmt.Printf("DEBUG HTTP STARTING REQUEST: %s\n", url)
	}

	httpStart := time.Now()
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error making http request:")
		fmt.Println(err)
		return err
	}

	if c.useAuthToken {
		c.addAuthToken(req)
	}

	if resp, err = c.httpClient.Do(req); err != nil {
		fmt.Println("Error getting HTTP Response:")
		fmt.Println(err)
		return err
	}

	httpElapsed := time.Since(httpStart)

	if c.debugFlags&HubClientDebugTimings != 0 {
		fmt.Printf("DEBUG HTTP ELAPSED TIME: %d ms.   -- Request: %s\n", (httpElapsed / 1000 / 1000), url)
	}

	if resp.StatusCode != expectedStatusCode { // Should this be a list at some point?
		fmt.Printf("Got a %d response instead of a 200.\n", resp.StatusCode)
		return fmt.Errorf("got a %d response instead of a 200", resp.StatusCode)
	}

	return c.processResponse(resp, result)
}

func (c *Client) addAuthToken(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
}
