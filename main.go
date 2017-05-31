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

type UsersService interface {
	GetTokenFor(nickname string) (string, error)
}

type RepositoriesService interface {
	Track(ctx context.Context, repoFullName, callbackURL string) error
	Untrack(ctx context.Context, repoFullName, callbackURL string) error
}

type RepositoryHandlers struct {
	repositories RepositoriesService
}

func (handlers *RepositoryHandlers) Track(w http.ResponseWriter, req *http.Request) {
	token, err := handlers.users.GetToken(username)

	if err != nil {

	}

	repositories := github.NewGithubRepositories(token)

	err = repositories.Track(
		context.Background(),
		repoFullName,
		fmt.Sprintf("https://blamewarrior.com/%s/webhook", repoFullName),
	)

	if err != nil {

	}

}

func (handlers *RepositoryHandlers) Untrack(w http.ResponseWriter, req *http.Request) {
	token, err := handlers.users.GetToken(username)
}

func NewRepositoryHandlers(usersClient *users.Client) {
	return &RepositoryHandlers{usersClient}
}

func main() {

	usersClient := users.NewClient()

	repositoryHandlers := NewRepositoryHandlers(usersClient)

	mux := pat.New()

	mux.Post("/track/:repository_id", http.HandlerFunc(repositoryHandlers.Track))
	mux.Post("/untrack/:repository_id", http.HandlerFunc(repositoryHandlers.Untrack))

	http.Handle("/", mux)

	log.Printf("blamewarrior users is running on 8080 port")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
