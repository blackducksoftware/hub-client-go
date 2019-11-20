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
	log "github.com/sirupsen/logrus"
)

func (c *Client) ListProjectVersionComponents(link hubapi.ResourceLink) (*hubapi.BomComponentList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList hubapi.BomComponentList
	err := c.GetPage(link.Href, nil, &bomList)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error while trying to get Project Version Component list")
	}

	return &bomList, nil
}

// TODO: Should this be used?
func (c *Client) ListProjectVersionVulnerableComponents(link hubapi.ResourceLink) (*hubapi.BomVulnerableComponentList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList hubapi.BomVulnerableComponentList
	err := c.GetPage(link.Href, nil, &bomList)

	if err != nil {
		return nil, AnnotateHubClientError(err, "Error trying to retrieve vulnerable components list")
	}

	return &bomList, nil
}

func (c *Client) CountProjectVersionVulnerableComponents(link hubapi.ResourceLink) (int, error) {
	return c.Count(link.Href)
}

func (c *Client) ListAllProjectVersionVulnerableComponents(link hubapi.ResourceLink) ([]hubapi.BomVulnerableComponent, error) {

	totalCount, err := c.Count(link.Href)

	if err != nil {
		log.Errorf("Error trying to retrieve vulnerable components list: %+v.", err)
		return nil, AnnotateHubClientError(err, "Error trying to retrieve project's vulnerable components")
	}

	result := make([]hubapi.BomVulnerableComponent, 0, totalCount)

	err = c.ForAllPages(nil, func(options *hubapi.GetListOptions) (int, error) {
		var bomPage hubapi.BomVulnerableComponentList
		err = c.GetPage(link.Href, options, &bomPage)
		if err != nil {
			log.Errorf("Error trying to retrieve vulnerable components list: %+v.", err)
			return 0, err
		}

		result = append(result, bomPage.Items...)

		return bomPage.TotalCount, nil
	})

	return result, nil
}
