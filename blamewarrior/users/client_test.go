package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blamewarrior/hooks/blamewarrior/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetToken(t *testing.T) {
	testAPIEndpoint, mux, teardown := setup()

	defer teardown()

	mux.HandleFunc("/users/blamewarrior", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(userResponse))
	})

	client := users.NewClient()
	client.BaseURL = testAPIEndpoint

	token, err := client.GetTokenFor("blamewarrior")

	require.NoError(t, err)

	assert.Equal(t, "test_token", token)

}

func setup() (baseURL string, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}

const userResponse = `{"token": "test_token"}`
