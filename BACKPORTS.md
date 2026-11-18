# Security backports onto Forgejo v11.0.16

This tree stays on the abandoned v11 line and only takes security fixes that are
already contained in a **published upstream tag**. Commits that exist only on
release branches, or only as `-dev` tags, are not applied.

## Latest upstream release checked

| Line | Latest published tag | Date | Commit |
| --- | --- | --- | --- |
| v17 (newest stable) | **v17.0.2** | 2026-11-18 | `eea4b0052f24f5ee986b46c00ebf30bdbbf654cb` |
| v16 | v16.0.7 | 2026-10-26 | `44861ff4ed2bdcbf2607d250600203a75ed94c59` |
| v15 | v15.0.12 | 2026-11-18 | `94c6cd65f5c25043e5c9fdb9328d1b8d3f89941a` |
| v14 | v14.0.5 | 2026-04-29 | already older than v11.0.16 |
| v13 | v13.0.5 | 2026-01-16 | already older than v11.0.16 |
| v12 | v12.0.4 | 2025-09-19 | already older than v11.0.16 |

v11.0.16 itself is `0b2b072dcecbde8a53770855fd9c57cc54f276c3` (2026-07-09), the
same day as v15.0.4. v12–v14 security releases were merged into v11 before it
was abandoned, so they are not repeated here.

v17.0.0 (`0d791421f85ca13a63b665053d81ea765fbbc571`, 2026-10-15) was checked
after it was published. Its `[SECURITY]` commits are `golang.org/x/net` v0.60.0,
toolchain go1.26.9 (both already backported from v15.0.10 / v16.0.6), and
`sharp` 0.35.5 (already skipped; not a v11 dependency). Nothing new to apply.

v16 has no published tag after v16.0.7. v17.0.2 and v15.0.12 (both 2026-11-18)
carry the same security set. The v15 commit is the one named below when both
lines have the fix.

Future runs should start after **v17.0.2** / **v16.0.7** / **v15.0.12** and
ignore branch tips until those commits appear in a new tag.

## Version scheme

`v11.0.16-secure.N`, starting at `v11.0.16-secure.1`.

The GitHub release asset is `forgejo-<version>-linux-amd64` with the leading
`v` removed, for example `forgejo-11.0.16-secure.1-linux-amd64`.

## Backported

Each row names the v15 commit when that line has the fix (closest to v11). The
v16 commit is the same change on the newest stable line. Patches were applied
by hand: v15 files have diverged (API context helpers, migration registration,
authorization reducer), so a straight cherry-pick does not apply.

