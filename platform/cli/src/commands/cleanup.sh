#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Cleanup does not accept arguments." 2
drift_loader_start "Cleaning workspace"
drift_queue_scene "$DRIFT_ICON_CLEANUP" "$DRIFT_CLEANUP_PHRASE"

cd "$DRIFT_PROJECT_ROOT"
find . -type d \( -name node_modules -o -name .next -o -name .turbo -o -name dist -o -name coverage \) -prune -exec rm -rf {} + 2>/dev/null || true
drift_print_success "Generated workspace state removed."
