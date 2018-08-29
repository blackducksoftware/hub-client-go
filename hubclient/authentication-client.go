package hubclient

import (
	"fmt"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	log "github.com/sirupsen/logrus"
)

// Public Key is retrieved with global setting API in Hub.

func (c *Client) GetPublicKey() (*hubapi.GlobalSetting, error) {
	var publicKey hubapi.GlobalSetting
	publicKeyURL := fmt.Sprintf("%s/api/public-key", c.baseURL)
	err := c.HttpGetJSON(publicKeyURL, &publicKey, 200)
	if err != nil {
		log.Errorf("Error while trying to get public key: %+v.\n", err)
		return nil, err
	}

	return &publicKey, nil
}
