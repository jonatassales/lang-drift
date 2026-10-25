# 0028: Start the platform site from a clone of the website

**Spec:** none yet; a direct owner request on 2026-10-25. The platform spec follows from grilling.

**What to build:** the LangDrift platform site is built in `apps/website`, which starts as a clone of today's langdrift.md. That site moves to `apps/website-old` and stays as the reference until the platform site takes over the domain. The clone copies every file Git tracks, nothing generated. The two apps differ only in their package names and dev ports: `website` on 3000, `website-old` on 2000. Two root scripts start them in dev mode through Turborepo, so the packages they depend on build first: `pnpm website` and `pnpm website-old`. The lockfile registers both workspace packages, and the dead-code allow list covers the formatters and `react-dom` in both.

**Blocked by:** none.

**Status:** ready-for-human

### Acceptance

- [x] `apps/website-old` holds today's website, renamed `website-old`, with `dev` pinned to port 2000.
- [x] `apps/website` holds the 164 tracked files of that website, named `website`, with `dev` on the default port 3000.
- [x] `pnpm website` and `pnpm website-old` exist in the root `package.json` and `.agents/context.md` lists them as the only root scripts that wrap `turbo run`.
- [x] CI ignores `apps/website-old`: the Turborepo tasks filter it out, Biome ignores it and the dead-code allow list covers it, so `drift check dead-code` passes.
- [ ] Pages render locally: each app needs `CRYPTO_SECRET` in its own `.env.local`.
