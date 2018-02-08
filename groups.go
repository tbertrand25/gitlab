package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (c glClient) GetGroup(path string) (glGroup, error) {
	URL := c.baseURL + "/groups/" + strings.Replace(path, "/", "%2F", -1)

	params := make(map[string]string)
	params["private_token"] = c.auth

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	var group glGroup
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&group)
	if err != nil {
		fmt.Println(err)
	}

	return group, nil
}

func (c glClient) GetSubprojects(groupPath string) ([]glProject, error) {
	group, err := c.GetGroup(groupPath)
	if err != nil {
		fmt.Println(err)
	}

	var projects []glProject

	URL := c.baseURL + "/groups/" + strconv.Itoa(group.ID) + "/projects"

	params := make(map[string]string)
	params["private_token"] = c.auth

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&projects)
	if err != nil {
		fmt.Println(err)
	}

	URL = c.baseURL + "/groups/" + strconv.Itoa(group.ID) + "/subgroups"

	var subgroups []glGroup
	resp, err = c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&subgroups)
	if err != nil {
		fmt.Println(err)
	}

	var p []glProject
	for _, subgroup := range subgroups {
		p, err = c.GetSubprojects(subgroup.FullPath)
		if err != nil {
			fmt.Println(err)
		}
		projects = append(projects, p...)
	}

	return projects, nil
}
