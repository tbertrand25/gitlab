package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type glMergeRequest struct {
	ID        int
	ProjectID int    `json:"project_id"`
	Source    string `json:"source_branch"`
	Target    string `json:"target_branch"`
}

func (c glClient) CreateMergeRequest(title string, path string, source string, target string, asssignee string) (glMergeRequest, error) {
	project, err := c.GetProject(path)
	if err != nil {
		fmt.Println(err)
	}

	URL := c.baseURL + "/projects/" + strconv.Itoa(project.ID) + "/merge_requests"

	user, err := c.GetUser(asssignee)
	if err != nil {
		fmt.Println(err)
	}

	params := make(map[string]string)
	params["private_token"] = c.auth
	params["title"] = title
	params["source_branch"] = source
	params["target_branch"] = target
	params["assignee_id"] = user.ID

	resp, err := c.Do("POST", URL, params)
	if err != nil {
		fmt.Println(err)
	}

	var mr glMergeRequest
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&mr)
	if err != nil {
		fmt.Println(err)
	}

	return mr, nil
}
