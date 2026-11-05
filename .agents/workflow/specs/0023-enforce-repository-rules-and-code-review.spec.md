# Enforce repository rules and code-review

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

Coding agents and developers modifying code in the Tao monorepo may inadvertently introduce plain CSS stylesheets, invalid directory structures, or bypass the two-axis code review process. While `.claude/rules/*.rule.md` documents architectural standards (including CSS bans, ratchet limits on legacy stylesheets, Next.js tree layout, and POSIX shell scripts for tooling), adherence was previously reliant on voluntary agent compliance without hard behavioral or mechanical gates.

## Solution

1. **Rule-level Enforcement:**
   - Add Rule 4 to `.claude/rules/00-general.rule.md` requiring mandatory two-axis `code-review` before every commit or delivery.
   - Update `.claude/skills/code-review/SKILL.md` to explicitly target all 10 rule files in `.claude/rules/*.rule.md` as the authoritative source of standards.
   - Update `.claude/skills/implement/SKILL.md` to require calling `code-review` and resolving all findings before delivery.

2. **Mechanical Guardrail (`drift check rules`):**
   - Implement `platform/cli/src/commands/check-rules.sh` to mechanically verify:
     - `00-general`: Tooling/platform scripts are POSIX shell (`.sh`); only `pnpm-lock.yaml` is present.
     - `01-filesystem`: Workflow markdown files have approved type suffixes (`.spec.md`, `.ticket.md`, `.adr.md`, `.research.md`) and rule documents end in `.rule.md`.
     - `02-frontend` & `03-tailwind`: Plain CSS whitelist is strictly enforced (rejects any unapproved `.css` file) and ratchet ceilings are enforced (legacy CSS files cannot increase in line count).
     - `04-code`: Centralized `env.ts` / `getEnv.ts` modules are forbidden.
     - `05-react`: Shared UI barrels (`packages/web`) never use `export * from`.
     - `06-nextjs`: Forbidden utility directories (`utils`, `helpers`, `common`, `misc`) and Pages Router trees are rejected; private app manifests carry no `version`.
     - `07-jotai`: Shared atoms in `atoms/` end with `.atom.ts` and contain no `'use client'`.

3. **Pre-Commit and CI Integration:**
   - Wire `drift check rules` into `platform/cli/src/commands/git-pre-commit.sh` so local commits fail if rules are broken.
   - Wire `platform/cli/drift check rules` into the CI `check` job in `.github/workflows/validate.yml`.
   - Wire `check biome` into `platform/cli/src/commands/check-biome.sh` and CI.

## Testing Decisions

- **Seam 1 — Drift Rules Check:** `platform/cli/drift check rules` executes and passes cleanly across all workspaces.
- **Seam 2 — Pre-commit execution:** `platform/cli/drift git pre-commit` executes `lint-staged` and `check rules`.
- **Seam 3 — CI Workflow Lint:** `platform/cli/drift check workflows` passes.
- **Seam 4 — Typecheck and Biome Check:** `pnpm turbo run check` and `pnpm turbo run typecheck` pass with zero errors.
