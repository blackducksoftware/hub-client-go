package hubclient

import (
	"crypto/rsa"
	"fmt"

	"github.com/blackducksoftware/hub-client-go/hubapi"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// Public Key is retrieved with global setting API in Hub.

func (c *Client) GetPublicKey() (*rsa.PublicKey, error) {
	var publicKey hubapi.GlobalSetting
	publicKeyURL := fmt.Sprintf("%s/api/public-key", c.baseURL)
	err := c.HttpGetJSON(publicKeyURL, &publicKey, 200)
	if err != nil {
		log.Errorf("Error while trying to get public key: %+v.\n", err)
		return nil, err
	}

	key, err := LoadKeyFromString(publicKey.Value)
	if err != nil {
		log.Errorf("Could not load key from the public key string", err)
		return nil, err
	}
	return key, nil
}

func LoadKeyFromString(publicKeyString string) (*rsa.PublicKey, error) {
	key := "-----BEGIN PUBLIC KEY-----\n" + publicKeyString + "\n-----END PUBLIC KEY-----"
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		log.Errorf("Cannot parse the public key string into a key", err)
		return nil, err
	}
	return rsaKey, nil
}
