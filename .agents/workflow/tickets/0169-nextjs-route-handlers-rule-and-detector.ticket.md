# 0169: Next.js route handlers own their implementation directly

**Spec:** .agents/workflow/specs/0028-route-handlers-own-implementation.spec.md

**What to build:** Add Rule 9 to `.claude/rules/06-nextjs.rule.md` requiring App Router API route handlers (`route.ts`) to declare their HTTP method functions directly in `route.ts` rather than delegating or re-exporting helper functions. Clarify `.claude/rules/00-general.rule.md` to establish that each exported HTTP method in `route.ts` is an entrypoint. Update detector line 13 in `platform/cli/src/quality/rule-detectors.tsv` to exclude `route.ts` from the "More than one function in a TypeScript file" check (`*.ts,*.tsx,!**/route.ts`). Add detector(s) in `rule-detectors.tsv` to forbid hollow route re-exports (`export const POST = ...` or `export { ... as POST }`) and rogue `*.route.ts` files.

**Blocked by:** none

**Status:** ready-for-agent

- [x] `.claude/rules/06-nextjs.rule.md` includes Rule 9 (Route handlers own their implementation).
- [x] `.claude/rules/00-general.rule.md` clarifies that Next.js route handlers export multiple HTTP method entrypoints directly.
- [x] `platform/cli/src/quality/rule-detectors.tsv` excludes `route.ts` from `More than one function in a TypeScript file`.
- [x] Quality detector for `06-nextjs` detects hollow re-exports and `*.route.ts` files.
- [x] `platform/cli/drift quality check` runs cleanly.

### Verification

- `platform/cli/drift quality check`
