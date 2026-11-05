# 0295: the console's runtime form keeps one responsibility per file

**Spec:** .agents/workflow/specs/0042-state-rules-compliance.spec.md

**What to build:** The runtime form gets its folder with a types file, a defaults function from the runtime being edited, maps for a configuration field, an applied recommendation and a component toggle, a request map (`PATCH` with the revision for an existing runtime, `POST` with the product id for a new one, the optional budget and native runtime id only when set), a Zod schema replacing the ad-hoc parse, and one handler per event. The console adds `zod`.

**Blocked by:** 0230, 0242.

**Status:** ready-for-agent

- [ ] An invalid build keeps today's alert.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Category:** `core-domain`. Which mutation a runtime form sends, and with which revision, is the console's contract with the runtime record.

**Paths:**
- `apps/console/src/app/[locale]/(menu)/runtimes/runtime-form/runtime-request.map.ts`
- `apps/console/src/app/[locale]/(menu)/runtimes/runtime-form/runtime-request.map.test.ts`

**Boundary:** `toRuntimeRequest(state: RuntimeFormState, submission: RuntimeSubmission, runtime?: Runtime): ConsoleMutation`, called by the submit handler, where `RuntimeFormState` and `RuntimeSubmission` are declared by the agent in `form.types.ts` and `Runtime` and the mutation type stay in the console's types.

**Goal:** Build the mutation from the form: `PATCH` with the revision for an existing runtime, `POST` with the product id for a new one, the selected components, and the optional budget and native runtime id only when set.

- [ ] Tests cover a new runtime, an existing runtime, and each optional field set and unset.
- [ ] The request it builds equals today's for the same inputs.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
