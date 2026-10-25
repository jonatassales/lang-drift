# 0028: Rainbow color group

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** `@langdrift/tokens` gains a fixed color group, `rainbow`, with six ordered stops that are the same in both color schemes and never reference the brand, exposed to Tailwind as `rainbow-1` to `rainbow-6` like the status, classification and drift groups. The header's demo button highlight is its first consumer.

**Blocked by:** none.

**Status:** ready-for-agent

### Tokens

- [ ] `packages/tokens/shared/rainbow.json` declares `lg.rainbow` (`$type: color`) with the stops `1` to `6`: `#F43F5E`, `#F97316`, `#FACC15`, `#22C55E`, `#3B82F6`, `#6D0AD6`.
- [ ] `packages/tokens/build.js` lists `rainbow` among the fixed groups, so `dist/website.css` carries `--lg-rainbow-<n>` under `:root` and `--color-rainbow-<n>` in `@theme inline`.

### Tests

- [ ] `packages/tokens/tokens.test.js` fixes the six stops in light and asserts that dark declares none, beside the classification test.

### Documents

- [ ] `.agents/context.md` names the group in "Design and tokens", and the glossary gains **Rainbow color**.

### Verification

- `pnpm turbo run build test --filter=@langdrift/tokens`
- `grep -c rainbow packages/tokens/dist/website.css` prints 12.
