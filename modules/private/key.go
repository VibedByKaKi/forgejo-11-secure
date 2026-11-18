// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"context"
	"fmt"

	"forgejo.org/modules/setting"
)

// UpdatePublicKeyInRepo update public key and if necessary deploy key updates
func UpdatePublicKeyInRepo(ctx context.Context, keyID, repoID int64) error {
	// Ask for running deliver hook and test pull request tasks.
	reqURL := setting.LocalURL + fmt.Sprintf("api/internal/ssh/%d/update/%d", keyID, repoID)
	req := newInternalRequest(ctx, reqURL, "POST")
	_, extra := requestJSONResp(req, &ResponseText{})
	return extra.Error
}

// AuthorizedPublicKeyByFingerprint finds a public key via its fingerprint.
func AuthorizedPublicKeyByFingerprint(ctx context.Context, fingerprint string) (*ResponseText, ResponseExtra) {
	reqURL := setting.LocalURL + "api/internal/ssh/authorized_keys"
	req := newInternalRequest(ctx, reqURL, "POST")
	req.Param("content", fingerprint)
	return requestJSONResp(req, &ResponseText{})
}
