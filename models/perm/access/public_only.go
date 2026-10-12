// Copyright 2026 The Forgejo Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package access

import (
	"context"

	repo_model "forgejo.org/models/repo"
)

// PublicOnlyContextKey marks an API request whose token is limited to public resources.
// v11 does not have the authorization-reducer framework added in later Forgejo releases.
// This flag carries the public-only restriction into permission checks that would
// otherwise use the user's full access, which is the part of the v15.0.7 / v16.0.3
// pull-request update fix that applies on v11.
type publicOnlyContextKey struct{}

// PublicOnlyContextKey is the context key set for public-only API tokens.
var PublicOnlyContextKey publicOnlyContextKey

// IsPublicOnly reports whether ctx was marked as a public-only API token request.
func IsPublicOnly(ctx context.Context) bool {
	v, _ := ctx.Value(PublicOnlyContextKey).(bool)
	return v
}

// RepoVisibleToPublicOnly reports whether a public-only token may use this repository.
// Matches the later PublicReposAuthorizationReducer: the owner must be public and the
// repository itself must not be private.
func RepoVisibleToPublicOnly(ctx context.Context, repo *repo_model.Repository) (bool, error) {
	if repo == nil {
		return false, nil
	}
	if err := repo.LoadOwner(ctx); err != nil {
		return false, err
	}
	return !repo.IsPrivate && repo.Owner.Visibility.IsPublic(), nil
}
