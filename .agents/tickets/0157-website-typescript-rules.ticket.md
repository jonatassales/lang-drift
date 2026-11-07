# 0157: The website follows the TypeScript rules

**Spec:** .agents/specs/0027-website-rules-compliance.spec.md

**What to build:** The translator is `translate` everywhere. Object shapes declared as `type` become interfaces. Non-trivial inline object types in signatures, the route props' `Promise<{ locale: string }>` included, become named interfaces declared above their function. Type parameters named `T` get a descriptive name. Every environment variable is read with `process.env.NAME` at the operation that uses it, with an explicit check that treats a missing value as "not configured" as today, and without a helper or `?? ''`.

**Blocked by:** 0156

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] No translator is named `t`.
- [ ] No object shape is a `type` alias, and no signature carries a non-trivial inline object type.
- [ ] No type parameter is a single letter.
- [ ] No environment variable is read through a helper or with a fallback.
- [ ] A missing service value still reads as "not configured", and the tests prove it.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check typescript` and `platform/cli/drift quality check general`
