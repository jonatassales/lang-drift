# 0117: ex's API clients keep credentials and bodies bounded

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** The judge's Anthropic client:

- refuses every redirect;
- reads at most 4 MiB of any response;
- refuses a base URL that is not HTTPS (HTTP on loopback excepted) before any request;
- reports provider failures as status and request id only, in warnings and in the judge's stop reason.

The console client reports a malformed console URL without quoting it. The CLI README states the SDK's credential sources, the redirect refusal and the order in which a credential's private file and record are written.

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] A redirecting endpoint never receives the key on the second host.
- [ ] A response over the bound fails with a bounded read.
- [ ] A provider error body never appears in an error, a warning or a stop reason.
- [ ] A non-HTTPS base URL outside loopback is refused before any request.
- [ ] A malformed console URL's error does not contain the URL or any userinfo.
- [ ] `go test -race ./...` passes in `packages/ex`.

### Verification

- `go test -race ./...` in `packages/ex`, with `httptest` servers for each criterion.
