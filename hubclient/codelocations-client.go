package hubclient

import (
	log "github.com/sirupsen/logrus"
)

func (c *Client) ListCodeLocations(link ResourceLink) (*CodeLocationList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var codeLocationList CodeLocationList
	err := c.httpGetJSON(link.Href, &codeLocationList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve code location list: %+v.\n", err)
		return nil, err
	}

	return &codeLocationList, nil
}

func (c *Client) GetCodeLocation(link ResourceLink) (*CodeLocation, error) {

	var codeLocation CodeLocation
	err := c.httpGetJSON(link.Href, &codeLocation, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a code location: %+v.\n", err)
		return nil, err
	}

	return &codeLocation, nil
}

func (c *Client) ListScanSummaries(link ResourceLink) (*ScanSummaryList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var scanSummaryList ScanSummaryList
	err := c.httpGetJSON(link.Href, &scanSummaryList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve scan summary list: %+v.\n", err)
		return nil, err
	}

	return &scanSummaryList, nil
}

func (c *Client) GetScanSummary(link ResourceLink) (*ScanSummary, error) {

	var scanSummary ScanSummary
	err := c.httpGetJSON(link.Href, &scanSummary, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a scan summary: %+v.\n", err)
		return nil, err
	}

	return &scanSummary, nil
}
