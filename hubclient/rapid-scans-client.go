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

func (c *Client) PollRapidScanResults(rapidScanEndpoint string, interval, timeout time.Duration, pageLimit int) (error, *hubapi.RapidScanResult) {
	offset := 0
	ticker := time.NewTicker(interval)
	timeoutTimer := time.NewTimer(timeout)
	var body string

	for {
		select {
		case <-timeoutTimer.C:
			ticker.Stop()
			return errors.New(fmt.Sprintf("The polling for rapid scan result timed out: %s", rapidScanEndpoint)), nil
		case <-ticker.C:
			err, statusCode := c.fetchResults(rapidScanEndpoint, offset, pageLimit, &body)

			if err != nil {
				ticker.Stop()
				timeoutTimer.Stop()
				log.Errorf("Error fetching rapid scan result", err)
				return err, nil
			}

			if statusCode == http.StatusOK {
				ticker.Stop()
				timeoutTimer.Stop()

				err, result := parseBody(body)

				if err != nil {
					return err, result
				}

				//read all pages of the result
				for result.Count > len(result.Components) {
					//increase offset to fetch the next page of results
					offset += pageLimit
					err, statusCode := c.fetchResults(rapidScanEndpoint, offset, pageLimit, &body)

					if err != nil {
						log.Errorf("Error fetching rapid scan result", err)
						return err, result
					}

					if statusCode != http.StatusOK {
						log.Errorf("Error fetching subsequent pages of a rapid scan result. Code: %d", statusCode)
						return err, result
					}

					err, pagedResult := parseBody(body)
					if err != nil {
						return err, result
					}

					result.Components = append(result.Components, pagedResult.Components...)
				}

				return nil, result
			}
		}
	}
}

func parseBody(body string) (error, *hubapi.RapidScanResult) {
	var pagedResult *hubapi.RapidScanResult
	err := json.Unmarshal([]byte(body), &pagedResult)

	if err != nil {
		log.Errorf("Error parsing rapid scan result", err)
		return err, nil
	}

	return nil, pagedResult
}

func (c *Client) fetchResults(rapidScanEndpoint string, offset int, limit int, body *string) (error, int) {
	url := rapidScanEndpoint + apiFullResults + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	err, statusCode := c.HttpGetString(url, body, []int{http.StatusOK, http.StatusNotFound}, hubapi.ContentTypeRapidScanResults)
	return err, statusCode
}
