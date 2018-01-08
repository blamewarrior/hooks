package collaborators_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blamewarrior/hooks/blamewarrior/collaborators"

	gh "github.com/blamewarrior/hooks/github"
)

func TestListCollaborator(t *testing.T) {
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

		collaborators, err := client.ListCollaborator("blamewarrior/test_repo")

		assert.Equal(t, result.ResponseError, err)
		assert.Equal(t, result.Collaborators, collaborators)

		teardown()
	}
}

func TestAddCollaborator(t *testing.T) {
	results := []struct {
		ResponseStatus int
		Collaborator   *gh.Collaborator
		ResponseError  error
	}{
		{ResponseStatus: http.StatusCreated, Collaborator: &gh.Collaborator{Id: 123, Admin: true, Login: "test_user"}, ResponseError: nil},
		{ResponseStatus: http.StatusInternalServerError, Collaborator: nil, ResponseError: errors.New(
			"Unable to add collaborator for blamewarrior/test_repo, status_code=500",
		)},
	}

	for _, result := range results {
		testAPIEndpoint, mux, teardown := setup()

		mux.HandleFunc("/blamewarrior/test_repo/collaborators", func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)
			w.WriteHeader(result.ResponseStatus)
			fmt.Fprint(w, `[{
                    "id": 123,
                    "admin": true,
                    "login": "test_user"
                }]`)
		})

		client := collaborators.NewClient()
		client.BaseURL = testAPIEndpoint

		err := client.AddCollaborator("blamewarrior/test_repo", result.Collaborator)
		assert.Equal(t, result.ResponseError, err)

		teardown()
	}
}

func TestEditCollaborator(t *testing.T) {
	results := []struct {
		ResponseStatus int
		Collaborator   *gh.Collaborator
		ResponseError  error
	}{
		{ResponseStatus: http.StatusOK, Collaborator: &gh.Collaborator{Id: 123, Admin: true, Login: "test_user"}, ResponseError: nil},
		{ResponseStatus: http.StatusInternalServerError, Collaborator: nil, ResponseError: errors.New(
			"Unable to edit collaborator for blamewarrior/test_repo, status_code=500",
		)},
	}

	for _, result := range results {
		testAPIEndpoint, mux, teardown := setup()

		mux.HandleFunc("/blamewarrior/test_repo/collaborators", func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "PUT", r.Method)
			w.WriteHeader(result.ResponseStatus)
			fmt.Fprint(w, `[{
                    "id": 123,
                    "admin": true,
                    "login": "test_user"
                }]`)
		})

		client := collaborators.NewClient()
		client.BaseURL = testAPIEndpoint

		err := client.EditCollaborator("blamewarrior/test_repo", result.Collaborator)
		assert.Equal(t, result.ResponseError, err)

		teardown()
	}
}

func TestDeleteCollaborator(t *testing.T) {
	results := []struct {
		ResponseStatus int
		Collaborator   *gh.Collaborator
		ResponseError  error
	}{
		{ResponseStatus: http.StatusNoContent, Collaborator: &gh.Collaborator{Id: 123, Admin: true, Login: "test_user"}, ResponseError: nil},
		{ResponseStatus: http.StatusInternalServerError, Collaborator: &gh.Collaborator{Id: 123, Admin: true, Login: "test_user"}, ResponseError: errors.New(
			"Unable to delete collaborator for blamewarrior/test_repo, status_code=500",
		)},
	}

	for _, result := range results {
		testAPIEndpoint, mux, teardown := setup()

		mux.HandleFunc("/blamewarrior/test_repo/collaborators/"+result.Collaborator.Login, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "DELETE", r.Method)
			w.WriteHeader(result.ResponseStatus)
		})

		client := collaborators.NewClient()
		client.BaseURL = testAPIEndpoint

		err := client.DeleteCollaborator("blamewarrior/test_repo", result.Collaborator.Login)
		assert.Equal(t, result.ResponseError, err)

		teardown()
	}
}

func setup() (baseURL string, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}
