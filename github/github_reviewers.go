package github

import (
	"net/http"

	gh "github.com/google/go-github/github"
)

type Reviewers interface {
	RequestReviewers(repositoryFullName string, reviewers []Collaborator) (err error)
	ReviewComments(commentsURL string) (string, error)
}

type GithubReviewers struct {
	c http.Client
}

func (service *GithubReviewers) RequestReviewers(ctx Context, ghPRHook *GithubPullRequestHook) (err error) {
	repositoryFullName := ghPRHook.Repository
	reviewers := ghPRHook.RequestedReviewers

	number := ghPRHook.PullRequest.Number

	reviewersRequest := new(gh.ReviewersRequest)

	for reviewer := range reviewers {
		reviewersRequest.Reviewers = append(reviewersRequest.Reviewers, reviewer.Login)
	}

	api := initAPIClient(ctx)

	owner, repo := SplitRepositoryName(repositoryFullName)

	api.Reviewers.RequestReviewers(ctx, owner, repo, number, reviewersRequest)
	return nil
}

func (service *GithubReviewers) ReviewComments(ctx Context, ghPRHook *GithubPullRequestHook) (string, error) {
	repositoryFullName := ghPRHook.Repository

	number := ghPRHook.PullRequest.Number

	api := initAPIClient(ctx)

	owner, repo := SplitRepositoryName(repositoryFullName)

	api.PullRequests.ListComments(ctx, owner, repo, number)
}
