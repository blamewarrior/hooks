package github

import gh "github.com/google/go-github/github"

type GithubPullRequestUser struct {
	Id int `json:"id"`
}

type GithubPullRequestHook struct {
	PullRequest gh.PullRequest `json:"pull_request"`

	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`

	RequestedReviewer *struct {
		Id int `json: "id"`
	} `json:"requested_reviewer"`
}
