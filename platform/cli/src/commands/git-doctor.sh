#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift git doctor" 2
drift_need git
drift_git_checkout || drift_die "Not a Git checkout. Run git init first."

status=0
hooks_path=$(git -C "$DRIFT_PROJECT_ROOT" config core.hooksPath || true)
if [ "$hooks_path" = "$DRIFT_HOOKS_DIR/_" ]; then
  drift_print_success "Hooks installed (core.hooksPath=$hooks_path)"
else
  drift_print_error "Hooks not installed. Run: drift git setup"
  status=1
fi
for hook in pre-commit commit-msg; do
  if [ -f "$DRIFT_PROJECT_ROOT/$DRIFT_HOOKS_DIR/$hook" ]; then
    drift_print_success "$DRIFT_HOOKS_DIR/$hook"
  else
    drift_print_error "$DRIFT_HOOKS_DIR/$hook missing"
    status=1
  fi
done
exit "$status"
