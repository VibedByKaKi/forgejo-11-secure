// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2026 The Forgejo Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"testing"

	"forgejo.org/models/db"
	"forgejo.org/models/unittest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate_InsertReleases(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	a := &Attachment{
		UUID: "a0eebc91-9c0c-4ef7-bb6e-6bb9bd380a12",
	}
	r := &Release{
		Attachments: []*Attachment{a},
	}

	err := InsertReleases(db.DefaultContext, r)
	require.NoError(t, err)
}

func TestReleaseLoadRepo(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	release := unittest.AssertExistsAndLoadBean(t, &Release{ID: 1})
	assert.Nil(t, release.Repo)

	require.NoError(t, release.LoadRepo(db.DefaultContext))

	assert.EqualValues(t, 1, release.Repo.ID)
}

func TestReleaseDisplayName(t *testing.T) {
	release := Release{TagName: "TagName"}

	assert.Empty(t, release.DisplayName())

	release.IsTag = true
	assert.Equal(t, "TagName", release.DisplayName())

	release.Title = "Title"
	assert.Equal(t, "Title", release.DisplayName())
}

func TestLoadAttributes(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	t.Run("Normal", func(t *testing.T) {
		release := unittest.AssertExistsAndLoadBean(t, &Release{ID: 1})

		require.NoError(t, release.LoadAttributes(t.Context()))

		if assert.NotNil(t, release.Publisher) {
			assert.EqualValues(t, 2, release.Publisher.ID)
		}
		if assert.NotNil(t, release.Repo) {
			assert.EqualValues(t, 1, release.Repo.ID)
		}
		if assert.NotNil(t, release.ArchiveDownloadCount) {
			assert.Empty(t, release.ArchiveDownloadCount.Zip)
			assert.Empty(t, release.ArchiveDownloadCount.TarGz)
		}
		if assert.Len(t, release.Attachments, 1) {
			assert.EqualValues(t, 9, release.Attachments[0].ID)
		}
	})

	t.Run("Not from DB", func(t *testing.T) {
		release := &Release{
			RepoID:      3,
			PublisherID: 3,
		}

		require.NoError(t, release.LoadAttributes(t.Context()))

		assert.Nil(t, release.ArchiveDownloadCount)
		assert.Nil(t, release.Attachments)
		if assert.NotNil(t, release.Publisher) {
			assert.EqualValues(t, 3, release.Publisher.ID)
		}
		if assert.NotNil(t, release.Repo) {
			assert.EqualValues(t, 3, release.Repo.ID)
		}
	})
}
