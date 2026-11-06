# 0154: The early-access form's styles become utilities

**Spec:** .agents/specs/0027-website-rules-compliance.spec.md

**What to build:** `early-access.css` becomes utilities in the early-access form, and is deleted. The form is paused and not rendered, but restoring it brings no stylesheet back. The 900px breakpoint becomes `md` and the 600px one `sm`. The field error and the invalid input's border take the status danger role.

**Blocked by:** 0150

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `early-access.css` is deleted and nothing imports it.
- [ ] The form's label, input, button, feedback and error carry their former styles as utilities, hover, placeholder, invalid and disabled states included.
- [ ] No palette literal remains: the error colors use the status danger role.
- [ ] The section keeps its look, and the captures match the baseline.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check tailwind`
- Built-in browser captures compared with the baseline.
