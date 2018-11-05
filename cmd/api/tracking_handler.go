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

	"github.com/blamewarrior/hooks/blamewarrior/collaborators"
	"github.com/blamewarrior/hooks/github"
	"github.com/go-redis/redis"
)

type TrackingHandler struct {
	hostname string

	redisClient *redis.Client

	repositories  github.Repositories
	collaborators collaborators.Client
}

func (handler *TrackingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	username := req.URL.Query().Get(":username")
	repo := req.URL.Query().Get(":repo")

	fullName := fmt.Sprintf("%s/%s", username, repo)

	trackingAction := req.URL.Query().Get(":action")

	err := handler.DoAction(fullName, trackingAction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}

}

func (handler *TrackingHandler) DoAction(repoFullName, action string) (err error) {

	switch action {
	case "track":
		err = handler.collaborators.FetchCollaborators(repoFullName)

		if err != nil {
			if err = handler.redisClient.LPush(
				"fetch_collaborators",
				fmt.Sprintf(`{"type":"fetch", repo_full_name:"%s"}`, repoFullName),
			).Err(); err != nil {
				return err
			}
			return err
		}

		err = handler.repositories.Track(
			github.Context{},
			repoFullName,
			fmt.Sprintf("https://%s/%s/webhook", handler.hostname, repoFullName),
		)
		return
	case "untrack":
		err = handler.repositories.Untrack(
			github.Context{},
			repoFullName,
			fmt.Sprintf("https://%s/%s/webhook", handler.hostname, repoFullName),
		)

		return
	default:
		return fmt.Errorf("Unsupported action %s", action)
	}
}

func NewTrackingHandler(hostname string, repositories github.Repositories, redisClient *redis.Client, collaborators collaborators.Client) *TrackingHandler {
	return &TrackingHandler{
		hostname:      hostname,
		repositories:  repositories,
		redisClient:   redisClient,
		collaborators: collaborators,
	}
}
