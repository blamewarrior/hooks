package blamewarrior

import (
	"time"
)

type PullRequest struct {
	Id             int        `json:"id"`
	HTMLURL        string     `json:"html_url"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	RepositoryName string     `json:"repository_name"`
	ReviewerIds    []int      `json:"reviewer_ids"`
	Number         int        `json:"number"`
	State          string     `json:"state"`
	CreatedAt      *time.Time `json:"opened_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	OwnerId        int        `json:"owner_id"`
	Commits        int        `json:"commits"`
	Additions      int        `json:"additions"`
	Deletions      int        `json:"deletions"`
}
