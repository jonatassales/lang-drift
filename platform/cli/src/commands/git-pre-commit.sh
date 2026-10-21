#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift git pre-commit" 2
cd "$DRIFT_PROJECT_ROOT"

# Biome checks the staged files through lint-staged.
[ -x node_modules/.bin/lint-staged ] || drift_die "lint-staged is missing. Run pnpm install."
exec pnpm exec lint-staged --config "$DRIFT_LINT_STAGED_CONFIG"
