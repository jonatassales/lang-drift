# 0181: Generate CRYPTO_SECRET on the fly in development

**What to build:** Remove the hardcoded `CRYPTO_SECRET` from `apps/website/.env.development`. In `apps/website/src/factories/cookie-service.factory.ts`, inspect the environment through `process.env`: when `CRYPTO_SECRET` is not set and the environment is development or test, generate a canonical 32-byte base64 secret on the fly and persist it in memory for the process lifetime. In production, preserve the requirement that `CRYPTO_SECRET` must be explicitly configured. Add unit tests for the factory in `cookie-service.factory.unit.ts` and update documentation.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Remove `CRYPTO_SECRET` from `apps/website/.env.development`.
- [x] Update `apps/website/src/factories/cookie-service.factory.ts` to generate `CRYPTO_SECRET` on the fly in development/test when unset.
- [x] Add unit tests in `apps/website/src/factories/cookie-service.factory.unit.ts`.
- [x] Update `apps/website/README.md` and `.agents/context.md` variable contracts.
- [x] Verify all tests and checks pass with `pnpm turbo run check typecheck test build`.
