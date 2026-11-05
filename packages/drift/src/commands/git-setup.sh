#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift git setup" 2
drift_git_checkout || drift_die "Not a Git checkout. Run git init first."
[ -d "$DRIFT_PROJECT_ROOT/node_modules/husky" ] || drift_die "Husky is not installed. Run pnpm install."

drift_loader_start "Configuring hooks"
(cd "$DRIFT_PROJECT_ROOT" && pnpm exec husky "$DRIFT_HOOKS_DIR")
drift_print_success "Git hooks configured."
