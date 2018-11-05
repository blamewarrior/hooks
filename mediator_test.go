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

package hooks_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/blamewarrior/hooks"
	bw "github.com/blamewarrior/hooks/blamewarrior"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gh "github.com/blamewarrior/hooks/github"
)

type PayloadServiceMock struct {
	mock.Mock
}

func (m *PayloadServiceMock) Save(payload string) (err error) {
	args := m.Called()
	return args.Error(0)
}

func (m *PayloadServiceMock) List(limit int64) (result []string, err error) {
	return
}

func (m *PayloadServiceMock) Delete(payload string) (err error) {
	return
}

type WebClientMock struct {
	mock.Mock
}

func (m *WebClientMock) ProcessPullRequest(pullRequest *bw.PullRequest) (err error) {
	args := m.Called(pullRequest)
	return args.Error(0)
}

type CollaboratorsClientMock struct {
	mock.Mock
}

func (m *CollaboratorsClientMock) FetchCollaborators(repositoryFullName string) error {
	args := m.Called(repositoryFullName)
	return args.Error(0)
}

func (m *CollaboratorsClientMock) ListCollaborator(repositoryFullName string) ([]gh.Collaborator, error) {
	args := m.Called(repositoryFullName)
	return args.Get(0).([]gh.Collaborator), args.Error(1)
}

func (m *CollaboratorsClientMock) AddCollaborator(repositoryFullName string, collaborator *gh.Collaborator) error {
	args := m.Called(repositoryFullName, collaborator)
	return args.Error(0)
}

func (m *CollaboratorsClientMock) EditCollaborator(repositoryFullName string, collaborator *gh.Collaborator) error {
	args := m.Called(repositoryFullName, collaborator)
	return args.Error(0)
}

func (m *CollaboratorsClientMock) DeleteCollaborator(repositoryFullName, login string) error {
	args := m.Called(repositoryFullName, login)
	return args.Error(0)
}

type ReviewersServiceMock struct {
	mock.Mock
}

func (m *ReviewersServiceMock) RequestReviewers(ctx gh.Context, repoFullName string, pullNumber int, reviewers []gh.Collaborator) (err error) {
	args := m.Called(ctx, repoFullName, pullNumber, reviewers)
	return args.Error(0)
}

func (m *ReviewersServiceMock) ReviewComments(ctx gh.Context, repoFullName string, pullNumber int) ([]gh.ReviewComment, error) {
	args := m.Called(ctx, repoFullName, pullNumber)
	return args.Get(0).([]gh.ReviewComment), args.Error(1)
}

func TestHooksMediator_Mediate_HooksWithAssignedReviewers(t *testing.T) {

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	webClientMock := new(WebClientMock)

	// layout := "2006-01-02T15:04:00Z"
	str := "2015-05-05T23:40:27Z"
	createdAt, err := time.Parse(time.RFC3339, str)

	require.NoError(t, err)

	pullRequest := &bw.PullRequest{Id: 34778301,
		HTMLURL:        "https://github.com/blamewarrior_user/public-repo/pull/1",
		Title:          "Update the README with new information",
		Body:           "This is a pretty simple change that we need to pull into master.",
		RepositoryName: "blamewarrior_user/public-repo",
		Reviewers: []gh.Collaborator{gh.Collaborator{Id: 6752318,
			Login: "blamewarrior_second_user",
			Admin: false}},
		Number:    1,
		State:     "open",
		CreatedAt: &createdAt,
		OwnerId:   6752317,
		Commits:   1,
		Additions: 1,
		Deletions: 1,
	}

	webClientMock.On("ProcessPullRequest", pullRequest).Return(nil)

	reviewersService := new(ReviewersServiceMock)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("pull_request", []byte(pullRequestHookPayloadWithAssignedReviewers))

	webClientMock.AssertExpectations(t)

}

