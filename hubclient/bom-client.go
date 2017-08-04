package hubclient

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (c *Client) ListProjectVerionComponents(link ResourceLink) (*BomComponentList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList BomComponentList
	err := c.httpGetJSON(link.Href+"?limit=2", &bomList, 200)

	if err != nil {
		log.Errorf("Error while trying to get Project Version Component list: %+v.\n", err)
		return nil, err
	}

	return &bomList, nil
}

func (c *Client) ListProjectVerionVulnerableComponents(link ResourceLink) (*BomVulnerableComponentList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList BomVulnerableComponentList
	err := c.httpGetJSON(link.Href+"?limit=2", &bomList, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &bomList, nil
}

func (c *Client) PageProjectVerionVulnerableComponents(link ResourceLink, offset uint32, limit uint32) (*BomVulnerableComponentList, error) {

	// Should we abstract list fetching like we did with a single Get?

	var bomList BomVulnerableComponentList
	url := fmt.Sprintf("%s?offset=%d&limit=%d", link.Href, offset, limit)
	err := c.httpGetJSON(url, &bomList, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &bomList, nil
}

func (c *Client) CountProjectVerionVulnerableComponents(link ResourceLink) (uint32, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList BomVulnerableComponentList
	err := c.httpGetJSON(link.Href+"?offset=0&limit=1", &bomList, 200)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return bomList.TotalCount, nil
}

func (c *Client) ListAllProjectVerionVulnerableComponents(link ResourceLink) ([]BomVulnerableComponent, error) {

	fmt.Printf("***** Getting total count.\n")
	//totalCount, err := c.CountProjectVerionVulnerableComponents(link)
	totalCount := uint32(100)
	fmt.Printf("***** Got total count: %d\n", totalCount)

	// if err != nil {
	// 	fmt.Printf("ERROR GETTING CONUT: %s\n", err)
	// }

	pageSize := uint32(100)
	result := make([]BomVulnerableComponent, totalCount, totalCount)

	for offset := uint32(0); offset < totalCount; offset += pageSize {

		fmt.Printf("***** Going to get vulnerable components. Offset: %d, Limit: %d \n", offset, pageSize)
		bomPage, err := c.PageProjectVerionVulnerableComponents(link, offset, pageSize)

		if err != nil {
			fmt.Println(err)
		}

		result = append(result, bomPage.Items...)
	}

	return result, nil
}
