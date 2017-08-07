package hubclient

import (
	log "github.com/sirupsen/logrus"
)

func (c *Client) GetExternalExtension(link ResourceLink) (*ExternalExtension, error) {

	var extension ExternalExtension
	err := c.httpGetJSON(link.Href, &extension, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve an external extension: %+v.\n", err)
		return nil, err
	}

	return &extension, nil
}

func (c *Client) UpdateExternalExtension(extension *ExternalExtension) error {

	err := c.httpPutJSON(extension.Meta.Href, &extension, ContentTypeExtensionJSON, 200)

	if err != nil {
		log.Errorf("Error trying to update an external extension: %+v.\n", err)
		return err
	}

	return nil
}
