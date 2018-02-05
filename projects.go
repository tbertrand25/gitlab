package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type glProject struct {
	ID                int
	PathWithNamespace string `json:"path_with_namespace"`
}

func (c glClient) SearchProjects(name string) []glProject {
	URL := c.baseURL + "/projects"
	params := make(map[string]string)
	params["private_token"] = c.auth
	params["search"] = name

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var projects []glProject

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&projects)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var nextProjects []glProject
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
			err = decoder.Decode(&nextProjects)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			projects = append(projects, nextProjects...)
		} else {
			break
		}
	}

	return projects
}

func (c glClient) GetProject(path string) (glProject, error) {
	URL := c.baseURL + "/projects/" + strings.Replace(path, "/", "%2F", -1)

	params := make(map[string]string)
	params["private_token"] = c.auth

	resp, err := c.Do("GET", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	var project glProject
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&project)
	if err != nil {
		fmt.Println(err)
		return glProject{}, err
	}

	return project, nil
}

func (c glClient) ProjectID(path string) (int, error) {

	project, err := c.GetProject(path)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return project.ID, nil
}

func (c glClient) CreateProject(path string) (glProject, error) {
	URL := c.baseURL + "/projects"

	elems := strings.Split(path, "/")
	name := elems[len(elems)-1]
	group := elems[len(elems)-2]

	potentialGroups := c.SearchGroups(group)
	groupID := -1
	for _, g := range potentialGroups {
		if g.FullPath == strings.Join(elems[:len(elems)-1], "/") {
			groupID = g.ID
		}
	}

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["name"] = name
	params["namespace_id"] = strconv.Itoa(groupID)

	resp, err := c.Do("POST", URL, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var project glProject
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&project)

	return project, nil
}
