// Copyright 2019 Synopsys, Inc.
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type BearerTokenResponse struct {
	BearerToken           string `json:"bearerToken"`
	ExpiresInMilliseconds int64  `json:"expiresInMilliseconds"`
}

func NewWithApiToken(baseURL string, apiToken string, debugFlags HubClientDebug, timeout time.Duration) (*Client, error) {
	client := createHttpClient(timeout)
	return NewWithApiTokenAndClient(baseURL, apiToken, debugFlags, client)
}

func NewWithApiTokenAndClient(baseURL string, apiToken string, debugFlags HubClientDebug, client *http.Client) (*Client, error) {
	if client == nil {
		client = createHttpClient(time.Minute)
	}

	url := baseURL + "/api/tokens/authenticate"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error creating API token request")
	}

	tokenValue := fmt.Sprintf("token %s", apiToken)

	req.Header.Add(HeaderNameAuthorization, tokenValue)

	resp, err := client.Do(req)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error logging in via API Token")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, HubClientErrorfWithStatusCode(resp.StatusCode, "got a %d response instead of a %d", resp.StatusCode, http.StatusOK)
	}

	csrf := resp.Header.Get(HeaderNameCsrfToken)
	if csrf == "" {
		return nil, newHubClientError(nil, resp, "CSRF token not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error reading body")
	}

	var bearerTokenResponse BearerTokenResponse

	err = json.Unmarshal(body, &bearerTokenResponse)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error decoding JSON")
	}

	if bearerTokenResponse.BearerToken == "" {
		return nil, newHubClientError(body, resp, "bearer token not found")
	}

	log.Debug("Logged in with auth token successfully")

	return &Client{
		httpClient:      client,
		baseURL:         baseURL,
		authToken:       bearerTokenResponse.BearerToken,
		useAuthToken:    true,
		csrfToken:       csrf,
		haveCsrfToken:   true,
		debugFlags:      debugFlags,
		headerOverrides: http.Header{},
	}, nil
}
