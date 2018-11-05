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
	"os"

	"github.com/blamewarrior/hooks"
	"github.com/blamewarrior/hooks/blamewarrior/collaborators"
	"github.com/blamewarrior/hooks/blamewarrior/tokens"
	"github.com/blamewarrior/hooks/blamewarrior/web"
	"github.com/blamewarrior/hooks/github"
	"github.com/bmizerany/pat"
	"github.com/go-redis/redis"
)

func main() {
	mux := pat.New()

	tokenClient := tokens.NewTokenClient()

	repositories := github.NewGithubRepositories(tokenClient)

	bwHost := os.Getenv("BW_HOST")
	if bwHost == "" {
		log.Fatal("missing bw host (expected to be passed via ENV['BW_HOST'])")
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	redisClient := redis.NewClient(opts)
	collaboratorsClient := collaborators.NewClient()

	mux.Post("/:action/:username/:repo", NewTrackingHandler(bwHost, repositories, redisClient, collaboratorsClient))

	payloadRepo := hooks.NewPayloadRepository(redisClient)

	webClient := web.NewClient()

	reviewersService := github.NewGithubReviewers(tokenClient)

	mediator := hooks.NewMediatorService(payloadRepo, webClient, collaboratorsClient, reviewersService)

	mux.Post("/:username/:repo/webhook", NewHooksPayloadHandler(mediator))

	http.Handle("/", mux)

	log.Printf("blamewarrior users is running on 8080 port")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
