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
	"net/http"
	"strconv"
	"time"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
)

const (
	headerBdMode          = "X-BD-MODE"
	headerBdDocumentCount = "X-BD-DOCUMENT-COUNT"
	bdModeAppend          = "append"
	bdModeFinish          = "finish"
)

type ChunkIterator interface {
	hasNext() bool
	next() string
}

func (c *Client) StartRapidScan(bdioHeaderContent string) (error, string) {
	rapidScansURL := hubapi.BuildUrl(c.baseURL, hubapi.DeveloperScansApi)
	bdioUploadEndpoint, err := c.HttpPostString(rapidScansURL, bdioHeaderContent, hubapi.ContentTypeRapidScanRequest, http.StatusCreated)

	if err != nil {
		log.Error("Error kicking off a rapid scan.", err)
		return err, ""
	}

	return nil, bdioUploadEndpoint
}

func (c *Client) UploadBdioFiles(bdioUploadEndpoint string, bdioContents []string) error {
	iterator := NewArrayChunkIterator(bdioContents)
	return c.UpdateBdioFilesByChunk(bdioUploadEndpoint, len(bdioContents), &iterator)
}

func (c *Client) UpdateBdioFilesByChunk(bdioUploadEndpoint string, chunkcount int, iterator ChunkIterator) error {
	header := http.Header{}
	header.Add(headerBdMode, bdModeAppend)
	header.Add(headerBdDocumentCount, strconv.Itoa(chunkcount))

	for iterator.hasNext() {
		bdioContent := iterator.next()
		err := c.HttpPutStringWithHeader(bdioUploadEndpoint, bdioContent, hubapi.ContentTypeRapidScanRequest, http.StatusAccepted, header)
		if err != nil {
			log.Error("Error uploading bdio files.", err)
			return err
		}
	}

	header.Set(headerBdMode, bdModeFinish)
	err := c.HttpPutStringWithHeader(bdioUploadEndpoint, "", hubapi.ContentTypeRapidScanRequest, http.StatusAccepted, header)
	if err != nil {
		log.Error("Error uploading bdio files.", err)
		return err
	}

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
				log.Error("Error fetching rapid scan result", err)
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
						log.Error("Error fetching rapid scan result", err)
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

func (c *Client) FetchResults(rapidScanEndpoint string, offset int, pageLimit int) (err error, httpStatus int, result *hubapi.RapidScanResult) {
	var body string
	err, statusCode := c.fetchResults(rapidScanEndpoint, offset, pageLimit, &body)
	if err != nil || statusCode != http.StatusOK {
		return err, statusCode, nil
	}
	err, result = parseBody(body)
	return err, statusCode, result
}

func parseBody(body string) (error, *hubapi.RapidScanResult) {
	var pagedResult *hubapi.RapidScanResult
	err := json.Unmarshal([]byte(body), &pagedResult)

	if err != nil {
		log.Error("Error parsing rapid scan result", err)
		return err, nil
	}

	return nil, pagedResult
}

func (c *Client) fetchResults(rapidScanEndpoint string, offset int, limit int, body *string) (error, int) {
	url := hubapi.BuildUrl(rapidScanEndpoint, hubapi.FullResultsApi)
	params := make(map[string]string)
	params["offset"] = strconv.Itoa(offset)
	params["limit"] = strconv.Itoa(limit)
	url = hubapi.AddParameters(url, params)
	err, statusCode := c.HttpGetString(url, body, []int{http.StatusOK, http.StatusNotFound}, hubapi.ContentTypeRapidScanResults)
	return err, statusCode
}
