# 0174: Five-second auto-redirect on sent confirmation page

**Spec:** .agents/workflow/specs/0029-resend-book-a-demo-integration.spec.md

**What to build:** Create a client component `apps/website/src/app/[locale]/(pages)/early-access-sent/auto-redirect.tsx` that sets a 5000ms timer to redirect visitors to the home page (`/`) via `useRouter` from `@i18n/navigation`, clearing the timer on unmount. Include it in `apps/website/src/app/[locale]/(pages)/early-access-sent/page.tsx`. Add unit tests in `auto-redirect.unit.tsx`.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `auto-redirect.tsx` implements 5-second automatic redirection with cleanup on unmount.
- [x] `apps/website/src/app/[locale]/(pages)/early-access-sent/page.tsx` renders `AutoRedirect`.
- [x] `auto-redirect.unit.tsx` tests the 5-second timeout and router navigation.

### Verification

- `pnpm turbo run test --filter=website`
- `platform/cli/drift check biome`
- `platform/cli/drift quality check 00-general 05-react`
