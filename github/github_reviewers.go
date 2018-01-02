/*
   Copyright (C) 2017 The BlameWarrior Authors.
   This file is a part of BlameWarrior service.
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package github

import (
	"github.com/blamewarrior/hooks/blamewarrior/tokens"
	gh "github.com/google/go-github/github"
)

type Reviewers interface {
	RequestReviewers(ctx Context) (err error)
	ReviewComments(ctx Context) ([]ReviewComment, error)
}

type GithubReviewers struct {
	tokenClient     tokens.Client
	pullRequestHook *GithubPullRequestHook
}

func NewGithubReviewers(tokenClient tokens.Client, hook *GithubPullRequestHook) *GithubReviewers {
	return &GithubReviewers{tokenClient, hook}
}

func (service *GithubReviewers) RequestReviewers(ctx Context) (err error) {
	repositoryFullName := service.pullRequestHook.Repository.FullName
	reviewers := service.pullRequestHook.RequestedReviewers
	number := service.pullRequestHook.PullRequest.Number

	reviewersRequest := gh.ReviewersRequest{}

	for _, reviewer := range reviewers {
		reviewersRequest.Reviewers = append(reviewersRequest.Reviewers, reviewer.Login)
	}

	api, err := initAPIClient(ctx, service.tokenClient)
	if err != nil {
		return err
	}

	owner, repo := SplitRepositoryName(repositoryFullName)

	_, _, err = api.PullRequests.RequestReviewers(ctx, owner, repo, *number, reviewersRequest)

	if err != nil {
		return err
	}

	return nil
}

func (service *GithubReviewers) ReviewComments(ctx Context) ([]ReviewComment, error) {
	repositoryFullName := service.pullRequestHook.Repository.FullName

	number := service.pullRequestHook.PullRequest.Number

	api, err := initAPIClient(ctx, service.tokenClient)
	if err != nil {
		return nil, err
	}

	owner, repo := SplitRepositoryName(repositoryFullName)

	ghComments, _, err := api.PullRequests.ListComments(ctx, owner, repo, *number, nil)
	if err != nil {
		return nil, err
	}

	reviewComments := make([]ReviewComment, 0)

	for _, comment := range ghComments {
		reviewComment := ReviewComment(*comment)
		reviewComments = append(reviewComments, reviewComment)
	}

	return reviewComments, nil
}
