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
		for _, reviewer := range requestedReviewers {
			pullRequest.ReviewerIds = append(pullRequest.ReviewerIds, reviewer.Id)
		}
	}

	return pullRequest
}
