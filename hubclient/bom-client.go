package hubclient

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (c *Client) ListProjectVersionComponents(link ResourceLink) (*BomComponentList, error) {

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

// TODO: Should this be used?
func (c *Client) ListProjectVerionVulnerableComponents(link ResourceLink) (*BomVulnerableComponentList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var bomList BomVulnerableComponentList
	err := c.httpGetJSON(link.Href+"?limit=2", &bomList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve vulnerable components list: %+v.\n", err)
		return nil, err
	}

	return &bomList, nil
}

func (c *Client) PageProjectVersionVulnerableComponents(link ResourceLink, offset uint32, limit uint32) (*BomVulnerableComponentList, error) {

	// Should we abstract list fetching like we did with a single Get?

	var bomList BomVulnerableComponentList
	url := fmt.Sprintf("%s?offset=%d&limit=%d", link.Href, offset, limit)
	err := c.httpGetJSON(url, &bomList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve vulnerable components page: %+v.\n", err)
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
		log.Errorf("Error trying to retrieve count of vulnerable components: %+v.\n", err)
		return 0, err
	}

	return bomList.TotalCount, nil
}

func (c *Client) ListAllProjectVerionVulnerableComponents(link ResourceLink) ([]BomVulnerableComponent, error) {

	log.Debugf("***** Getting total count.\n")
	//totalCount, err := c.CountProjectVerionVulnerableComponents(link)
	totalCount := uint32(100)
	log.Debugf("***** Got total count: %d\n", totalCount)

	// if err != nil {
	// 	log.Debugf("ERROR GETTING COUNT: %s\n", err)
	// }

	pageSize := uint32(100)
	result := make([]BomVulnerableComponent, totalCount, totalCount)

	for offset := uint32(0); offset < totalCount; offset += pageSize {

		log.Debugf("***** Going to get vulnerable components. Offset: %d, Limit: %d \n", offset, pageSize)
		bomPage, err := c.PageProjectVersionVulnerableComponents(link, offset, pageSize)

		if err != nil {
			log.Errorf("Error trying to retrieve vulnerable components list: %+v.\n", err)
		}

		result = append(result, bomPage.Items...)
	}

	return result, nil
}
