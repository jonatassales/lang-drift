# 0223: the console judges at `/api/judge`

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** A request carrying the deployment's secret and a question of the contract gets a judged answer from the platform's own model. The console gains the route `/api/judge`: a route handler that owns its `POST`, with one interceptor per concern beside it (method, bearer secret compared in constant time, content type, body size, the question's schema), a provider adapter under the console's `infra/` folder (new, as the Next tree names it) that reads the platform's key from the environment at the call and refuses when it is missing, validation of the model's text against the kind's schema before answering, the cost of the call computed from the usage at the platform's list prices and marked illustrative, and one observability line per judgment with usage and cost and no source text. 401 without the secret, 400 for a question outside the contract, 503 when the model is unavailable or refuses to start, 200 with the response of the contract. The console accepts only the models it offers, named in its environment, and refuses another with 400 naming the offered set. `GET /api/judge` with the secret answers 204 and without it 401: the connection check `dft agent connect` performs.

**Blocked by:** 0222.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `POST /api/judge` without the bearer answers 401; with a body outside the request schema, 400 naming the field; with the platform's key missing or the provider unreachable, 503; with a valid question, 200 with a response that validates against the contract's response schema.
- [x] The response's answer text validates against the kind's schema; a model text that does not is answered as a refusal, never forwarded.
- [x] The response carries the tokens used and a cost in cents with provenance `console` and `illustrative: true`, computed from the usage at the list prices the console holds for its model.
- [x] `GET /api/judge` answers 204 with the secret and 401 without; a `POST` naming a model the console does not offer answers 400 with the offered set.
- [x] The route's tests run with a fake provider and the contract's samples; no test calls a provider.
- [x] The console README documents `RUNTIME_SECRET` and the platform's provider key by name, and the owner's step of setting both on the Vercel project.

## Dev part

**Category:** `core-domain`. Turning what the model said into a judgment the platform stands behind, or into a refusal, and what that judgment costs, is the logic the owner must be able to defend to a consumer reading their ledger.

**Paths:**
- `apps/console/src/app/api/judge/judge-response.compute.ts`
- `apps/console/src/app/api/judge/judge-response.compute.test.ts`
- `apps/console/src/app/api/judge/judgment-cost.compute.ts` and `has-token-counts.compute.ts`, the pricing and the count check it calls, each in its own file under the one-function rule

**Boundary:** `judgeResponse(result: ProviderResult, kind: QuestionKind, prices: ModelPrices): JudgeResponse | JudgeRefusal`, called by the route after the provider answers; `ProviderResult` (text, usage, stop reason), `QuestionKind`, `ModelPrices`, `JudgeResponse` (the validated answer text, usage, cost in cents, provenance, illustrative, judge version) and `JudgeRefusal` (the reason) are declared in `judge.types.ts` beside it by the agent, which also writes the route, the interceptors, the provider adapter and the Zod schemas of the contract.

**Goal:** Decide what the route answers from what the provider returned: parse the text as JSON and validate it against the kind's schema, turning anything else, a truncated stop or an empty text into a refusal that names the reason; price the usage at the list prices in cents per million tokens, input and output apart, rounded up to the cent, never zero when tokens were used, marked illustrative; refuse a usage with missing or negative counts.

- [x] The contract's sample answers produce responses whose cost the owner computes by hand; one token costs one cent; a usage of zero tokens with a non-empty text is refused.
- [x] A text outside the kind's schema, a truncated stop reason and an empty text each produce a refusal naming the cause, and no response carries unvalidated text.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("you write"), so the agent wrote `judgeResponse` and its test.

## Verification

- `pnpm turbo run build typecheck test --filter=console`
- `drift quality check 06-nextjs`
