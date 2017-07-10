package hubclient

import (
	"fmt"
	"net/url"
)

func (c *Client) Login(username string, password string) error {

	loginURL := fmt.Sprintf("%s/j_spring_security_check", c.baseURL)
	formValues := url.Values{
		"j_username": {username},
		"j_password": {password},
	}

	resp, err := c.httpClient.PostForm(loginURL, formValues)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if resp.StatusCode != 204 {
		fmt.Printf("Login: Got a %d reponse instead of a 204.\n", resp.StatusCode)
		return fmt.Errorf("got a %d response instead of a 204", resp.StatusCode)
	}

	fmt.Println("Login: Successfully authenticated")

	return nil
}
