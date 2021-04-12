// Copyright 2021 Synopsys, Inc.
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
	"strconv"
)

const (
	apiDeveloperScans     = "/api/developer-scans"
	headerBdMode          = "X-BD-MODE"
	headerBdDocumentCount = "X-BD-DOCUMENT-COUNT"
)

func (c *Client) StartRapidScan(bdioHeaderContent string) (error, string) {
	rapidScansURL := c.baseURL + apiDeveloperScans
	bdioUploadEndpoint, err := c.HttpPostJSON(rapidScansURL, bdioHeaderContent, hubapi.ContentTypeRapidScan, 201)

	if err != nil {
		log.Errorf("Error kicking off a rapid scan.", err)
		return err, ""
	}

	return nil, bdioUploadEndpoint
}

func (c *Client) UploadBdioFiles(bdioUploadEndpoint string, bdioContents []string) error {
	c.AddHeaderValue(headerBdMode, "append")
	c.AddHeaderValue(headerBdDocumentCount, strconv.Itoa(len(bdioContents)))

	for _, bdioContent := range bdioContents {
		err := c.HttpPutJSON(bdioUploadEndpoint, bdioContent, hubapi.ContentTypeRapidScan, 202)
		if err != nil {
			log.Errorf("Error uploading bdio files.", err)
			return err
		}
	}

	c.SetHeaderValue(headerBdMode, "finish")
	err := c.HttpPutJSON(bdioUploadEndpoint, "", hubapi.ContentTypeRapidScan, 202)
	if err != nil {
		log.Errorf("Error uploading bdio files.", err)
		return err
	}

	c.DeleteHeaderValue(headerBdMode)
	c.DeleteHeaderValue(headerBdDocumentCount)

	return nil
}
