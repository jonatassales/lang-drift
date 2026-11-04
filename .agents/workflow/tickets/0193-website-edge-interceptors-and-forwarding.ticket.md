# 0193: website edge: interceptors and forwarding

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** The website's Book a Demo route stops sending email and forwards an accepted request to the `demo-requests` service, while every refusal it performs today stays. The route's domain gains an `interceptors/` folder with one interceptor per concern (method, origin, content type, body size, bounded body and strict JSON, schema), each answering the refusal itself or yielding its value, reusing the form's Zod schema and the existing request helpers. An adapter posts the accepted submission to the service with the bearer token, a bounded timeout, no redirects and no cache; 202 becomes 200 `OK`, and a missing URL or token, a network failure, a timeout or any other status becomes 503 `UNAVAILABLE`, the paused state the form already renders. The website's Resend adapters, their tests and the `resend` dependency go; the README and the development environment file carry the two new keys and stop saying the form answers that requests are paused.

**Blocked by:** 0186.

**Status:** ready-for-agent

- [x] Every existing refusal (405, 403, 415, 413, oversized or malformed body, 400 with field errors) answers exactly as before, each from its own interceptor.
- [x] An accepted submission is posted to `<DEMO_REQUESTS_URL>/v1/demo-requests` with `Authorization: Bearer <DEMO_REQUESTS_TOKEN>` and the submission as JSON, and the route answers 200 `OK` on 202.
- [x] With either key unset, with a rejected fetch, a timeout, or any upstream status but 202, the route answers 503 `UNAVAILABLE` with `Retry-After: 30`, and the form shows the return-soon notice.
- [x] `DELIVERY_FAILED` is no longer a reply code; `apps/website/src/infra/resend` and the `resend` dependency are gone; the early-access provider-configuration test keeps passing.
- [x] The route test covers forwarding, the mapping of every upstream outcome and the missing keys with the transport stubbed; the adapter has its own unit test.
- [x] The website README's environment table lists `DEMO_REQUESTS_URL` and `DEMO_REQUESTS_TOKEN` and no longer lists `RESEND_API_KEY`; its prose about the form and about mail is true again.

### Verification

- `pnpm turbo run build typecheck test --filter=website`
- `pnpm turbo run check`
- `platform/cli/drift check biome`
- `platform/cli/drift quality dead-code`
- `platform/cli/drift quality check`
- Manual: with the services up, submit the form and land on the sent page; with them down, see the return-soon notice.
