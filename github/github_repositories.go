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
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	gh "github.com/google/go-github/github"
)

var (
	ErrRateLimitReached = errors.New("GitHub API request rate limit reached")
	ErrNoSuchRepository = errors.New("no such repository")
)

type GithubRepositories struct {
	// BaseURL overrides GitHub API endpoint and is intended for use in tests.
	BaseURL *url.URL

	token      string
	httpClient *http.Client
}

// NewClient returns a new copy of github repositories service that uses given http.Client
// to make GitHub API requests.
func NewGithubRepositories(httpClient *http.Client, token string) *GithubRepositories {
	return &GithubRepositories{httpClient: httpClient, token: token}
}

// Tracks pull requests sets up "pull_request" event to be sent to callback
func (service *GithubRepositories) TrackPullRequests(repoFullName, callbackURL string) (err error) {
	owner, name := SplitRepositoryName(repoFullName)

	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: service.token})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	api := gh.NewClient(httpClient)
	if service.BaseURL != nil {
		api.BaseURL = service.BaseURL
	}

	hook := &gh.Hook{
		Name:   new(string),
		Active: new(bool),
		Events: []string{"pull_request"},
		Config: map[string]interface{}{
			"url":          callbackURL,
			"content_type": "json",
		},
	}

	*hook.Name = "web"
	*hook.Active = true

	_, _, err = api.Repositories.CreateHook(ctx, owner, name, hook)
	return err
}

// SplitRepositoryName splits full GitHub repository name into owner and name parts.
func SplitRepositoryName(fullName string) (owner, repo string) {
	sep := strings.IndexByte(fullName, '/')
	if sep <= 0 || sep == len(fullName)-1 {
		return "", ""
	}

	return fullName[0:sep], fullName[sep+1:]
}
