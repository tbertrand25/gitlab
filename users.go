package gitlab

import (
	"encoding/json"
	"fmt"
)

type glUser struct {
	ID       int
	UserName string
	Name     string
}

func (c glClient) SearchUsers(username string) ([]glUser, error) {
	URL := c.baseURL + "/users"

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["search"] = username

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	var users []glUser
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&users)
	if err != nil {
		fmt.Println(err)
	}

	return users, nil
}
