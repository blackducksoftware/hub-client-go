package hubclient

import "fmt"

func (c *Client) GetExternalExtension(link ResourceLink) (*ExternalExtension, error) {

	var extension ExternalExtension
	err := c.httpGetJSON(link.Href, &extension, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &extension, nil
}

func (c *Client) UpdateExternalExtension(extension *ExternalExtension) error {

	err := c.httpPutJSON(extension.Meta.Href, &extension, ContentTypeExtensionJSON, 200)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
