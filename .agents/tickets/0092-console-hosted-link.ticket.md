# 0092: The console reads a hosted extub

**Spec:** .agents/specs/0015-hosted-extub.spec.md

**What to build:** with `EXTUB_TOKEN` and an `https` `EXTUB_URL` set, the console forwards to the hosted extub from its Vercel host, sends the token, and shows a refused change as the preview's rule. Without the token the loopback guards hold as today.

**Blocked by:** 0091.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] One pure rule decides whether a host and an `EXTUB_URL` are allowed, from whether a token is set; `extubRequest` and the server reads use it.
- [ ] In hosted mode the proxy sends `Authorization: Bearer <token>` and never forwards browser credentials.
- [ ] The `errors` catalog gains `read_only_preview`, `token_required` and `stream_unavailable` in the four locales.
- [ ] `EXTUB_TOKEN` joins Turbo's pass-through environment; the README lists the variables and the owner's steps.
- [ ] Node tests for the rule.

### Verification

- `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`; the console against a hosted-mode extub on loopback shows data and refuses a change with the preview's message.
