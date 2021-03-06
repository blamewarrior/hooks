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

package hooks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	bw "github.com/blamewarrior/hooks/blamewarrior"
	"github.com/blamewarrior/hooks/blamewarrior/collaborators"
	"github.com/blamewarrior/hooks/blamewarrior/tokens"
	"github.com/blamewarrior/hooks/blamewarrior/web"
	gh "github.com/blamewarrior/hooks/github"
)

var SendingError = fmt.Errorf("sending error")

type Mediator interface {
	Mediate(event string, payload []byte) (err error)
}

type MediatorService struct {
	ConsumerBaseURL string
	c               *http.Client

	payloads Payloads

	webClient           web.Client
	collaboratorsClient collaborators.Client
	tokenClient         tokens.Client
	reviewersClient     gh.Reviewers
}

func NewMediatorService(
	payloads Payloads, webClient web.Client,
	collaboratorsClient collaborators.Client, reviewers gh.Reviewers) *MediatorService {
	return &MediatorService{
		payloads:            payloads,
		webClient:           webClient,
		collaboratorsClient: collaboratorsClient,
		reviewersClient:     reviewers,
	}
}

func (service *MediatorService) Mediate(event string, payload []byte) (err error) {

	switch event {
	case "pull_request":
		err = service.handlePullRequestPayload(payload)
	case "member":
		err = service.handleMemberPayload(payload)
	}

	if err != nil {
		if err = service.payloads.Save(string(payload)); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (service *MediatorService) handleMemberPayload(payload []byte) (err error) {
	ghMemberHook := new(gh.GithubMemberHook)

	if err = json.Unmarshal(payload, &ghMemberHook); err != nil {
		return err
	}

	memberRepositoryName := ghMemberHook.Repository.FullName
	member := ghMemberHook.Member

	switch ghMemberHook.Action {
	case "added":
		err = service.collaboratorsClient.AddCollaborator(memberRepositoryName, &member)
	case "edited":
		err = service.collaboratorsClient.EditCollaborator(memberRepositoryName, &member)
	case "deleted":
		err = service.collaboratorsClient.DeleteCollaborator(memberRepositoryName, ghMemberHook.Member.Login)
	}
	return err
}

func (service *MediatorService) handlePullRequestPayload(payload []byte) (err error) {
	ghPullRequestHook := new(gh.GithubPullRequestHook)
	if err = json.Unmarshal(payload, &ghPullRequestHook); err != nil {
		return err
	}

	hookRepositoryName := ghPullRequestHook.Repository.FullName

	pullRequest := bw.NewPullRequestFromGithubHook(ghPullRequestHook)

	if len(pullRequest.Reviewers) == 0 {
		listCollaborators, err := service.collaboratorsClient.ListCollaborator(hookRepositoryName)

		reviewer := service.pickCollaboratorFrom(listCollaborators)
		reviewers := []gh.Collaborator{*reviewer}

		pullRequest.Reviewers = reviewers

		if err = service.reviewersClient.RequestReviewers(gh.Context{},
			hookRepositoryName,
			pullRequest.Number,
			reviewers,
		); err != nil {
			return err
		}
	}

	if pullRequest.State != "open" {
		comments, err := service.reviewersClient.ReviewComments(gh.Context{},
			hookRepositoryName,
			pullRequest.Number,
		)
		if err != nil {
			return err
		}

		pullRequest.ReviewComments = comments
	}

	err = service.webClient.ProcessPullRequest(pullRequest)

	if err != nil {
		return err
	}

	return nil
}

func (service *MediatorService) pickCollaboratorFrom(collaborators []gh.Collaborator) *gh.Collaborator {
	admins := make([]gh.Collaborator, 0)

	for _, collaborator := range collaborators {
		if collaborator.Admin {
			admins = append(admins, collaborator)
		}
	}

	rand.Seed(time.Now().Unix())

	n := rand.Int() % len(admins)

	return &admins[n]
	return nil
}
