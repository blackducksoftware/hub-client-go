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
	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
)

const apiCurrentUser = "/api/current-user"

func (c *Client) CreateApiToken(name, description string, readOnly bool) (location string, token string, err error) {

	tokensUrl := c.baseURL + apiCurrentUser + "/tokens"

	tokenRequest := &hubapi.ApiToken{
		Name:        name,
		Description: description,
		Scopes:      []string{"read"},
	}

	if !readOnly {
		tokenRequest.Scopes = append(tokenRequest.Scopes, "write")
	}

	var tokenResponse hubapi.CreateApiTokenResponse
	location, err = c.HttpPostJSONExpectResult(tokensUrl, tokenRequest, &tokenResponse, "application/json", 201)

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
	tokensUrl := c.baseURL + apiCurrentUser + "/tokens"
	var apiTokenList hubapi.ApiTokenList

	err := c.GetPage(tokensUrl, options, &apiTokenList)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to retrieve api tokens list")
	}

	return &apiTokenList, nil
}

func (c *Client) GetCurrentUser() (response *hubapi.CurrentUserResponse, err error) {
	currentUserUrl := c.baseURL + apiCurrentUser

	response = &hubapi.CurrentUserResponse{}

	err = c.HttpGetJSON(currentUserUrl, response, 200)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to get current user")
	}

	return response, nil
}
