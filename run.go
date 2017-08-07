package main

import "bitbucket.org/bdsengineering/go-hub-client/hubclient"
import "fmt"

func main() {

	client, err := hubclient.NewWithSession("https://localhost", hubclient.HubClientDebugTimings)

	if err != nil {
		fmt.Printf("Error making hub client: %s\n", err)
	}

	err = client.Login("sysadmin", "blackduck")

	if err != nil {
		fmt.Printf("Error logging into hub: %s\n", err)
	}

	// Projects

	fmt.Println("\n\n===========================================================")
	projects, err := client.ListProjects()

	if err != nil {
		fmt.Printf("Error getting project list: %s\n", err)
	}

	fmt.Printf("Projects: %+v\n", projects)

	// Versions

	fmt.Println("\n\n===========================================================")
	versionLink, err := projects.Items[0].GetProjectVersionLink()

	if err != nil {
		fmt.Printf("Error getting project versions link: %s\n", err)
	}

	versions, err := client.ListProjectVersions(*versionLink)

	if err != nil {
		fmt.Printf("Error getting versions list: %s\n", err)
	}

	fmt.Printf("Versions: %+v\n", versions)

	// Components

	fmt.Println("\n\n===========================================================")
	versionCompLink, err := versions.Items[0].GetComponentsLink()

	if err != nil {
		fmt.Printf("Error getting project verison components link: %s\n", err)
	}

	components, err := client.ListProjectVersionComponents(*versionCompLink)

	if err != nil {
		fmt.Printf("Error getting version component list: %s\n", err)
	}

	fmt.Printf("Components: %+v\n", components)

	// Vulnerable Components

	fmt.Println("\n\n===========================================================")
	versionVulnCompLink, err := versions.Items[0].GetVulnerableComponentsLink()

	if err != nil {
		fmt.Printf("Error getting project verion vuln components link: %s\n", err)
	}

	client.ListAllProjectVerionVulnerableComponents(*versionVulnCompLink)
	// vulnComponents, err := client.ListAllProjectVerionVulnerableComponents(*versionVulnCompLink)

	// 	if err != nil {
	// 		fmt.Printf("Error getting version vuln component list: %s\n", err)
	// 	}

	// 	fmt.Printf("Vuln Components: %+v\n", vulnComponents)
}
