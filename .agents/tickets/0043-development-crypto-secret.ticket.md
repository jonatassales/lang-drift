# 0043: Track a development-only CRYPTO_SECRET

**Spec:** None

**What to build:** a fresh checkout runs the website without creating `.env.local`. The tracked `apps/website/.env.development` carries a development-only `CRYPTO_SECRET`, so the locale proxy works out of the box; deployed environments keep setting their own key.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `.env.development` holds a 32-byte canonical base64 key, with a comment saying it is development-only and that no deployed environment may reuse it.
- [ ] The website README and `.agents/context.md` describe the tracked development key instead of saying the file never holds a value or that `.env.local` is required locally.
- [ ] `CRYPTO_SECRET` stays pass-through only in Turborepo; nothing reads it as `NEXT_PUBLIC_`.

### Verification

- The dev server starts without `.env.local`; `/` redirects to `/en/` and the home page renders.
- `pnpm turbo run check`.
