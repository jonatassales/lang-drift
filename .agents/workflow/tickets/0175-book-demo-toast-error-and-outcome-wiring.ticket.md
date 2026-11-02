# 0175: Toast error feedback and outcome wiring in BookDemoForm

**Spec:** .agents/workflow/specs/0029-resend-book-a-demo-integration.spec.md

**What to build:** Create an accessible, styled Toast primitive in `packages/web/src/ui/toast/` following design tokens. In `apps/website/src/app/lib/components/book-a-demo/form.tsx`, on submission success (HTTP 200), navigate to `/early-access-sent` using `router.push('/early-access-sent')`. On submission failure, trigger the Toast displaying the error message (`copy.form.failed`). Add unit tests in `packages/web/src/ui/toast/toast.unit.tsx` and update `apps/website/src/app/lib/components/book-a-demo/form.unit.tsx`.

**Blocked by:** 0173, 0174

**Status:** ready-for-agent

- [ ] Toast component created in `packages/web/src/ui/toast/` with props interface, accessible roles, and styling.
- [ ] Exported from `packages/web` exports and index.
- [ ] `BookDemoForm` navigates to `/early-access-sent` on successful submission.
- [ ] `BookDemoForm` triggers Toast on failure.
- [ ] Unit tests in `toast.unit.tsx` and `form.unit.tsx`.

### Verification

- `pnpm turbo run check typecheck test build`
- `platform/cli/drift quality check`
