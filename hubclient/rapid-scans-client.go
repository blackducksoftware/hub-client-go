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
	"net/http"
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
)

func (c *Client) StartRapidScan(bdioHeaderContent string) (error, string) {
	rapidScansURL := c.baseURL + apiDeveloperScans
	bdioUploadEndpoint, err := c.HttpPostString(rapidScansURL, bdioHeaderContent, hubapi.ContentTypeRapidScanRequest, http.StatusCreated)

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
		err := c.HttpPutString(bdioUploadEndpoint, bdioContent, hubapi.ContentTypeRapidScanRequest, http.StatusAccepted)
		if err != nil {
			log.Errorf("Error uploading bdio files.", err)
			return err
		}
	}

	c.SetHeaderValue(headerBdMode, bdModeFinish)
	err := c.HttpPutString(bdioUploadEndpoint, "", hubapi.ContentTypeRapidScanRequest, http.StatusAccepted)
	if err != nil {
		log.Errorf("Error uploading bdio files.", err)
		return err
	}

	c.DeleteHeaderValue(headerBdMode)
	c.DeleteHeaderValue(headerBdDocumentCount)

	return nil
}

func (c *Client) PollRapidScanResults(rapidScanEndpoint string, interval, timeout time.Duration) (error, *hubapi.RapidScanResult) {
	url := rapidScanEndpoint + apiFullResults
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
			err := c.HttpGetString(url, &body, http.StatusOK, hubapi.ContentTypeRapidScanResults)
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
