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
	"testing"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetComponentVersion(t *testing.T) {
	client := createTestClient(t)
	assert.NotNil(t, client, "unable to get client")

	option := "maven:org.apache.commons:commons-collections4:4.0"
	listOptions := &hubapi.GetListOptions{Q: &option}

	componentList, err := client.ListComponents(listOptions)
	assert.NoError(t, err)
	assert.True(t, len(componentList.Items) > 0, "Expected at least one componentlist item")

	item := componentList.Items[0]

	link := hubapi.ResourceLink{Href: item.Version}

	componentVersion, err := client.GetComponentVersion(link)
	assert.NoError(t, err)
	assert.NotEmpty(t, componentVersion.Type)
	assert.NotEmpty(t, componentVersion.ApprovalStatus)
	assert.NotEmpty(t, componentVersion.Source)
	assert.NotEmpty(t, componentVersion.VersionName)
	assert.NotNil(t, componentVersion.ReleasedOn)
	assert.NotNil(t, componentVersion.License)
	assert.NotNil(t, componentVersion.Meta)
	assert.True(t, len(componentVersion.Meta.Links) > 0)
	assert.NotEmpty(t, componentVersion.Meta.Href)
	assert.True(t, len(componentVersion.Meta.Allow) > 0)
}

func TestClient_GetComponentVersionFromVariant(t *testing.T) {

	client := createTestClient(t)
	assert.NotNil(t, client, "unable to get client")

	option := "maven:org.apache.commons:commons-collections4:4.0"
	listOptions := &hubapi.GetListOptions{Q: &option}

	componentList, err := client.ListComponents(listOptions)
	assert.NoError(t, err)
	assert.True(t, len(componentList.Items) > 0, "Expected at least one componentlist item")

	componentVariant := componentList.Items[0]

	componentVersion, err := client.GetComponentVersion(hubapi.ResourceLink{Href: componentVariant.Version})
	assert.NoError(t, err)
	assert.NotEmpty(t, componentVersion.Type)
	assert.NotEmpty(t, componentVersion.ApprovalStatus)
	assert.NotEmpty(t, componentVersion.Source)
	assert.NotEmpty(t, componentVersion.VersionName)
	assert.NotNil(t, componentVersion.ReleasedOn)
	assert.NotNil(t, componentVersion.License)
	assert.NotNil(t, componentVersion.Meta)
	assert.True(t, len(componentVersion.Meta.Links) > 0)
	assert.NotEmpty(t, componentVersion.Meta.Href)
	assert.True(t, len(componentVersion.Meta.Allow) > 0)
}
