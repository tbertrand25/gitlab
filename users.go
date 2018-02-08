package gitlab

import (
	"encoding/json"
	"fmt"
)

type glUser struct {
	ID   int
	Name string
}

func (c glClient) GetUsers() ([]glUser, error) {
	URL := c.baseURL + "/users"

	params := make(map[string]string)
	params["private_token"] = c.auth

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
