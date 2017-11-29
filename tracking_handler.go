package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/github"
)

type TrackingHandler struct {
	hostname string
}

func (handler *TrackingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	username := req.URL.Query().Get(":username")
	repo := req.URL.Query().Get(":repo")

	fullName := fmt.Sprintf("%s/%s", username, repo)

	token := req.Header.Get("X-Token")

	if token == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Invalid token")
		return
	}

	repositories := github.NewGithubRepositories(token)

	trackingAction := req.URL.Query().Get(":action")

	err := handler.DoAction(repositories, fullName, trackingAction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}

}

func (handler *TrackingHandler) DoAction(repos github.RepositoriesService, repoFullName, action string) (err error) {

	switch action {
	case "track":
		err = repos.Track(
			context.Background(),
			repoFullName,
			fmt.Sprintf("https://%s/%s/webhook", handler.hostname, repoFullName),
		)
		return
	case "untrack":
		err = repos.Untrack(
			context.Background(),
			repoFullName,
			fmt.Sprintf("https://%s/%s/webhook", handler.hostname, repoFullName),
		)

		return
	default:
		return fmt.Errorf("Unsupported action %s", action)
	}
}

func NewTrackingHandler(hostname string) *TrackingHandler {
	return &TrackingHandler{
		hostname: hostname,
	}
}
