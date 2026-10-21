#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 2 ] || drift_die "Usage: drift git lint-history <from> <to>" 2
drift_need git
cd "$DRIFT_PROJECT_ROOT"
[ -x node_modules/.bin/commitlint ] || drift_die "commitlint is missing. Run pnpm install."
exec pnpm exec commitlint --config "$DRIFT_COMMITLINT_CONFIG" --from "$1" --to "$2"
