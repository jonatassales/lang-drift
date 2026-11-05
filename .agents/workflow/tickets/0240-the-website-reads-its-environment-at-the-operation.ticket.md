# 0240: the website reads its environment where it uses it

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** The early-access registry adapter reads `EARLY_ACCESS_REGISTRY_URL`, `EARLY_ACCESS_REGISTRY_KEY` and `EARLY_ACCESS_PUBLIC_ORIGIN` at the fetch, with the credential-free origin check and the key-length check beside the reads. The Turnstile verification reads `TURNSTILE_SECRET_KEY`, `TURNSTILE_VERIFY_URL` and `TURNSTILE_EXPECTED_HOSTNAME` at the fetch, with the hostname shape, the official-endpoint-in-production and the no-test-credentials-in-production checks beside them. The mode is derived at each of those operations from `NODE_ENV` and `EARLY_ACCESS_MODE`; the three reader modules are deleted. The cookie service factory requires `CRYPTO_SECRET` in every environment and invents nothing; the README documents every key and that a fresh checkout takes the value from the developer's untracked `.env`.

**Blocked by:** 0232, 0239.

**Status:** ready-for-agent. Its early-access items are superseded by spec 0046, which deletes that code. 2026-11-04.

- [ ] `read-registry.compute.ts`, `read-turnstile.compute.ts` and `early-access-mode.compute.ts` are gone, and each variable is read with `process.env.NAME` at the operation that uses it.
- [ ] The cookie factory throws without `CRYPTO_SECRET` in every `NODE_ENV`, with unit tests.
- [ ] The Turnstile verification answers `unavailable` for a non-official endpoint in production and `verified` for the official one; the registry adapter refuses a registry URL with credentials.
- [ ] `apps/website/README.md` lists every key and no value, and neither the README nor `.agents/context.md` says the secret is generated on the fly.

## Dev part

**Override:** No part qualifies: the reads and the checks beside them move to their operations as the spec writes them, and the cookie factory's new failure is one guard; none holds a problem of its own (policy, Dev part). The new failure is behaviour, covered by the criteria. 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift env validate`
