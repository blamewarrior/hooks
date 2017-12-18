package collaborators_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blamewarrior/hooks/blamewarrior/collaborators"

	gh "github.com/blamewarrior/hooks/github"
)

func TestGetCollaborators(t *testing.T) {
	results := []struct {
		ResponseStatus int
		Collaborators  []gh.Collaborator
		ResponseError  error
	}{
		{ResponseStatus: http.StatusOK, Collaborators: []gh.Collaborator{{Id: 123, Admin: true, Login: "test_user"}}, ResponseError: nil},
		{ResponseStatus: http.StatusNotFound, Collaborators: nil, ResponseError: errors.New(
			"Unable to get collaborators for blamewarrior/test_repo, status_code=404",
		)},
	}

	for _, result := range results {
		testAPIEndpoint, mux, teardown := setup()

		mux.HandleFunc("/blamewarrior/test_repo/collaborators", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(result.ResponseStatus)
			fmt.Fprint(w, `[{
                    "id": 123,
                    "admin": true,
                    "login": "test_user"
                }]`)
		})

		client := collaborators.NewClient()
		client.BaseURL = testAPIEndpoint

		collaborators, err := client.GetCollaborators("blamewarrior/test_repo")

		assert.Equal(t, result.ResponseError, err)
		assert.Equal(t, result.Collaborators, collaborators)

		teardown()
	}
}

func setup() (baseURL string, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}
