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
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	gh "github.com/google/go-github/github"
)

var (
	ErrRateLimitReached = errors.New("GitHub API request rate limit reached")
	ErrNoSuchRepository = errors.New("no such repository")
)

type tokenService interface {
	GetToken() (string, error)
}

type GithubRepositories struct {
	// BaseURL overrides GitHub API endpoint and is intended for use in tests.
	BaseURL *url.URL

	ts tokenService

	user string
}

// NewClient returns a new copy of github repositories service that uses given http.Client
// to make GitHub API requests.
func NewGithubRepositories(ts tokenService) *GithubRepositories {
	return &GithubRepositories{ts: ts}
}

// Tracks pull requests sets up "pull_request" event to be sent to callback
func (service *GithubRepositories) Track(ctx context.Context, repoFullName, callbackURL string) (err error) {
	owner, name := SplitRepositoryName(repoFullName)

	api, err := service.initAPIClient(ctx)

	if err != nil {
		return err
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

func (service *GithubRepositories) Untrack(ctx context.Context, repoFullName, callbackURL string) (err error) {
	owner, name := SplitRepositoryName(repoFullName)

	api, err := service.initAPIClient(ctx)

	if err != nil {
		return err
	}

	hooks, _, err := api.Repositories.ListHooks(ctx, owner, name, nil)

	if err != nil {
		return err
	}

	for _, hook := range hooks {

		configURL := hook.Config["url"]

		if strings.Index(*hook.URL, repoFullName) != -1 &&
			configURL == callbackURL {
			_, err = api.Repositories.DeleteHook(ctx, owner, name, *hook.ID)
			return err
		}
	}

	return fmt.Errorf("Hook not found")
}

func (service *GithubRepositories) initAPIClient(ctx context.Context) (*gh.Client, error) {

	token, err := service.ts.GetToken()

	if err != nil {
		return nil, fmt.Errorf("unable to get token to init API client: %s", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	api := gh.NewClient(oauthClient)
	if service.BaseURL != nil {
		api.BaseURL = service.BaseURL
	}

	return api, nil

}

// SplitRepositoryName splits full GitHub repository name into owner and name parts.
func SplitRepositoryName(fullName string) (owner, repo string) {
	sep := strings.IndexByte(fullName, '/')
	if sep <= 0 || sep == len(fullName)-1 {
		return "", ""
	}

	return fullName[0:sep], fullName[sep+1:]
}
