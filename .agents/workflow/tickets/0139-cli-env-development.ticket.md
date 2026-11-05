# 0139: Drift CLI and documentation support .env.development

**Spec:** None

**What to build:** The `drift env setup` and `drift env validate` commands recognize `.env.development` as the tracked environment template instead of `.env.example`, creating and validating `.env`. `.gitignore` ignores only `.env`, and documentation in `apps/website/README.md` and `.agents/context.md` reflects this contract.

**Blocked by:** None

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `platform/cli/src/core/common.sh` discovers `.env.development` instead of `.env.example`.
- [x] `platform/cli/src/commands/env-setup.sh` and `platform/cli/src/commands/env-validate.sh` target `.env` from `.env.development`.
- [x] `platform/cli/src/commands/help.sh` describes `env setup` and `env validate` using `.env.development`.
- [x] `.gitignore` ignores only `.env`.
- [x] `apps/website/README.md` and `.agents/context.md` document the `.env.development` template and `drift env` commands.

### Verification

- `drift env validate` reports missing `.env` when not present.
- `drift env setup` creates `.env` with 0600 permissions.
- `drift quality check`, `drift check biome`, `drift check dead-code`, and `drift check workflows` pass.
