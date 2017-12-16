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

package blamewarrior

import (
	"time"

	gh "github.com/blamewarrior/hooks/github"
)

type Collaborator struct {
	Id int `json:"id"`
}

type PullRequest struct {
	Id             int            `json:"id"`
	HTMLURL        string         `json:"html_url"`
	Title          string         `json:"title"`
	Body           string         `json:"body"`
	RepositoryName string         `json:"repository_name"`
	Reviewers      []Collaborator `json:"reviewers"`
	Number         int            `json:"number"`
	State          string         `json:"state"`
	CreatedAt      *time.Time     `json:"opened_at"`
	ClosedAt       *time.Time     `json:"closed_at"`
	OwnerId        int            `json:"owner_id"`
	Commits        int            `json:"commits"`
	Additions      int            `json:"additions"`
	Deletions      int            `json:"deletions"`
}

func (pullRequet *PullRequest) Valid() *Validator {
	v := new(Validator)

	v.MustNotBeZero(pullRequet.Id, "id must not be empty")
	v.MustNotBeZero(pullRequet.OwnerId, "owner_id must not be empty")
	v.MustNotBeZero(pullRequet.Number, "number must not be empty")
	v.MustNotBeEmpty(pullRequet.Title, "title must not be empty")
	v.MustNotBeEmpty(pullRequet.RepositoryName, "repository name must not be empty")
	v.MustNotBeEmpty(pullRequet.State, "state must not be empty")

	reviewers := make([]interface{}, len(pullRequet.Reviewers))
	for i, v := range pullRequet.Reviewers {
		reviewers[i] = v
	}
	v.MustNotBeZeroLength(reviewers, "reviewers must not be empty")

	return v
}

func NewPullRequestFromGithubHook(ghPullRequestHook *gh.GithubPullRequestHook) *PullRequest {
	pullRequest := &PullRequest{
		Id:             *ghPullRequestHook.PullRequest.ID,
		HTMLURL:        *ghPullRequestHook.PullRequest.HTMLURL,
		Title:          *ghPullRequestHook.PullRequest.Title,
		Body:           *ghPullRequestHook.PullRequest.Body,
		Number:         *ghPullRequestHook.PullRequest.Number,
		State:          *ghPullRequestHook.PullRequest.State,
		CreatedAt:      ghPullRequestHook.PullRequest.CreatedAt,
		ClosedAt:       ghPullRequestHook.PullRequest.ClosedAt,
		Commits:        *ghPullRequestHook.PullRequest.Commits,
		Additions:      *ghPullRequestHook.PullRequest.Additions,
		Deletions:      *ghPullRequestHook.PullRequest.Deletions,
		RepositoryName: ghPullRequestHook.Repository.FullName,
		OwnerId:        *ghPullRequestHook.PullRequest.User.ID,
	}

	requestedReviewers := ghPullRequestHook.RequestedReviewers

	if len(requestedReviewers) > 0 {
		pullRequest.Reviewers = make([]Collaborator, 0)
		for _, reviewer := range requestedReviewers {
			reviewer := &Collaborator{Id: reviewer.Id}
			pullRequest.Reviewers = append(pullRequest.Reviewers, *reviewer)
		}
	}

	return pullRequest
}
