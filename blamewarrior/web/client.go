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

package web

import (
	"fmt"
	"net/http"

	bw "github.com/blamewarrior/hooks/blamewarrior"
)

type Client struct {
	BaseURL string
	c       *http.Client
}

func NewClient() *Client {
	client := &Client{
		BaseURL: "https://blamewarrior.com",
		c:       http.DefaultClient,
	}

	return client
}

func (client *Client) ProcessPullRequest(pullRequest *bw.PullRequest) (err error) {

	repositoryName := pullRequest.RepositoryName

	requestUrl := fmt.Sprintf("%s/api/%s/pull_requests/process", client.BaseURL, repositoryName)

	response, err := client.c.Get(requestUrl)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Impossible to process hook for %s", repositoryName)
	}

	return nil
}
