package hubclient

import (
	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (c *Client) StartSnippetScan(content io.Reader) (*hubapi.SnippetMatchResponse, error) {
	snippetMatchingURL := hubapi.BuildUrl(c.baseURL, hubapi.SnippetMatchingApi)
	var result hubapi.SnippetMatchResponse
	_, err := c.HttpPostRawJSONExpectResult(snippetMatchingURL, content, &result, "text/plain", http.StatusOK)
	if err != nil {
		log.Error("Error kicking off a snippet scan.", err)
		return nil, TraceHubClientError(err)
	}

	return &result, nil
}
