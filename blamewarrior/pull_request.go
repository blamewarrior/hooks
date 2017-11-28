package blamewarrior

import (
	"encoding/json"
	"time"
)

type PullRequest struct {
	Id             int        `json:"id,omitempty"`
	HTMLURL        string     `json:"html_url,omitempty"`
	Title          string     `json:"title,omitempty"`
	Body           string     `json:"body,omitempty"`
	RepositoryName string     `json:"repository_name,omitempty"`
	ReviewerId     int        `json:"reviewer_id,omitempty"`
	Number         int        `json:"number,omitempty"`
	State          string     `json:"state,omitempty"`
	CreatedAt      *time.Time `json:"opened_at,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	OwnerId        int        `json:"owner_id,omitempty"`
	Commits        int        `json:"commits,omitempty"`
	Additions      int        `json:"additions,omitempty"`
	Deletions      int        `json:"deletions,omitempty"`
}

func (pr *PullRequest) UnmarshalJSON(data []byte) error {
	type Alias PullRequest

	type User struct {
		Id int `json:"id"`
	}

	githubPullRequest := &struct {
		PullRequest struct {
			User  *User  `json: "user,omitempty"`
			State string `json:"state"`
		} `json:"pull_request"`

		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`

		RequestedReviewer *struct {
			Id int `json: "id"`
		} `json:"requested_reviewer"`
		*Alias
	}{
		Alias: (*Alias)(pr),
	}
	if err := json.Unmarshal(data, &githubPullRequest); err != nil {
		return err
	}

	pr.OwnerId = githubPullRequest.PullRequest.User.Id
	pr.State = githubPullRequest.PullRequest.State
	pr.RepositoryName = githubPullRequest.Repository.FullName

	if githubPullRequest.RequestedReviewer != nil {
		pr.ReviewerId = githubPullRequest.RequestedReviewer.Id
	}
	return nil
}
