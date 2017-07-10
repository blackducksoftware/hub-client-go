package hubclient

import "fmt"

// What about continuation for these?
// Should we have something where user can pass in an optional continuation/next placeholder?
// Or maybe that is something more for RX?
// Or maybe a special return type that can keep querying for all of them when it runs out?
// Is there any iterator type in GoLang?

func (c *Client) ListProjects() (*ProjectList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	projectsURL := fmt.Sprintf("%s/api/projects", c.baseURL)

	var projectList ProjectList
	err := c.httpGetJSON(projectsURL, &projectList, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &projectList, nil
}

func (c *Client) GetProject(link ResourceLink) (*Project, error) {

	var project Project
	err := c.httpGetJSON(link.Href, &project, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &project, nil
}

func (c *Client) ListProjectVerions(link ResourceLink) (*ProjectVersionList, error) {

	// Need offset/limit
	// Should we abstract list fetching like we did with a single Get?

	var versionList ProjectVersionList
	err := c.httpGetJSON(link.Href, &versionList, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &versionList, nil
}

func (c *Client) GetProjectVerion(link ResourceLink) (*ProjectVersion, error) {

	var projectVersion ProjectVersion
	err := c.httpGetJSON(link.Href, &projectVersion, 200)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &projectVersion, nil
}
