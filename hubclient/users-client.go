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
	"github.com/blackducksoftware/hub-client-go/hubapi"
)

const apiUsers = "/api/users"

// TODO: This API should also be returning a location
func (c *Client) CreateUser(userRequest *hubapi.UserRequest) (*hubapi.User, error) {
	usersURL := c.baseURL + apiUsers

	var result hubapi.User
	_, err := c.HttpPostJSONExpectResult(usersURL, userRequest, &result, "application/json", 201)

	if err != nil {
		return nil, TraceHubClientError(err)
	}

	// TODO: Warn once user creation returns a location.
	// if location == "" {
	// 	log.Warnf("Did not get a location header back for user creation")
	// }

	return &result, err
}

func (c *Client) ListUsers(options *hubapi.GetListOptions) (*hubapi.UserList, error) {
	usersURL := c.baseURL + "/api/users"

	var userList hubapi.UserList
	err := c.GetPage(usersURL, options, &userList)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to retrieve user list")
	}

	return &userList, nil
}

func (c *Client) GetUser(link hubapi.ResourceLink) (*hubapi.User, error) {
	var user hubapi.User
	err := c.HttpGetJSON(link.Href, &user, 200)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to retrieve a user")
	}

	return &user, nil
}
