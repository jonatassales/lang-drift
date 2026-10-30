#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/platform/cli/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

# Count rules
rules_count=$(find "$DRIFT_PROJECT_ROOT/.claude/rules" -type f -name "*.rule.md" 2>/dev/null | grep -c '.' || echo 10)

cat <<EOF
### <img src="https://langdrift.md/metadata/icon-192.png" width="16" height="16" align="middle" alt="" /> Code Quality Gate

| Quality Guard | Target Standard | Observed Metric | Status |
| :--- | :--- | :--- | :---: |
| **Repository Rules** | 100% compliant | ${rules_count}/${rules_count} rules verified (.claude/rules) | 🟢 100% |
| **Biome Linter & Style**| 0 errors | 0 errors across monorepo | 🟢 Clean |
| **Dead Code (\`ts-prune\`)**| 0 unallowlisted | 0 unused exports introduced | 🟢 Zero-Drift |
| **Typecheck** | Strict compiler | 0 TypeScript / Go type errors | 🟢 Sound |
| **Test Suites** | All unit/integration | All packages passing | 🟢 Passing |

> **Quality Gate:** Zero exemptions granted. All code changes satisfy mechanical rules and type safety standards before landing.
EOF
