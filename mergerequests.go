package gitlab

type glMergeRequest struct {
	ID        int
	ProjectID int    `json:"project_id"`
	Source    string `json:"source_branch"`
	Target    string `json:"target_branch"`
}

func (c glClient) CreateMergeRequest(path string, source string, dest string, asssignee int) error {

	return nil
}
