# 0222: one judge contract for the CLI and the console

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** The three questions the planner asks and the answer it expects have one definition both sides are tested against. The core package gains a `judge` surface: a JSON Schema per question kind (`sources`, `episodes`, `suggestion`) for the answer's JSON, one schema for the request (kind, planner model, system text, prompt, token bound and prompt version) and one for the response (the validated answer text, the tokens used, whether the model refused, the cost in cents with its provenance and the judge's version), plus sample requests and answers. The Go surface parses and validates an answer of each kind against its schema; the schemas are plain files the console can read in its tests. The CLI's planner does not change yet.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] Under the core package, a `judge` folder holds the request, response and three answer schemas and one sample per kind; every sample validates against its schema.
- [ ] The Go module exposes the request and response types and a parser that turns an answer's text into the typed answer of its kind, refusing text outside the schema with the field that broke it.
- [ ] The core README names the surface and the rule that the CLI and the console both test against these files.
- [ ] `go vet` and `go test -race` pass in the core module.

## Dev part

**Category:** `core-domain`. What a judgment is, field by field, is the contract the whole feature stands on; the owner should be able to explain every field of an answer.

**Paths:**
- `lib/core/contracts/go/judge/answer.go`
- `lib/core/contracts/go/judge/answer_test.go`

**Boundary:** `judge.ParseAnswer(kind Kind, text string) (Answer, error)`, called by the CLI's planner, where `Kind` and the per-kind answer structs (`SourcesAnswer`, `EpisodesAnswer`, `SuggestionAnswer`) are declared in `types.go` beside it by the agent; the agent also writes the schema files, the samples and the request and response types.

**Goal:** Parse an answer's JSON into the typed answer of its kind and refuse anything else: a missing field, a verdict outside the allowed set, a confidence outside 0 to 1, an effort outside `low`, `medium`, `high`, a reason that is empty. The planner's string kinds map one to one onto the surface's `Kind`.

- [ ] Each sample answer parses into its kind with every field populated.
- [ ] Each refusal names the field and the rule it broke, and the tests cover one refusal per rule.

## Verification

- `cd lib/core/contracts/go && GOTOOLCHAIN=local go test -race ./...`
- `pnpm turbo run test typecheck --filter @langdrift/core`
