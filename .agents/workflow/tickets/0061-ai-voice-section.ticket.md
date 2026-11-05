# 0061: The AI Voice section plays the orb on the home

**Spec:** .agents/workflow/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** a fourth section in the same skeleton for AI Voice: its icon and name at the right, the voice title and description, the orb with the sample question and its scripted answer, no link because the product has no page. The preview reads correctly in the light scheme.

**Blocked by:** 0060 (The home section shell and the Vision, Microruntimes and FinOps sections).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `VoicePreview` is trimmed to the orb, the question row and the answer card; its kicker, heading and paragraph leave, since the section renders the voice title and description from the home catalog.
- [ ] Its styling is role utilities: borders in the border role, muted text in the muted role, the button in the primary role, the answer card's fade and slide still under reduced motion; its stylesheet is deleted and its line leaves the frontend rule's shrink-only list.
- [ ] The kicker's dark modifier class leaves the global stylesheet, and the voice preview's two lines leave the dead-code allow list.
- [ ] The home renders the AI Voice section fourth, with the sample answer playing on request and the accessible status announcing each state.
- [ ] The voice preview test keeps every state and timer case.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check dead-code`.
- `/en/` in light and dark: the orb and the answer card read in both schemes; the answer appears without sliding under reduced motion.
