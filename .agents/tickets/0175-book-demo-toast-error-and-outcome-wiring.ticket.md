# 0175: Toast error feedback and outcome wiring in BookDemoForm

**Spec:** .agents/specs/0029-resend-book-a-demo-integration.spec.md

**What to build:** Create an accessible, styled Toast primitive in `packages/web/src/ui/toast/` following design tokens. In `apps/website/src/app/lib/components/book-a-demo/form.tsx`, on submission success (HTTP 200), navigate to `/early-access-sent` using `router.push('/early-access-sent')`. On submission failure, trigger the Toast displaying the error message (`copy.form.failed`). Add unit tests in `packages/web/src/ui/toast/toast.unit.tsx` and update `apps/website/src/app/lib/components/book-a-demo/form.unit.tsx`.

**Blocked by:** 0173, 0174

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Toast component created in `packages/web/src/ui/toast/` with props interface, accessible roles, and styling.
- [x] Exported from `packages/web` exports and index.
- [x] `BookDemoForm` navigates to `/early-access-sent` on successful submission.
- [x] `BookDemoForm` triggers Toast on failure.
- [x] Unit tests in `toast.unit.tsx` and `form.unit.tsx`.

### Verification

- `pnpm turbo run check typecheck test build`
- `platform/cli/drift quality check`
