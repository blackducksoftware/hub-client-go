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
