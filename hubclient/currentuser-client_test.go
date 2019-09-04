package hubclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestClient(t *testing.T) *Client {
	var serverInfo = struct {
		Url      string
		User     string
		Password string
	}{
		Url:      "https://localhost",
		User:     "sysadmin",
		Password: "blackduck",
	}

	if _, err := os.Stat("testserver.json"); err == nil {
		file, err := ioutil.ReadFile("testserver.json")
		assert.NotNil(t, file)
		assert.NoError(t, err)

		fileBytes := rot13rot5(string(file))

		err = json.Unmarshal([]byte(fileBytes), &serverInfo)
		assert.NoError(t, err)
	}

	client, err := NewWithSession(serverInfo.Url, HubClientDebugTimings, 5*time.Second)
	assert.NoError(t, err, "unable to instantiate client: %s", err)

	err = client.Login(serverInfo.User, serverInfo.Password)
	assert.NoError(t, err, "unable to login client: %s", err)

	return client
}

func TestClient_CreateAndDeleteApiToken(t *testing.T) {
	client := createTestClient(t)
	assert.NotNil(t, client, "unable to get client")

	tokenName := "testToken_" + randomString(10)
	tokenDesc := "testTokenDesc_" + randomString(20)

	for _, ro := range []bool{false, true} {
		location, token, err := client.CreateApiToken(tokenName, tokenDesc, ro)

		assert.NoError(t, err, "unable to create a token: %s readonly=%v", err, ro)
		assert.NotEmpty(t, token, "empty token created")

		err = client.DeleteApiToken(location)
		assert.NoError(t, err, "unable to delete a token")
	}
}

func TestClient_ListApiTokens(t *testing.T) {
	client := createTestClient(t)

	currentVersion, err := client.ListApiTokens(nil)
	assert.NoError(t, err, "unable to get current version")

	t.Logf("successfully fetched current version %+v", currentVersion)
}
