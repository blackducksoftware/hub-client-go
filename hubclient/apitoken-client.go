package hubclient

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type BearerTokenResponse struct {
	BearerToken           string `json:"bearerToken"`
	ExpiresInMilliseconds int64  `json:"expiresInMilliseconds"`
}

func NewWithApiToken(baseURL string, apiToken string, debugFlags HubClientDebug, timeout time.Duration) (*Client, error) {
	url := fmt.Sprintf("%s/api/tokens/authenticate", baseURL)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error creating API token request")
	}

	tokenValue := fmt.Sprintf("token %s", apiToken)

	req.Header.Add(HeaderNameAuthorization, tokenValue)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, AnnotateHubClientError(err, "error logging in via API Token")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, HubClientErrorf("got a %d response instead of a %d", resp.StatusCode, http.StatusOK)
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
		httpClient:    client,
		baseURL:       baseURL,
		authToken:     bearerTokenResponse.BearerToken,
		useAuthToken:  true,
		csrfToken:     csrf,
		haveCsrfToken: true,
		debugFlags:    debugFlags,
	}, nil

}

func init() {
	log.SetLevel(log.DebugLevel)
}