func TestHooksMediator_Mediate_HooksWithoutAssignedReviewers(t *testing.T) {

	collaborators := []gh.Collaborator{
		{
			Id:    123,
			Login: "admin_user",
			Admin: true,
		},
	}

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	collaboratorsClientMock.On("ListCollaborator", "blamewarrior_user/public-repo").Return(collaborators, nil)

	webClientMock := new(WebClientMock)

	// layout := "2006-01-02T15:04:00Z"
	str := "2015-05-05T23:40:27Z"
	createdAt, err := time.Parse(time.RFC3339, str)

	require.NoError(t, err)

	pullRequest := &bw.PullRequest{Id: 34778301,
		HTMLURL:        "https://github.com/blamewarrior_user/public-repo/pull/1",
		Title:          "Update the README with new information",
		Body:           "This is a pretty simple change that we need to pull into master.",
		RepositoryName: "blamewarrior_user/public-repo",
		Reviewers: []gh.Collaborator{gh.Collaborator{
			Id:    123,
			Login: "admin_user",
			Admin: true},
		},
		Number:    1,
		State:     "open",
		CreatedAt: &createdAt,
		OwnerId:   6752317,
		Commits:   1,
		Additions: 1,
		Deletions: 1,
	}

	webClientMock.On("ProcessPullRequest", pullRequest).Return(nil)

	reviewersService := new(ReviewersServiceMock)
	reviewersService.On("RequestReviewers",
		gh.Context{},
		"blamewarrior_user/public-repo",
		1,
		collaborators,
	).Return(nil)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("pull_request", []byte(pullRequestHookPayloadWithoutAssignedReviewers))

	reviewersService.AssertExpectations(t)

}

func TestHooksMediator_Mediate_ClosedPullRequest(t *testing.T) {
	commentBody := "great stuff"
	comments := []gh.ReviewComment{
		{Body: &commentBody},
	}

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	webClientMock := new(WebClientMock)

	// layout := "2006-01-02T15:04:00Z"
	str := "2015-05-05T23:40:27Z"
	createdAt, err := time.Parse(time.RFC3339, str)

	require.NoError(t, err)

	pullRequest := &bw.PullRequest{Id: 34778301,
		HTMLURL:        "https://github.com/blamewarrior_user/public-repo/pull/1",
		Title:          "Update the README with new information",
		Body:           "This is a pretty simple change that we need to pull into master.",
		RepositoryName: "blamewarrior_user/public-repo",
		Reviewers: []gh.Collaborator{gh.Collaborator{Id: 6752318,
			Login: "blamewarrior_second_user",
			Admin: false}},
		Number:         1,
		State:          "closed",
		CreatedAt:      &createdAt,
		OwnerId:        6752317,
		Commits:        1,
		Additions:      1,
		Deletions:      1,
		ReviewComments: comments[0:],
	}

	webClientMock.On("ProcessPullRequest", pullRequest).Return(nil)

	reviewersService := new(ReviewersServiceMock)
	reviewersService.On("ReviewComments", gh.Context{}, "blamewarrior_user/public-repo", 1).Return(comments[0:], nil)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("pull_request", []byte(closedPullRequestHookPayload))

	reviewersService.AssertExpectations(t)
}

func TestMediator_Mediate_AddColaborator(t *testing.T) {

	collaborator := &gh.Collaborator{
		Id:    583231,
		Login: "octocat",
		Admin: false,
	}

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	collaboratorsClientMock.On("AddCollaborator", "baxterthehacker/public-repo", collaborator).Return(nil)

	webClientMock := new(WebClientMock)

	reviewersService := new(ReviewersServiceMock)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("member", []byte(fmt.Sprintf(pullRequestPayloadWithMember, "added")))

	collaboratorsClientMock.AssertExpectations(t)
}

func TestMediator_Mediate_EditColaborator(t *testing.T) {

	collaborator := &gh.Collaborator{
		Id:    583231,
		Login: "octocat",
		Admin: false,
	}

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	collaboratorsClientMock.On("EditCollaborator", "baxterthehacker/public-repo", collaborator).Return(nil)

	webClientMock := new(WebClientMock)

	reviewersService := new(ReviewersServiceMock)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("member", []byte(fmt.Sprintf(pullRequestPayloadWithMember, "edited")))

	collaboratorsClientMock.AssertExpectations(t)
}

