# 0281: the website files carry the suffix of their responsibility

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Every TypeScript file in the website that the earlier splits did not already rename takes the suffix of its responsibility, by the spec's calls: metadata companions `.map` with the framework stem, I/O functions `.adapter`, pure ones `.compute`, commands `.action`, handles `.controller`, fresh instances `.factory`, error classes `.type`, module-level mutable objects `.state`, constants `.data`; a re-export-only module that is not `index.ts` is dissolved unless it is a package export entry.

**Blocked by:** 0239, 0240, 0279, 0277.

**Status:** ready-for-agent

- [ ] The file-suffix detector reports nothing in the website.
- [ ] No import path breaks and no export is renamed beyond the spec's calls.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check 04-typescript`
