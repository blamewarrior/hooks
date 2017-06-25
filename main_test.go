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

package main_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	main "github.com/blamewarrior/hooks"
)

type RepositoriesServiceMock struct {
	mock.Mock
}

func (m *RepositoriesServiceMock) Track(ctx context.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func (m *RepositoriesServiceMock) Untrack(ctx context.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func TestTrackingHandler_DoAction(t *testing.T) {
	reposService := new(RepositoriesServiceMock)

	reposService.On(
		"Track",
		context.Background(),
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	reposService.On(
		"Untrack",
		context.Background(),
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	handler := main.NewTrackingHandler("blamewarrior.com")

	suits := []struct {
		Action string
		Err    error
	}{
		{
			"track",
			nil,
		},
		{
			"untrack",
			nil,
		},

		{
			"custom",
			errors.New("Unsupported action custom"),
		},
	}

	for _, suits := range suits {
		err := handler.DoAction(reposService, "blamewarrior/hooks", suits.Action)
		assert.Equal(t, suits.Err, err)
	}

}

const GithubPullRequestHookBody = `
  {
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
  "requested_reviewer": {
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
  "installation": {
    "id": 234
  }
}
`
