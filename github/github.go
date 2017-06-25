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
package github

type GithubUser struct {
	Login     string `json:"login"`
	ID        string `json:"id"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
}

type PullRequestPayload struct {
	Number int    `json:"number"`
	Body   string `json:"body"`

	PullRequest struct {
		State string `json:"state"`
		Title string `json:"title"`

		User GithubUser `json:"user"`
	} `json:"pull_request"`

	RequestedReviewer GithubUser `json:"requested_reviewer"`
}
