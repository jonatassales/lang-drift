# 0028: Clone the website as the platform-site sandbox

**Spec:** none yet; a direct owner request on 2026-10-25. The platform spec follows from grilling.

**What to build:** `apps/website-lab` is a clone of `apps/website`, where the LangDrift platform site is built while `apps/website` keeps serving today's langdrift.md. The clone copies every file Git tracks in the website, nothing generated, and differs only in its package name, `website-lab`, and its dev port, 2000. Two root scripts start the apps in dev mode through Turborepo, so the packages they depend on build first: `pnpm website` on port 3000 and `pnpm website-lab` on port 2000. The lockfile registers the new workspace package.

**Blocked by:** none.

**Status:** ready-for-human

### Acceptance

- [x] `apps/website-lab` holds the 164 files Git tracks in `apps/website`, with `package.json` renamed and `dev` pinned to port 2000.
- [x] `pnpm website` and `pnpm website-lab` exist in the root `package.json` and `.agents/context.md` lists them as the only root scripts that wrap `turbo run`.
- [x] `pnpm website-lab` starts Next.js on `http://localhost:2000`.
- [ ] Pages render locally: they need `CRYPTO_SECRET` in `apps/website-lab/.env.local`, as for the website.
