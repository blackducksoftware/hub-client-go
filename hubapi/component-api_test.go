package hubapi_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	"github.com/stretchr/testify/assert"
)

func TestReadList(t *testing.T) {

	content, err := ioutil.ReadFile("testdata/componentList.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	var items hubapi.ComponentList

	err = json.Unmarshal([]byte(content), &items)
	assert.NoError(t, err)
	assert.True(t, items.TotalCount > 0)
	assert.Len(t, items.Items, 1)

	component := items.Items[0]

	assert.NotEmpty(t, component.ComponentName)
	assert.NotEmpty(t, component.Component)
	assert.NotEmpty(t, component.Version)
}

func TestReadComponentVersionList(t *testing.T) {

	content, err := ioutil.ReadFile("testdata/componentVersionList.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	var items hubapi.ComponentVersionList

	err = json.Unmarshal([]byte(content), &items)
	assert.NoError(t, err)
	assert.True(t, items.TotalCount > 0)
	assert.Len(t, items.Items, 10)

	component := items.Items[0]

	assert.NotEmpty(t, component.VersionName)
	assert.NotEmpty(t, component.Source)
	assert.NotEmpty(t, component.Type)
}

func TestReadComponentVersionOriginList(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/componentVersionOriginList.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	var items hubapi.ComponentVersionOriginList

	err = json.Unmarshal([]byte(content), &items)
	assert.NoError(t, err)
	assert.Len(t, items.Items, 10)
	assert.True(t, items.TotalCount > 0)

	item := items.Items[1]

	assert.NotEmpty(t, item.VersionName)
	assert.NotEmpty(t, item.Origin)
	assert.NotEmpty(t, item.OriginID)

	assert.True(t, len(item.Meta.Links) > 5)
	resLink, err := item.Meta.FindLinkByRel("version")
	assert.NoError(t, err)
	assert.NotNil(t, resLink)

	resLink, err = item.Meta.FindLinkByRel("vulnerabilities")
	assert.NoError(t, err)
	assert.NotNil(t, resLink)
}
