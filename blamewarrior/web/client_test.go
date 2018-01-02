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

package web_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	bw "github.com/blamewarrior/hooks/blamewarrior"
	"github.com/blamewarrior/hooks/blamewarrior/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gh "github.com/blamewarrior/hooks/github"
)

func TestProcessPullRequest(t *testing.T) {

	results := []struct {
		ResponseStatus int
		ResponseError  error
	}{
		{ResponseStatus: http.StatusNoContent, ResponseError: nil},
		{ResponseStatus: http.StatusNotFound, ResponseError: errors.New("Impossible to process hook for blamewarrior/test_repo, status_code=404")},
	}

	for _, result := range results {
		testAPIEndpoint, mux, teardown := setup()

		mux.HandleFunc("/api/blamewarrior/test_repo/pull_requests/process", func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)

			bodyBytes, _ := ioutil.ReadAll(r.Body)

			require.Equal(t,
				"{\"id\":123,\"html_url\":\"\",\"title\":\"bug fixes\",\"body\":\"\",\"repository_name\":\"blamewarrior/test_repo\",\"reviewers\":[{\"id\":2,\"login\":\"test_user\",\"admin\":true}],\"number\":12,\"state\":\"open\",\"opened_at\":null,\"closed_at\":null,\"owner_id\":1,\"commits\":0,\"additions\":0,\"deletions\":0,\"review_comments\":null}",
				string(bodyBytes))

			w.WriteHeader(result.ResponseStatus)
		})

		client := web.NewClient()
		client.BaseURL = testAPIEndpoint

		pullRequest := &bw.PullRequest{
			Id:             123,
			OwnerId:        1,
			Number:         12,
			Title:          "bug fixes",
			State:          "open",
			RepositoryName: "blamewarrior/test_repo",
			Reviewers: []gh.Collaborator{
				{Id: 2, Admin: true, Login: "test_user"},
			},
		}

		err := client.ProcessPullRequest(pullRequest)

		assert.Equal(t, result.ResponseError, err)

		teardown()
	}
}

func setup() (baseURL string, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}
