package hubclient

import (
	"fmt"

	"bitbucket.org/bdsengineering/go-hub-client/hubapi"

	log "github.com/sirupsen/logrus"
)

// What about continuation for these?
// Should we have something where user can pass in an optional continuation/next placeholder?
// Or maybe that is something more for RX?
// Or maybe a special return type that can keep querying for all of them when it runs out?
// Is there any iterator type in GoLang?

func (c *Client) ListProjects() (*hubapi.ProjectList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	projectsURL := fmt.Sprintf("%s/api/projects", c.baseURL)

	var projectList hubapi.ProjectList
	err := c.httpGetJSON(projectsURL, &projectList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve project list: %+v.", err)
		return nil, err
	}

	return &projectList, nil
}

func (c *Client) GetProject(link hubapi.ResourceLink) (*hubapi.Project, error) {

	var project hubapi.Project
	err := c.httpGetJSON(link.Href, &project, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a project: %+v.", err)
		return nil, err
	}

	return &project, nil
}

func (c *Client) ListProjectVersions(link hubapi.ResourceLink) (*hubapi.ProjectVersionList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var versionList hubapi.ProjectVersionList
	err := c.httpGetJSON(link.Href, &versionList, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve project version list list: %+v.", err)
		return nil, err
	}

	return &versionList, nil
}

func (c *Client) GetProjectVersion(link hubapi.ResourceLink) (*hubapi.ProjectVersion, error) {

	var projectVersion hubapi.ProjectVersion
	err := c.httpGetJSON(link.Href, &projectVersion, 200)

	if err != nil {
		log.Errorf("Error trying to retrieve a project version: %+v.", err)
		return nil, err
	}

	return &projectVersion, nil
}
