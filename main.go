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
	"fmt"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/blamewarrior/users"
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

	tokenService := users.NewClient(username)

	repositories := github.NewGithubRepositories(tokenService)

	trackingAction := req.URL.Query().Get(":action")

	err := handler.DoAction(repositories, fullName, trackingAction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "GET", req.RequestURI, http.StatusInternalServerError, err)
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

func main() {

	mux := pat.New()

	mux.Post("/:action/:repo_full_name", NewTrackingHandler("blamewarrior.com"))

	http.Handle("/", mux)

	log.Printf("blamewarrior users is running on 8080 port")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
