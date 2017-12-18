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

package collaborators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	gh "github.com/blamewarrior/hooks/github"
)

type Client interface {
	GetCollaborators(repositoryFullName string) ([]gh.Collaborator, error)
}

type CollaboratorsClient struct {
	BaseURL string
	c       *http.Client
}

func NewClient() *CollaboratorsClient {
	client := &CollaboratorsClient{
		BaseURL: "https://blamewarrior.com",
		c:       http.DefaultClient,
	}

	return client
}

func (client *CollaboratorsClient) GetCollaborators(repositoryFullName string) ([]gh.Collaborator, error) {

	requestUrl := fmt.Sprintf("%s/%s/collaborators", client.BaseURL, repositoryFullName)

	response, err := client.c.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to get collaborators for %s, status_code=%d", repositoryFullName, response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Unable to get response for %s", repositoryFullName)
	}

	collaborators := make([]gh.Collaborator, 0)

	if err = json.Unmarshal(bodyBytes, &collaborators); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal income json for %s, income=%s", repositoryFullName, string(bodyBytes))
	}

	return collaborators, nil
}
