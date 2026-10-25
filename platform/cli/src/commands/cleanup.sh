#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Cleanup does not accept arguments." 2
drift_loader_start "Cleaning workspace"
drift_queue_scene "$DRIFT_ICON_CLEANUP" "$DRIFT_CLEANUP_PHRASE"

cd "$DRIFT_PROJECT_ROOT"
removable=$(mktemp)
trap 'rm -f "$removable"; drift_flush_results' EXIT

# Generated state first, so each .gitkeep is judged on what remains beside it.
drift_find_generated_state >"$removable"
drift_remove_paths "$removable"
drift_find_redundant_gitkeeps >"$removable"
drift_remove_paths "$removable"

if [ "$drift_removed_count" -eq 0 ]; then
  drift_print_info "Nothing to remove."
else
  drift_print_success "Generated workspace state removed."
fi
