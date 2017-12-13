/*
   Copyright (C) 2016 The BlameWarrior Authors.
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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blamewarrior/hooks"
)

var SendingError = fmt.Errorf("sending error")

type Mediator interface {
	Mediate(payload string) (err error)
}

type MediatorService struct {
	ConsumerBaseURL string
	c               *http.Client

	payloads hooks.Payloads
}

func (service *MediatorService) Mediate(event string, payload []byte) (err error) {

	switch event {
	case "pull_request":
		ghPullRequestHook := new(gh.GithubPullRequestHook)
		err = json.Unmarshal(payload, &ghPullRequestHook)

		hookRepositoryName := ghPullRequestHook.Repository.FullName

		pullRequest := bw.NewPullRequestFromGithubHook(ghPullRequestHook)

		listCollaborators := collaborators.GetCollaboratorsFor(hookRepositoryName)

		collaboratorIds := service.pickCollaborators(listCollaborators)

	}

	err = service.send(result.ValueBytes)

	if err != nil {

	}

	return nil
}

func (service *MediatorService) send(payload []byte) (err error) {

	response, err := service.c.Post(
		service.ConsumerBaseURL+"/process_hook",
		"application/json",
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return SendingError
	}
}

func (service *MediatorService) pickCollaborators(collaborators []Collaborator) error {
	return nil
}
