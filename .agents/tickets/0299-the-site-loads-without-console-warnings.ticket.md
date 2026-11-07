# 0299: the site loads without console warnings

**Spec:** .agents/specs/0046-the-website-focus-pass.spec.md

**What to build:** A visitor opens any page and the browser console stays clean. The Permissions-Policy of the website and the console drops `browsing-topics=()`. The Microruntimes workflow card's selects take the names `model` and `effort`, and its human-stage checkboxes `human-<stage>`. The vendored tabs and accordion draw with Tailwind utilities and lose their stylesheets, as ticket 0273 of spec 0038 sets out, so no stylesheet is prefetched for a page that does not use it. The `eval` report goes with 0300, which takes the only client zod schema out of the bundle; the Content Security Policy keeps `'unsafe-eval'` for development only.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Both apps send `Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=(), usb=()`, asserted by the website's Next configuration test; the README follows.
- [x] Every `select` and checkbox of the Microruntimes workflow card has a `name`.
- [x] Ticket 0273's criteria hold: `animated-tabs.css` and `basic-accordion.css` and their imports are gone, and the components' tests assert their utilities.
- [x] A full load of `/`, `/pricing` and each platform page in the built-in browser shows no Permissions-Policy, preload or form-field issue.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `pnpm turbo run test --filter=console`