func TestMediator_Mediate_DeleteColaborator(t *testing.T) {

	collaborator := &gh.Collaborator{
		Id:    583231,
		Login: "octocat",
		Admin: false,
	}

	payloadServiceMock := new(PayloadServiceMock)

	payloadServiceMock.On("Save", "").Return(nil)

	collaboratorsClientMock := new(CollaboratorsClientMock)

	collaboratorsClientMock.On("DeleteCollaborator", "baxterthehacker/public-repo", collaborator.Login).Return(nil)

	webClientMock := new(WebClientMock)

	reviewersService := new(ReviewersServiceMock)

	m := hooks.NewMediatorService(payloadServiceMock, webClientMock, collaboratorsClientMock, reviewersService)
	m.Mediate("member", []byte(fmt.Sprintf(pullRequestPayloadWithMember, "deleted")))

	collaboratorsClientMock.AssertExpectations(t)
}

const (
	pullRequestHookPayloadWithAssignedReviewers = `{
  "action": "opened",
  "number": 1,
  "pull_request": {
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1",
    "id": 34778301,
    "html_url": "https://github.com/blamewarrior_user/public-repo/pull/1",
    "diff_url": "https://github.com/blamewarrior_user/public-repo/pull/1.diff",
    "patch_url": "https://github.com/blamewarrior_user/public-repo/pull/1.patch",
    "issue_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1",
    "number": 1,
    "state": "open",
    "locked": false,
    "title": "Update the README with new information",
    "user": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "body": "This is a pretty simple change that we need to pull into master.",
    "created_at": "2015-05-05T23:40:27Z",
    "updated_at": "2015-05-05T23:40:27Z",
    "closed_at": null,
    "merged_at": null,
    "merge_commit_sha": null,
    "assignee": null,
    "milestone": null,
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits",
    "review_comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments",
    "review_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
    "head": {
      "label": "blamewarrior_user:changes",
      "ref": "changes",
      "sha": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "base": {
      "label": "blamewarrior_user:master",
      "ref": "master",
      "sha": "9049f1265b7d61be4a8904a9a27120d2064dab3b",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "_links": {
      "self": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1"
      },
      "html": {
        "href": "https://github.com/blamewarrior_user/public-repo/pull/1"
      },
      "issue": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1"
      },
      "comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments"
      },
      "review_comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments"
      },
      "review_comment": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}"
      },
      "commits": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits"
      },
      "statuses": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c"
      }
    },
    "merged": false,
    "mergeable": null,
    "mergeable_state": "unknown",
    "merged_by": null,
    "comments": 0,
    "review_comments": 0,
    "commits": 1,
    "additions": 1,
    "deletions": 1,
    "changed_files": 1
  },
  "repository": {
    "id": 35129377,
    "name": "public-repo",
    "full_name": "blamewarrior_user/public-repo",
    "owner": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "private": false,
    "html_url": "https://github.com/blamewarrior_user/public-repo",
    "description": "",
    "fork": false,
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
    "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
    "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
    "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
    "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
    "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
    "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
    "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
    "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
    "created_at": "2015-05-05T23:40:12Z",
    "updated_at": "2015-05-05T23:40:12Z",
    "pushed_at": "2015-05-05T23:40:26Z",
    "git_url": "git://github.com/blamewarrior_user/public-repo.git",
    "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
    "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
    "svn_url": "https://github.com/blamewarrior_user/public-repo",
    "homepage": null,
    "size": 0,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": null,
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": true,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 1,
    "forks": 0,
    "open_issues": 1,
    "watchers": 0,
    "default_branch": "master"
  },
  "requested_reviewers": [{
    "login": "blamewarrior_second_user",
    "id": 6752318,
    "avatar_url": "https://avatars.githubusercontent.com/u/6752318?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/blamewarrior_second_user",
    "html_url": "https://github.com/blamewarrior_second_user",
    "followers_url": "https://api.github.com/users/blamewarrior_second_user/followers",
    "following_url": "https://api.github.com/users/blamewarrior_second_user/following{/other_user}",
    "gists_url": "https://api.github.com/users/blamewarrior_second_user/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/blamewarrior_second_user/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/blamewarrior_second_user/subscriptions",
    "organizations_url": "https://api.github.com/users/blamewarrior_second_user/orgs",
    "repos_url": "https://api.github.com/users/blamewarrior_second_user/repos",
    "events_url": "https://api.github.com/users/blamewarrior_second_user/events{/privacy}",
    "received_events_url": "https://api.github.com/users/blamewarrior_second_user/received_events",
    "type": "User",
    "site_admin": false
  }],
  "installation": {
    "id": 234
  }
}`

	pullRequestHookPayloadWithoutAssignedReviewers = `{
  "action": "opened",
  "number": 1,
  "pull_request": {
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1",
    "id": 34778301,
    "html_url": "https://github.com/blamewarrior_user/public-repo/pull/1",
    "diff_url": "https://github.com/blamewarrior_user/public-repo/pull/1.diff",
    "patch_url": "https://github.com/blamewarrior_user/public-repo/pull/1.patch",
    "issue_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1",
    "number": 1,
    "state": "open",
    "locked": false,
    "title": "Update the README with new information",
    "user": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "body": "This is a pretty simple change that we need to pull into master.",
    "created_at": "2015-05-05T23:40:27Z",
    "updated_at": "2015-05-05T23:40:27Z",
    "closed_at": null,
    "merged_at": null,
    "merge_commit_sha": null,
    "assignee": null,
    "milestone": null,
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits",
    "review_comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments",
    "review_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
    "head": {
      "label": "blamewarrior_user:changes",
      "ref": "changes",
      "sha": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "base": {
      "label": "blamewarrior_user:master",
      "ref": "master",
      "sha": "9049f1265b7d61be4a8904a9a27120d2064dab3b",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "_links": {
      "self": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1"
      },
      "html": {
        "href": "https://github.com/blamewarrior_user/public-repo/pull/1"
      },
      "issue": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1"
      },
      "comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments"
      },
      "review_comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments"
      },
      "review_comment": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}"
      },
      "commits": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits"
      },
      "statuses": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c"
      }
    },
    "merged": false,
    "mergeable": null,
    "mergeable_state": "unknown",
    "merged_by": null,
    "comments": 0,
    "review_comments": 0,
    "commits": 1,
    "additions": 1,
    "deletions": 1,
    "changed_files": 1
  },
  "repository": {
    "id": 35129377,
    "name": "public-repo",
    "full_name": "blamewarrior_user/public-repo",
    "owner": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "private": false,
    "html_url": "https://github.com/blamewarrior_user/public-repo",
    "description": "",
    "fork": false,
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
    "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
    "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
    "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
    "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
    "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
    "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
    "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
    "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
    "created_at": "2015-05-05T23:40:12Z",
    "updated_at": "2015-05-05T23:40:12Z",
    "pushed_at": "2015-05-05T23:40:26Z",
    "git_url": "git://github.com/blamewarrior_user/public-repo.git",
    "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
    "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
    "svn_url": "https://github.com/blamewarrior_user/public-repo",
    "homepage": null,
    "size": 0,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": null,
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": true,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 1,
    "forks": 0,
    "open_issues": 1,
    "watchers": 0,
    "default_branch": "master"
  },
  "requested_reviewers": [],
  "installation": {
    "id": 234
  }
}`

	closedPullRequestHookPayload = `{
  "action": "closed",
  "number": 1,
  "pull_request": {
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1",
    "id": 34778301,
    "html_url": "https://github.com/blamewarrior_user/public-repo/pull/1",
    "diff_url": "https://github.com/blamewarrior_user/public-repo/pull/1.diff",
    "patch_url": "https://github.com/blamewarrior_user/public-repo/pull/1.patch",
    "issue_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1",
    "number": 1,
    "state": "closed",
    "locked": false,
    "title": "Update the README with new information",
    "user": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "body": "This is a pretty simple change that we need to pull into master.",
    "created_at": "2015-05-05T23:40:27Z",
    "updated_at": "2015-05-05T23:40:27Z",
    "closed_at": null,
    "merged_at": null,
    "merge_commit_sha": null,
    "assignee": null,
    "milestone": null,
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits",
    "review_comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments",
    "review_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
    "head": {
      "label": "blamewarrior_user:changes",
      "ref": "changes",
      "sha": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "base": {
      "label": "blamewarrior_user:master",
      "ref": "master",
      "sha": "9049f1265b7d61be4a8904a9a27120d2064dab3b",
      "user": {
        "login": "blamewarrior_user",
        "id": 6752317,
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "gravatar_id": "",
        "url": "https://api.github.com/users/blamewarrior_user",
        "html_url": "https://github.com/blamewarrior_user",
        "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
        "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
        "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
        "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
        "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
        "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
        "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
        "type": "User",
        "site_admin": false
      },
      "repo": {
        "id": 35129377,
        "name": "public-repo",
        "full_name": "blamewarrior_user/public-repo",
        "owner": {
          "login": "blamewarrior_user",
          "id": 6752317,
          "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
          "gravatar_id": "",
          "url": "https://api.github.com/users/blamewarrior_user",
          "html_url": "https://github.com/blamewarrior_user",
          "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
          "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
          "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
          "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
          "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
          "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
          "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
          "type": "User",
          "site_admin": false
        },
        "private": false,
        "html_url": "https://github.com/blamewarrior_user/public-repo",
        "description": "",
        "fork": false,
        "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
        "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
        "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
        "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
        "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
        "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
        "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
        "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
        "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
        "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
        "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
        "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
        "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
        "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
        "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
        "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
        "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
        "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
        "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
        "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
        "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
        "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
        "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
        "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
        "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
        "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
        "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
        "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
        "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
        "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
        "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
        "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
        "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
        "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
        "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
        "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
        "created_at": "2015-05-05T23:40:12Z",
        "updated_at": "2015-05-05T23:40:12Z",
        "pushed_at": "2015-05-05T23:40:26Z",
        "git_url": "git://github.com/blamewarrior_user/public-repo.git",
        "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
        "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
        "svn_url": "https://github.com/blamewarrior_user/public-repo",
        "homepage": null,
        "size": 0,
        "stargazers_count": 0,
        "watchers_count": 0,
        "language": null,
        "has_issues": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": true,
        "forks_count": 0,
        "mirror_url": null,
        "open_issues_count": 1,
        "forks": 0,
        "open_issues": 1,
        "watchers": 0,
        "default_branch": "master"
      }
    },
    "_links": {
      "self": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1"
      },
      "html": {
        "href": "https://github.com/blamewarrior_user/public-repo/pull/1"
      },
      "issue": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1"
      },
      "comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/1/comments"
      },
      "review_comments": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/comments"
      },
      "review_comment": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/comments{/number}"
      },
      "commits": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls/1/commits"
      },
      "statuses": {
        "href": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c"
      }
    },
    "merged": false,
    "mergeable": null,
    "mergeable_state": "unknown",
    "merged_by": null,
    "comments": 0,
    "review_comments": 0,
    "commits": 1,
    "additions": 1,
    "deletions": 1,
    "changed_files": 1
  },
  "repository": {
    "id": 35129377,
    "name": "public-repo",
    "full_name": "blamewarrior_user/public-repo",
    "owner": {
      "login": "blamewarrior_user",
      "id": 6752317,
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/blamewarrior_user",
      "html_url": "https://github.com/blamewarrior_user",
      "followers_url": "https://api.github.com/users/blamewarrior_user/followers",
      "following_url": "https://api.github.com/users/blamewarrior_user/following{/other_user}",
      "gists_url": "https://api.github.com/users/blamewarrior_user/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/blamewarrior_user/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/blamewarrior_user/subscriptions",
      "organizations_url": "https://api.github.com/users/blamewarrior_user/orgs",
      "repos_url": "https://api.github.com/users/blamewarrior_user/repos",
      "events_url": "https://api.github.com/users/blamewarrior_user/events{/privacy}",
      "received_events_url": "https://api.github.com/users/blamewarrior_user/received_events",
      "type": "User",
      "site_admin": false
    },
    "private": false,
    "html_url": "https://github.com/blamewarrior_user/public-repo",
    "description": "",
    "fork": false,
    "url": "https://api.github.com/repos/blamewarrior_user/public-repo",
    "forks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/forks",
    "keys_url": "https://api.github.com/repos/blamewarrior_user/public-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/blamewarrior_user/public-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/blamewarrior_user/public-repo/teams",
    "hooks_url": "https://api.github.com/repos/blamewarrior_user/public-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/blamewarrior_user/public-repo/events",
    "assignees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/blamewarrior_user/public-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/tags",
    "blobs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/blamewarrior_user/public-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/blamewarrior_user/public-repo/languages",
    "stargazers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/blamewarrior_user/public-repo/subscription",
    "commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/blamewarrior_user/public-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/blamewarrior_user/public-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/blamewarrior_user/public-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/blamewarrior_user/public-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/blamewarrior_user/public-repo/merges",
    "archive_url": "https://api.github.com/repos/blamewarrior_user/public-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/blamewarrior_user/public-repo/downloads",
    "issues_url": "https://api.github.com/repos/blamewarrior_user/public-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/blamewarrior_user/public-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/blamewarrior_user/public-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/blamewarrior_user/public-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/blamewarrior_user/public-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/blamewarrior_user/public-repo/releases{/id}",
    "created_at": "2015-05-05T23:40:12Z",
    "updated_at": "2015-05-05T23:40:12Z",
    "pushed_at": "2015-05-05T23:40:26Z",
    "git_url": "git://github.com/blamewarrior_user/public-repo.git",
    "ssh_url": "git@github.com:blamewarrior_user/public-repo.git",
    "clone_url": "https://github.com/blamewarrior_user/public-repo.git",
    "svn_url": "https://github.com/blamewarrior_user/public-repo",
    "homepage": null,
    "size": 0,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": null,
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": true,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 1,
    "forks": 0,
    "open_issues": 1,
    "watchers": 0,
    "default_branch": "master"
  },
  "requested_reviewers": [{
    "login": "blamewarrior_second_user",
    "id": 6752318,
    "avatar_url": "https://avatars.githubusercontent.com/u/6752318?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/blamewarrior_second_user",
    "html_url": "https://github.com/blamewarrior_second_user",
    "followers_url": "https://api.github.com/users/blamewarrior_second_user/followers",
    "following_url": "https://api.github.com/users/blamewarrior_second_user/following{/other_user}",
    "gists_url": "https://api.github.com/users/blamewarrior_second_user/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/blamewarrior_second_user/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/blamewarrior_second_user/subscriptions",
    "organizations_url": "https://api.github.com/users/blamewarrior_second_user/orgs",
    "repos_url": "https://api.github.com/users/blamewarrior_second_user/repos",
    "events_url": "https://api.github.com/users/blamewarrior_second_user/events{/privacy}",
    "received_events_url": "https://api.github.com/users/blamewarrior_second_user/received_events",
    "type": "User",
    "site_admin": false
  }],
  "installation": {
    "id": 234
  }
}`

	pullRequestPayloadWithMember = `
  {
  "action": "%s",
  "member": {
    "login": "octocat",
    "id": 583231,
    "avatar_url": "https://avatars.githubusercontent.com/u/583231?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  },
  "repository": {
    "id": 35129377,
    "name": "public-repo",
    "full_name": "baxterthehacker/public-repo"
  }
}
`
)
