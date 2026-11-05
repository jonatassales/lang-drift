# 0235: drift measures the Next.js tree and route handlers

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** `drift quality check 06-nextjs` gains: a tracked-path detector for a `lib`, `components`, `views`, `ui`, `hooks` or `data` directory directly under `src/app`; an awk detector for an HTTP method in `route.ts` whose body is one `return` of another function; an awk detector for a function in `route.ts` that is not an HTTP method; a paths detector for a file under `template/formatters/` that is not a `.fmt` file; a line detector for `headers.set` or `cookies.set` in `src/proxy.ts`. With the spec's other tickets landed, it reports nothing under `apps/`.

**Blocked by:** 0229, 0230, 0231, 0232, 0233, 0234.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] Each new detector, run against the commit before 0229 merged, reports that tree's violations, recorded as local evidence under `.audits/`.
- [ ] `platform/cli/drift quality check 06-nextjs` reports nothing under `apps/`.
- [ ] The review-enforced classes (a computed internal `href`, a module's owner, a framework file's meaning) are listed in the quality runner's notes.

## Dev part

**Override:** No part qualifies: pattern rows plus two short awk programs that each match one shape of `route.ts` body; unlike 0236's condition reader, neither has a difficulty of its own (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check 06-nextjs`
- `platform/cli/drift quality coverage 06-nextjs`
