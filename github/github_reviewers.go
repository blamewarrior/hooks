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
	RequestReviewers(ctx Context, repoFullName string, pullNumber int, reviewers []Collaborator) (err error)
	ReviewComments(ctx Context, repoFullName string, pullNumber int) ([]ReviewComment, error)
}

type GithubReviewers struct {
	tokenClient tokens.Client
}

func NewGithubReviewers(tokenClient tokens.Client) *GithubReviewers {
	return &GithubReviewers{tokenClient}
}

func (service *GithubReviewers) RequestReviewers(ctx Context, repoFullName string, pullNumber int, reviewers []Collaborator) (err error) {

	reviewersRequest := gh.ReviewersRequest{}

	for _, reviewer := range reviewers {
		reviewersRequest.Reviewers = append(reviewersRequest.Reviewers, reviewer.Login)
	}

	api, err := initAPIClient(ctx, service.tokenClient)
	if err != nil {
		return err
	}

	owner, repo := SplitRepositoryName(repoFullName)

	_, _, err = api.PullRequests.RequestReviewers(ctx, owner, repo, pullNumber, reviewersRequest)

	if err != nil {
		return err
	}

	return nil
}

func (service *GithubReviewers) ReviewComments(ctx Context, repoFullName string, pullNumber int) ([]ReviewComment, error) {

	api, err := initAPIClient(ctx, service.tokenClient)
	if err != nil {
		return nil, err
	}

	owner, repo := SplitRepositoryName(repoFullName)

	ghComments, _, err := api.PullRequests.ListComments(ctx, owner, repo, pullNumber, nil)
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
