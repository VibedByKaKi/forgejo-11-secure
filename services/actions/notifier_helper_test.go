// Copyright 2024 The Forgejo Authors
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	actions_model "forgejo.org/models/actions"
	"forgejo.org/models/db"
	issues_model "forgejo.org/models/issues"
	repo_model "forgejo.org/models/repo"
	"forgejo.org/models/unittest"
	user_model "forgejo.org/models/user"
	webhook_module "forgejo.org/modules/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SkipPullRequestEvent(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	repoID := int64(1)
	commitSHA := "1234"

	// event is not webhook_module.HookEventPullRequestSync, never skip
	assert.False(t, SkipPullRequestEvent(db.DefaultContext, webhook_module.HookEventPullRequest, repoID, commitSHA))

	// event is webhook_module.HookEventPullRequestSync but there is nothing in the ActionRun table, do not skip
	assert.False(t, SkipPullRequestEvent(db.DefaultContext, webhook_module.HookEventPullRequestSync, repoID, commitSHA))

	// there is a webhook_module.HookEventPullRequest event but the SHA is different, do not skip
	index := int64(1)
	run := &actions_model.ActionRun{
		Index:     index,
		Event:     webhook_module.HookEventPullRequest,
		RepoID:    repoID,
		CommitSHA: "othersha",
	}
	unittest.AssertSuccessfulInsert(t, run)
	assert.False(t, SkipPullRequestEvent(db.DefaultContext, webhook_module.HookEventPullRequestSync, repoID, commitSHA))

	// there already is a webhook_module.HookEventPullRequest with the same SHA, skip
	index++
	run = &actions_model.ActionRun{
		Index:     index,
		Event:     webhook_module.HookEventPullRequest,
		RepoID:    repoID,
		CommitSHA: commitSHA,
	}
	unittest.AssertSuccessfulInsert(t, run)
	assert.True(t, SkipPullRequestEvent(db.DefaultContext, webhook_module.HookEventPullRequestSync, repoID, commitSHA))
}

func TestForkReviewWorkflowsFromBase(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
	outsider := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})

	forkPR := &issues_model.PullRequest{
		HeadRepoID: 2,
		BaseRepoID: repo.ID,
		Issue: &issues_model.Issue{
			PosterID: outsider.ID,
			Poster:   outsider,
			RepoID:   repo.ID,
			Repo:     repo,
		},
	}
	sameRepoPR := &issues_model.PullRequest{
		HeadRepoID: repo.ID,
		BaseRepoID: repo.ID,
		Issue: &issues_model.Issue{
			PosterID: outsider.ID,
			Poster:   outsider,
			RepoID:   repo.ID,
			Repo:     repo,
		},
	}
	trustedPosterPR := &issues_model.PullRequest{
		HeadRepoID: 2,
		BaseRepoID: repo.ID,
		Issue: &issues_model.Issue{
			PosterID: owner.ID,
			Poster:   owner,
			RepoID:   repo.ID,
			Repo:     repo,
		},
	}

	fromBase, err := forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        owner,
		Event:       webhook_module.HookEventPullRequestReviewApproved,
		PullRequest: forkPR,
	})
	require.NoError(t, err)
	assert.True(t, fromBase)

	fromBase, err = forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        owner,
		Event:       webhook_module.HookEventPullRequestReviewComment,
		PullRequest: forkPR,
	})
	require.NoError(t, err)
	assert.True(t, fromBase)

	fromBase, err = forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        owner,
		Event:       webhook_module.HookEventPullRequestReviewApproved,
		PullRequest: sameRepoPR,
	})
	require.NoError(t, err)
	assert.False(t, fromBase)

	fromBase, err = forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        owner,
		Event:       webhook_module.HookEventPullRequest,
		PullRequest: forkPR,
	})
	require.NoError(t, err)
	assert.False(t, fromBase)

	fromBase, err = forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        owner,
		Event:       webhook_module.HookEventPullRequestReviewApproved,
		PullRequest: trustedPosterPR,
	})
	require.NoError(t, err)
	assert.False(t, fromBase)

	fromBase, err = forkReviewWorkflowsFromBase(db.DefaultContext, &notifyInput{
		Repo:        repo,
		Doer:        outsider,
		Event:       webhook_module.HookEventPullRequestReviewRejected,
		PullRequest: forkPR,
	})
	require.NoError(t, err)
	assert.False(t, fromBase)
}
