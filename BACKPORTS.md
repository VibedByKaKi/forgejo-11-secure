# Security backports onto Forgejo v11.0.16

This tree stays on the abandoned v11 line and only takes security fixes that are
already contained in a **published upstream tag**. Commits that exist only on
release branches, or only as `-dev` tags, are not applied.

## Latest upstream release checked

| Line | Latest published tag | Date | Commit |
| --- | --- | --- | --- |
| v16 (newest stable) | **v16.0.6** | 2026-10-12 | `1aa13944143b3bc4312607239a68da23512d5ce3` |
| v15 | v15.0.10 | 2026-10-12 | `17e1001d5e00e8ee3cbc1b6dd31be94ac38185f9` |
| v14 | v14.0.5 | 2026-04-29 | already older than v11.0.16 |
| v13 | v13.0.5 | 2026-01-16 | already older than v11.0.16 |
| v12 | v12.0.4 | 2025-09-19 | already older than v11.0.16 |

v11.0.16 itself is `0b2b072dcecbde8a53770855fd9c57cc54f276c3` (2026-07-09), the
same day as v15.0.4. v12–v14 security releases were merged into v11 before it
was abandoned, so they are not repeated here. There is no stable v17 tag
(`v17.0.0-dev` only).

Future runs should start after **v16.0.6** / **v15.0.10** and ignore branch tips
until those commits appear in a new tag.

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
| Go toolchain `go1.26.9` | `51056ebd74` | v15.0.10, v16.0.6 | `go 1.26.0` stays. `.github/workflows/release.yml` pins `GOTOOLCHAIN`, so that pin moves to go1.26.9 with the directive |

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
| Actions trust helper and the fork-event integration test | `cc039bc494` inside `17e1001d5e` (`getPullRequestCommitAndApproval`, `UpdateTrustedWithPullRequest`, `TestActionsForkPullRequest`) | Needs `services/actions/trust.go`, `ActionUser`, and `forgery.CreateRepository`, which v11 does not have. The review-event dispatch is ported above. Other non-sync fork events (assign, label, close, and so on) already call `WithPullRequest` on v11 but still detect workflows from the head, because that trust helper is absent |

## Adaptations

- Whole merge commits were not cherry-picked. Applicable hunks were copied onto the v11 tree.
- OpenID: upstream type-asserts `Session.Get` to `string` and panics when the key is missing. This port treats a non-string or empty value as an invalid state.
- `CanMaintainerWriteToBranch` checks the base repository only after `LoadBaseRepo`.
- CVE-2026-20896 removes a default. Installs that use reverse-proxy authentication must set `REVERSE_PROXY_TRUSTED_PROXIES` explicitly.
- Fork review events: v15 relies on `getPullRequestCommitAndApproval` after `WithPullRequest`. This tree copies that decision with actions write permission, because there is no explicit-trust user table. A poster who can write still runs the head workflow.
- Issue label delete: v15 looks up a name, then a repository id. v11 only has a numeric id, and organization labels are valid on an organization's issues, so the id is accepted from the repository or from that repository's organization.
- `DeleteLabel` callers (`routers/api/v1/org`, `routers/api/v1/repo`, `routers/web/org`, `routers/web/repo`, `services/f3/driver`) now use `DeleteLabelInOrg` or `DeleteLabelInRepo`.
- The GitHub release workflow hardcodes `go-version` and `GOTOOLCHAIN`. Both are go1.26.9 so the published binary is built with the toolchain in `go.mod`.
