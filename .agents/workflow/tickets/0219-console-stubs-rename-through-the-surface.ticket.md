# 0219: console-stubs is its name on every visible surface

**Spec:** .agents/workflow/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** Nothing the console, its operator or the stubs' client sees still says `extub`. The environment contract becomes `CONSOLE_STUBS_URL`, `CONSOLE_STUBS_TOKEN` and `CONSOLE_STUBS_PORT` in the console, the stubs and the shared development defaults; the console's proxy route becomes `/api/console-stubs/...` and its data files, link helper and base computation are named after console-stubs; the stubs answer with a `console-stubs` realm, service field and load-test target, expect a `console-stubs` data channel and log under that prefix; the catalog key for the unavailable notice is renamed in every served locale with the same copy. The READMEs and the context document describe the new contract. The owner recreates the three variables on Vercel under the new names before this merges.

**Blocked by:** 0218.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The console and the stubs read only `CONSOLE_STUBS_*` variables; a missing one fails with its new name; the development defaults and both READMEs document the new keys.
- [ ] `GET /api/console-stubs/<path>` proxies to the stubs and `/api/extub/...` no longer exists; the proxy-matcher, link, entity-id and load-test tests pass with the new names.
- [ ] The stubs answer `WWW-Authenticate: Bearer realm="console-stubs"` and `"service": "console-stubs"`, the load tester reports target `console-stubs-loopback`, and the WebRTC hub accepts the `console-stubs` data channel; their tests say so.
- [ ] The console catalogs carry `console_stubs_unavailable` with unchanged copy in every served locale and the parity test passes; no catalog carries `extub_unavailable`.
- [ ] `git grep -n -i 'extub'` finds only workflow history and CLI fixtures.
- [ ] Owner, before merge: `CONSOLE_STUBS_TOKEN` exists on the console and stubs Vercel projects and `CONSOLE_STUBS_URL` on the console project; after merge, the `EXTUB_*` variables are removed.

## Dev part

**Override:** No part qualifies: a rename of identifiers, routes and keys with no new behaviour (policy, Categories: never a rename). 2026-11-04.

## Verification

- `pnpm turbo run build typecheck test --filter=console-stubs --filter=console`
- `platform/cli/drift check localization`
- `pnpm console`, then `curl -i http://127.0.0.1:3001/api/console-stubs/health`
