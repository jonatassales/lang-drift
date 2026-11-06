# 0024: No container image

**Spec:** .agents/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the repository builds no image: the website's Dockerfile and the root `.dockerignore` are gone, and no document describes them.

**Blocked by:** none.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `apps/website/Dockerfile` and `.dockerignore` are deleted.
- [ ] AGENTS.md no longer describes a website image; ticket 0021 notes that its Dockerfile criterion is superseded.
