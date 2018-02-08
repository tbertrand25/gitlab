package gitlab

import (
	"encoding/json"
	"fmt"
)

type glUser struct {
	ID   string
	Name string
}

func (c glClient) GetUser(name string) (glUser, error) {
	URL := c.baseURL + "/users"

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["username"] = name

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	var user glUser
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&user)
	if err != nil {
		fmt.Println(err)
	}

	return user, nil
}
