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
	"fmt"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/blamewarrior/tokens"
	"github.com/blamewarrior/hooks/github"
)

type TrackingHandler struct {
	hostname string
}

func (handler *TrackingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	username := req.URL.Query().Get(":username")
	repo := req.URL.Query().Get(":repo")

	fullName := fmt.Sprintf("%s/%s", username, repo)

	tokenClient := tokens.NewTokenClient("blamewarrior")

	repositories := github.NewGithubRepositories(tokenClient)

	trackingAction := req.URL.Query().Get(":action")

	err := handler.DoAction(repositories, fullName, trackingAction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}

}

func (handler *TrackingHandler) DoAction(repos github.Repositories, repoFullName, action string) (err error) {

	switch action {
	case "track":
		err = repos.Track(
			github.Context{},
			repoFullName,
			fmt.Sprintf("https://%s/%s/webhook", handler.hostname, repoFullName),
		)
		return
	case "untrack":
		err = repos.Untrack(
			github.Context{},
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
