// Copyright (c) 2019 Synopsys, Inc.
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
	"fmt"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
)

func (c *Client) CreateApiToken(name, description string, readOnly bool) (location string, token string, err error) {

	tokensUrl := fmt.Sprintf("%s/api/current-user/tokens", c.baseURL)

	tokenRequest := &hubapi.ApiToken {
		Name: name,
		Description:description,
		Scopes: []string{"read"},
	}

	if !readOnly {
		tokenRequest.Scopes = append(tokenRequest.Scopes, "write")
	}

	var tokenResponse hubapi.CreateApiTokenResponse
	location, err = c.HttpPostJSONExpectResult(tokensUrl, tokenRequest, &tokenResponse,"application/json", 201)

	if err != nil {
		return location, "", TraceHubClientError(err)
	}

	if location == "" {
		log.Warn("Did not get a location header back from token creation")
	}

	return location, tokenResponse.Token, err
}

func (c *Client) DeleteApiToken(tokenUrl string) error {
	return c.HttpDelete(tokenUrl, "application/json", 204)
}

func (c *Client) ListApiTokens(options *hubapi.GetListOptions) (*hubapi.ApiTokenList, error) {

	tokensUrl := fmt.Sprintf("%s/api/current-user/tokens", c.baseURL)

	params := ""
	if options != nil {
		params = fmt.Sprintf("?%s", hubapi.ParameterString(options))
	}

	apiTokenListURL := fmt.Sprintf("%s%s", tokensUrl, params)

	var apiTokenList hubapi.ApiTokenList
	err := c.HttpGetJSON(apiTokenListURL, &apiTokenList, 200)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to retrieve api tokens list")
	}

	return &apiTokenList, nil
}

