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
	"fmt"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
)

const (
	remediatingApi = "/remediating"
)

func (c *Client) ListComponents(options *hubapi.GetListOptions) (*hubapi.ComponentList, error) {
	params := ""
	if options != nil {
		params = fmt.Sprintf("?%s", hubapi.ParameterString(options))
	}

	componentURL := fmt.Sprintf("%s/api/components%s", c.baseURL, params)

	var componentList hubapi.ComponentList
	err := c.HttpGetJSON(componentURL, &componentList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve component list: %+v.", err)
		return nil, err
	}

	return &componentList, nil
}

func (c *Client) GetComponent(link hubapi.ResourceLink) (*hubapi.Component, error) {
	var component hubapi.Component
	err := c.HttpGetJSON(link.Href, &component, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a component: %+v.", err)
		return nil, err
	}

	return &component, nil
}

func (c *Client) CreateComponent(componentRequest *hubapi.ComponentRequest) (string, error) {
	componentURL := fmt.Sprintf("%s/api/components", c.baseURL)
	location, err := c.HttpPostJSON(componentURL, componentRequest, "application/json", 201)

	if err != nil {
		return location, err
	}

	if location == "" {
		log.Warnf("Did not get a location header back for component creation")
	}

	return location, err
}

func (c *Client) DeleteComponent(componentURL string) error {
	return c.HttpDelete(componentURL, "application/json", 204)
}

func (c *Client) GetComponentVersion(link *hubapi.ResourceLink) (*hubapi.ComponentVersion, error) {
	var componentVersion hubapi.ComponentVersion
	err := c.HttpGetJSON(link.Href, &componentVersion, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a component: %+v.", err)
		return nil, err
	}

	return &componentVersion, nil
}

func (c *Client) GetComponentVersionFromVariant(componentVariant *hubapi.ComponentVariant) (*hubapi.ComponentVersion, error) {
	return c.GetComponentVersion(&hubapi.ResourceLink{Href: componentVariant.Version})
}

func (c *Client) GetRemediationForComponentVersion(componentVersion *hubapi.ComponentVersion) (*hubapi.ComponentRemediation, error) {
	var componentRemediation hubapi.ComponentRemediation
	err := c.HttpGetJSON(componentVersion.Meta.Href+remediatingApi, &componentRemediation, 200)
	if err != nil {
		log.Errorf("Error trying to retrieve component remediation: %+v.", err)
		return nil, err
	}
	return &componentRemediation, nil
}