| Fix | Upstream commit (in a tag) | Shipped in | Notes |
| --- | --- | --- | --- |
| Treat RFC 6598 CGNAT `100.64.0.0/10` as private in the host matcher (SSRF) | `24c2be351abf929e3ede80f696277a5153b4e996` | v15.0.6, v16.0.2 | Applied with the upstream tests |
| Org-mode `#+INCLUDE` must not read server files | `0c41d710877c10a736a83fbf773ac4a44b8fc4ac` (org-mode half) | v15.0.6, v16.0.2 | `ReadFile` returns empty and logs are discarded |
| Disable hooks while applying a diff patch (GHSA-rcr6-4jqh-j84m) | `5f7e2e5c003c066a865ea483e42809fa87d85eae` | v15.0.6, v16.0.2 | `core.hooksPath=/dev/null` on apply, write-tree, commit-tree, and the temp-repo push. Local hooks only; remote pre-receive hooks still run |
| Apply the diff patch outside the bare temporary repo | `8efc220692` inside `19b9b9d216bbfb501c18514bd1a8c980246ca3f7` | v15.0.9, v16.0.5 | `git --git-dir <bare> apply` with a separate working directory, so a malicious patch cannot rewrite repo config |
| OpenID association CSRF | `8769076e3e` inside `19b9b9d216` | v15.0.9, v16.0.5 | Session stores `openIDState`. A missing session value is a failure, not a panic |
| Collaborator access mode cannot be raised to owner | `34091413a2` inside `d4de9eb2a87c26b402fdd0259e079957f8cd2b4b` | v15.0.7, v16.0.3 | `mode >= AccessModeOwner` is rejected |
| Public-only API tokens cannot use private repositories when updating a pull request | public-only half of `698e3dca06` / `d4de9eb2a8` and `b04c25d3d2` / `06d641d5b8c8b5670908bb414e22db672a2fb893` | v15.0.7–v15.0.8, v16.0.3–v16.0.4 | Adapted. v11 has public-only tokens but not the authorization reducer. `PublicOnlyContextKey` is set from `tokenRequiresScopes` and checked in `IsUserAllowedToUpdate` and `CanMaintainerWriteToBranch` |
| Draft release attachments require release write permission | `106b85d25a46023a305cb1375d2aa0e0b46b2629` | v15.0.8, v16.0.4 | API `checkReleaseMatchRepo` and web `ServeAttachment`. v11.0.16 already hid draft releases on get/list |
| Organization admin team is not promoted to owner | `450597142beff3af7d571b778738542cac2cfcd7` | v15.0.8, v16.0.4 | `perm.AccessMode = team.AccessMode` |
| Template generation removes a recreated `.git` before `git init` | `dee760c9ff` inside `06d641d5b8` | v15.0.8, v16.0.4 | Both `generate_repo_commit.go` (go1.25) and the legacy file |
| CVE-2026-20896: drop `REVERSE_PROXY_TRUSTED_PROXIES = *` from Docker app.ini templates | `81a5b0d8403d1bbb01421bb75b70bbfeb59d0c73` | v16.0.0 through v16.0.5 | Breaking for operators who relied on the wildcard. Set `FORGEJO__security__REVERSE_PROXY_TRUSTED_PROXIES` to the proxy address. Upstream did not put this in v15 LTS; v11 still shipped the vulnerable default |
| `github.com/go-chi/chi/v5` v5.3.2 | `00ebec2c6c` (v5.3.0 security) then `5af242d43a` (v5.3.2, same tags) | v15.0.6+ / v16.0.5 go.mod | v11 was v5.2.4 |
| `github.com/klauspost/compress` v1.18.7 | `c700a1ba0c` | v15.0.6, v16.0.2 | |
| `golang.org/x/crypto` v0.56.0 | `462b8b4993` | v15.0.8, v16.0.4 | CVE-2026-78662 and CVE-2026-56855 (SSH DoS). Module requires Go ≥ 1.26 |
| `golang.org/x/image` v0.45.0 | `ff6d06accc` | v15.0.7, v16.0.3 | v11 already had v0.43.0 |
| `golang.org/x/net` v0.58.0 | `7335a4d31d` (v0.56.0 security); v0.58.0 is what v15.0.7 and v16.0.5 ship | v15.0.6+ | v11 was v0.55.0. Superseded by v0.60.0 below |
| `golang.org/x/text` v0.41.0 | `7bb6ba2197` (v0.39.0 security); v0.41.0 is what v16.0.5 ships | v15.0.6+ | v11 was v0.38.0 |
| `golang.org/x/mod` v0.40.0 (indirect) | `0776a1843d` | v15.0.7, v16.0.3 | |
| `golang.org/x/sync` v0.22.0, `x/sys` v0.47.0, `x/tools` v0.49.0 | transitive, same versions as v16.0.5 `go.mod` | v16.0.5 | Pulled in by the modules above |
| Go `go 1.26.0` / `toolchain go1.26.7` | `525fc5b43c`, `971568e264` | v15.0.7, v16.0.5 | Required to consume `x/crypto` v0.56.0. Superseded by go1.26.9 below |
| mermaid 11.16.1 | `40c92cc422` | v15.0.7, v16.0.3 | lockfile updated |
| postcss 8.5.23 | `583c447a11` | v15.0.6, v16.0.2 | lockfile updated. Intermediate 8.5.12 and 8.5.18 are included |
| svgo 4.1.0 | `4aa4c8793f` | v15.0.8, v16.0.4 | `import {optimize} from 'svgo'` is unchanged. Intermediate 4.0.2 is included |
| Fork `pull_request_review` and `pull_request_review_comment` run from the pull request base | `6ff8b8e14f` inside `17e1001d5e` | v15.0.10, v16.0.6 | Adapted. v15 only adds `WithPullRequest` because `getPullRequestCommitAndApproval` already chooses the base commit. v11 has no `ActionUser` trust table, so a fork review uses the base workflow when the reviewer can write actions and the poster cannot. Same-repo pull requests and a trusted poster still use the head. `IsForkPullRequest` is set, so repository secrets stay withheld |
| Tag creation does not load attachments for an unsaved release | `65ca1b0dbe` inside `17e1001d5e` | v15.0.10, v16.0.6 | `LoadAttributes` returns before archive counts and attachments when `ID == 0`. The second `LoadAttributes` call inside `createTag` is removed |
| Issue label delete cannot read another repository's label | `9623034655` inside `17e1001d5e` | v15.0.10, v16.0.6 | Adapted. v11's route is `DELETE .../labels/{id}` and returns 204, not the GitHub label body. The id is resolved with `GetLabelInRepoByID`, then `GetLabelInOrgByID` when the owner is an organization, so an organization label on an issue can still be removed |
| Org and repo label deletes no longer share one id check | `5e4dc852a4` inside `17e1001d5e` | v15.0.10, v16.0.6 | `DeleteLabel` is split into `DeleteLabelInOrg` and `DeleteLabelInRepo`. A type mismatch returns an error instead of comparing an org id with a repo id |
| `katex` 0.18.2 | `b6291844dc` | v15.0.10, v16.0.6 | Direct dependency. v11 was 0.16.21. Mermaid still nests katex 0.16.45; micromark-extension-math nests 0.16.47, same shape as the v15 lockfile |
| `golang.org/x/net` v0.60.0 | `50d17fc53d` | v15.0.10, v16.0.6 | Also pulls `x/crypto` v0.57.0, `x/sync` v0.23.0, `x/sys` v0.48.0, `x/text` v0.42.0, and `x/mod` v0.41.0 |
| Go toolchain `go1.26.9` | `51056ebd74` | v15.0.10, v16.0.6 | `go 1.26.0` stays. `.github/workflows/release.yml` pins `GOTOOLCHAIN`, so that pin moves to go1.26.9 with the directive. Superseded by go1.26.10 below |
| `pull_request_review` Actions use `refs/pull/%d/head`, not the review `commit_id` | `291cf01266` inside `753f7edc77` | v15.0.11, v16.0.7, v17.0.1 | `WithRef(review.CommitID)` is removed. `WithPullRequest` sets the ref. v11 has no `TestActionsForkPullRequest` harness, so `TestPullRequestReviewRefUsesPullHead` covers the ref contract |
| Visibility of users when listing watched and starred repositories | `16811fea1f`, `3f3e8bea89`, `dae5735aee`, `e0bfad0d6c` inside `753f7edc77` | v15.0.11, v16.0.7, v17.0.1 | API lists go through `SearchRepository` (`Actor` is the viewer, `Private` follows sign-in). Permissions use the viewer. `GetStarredRepos` / `GetWatchedRepos` are removed. Watched search excludes `WatchModeDont`. No authorization reducer on v11 |
| NuGet `X-NuGet-ApiKey` access tokens respect package scopes | `4ff23f1ffb` inside `753f7edc77` | v15.0.11, v16.0.7, v17.0.1 | v11 stores the scope on the request (`IsApiToken`, `ApiTokenScope`) instead of an auth-result `Scope()`. basic, oauth2, container, and conan already did. chef and HTTP signatures are not access tokens |
| Single quotes around pull-request checkout and merge commands | `e5145ff6a5` | v17.0.1 | v11 template had no quotes (`git checkout`, not `git switch`). A branch name could run a shell substitution when the instructions were pasted. Quotes match the v17 fix |
| Go toolchain `go1.26.10` | `5f7f338795` | v15.0.11, v16.0.7, v17.0.1 | `go 1.26.0` stays. `.github/workflows/release.yml` `go-version` and `GOTOOLCHAIN` move to go1.26.10 |
| `happy-dom` 21.0.4 | `115c57529a` | v15.0.11, v16.0.7, v17.0.1 | Frontend test devDependency. v11 was 20.8.9. Lockfile also takes `buffer-image-size` 0.6.4 and `ws` 8.22.0, same as the v15 bump |
| Search authorized SSH keys by fingerprint, then compare the key bytes | `94c6cd65f5` | v15.0.12, v17.0.2 | [v15 #15299](https://codeberg.org/forgejo/forgejo/pulls/15299) / [v17 #15298](https://codeberg.org/forgejo/forgejo/pulls/15298). Replaces the SQL `LIKE` on key text. The same commit includes the v17 fingerprint lookup ([#13614](https://codeberg.org/forgejo/forgejo/pulls/13614)) that v15 needed for the fix. `cmd/keys` stays on urfave/cli v2, parses the key, and sends `SHA256` fingerprint. A missing key returns `# No key found` |
| Primary email must already be activated | `94c6cd65f5` | v15.0.12, v17.0.2 | [v15 #15299](https://codeberg.org/forgejo/forgejo/pulls/15299) / [v17 #15298](https://codeberg.org/forgejo/forgejo/pulls/15298). `MakeEmailAddressPrimary` returns `ErrEmailNotActivated`. `ReplaceInactivePrimaryEmail` still uses the unexported helper |
| Commit `tree` / `parent` / `author` / `committer` headers are positional | `94c6cd65f5` | v15.0.12, v17.0.2 | [v15 #15299](https://codeberg.org/forgejo/forgejo/pulls/15299) / [v17 #15298](https://codeberg.org/forgejo/forgejo/pulls/15298). Later copies of those headers stay in the payload and do not replace the real commit. `encoding` and `change-id` still go to the payload |
| `/tree-list` and `/find` require the code unit | `94c6cd65f5` | v15.0.12, v17.0.2 | [v15 #15299](https://codeberg.org/forgejo/forgejo/pulls/15299) / [v17 #15298](https://codeberg.org/forgejo/forgejo/pulls/15298). `reqRepoCodeReader` on both routes |
| Public-only tokens cannot read private repositories in activity feeds or the heatmap | public-only half of `94c6cd65f5` | v15.0.12, v17.0.2 | [v15 #15299](https://codeberg.org/forgejo/forgejo/pulls/15299) / [v17 #15298](https://codeberg.org/forgejo/forgejo/pulls/15298). No authorization reducer. `IsPublicOnly` adds the public-repository subquery used by feeds and the heatmap |
| Debian packages parsed with `pault.ag/go/debian` v0.21.4 | `b6b96a2ca9` | v15.0.12, v17.0.2 | [v15 #15304](https://codeberg.org/forgejo/forgejo/pulls/15304) / [v17 #15302](https://codeberg.org/forgejo/forgejo/pulls/15302), backport of [#15242](https://codeberg.org/forgejo/forgejo/pulls/15242). `ParsePackage` takes `io.ReaderAt`. Tests still build archives with `blakesmith/ar`. The module requires `klauspost/compress` v1.20.3, `ulikunitz/xz` v0.5.17, `ProtonMail/go-crypto` v1.5.2, and `cloudflare/circl` v1.6.5 |
| `minimatch` 11.0.1 | `70f889ebb2` | v15.0.12, v17.0.2 | [v15 #15083](https://codeberg.org/forgejo/forgejo/pulls/15083) / [v17 #15085](https://codeberg.org/forgejo/forgejo/pulls/15085). `[SECURITY]`. Direct dependency. `import {minimatch} from 'minimatch'` is unchanged. v11 was 10.2.3 |
| `github.com/andybalholm/brotli` v1.2.7 | `0c83767960` | v15.0.12, v17.0.2 | [v15 #15158](https://codeberg.org/forgejo/forgejo/pulls/15158) / [v17 #15160](https://codeberg.org/forgejo/forgejo/pulls/15160). `[SECURITY]`. Indirect. v11 was v1.1.1. The v17 PR was closed as already included by the tag |
| `github.com/klauspost/compress` v1.20.3 | `2f15a8d128` (v1.20.2 `[SECURITY]`), then `b6b96a2ca9` | v15.0.12, v17.0.2 | [v15 #15159](https://codeberg.org/forgejo/forgejo/pulls/15159) / [v17 #15161](https://codeberg.org/forgejo/forgejo/pulls/15161) bump to v1.20.2. The debian parser commit in the same tags requires v1.20.3, which is what v15.0.12 and v17.0.2 ship. v11 was v1.18.7 |
| `golang.org/x/net` v0.62.0 | `44426ef404` | v15.0.12, v17.0.2 | [v15 #15063](https://codeberg.org/forgejo/forgejo/pulls/15063) / [v17 #15065](https://codeberg.org/forgejo/forgejo/pulls/15065). `[SECURITY]`. Also pulls `x/crypto` v0.58.0, `x/sync` v0.24.0, `x/sys` v0.49.0, `x/text` v0.43.0, and `x/tools` v0.51.0 |
| Go toolchain `go1.26.11` | `11c63cb45e` | v15.0.12, v17.0.2 | [v15 #15224](https://codeberg.org/forgejo/forgejo/pulls/15224) / [v17 #15225](https://codeberg.org/forgejo/forgejo/pulls/15225). The commit title is not `[SECURITY]`. Both published tags pin this toolchain, and `.github/workflows/release.yml` follows `go.mod`, same as go1.26.10. `go 1.26.0` stays |

## Skipped

| Fix | Upstream | Why |
| --- | --- | --- |
| Actions pre-execution stored XSS | org-mode commit `0c41d71087` also edits `models/actions/pre_execution_errors.go` | `PreExecutionError` was added in v13 and is not in v11 |
| Expand local reusable workflows from the base branch on `pull_request_target` | `841c28d1f4` inside `d4de9eb2a8` | Needs `ExpandLocalReusableWorkflows`, `WorkflowSourceCommit`, and migration registration that v11 does not have |
| Authorization reducer and repo-specific access tokens | remainder of `d4de9eb2a8` and `06d641d5b8` | `services/authz` and `GetUserRepoPermissionWithReducer` do not exist on v11. The public-only behavior was ported separately (see above) |
| sharp security bumps (v0.35.0 / v0.35.4 / v0.35.5) | `2c05f9af24`, `8fbae40d9c`, `cb23914d45` | `sharp` is not a v11 dependency. v15.0.10 and v16.0.6 ship 0.35.5 |
| `google.golang.org/grpc` indirect bumps | `9b6de2d5ae`, `72f8a67635`, `f2d33d5625` | grpc is not in v11 `go.mod` |
| `github.com/jackc/pgx` | later lines | not a v11 dependency |
| `github.com/go-ldap/ldap/v3` v3.4.14 | `77b5524f6e` | not marked `[SECURITY]` |
| goldmark v1.8.2 | later lines | not marked `[SECURITY]`; v11 already has v1.7.17 |
| go-git v5.17.1 / v5.18.0 | older `[SECURITY]` commits | v11.0.16 is already on v5.18.0 |
| happy-dom 20.8.9, minimatch 10.2.3, webpack 5.104.1 | commits before v11.0.16 | already in v11.0.16. Later webpack/minimatch/katex numbers in v16.0.5 are not marked `[SECURITY]` |
| 2026-01-08, 2026-05-12, and 2026-06-10 security patches | `52406dc6ea`, `9e51a55b63`, `a32151caf6` | already in v11.0.16 |
| Actions trust helper and the fork-event integration test | `cc039bc494` inside `17e1001d5e` (`getPullRequestCommitAndApproval`, `UpdateTrustedWithPullRequest`, `TestActionsForkPullRequest`) | Needs `services/actions/trust.go`, `ActionUser`, and `forgery.CreateRepository`, which v11 does not have. The review-event dispatch is ported above. Other non-sync fork events (assign, label, close, and so on) already call `WithPullRequest` on v11 but still detect workflows from the head, because that trust helper is absent. The v15.0.11 `expectedRef` extension of that same test is skipped for the same reason |
| Type confusion in `GetProjectByIDForOwner` | `520afcaf8e` | v17.0.1 only. The helper and `ensureProjectOwnedBy` are not in v11. Repo routes compare `RepoID` with the repository, and user/org routes compare `OwnerID` with the user. `CanBeAccessedByOwnerRepo` already rejects a user id used as a repository id |
| NuGet authorization reducer | reducer half of `4ff23f1ffb` | `services/authz` is not on v11. The scope return is ported above |
| Repo-specific access tokens on activity feeds and the heatmap | reducer half of `94c6cd65f5` / `eea4b0052f` | `services/authz`, `ctx.Reducer()`, and `createFineGrainedRepoAccessToken` are not on v11. The semgrep rule that requires `AuthorizationReducer` on `GetFeedsOptions` is the same feature. The public-only half is ported above |
| Replace `blakesmith/ar` with the `code.forgejo.org/crystal/ar` mirror | `3e36a8f11a` / `e02b117` | Not marked `[SECURITY]`. The author describes the govulncheck report as a request-level panic, not an exploitable issue. The debian parser change in the same tags removes production use of that library |

## Adaptations

- Whole merge commits were not cherry-picked. Applicable hunks were copied onto the v11 tree.
- OpenID: upstream type-asserts `Session.Get` to `string` and panics when the key is missing. This port treats a non-string or empty value as an invalid state.
- `CanMaintainerWriteToBranch` checks the base repository only after `LoadBaseRepo`.
- CVE-2026-20896 removes a default. Installs that use reverse-proxy authentication must set `REVERSE_PROXY_TRUSTED_PROXIES` explicitly.
- Fork review events: v15 relies on `getPullRequestCommitAndApproval` after `WithPullRequest`. This tree copies that decision with actions write permission, because there is no explicit-trust user table. A poster who can write still runs the head workflow.
- Issue label delete: v15 looks up a name, then a repository id. v11 only has a numeric id, and organization labels are valid on an organization's issues, so the id is accepted from the repository or from that repository's organization.
- `DeleteLabel` callers (`routers/api/v1/org`, `routers/api/v1/repo`, `routers/web/org`, `routers/web/repo`, `services/f3/driver`) now use `DeleteLabelInOrg` or `DeleteLabelInRepo`.
- The GitHub release workflow hardcodes `go-version` and `GOTOOLCHAIN`. Both are go1.26.10 so the published binary is built with the toolchain in `go.mod`.
- Watched and starred API lists use `SearchRepository` without `AuthorizationReducer`. Signed-in viewers who are not admins are limited by `AccessibleRepositoryCondition`. Admins still see every repository the search returns, which is the existing v11 search behavior.
- NuGet scope enforcement uses `store` `IsApiToken` / `ApiTokenScope`, which `reqPackageAccess` already reads. The v15 auth-result `Scope()` and `Reducer()` methods do not exist here.
- Pull-request CLI instructions keep `git checkout`. Only the quoting changes, from bare words to single quotes.
- SSH authorized-key lookup uses the fingerprint index that already exists on v11. `cmd/keys` stays on urfave/cli v2. LDAP sync compares stored key rows instead of searching by `LIKE`.
- `ReplaceInactivePrimaryEmail` calls the unexported primary-email helper, so replacing an inactive primary during registration still works. The settings form redirects when the address is missing or not activated.
- Commit parsing keeps `encoding` and `change-id` in the signature payload. A second `tree`, `parent`, `author`, or `committer` after the positional headers does not replace them.
- Activity feeds and the heatmap do not grow an `AuthorizationReducer` field. A public-only token adds a repository subquery from `PublicRepositoryFilter`, which matches `RepoVisibleToPublicOnly`.
- The GitHub release workflow `go-version` and `GOTOOLCHAIN` are go1.26.11, matching `go.mod`.
