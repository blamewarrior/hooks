package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/blamewarrior"
	"github.com/blamewarrior/hooks/github"
)

type HooksPayloadHandler struct{}

func (handler *HooksPayloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := handler.handlePayload(w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}
}

func (handler *HooksPayloadHandler) handlePayload(w http.ResponseWriter, req *http.Request) error {
	username := req.URL.Query().Get(":username")
	repo := req.URL.Query().Get(":repo")

	fullName := fmt.Sprintf("%s/%s", username, repo)

	respBytes, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))

	if err != nil {
		return err
	}

	if err := req.Body.Close(); err != nil {
		return err
	}

	event := req.Header.Get("X-GitHub-Event")

	switch event {
	case "pull_request":
		ghPullRequestHook := new(github.GithubPullRequestHook)
		err = json.Unmarshal(respBytes, &ghPullRequestHook)

		hookRepositoryName := ghPullRequestHook.Repository.FullName

		if hookRepositoryName != fullName {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("%s\t%s\t%v\t%s",
				"POST",
				req.RequestURI,
				http.StatusBadRequest,
				fmt.Sprintf("hooks repository name doesn't match with url pattern: %s != %s",
					hookRepositoryName, fullName),
			)
		}

		pullRequest := &blamewarrior.PullRequest{
			Id:             *ghPullRequestHook.PullRequest.ID,
			HTMLURL:        *ghPullRequestHook.PullRequest.HTMLURL,
			Title:          *ghPullRequestHook.PullRequest.Title,
			Body:           *ghPullRequestHook.PullRequest.Body,
			Number:         *ghPullRequestHook.PullRequest.Number,
			State:          *ghPullRequestHook.PullRequest.State,
			CreatedAt:      ghPullRequestHook.PullRequest.CreatedAt,
			ClosedAt:       ghPullRequestHook.PullRequest.ClosedAt,
			Commits:        *ghPullRequestHook.PullRequest.Commits,
			Additions:      *ghPullRequestHook.PullRequest.Additions,
			Deletions:      *ghPullRequestHook.PullRequest.Deletions,
			RepositoryName: hookRepositoryName,
			OwnerId:        *ghPullRequestHook.PullRequest.User.ID,
		}

		if ghPullRequestHook.RequestedReviewer != nil {
			pullRequest.ReviewerIds = []int{ghPullRequestHook.RequestedReviewer.Id}
		}

		if err != nil {
			return err
		}
	}
	return nil
}
