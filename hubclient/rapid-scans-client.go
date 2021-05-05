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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	apiDeveloperScans     = "/api/developer-scans"
	apiFullResults        = "/full-result"
	headerBdMode          = "X-BD-MODE"
	headerBdDocumentCount = "X-BD-DOCUMENT-COUNT"
	bdModeAppend          = "append"
	bdModeFinish          = "finish"
	OkStatusCode          = 200
	CreatedStatusCode     = 201
	AcceptedStatusCode    = 202
)

func (c *Client) StartRapidScan(bdioHeaderContent string) (error, string) {
	rapidScansURL := c.baseURL + apiDeveloperScans
	bdioUploadEndpoint, err := c.HttpPostString(rapidScansURL, bdioHeaderContent, hubapi.ContentTypeRapidScanRequest, CreatedStatusCode)

	if err != nil {
		log.Errorf("Error kicking off a rapid scan.", err)
		return err, ""
	}

	return nil, bdioUploadEndpoint
}

func (c *Client) UploadBdioFiles(bdioUploadEndpoint string, bdioContents []string) error {
	c.AddHeaderValue(headerBdMode, bdModeAppend)
	c.AddHeaderValue(headerBdDocumentCount, strconv.Itoa(len(bdioContents)))

	for _, bdioContent := range bdioContents {
		err := c.HttpPutString(bdioUploadEndpoint, bdioContent, hubapi.ContentTypeRapidScanRequest, AcceptedStatusCode)
		if err != nil {
			log.Errorf("Error uploading bdio files.", err)
			return err
		}
	}

	c.SetHeaderValue(headerBdMode, bdModeFinish)
	err := c.HttpPutString(bdioUploadEndpoint, "", hubapi.ContentTypeRapidScanRequest, AcceptedStatusCode)
	if err != nil {
		log.Errorf("Error uploading bdio files.", err)
		return err
	}

	c.DeleteHeaderValue(headerBdMode)
	c.DeleteHeaderValue(headerBdDocumentCount)

	return nil
}

func (c *Client) PollRapidScanResults(rapidScanEndpoint string, intervalInSeconds int, timeoutInSeconds int) (error, *hubapi.RapidScanResult) {
	url := rapidScanEndpoint + apiFullResults
	interval := time.Duration(intervalInSeconds) * time.Second
	timeout := time.Duration(timeoutInSeconds) * time.Second
	ticker := time.NewTicker(interval)
	timeoutTimer := time.NewTimer(timeout)

	var result *hubapi.RapidScanResult
	var body string

	for {
		select {
		case <-timeoutTimer.C:
			ticker.Stop()
			return errors.New(fmt.Sprintf("The polling for rapid scan result timed out: %s", rapidScanEndpoint)), result
		case <-ticker.C:
			err := c.HttpGetString(url, &body, OkStatusCode, hubapi.ContentTypeRapidScanResults)
			if err != nil || body != "" {
				ticker.Stop()
				timeoutTimer.Stop()
				if err != nil {
					log.Errorf("Error reading rapid scan result", err)
					return err, nil
				}

				err = json.Unmarshal([]byte(body), &result)

				if err != nil {
					log.Errorf("Error parsing rapid scan result", err)
				}

				return err, result
			}
		}
	}
}
