package gitlab

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func (c glClient) CreateBranch(path string, newBranch string, sourceBranch string) error {
	ID, _ := c.ProjectID(path)
	URL := c.baseURL + "/projects/" + strconv.Itoa(ID) + "/repository/branches"

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["branch"] = newBranch
	params["ref"] = sourceBranch

	_, err := c.Do("POST", URL, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func (c glClient) ProtectBranch(pathToBranch string, push int, merge int) {
	path := strings.Split(pathToBranch, "#")[0]
	branch := strings.Split(pathToBranch, "#")[1]
	ID, err := c.ProjectID(path)
	if err != nil {
		fmt.Println(err)
	}

	URL := c.baseURL + "/projects/" + strconv.Itoa(ID) + "/protected_branches/" + branch
	params := make(map[string]string)
	params["private_token"] = c.auth

	resp, _ := c.Do("DELETE", URL, params)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))

	URL = c.baseURL + "/projects/" + strconv.Itoa(ID) + "/protected_branches"

	params["name"] = branch
	params["push_access_level"] = strconv.Itoa(push)
	params["merge_access_level"] = strconv.Itoa(merge)

	resp, err = c.Do("POST", URL, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (c glClient) SetDefaultBranch(path string) error {
	project := strings.Split(path, "#")[0]
	branch := strings.Split(path, "#")[1]

	ID, err := c.ProjectID(project)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	URL := c.baseURL + "/projects/" + strconv.Itoa(ID)

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["default_branch"] = branch

	_, err = c.Do("PUT", URL, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
