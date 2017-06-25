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

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/blamewarrior"
	"github.com/blamewarrior/hooks/github"
	"github.com/bmizerany/pat"
)

type RepositoriesService interface {
	Track(ctx context.Context, repoFullName, callbackURL string) error
	Untrack(ctx context.Context, repoFullName, callbackURL string) error
}

type TrackingHandler struct {
	hostname string
}

func (handler *TrackingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	fullName := req.URL.Query().Get(":repo_full_name")
	username, _ := github.SplitRepositoryName(fullName)

	tokenService := blamewarrior.NewUsersClient(username)

	repositories := github.NewGithubRepositories(tokenService)

	trackingAction := req.URL.Query().Get(":action")

	err := handler.DoAction(repositories, fullName, trackingAction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}

}

func (handler *TrackingHandler) DoAction(repos RepositoriesService, repoFullName, action string) (err error) {

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

type HooksPayloadHandler struct{}

func (handler *HooksPayloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func (handler *HooksPayloadHandler) HandlePayload() error {
	fullName := req.URL.Query().Get(":repo_full_name")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	payload := new(github.PullRequestPayload)

	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return errors.New("Unable to unmarshal github hook: %s", body)
	}

	pullRequests = blamewarrior.NewPullRequestsClient()

	err = pullRequests.Handle(payload)

	if err != nil {
		return err
	}

	return nil
}

func main() {

	mux := pat.New()

	mux.Post("/:action/:repo_full_name", NewTrackingHandler("blamewarrior.com"))

	mux.Post("/:repo_full_name/webhook", new(HooksPayloadHandler))

	http.Handle("/", mux)

	log.Printf("blamewarrior users is running on 8080 port")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
