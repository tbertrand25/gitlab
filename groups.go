package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
)

type glGroup struct {
	ID       int
	FullPath string `json:"full_path"`
}

func (c glClient) SearchGroups(name string) []glGroup {
	var groups []glGroup

	URL := c.baseURL + "/groups"

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["search"] = name

	fmt.Println(params)

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&groups)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var nextGroups []glGroup
	for {
		linkMap := c.parseLinkHeader(resp.Header.Get("Link"))

		next, exists := linkMap["next"]
		if exists {
			resp, err = c.Do("GET", next, params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			decoder = json.NewDecoder(resp.Body)
			err = decoder.Decode(&nextGroups)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			groups = append(groups, nextGroups...)
		} else {
			break
		}
	}

	return groups
}
