package hubclient

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
