# 0155: The website's state barrel moves to src/app/state

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** The `@state` barrel moves from `src/app/lib/state` to `src/app/state`, as the Next.js rules lay the tree out, and `tsconfig.paths.json` points the alias there. Every importer keeps importing `@state`.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `src/app/state/index.ts` holds the barrel and `src/app/lib/state` is gone.
- [ ] `@state` resolves to the new path, and every importer still uses `@state`.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
