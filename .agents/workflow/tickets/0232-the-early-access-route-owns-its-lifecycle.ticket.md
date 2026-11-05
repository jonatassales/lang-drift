# 0232: the early-access route owns its lifecycle behind an explicit pause

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** The website's early-access `POST` first checks an explicit pause value, a `.data` constant beside the route, and answers the paused 503 as today. Past the pause it runs the interceptors (method, origin, content type, body size, body, schema) from an `interceptors/` folder beside the route, admits the email through the registry adapter, verifies the challenge, registers through the adapter, reads the receipt and replies. The registry calls live in one adapter under `src/infra/`, which takes the registry target (URL and key) as input, bounds the payload, composes the request signal with its own timeout and refuses redirects. `try/catch` wraps only the registry calls and the challenge verification, and the route's two discarded rejections are handled.

**Blocked by:** 0229.

**Status:** superseded by spec 0046, which removes the early-access code this ticket reshapes. 2026-11-04.

- [ ] `POST` answers the paused 503 while the pause value is on, and runs the pipeline when it is off.
- [ ] Each interceptor lives in its own file under the route's `interceptors/` folder, and the registry calls in one adapter under `src/infra/`.
- [ ] `try/catch` wraps only external calls, and no rejection is discarded.
- [ ] The route's unit tests call `POST` once with the pause on and, with the pause mocked off, for every refusal and the registry exchange with `fetch` stubbed; no module is left dead.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
