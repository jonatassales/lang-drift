# 0186: Nest rule and nested-workspace mechanics

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** The repository knows how a Nest app looks and where it lives before the first one lands. A new rule file for Nest apps states the tree, the suffixes (`.module`, `.controller` as a Nest controller class, `.strategy`, `.pipe`, `.filter`, `.token`, `.context`, `.guard` reused), injection by token with decorator metadata off, Zod through one validation pipe, no configuration module, catalogs under the app's `i18n/messages` folder, scripts that call `nest` directly, and unit tests that never boot Nest. The TypeScript suffix table gains `.interceptor` for an API-route domain's `interceptors/` folder. The workspace, the dead-code scan, the localization check and the launch configurations learn the nested `apps/microservices/*` layout, the context document lists it and defines "Event log", and spec 0029 is marked superseded by spec 0031. This is the prefactoring every other ticket of the spec builds on.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] `.claude/rules/11-nestjs.rule.md` exists with the decisions above and is written for the owner to promote.
- [x] The `.interceptor` row is in the suffix table of the TypeScript rules.
- [x] `pnpm-workspace.yaml` includes `apps/microservices/*`, and `pnpm install` succeeds with the folder absent and with an empty folder present.
- [x] `drift quality dead-code` scans `apps/microservices/*` workspaces and their `tsconfig.paths.json`.
- [x] `drift check localization` also finds catalogs under an app's `src/i18n/messages` folder and still passes.
- [x] The launch configurations include `demo-requests` on port 3402 and `mailer` on port 3403.
- [x] `.agents/context.md` lists `apps/microservices/*` in the repository section and adds "Event log" to the glossary.
- [x] Spec 0029 carries a line under its status saying spec 0031 supersedes its email dispatch decision.

### Verification

- `pnpm install`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
- `platform/cli/drift quality dead-code`
- `platform/cli/drift quality check`
