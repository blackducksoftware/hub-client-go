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
	"time"

	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
)

func TestNewWithApiToken(t *testing.T) {
	defer gock.Off()

	gock.New("http://server.com").
		Post("/api/tokens/authenticate").
		Reply(200).
		JSON(map[string]interface{}{"bearerToken": "testBearerToken", "expiresInMilliseconds": 7199998}).
		SetHeader(HeaderNameCsrfToken, "csrf-fake-value")

	gock.DisableNetworking()

	gock.Intercept()

	client, err := NewWithApiToken("http://server.com", "foo", 0, time.Second*10)

	if err != nil {
		t.Fatalf("Got error: %v", err)
	}

	assert.True(t, client.useAuthToken)
	assert.Equal(t, "testBearerToken", client.authToken)
	assert.True(t, client.haveCsrfToken)
	assert.Equal(t, "csrf-fake-value", client.csrfToken)
}

func TestNewWithApiTokenAndClient(t *testing.T) {
	defer gock.Off()

	gock.New("http://server.com").
		Post("/api/tokens/authenticate").
		Reply(200).
		JSON(map[string]interface{}{"bearerToken": "testBearerToken", "expiresInMilliseconds": 7199998}).
		SetHeader(HeaderNameCsrfToken, "csrf-fake-value")

	gock.DisableNetworking()

	gock.Intercept()

	httpClient := createHttpClient(time.Second * 10)

	client, err := NewWithApiTokenAndClient("http://server.com", "foo", 0, httpClient)

	if err != nil {
		t.Fatalf("Got error: %v", err)
	}

	assert.True(t, client.useAuthToken)
	assert.Equal(t, "testBearerToken", client.authToken)
	assert.True(t, client.haveCsrfToken)
	assert.Equal(t, "csrf-fake-value", client.csrfToken)
}
