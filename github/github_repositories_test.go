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

package github_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blamewarrior/hooks/github"

	api "github.com/google/go-github/github"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackRepositoryPullRequests(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/blamewarrior/hooks/hooks", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, "POST")

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var hook api.Hook
		require.NoError(t, json.Unmarshal(body, &hook), string(body))

		assert.Equal(t, *hook.Name, "web")
		assert.Contains(t, hook.Events, "pull_request")
		assert.Equal(t, hook.Config["url"], "https://example.com/blamewarrior/hooks/webhook")
		assert.True(t, *hook.Active)

		fmt.Fprint(w, `{"id":1}`)
	})

	githubRepos := github.NewGithubRepositories("token")

	githubRepos.BaseURL = baseURL

	callbackURL := "https://example.com/blamewarrior/hooks/webhook"

	err := githubRepos.Track(context.Background(), "blamewarrior/hooks", callbackURL)
	require.NoError(t, err)

}

func TestUntrackRepositoryPullRequests(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/blamewarrior/hooks/hooks", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, "GET")

		fmt.Fprint(w, hooksResponse)
	})

	mux.HandleFunc("/repos/blamewarrior/hooks/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, "DELETE")
	})

	githubRepos := github.NewGithubRepositories("token")

	githubRepos.BaseURL = baseURL

	callbackURL := "https://example.com/blamewarrior/hooks/webhook"

	err := githubRepos.Untrack(context.Background(), "blamewarrior/hooks", callbackURL)
	require.NoError(t, err)

}

func setup() (mux *http.ServeMux, baseURL *url.URL, teardownFn func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)

	return mux, url, server.Close
}

const (
	hooksResponse = `[
  {
    "id": 1,
    "url": "https://api.github.com/repos/blamewarrior/hooks/hooks/1",
    "test_url": "https://api.github.com/repos/blamewarrior/hooks/hooks/1/test",
    "ping_url": "https://api.github.com/repos/blamewarrior/hooks/hooks/1/pings",
    "name": "web",
    "events": [
      "pull_request"
    ],
    "active": true,
    "config": {
      "url": "https://example.com/blamewarrior/hooks/webhook",
      "content_type": "json"
    },
    "updated_at": "2011-09-06T20:39:23Z",
    "created_at": "2011-09-06T17:26:27Z"
  }
]`
)
