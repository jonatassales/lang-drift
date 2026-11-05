#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Reset does not accept arguments." 2
drift_need git
drift_need pnpm
drift_loader_start "Resetting workspace"
drift_queue_scene "$DRIFT_ICON_CLEANUP" "$DRIFT_RESET_PHRASE"

cd "$DRIFT_PROJECT_ROOT"
removable=$(mktemp)
trap 'rm -f "$removable"; drift_flush_results' EXIT

# Everything cleanup removes.
drift_find_generated_state >"$removable"
drift_remove_paths "$removable"
drift_find_redundant_gitkeeps >"$removable"
drift_remove_paths "$removable"

# Then what pnpm install recreates and the local configuration: the Git hook
# shims, every .env file Git does not track, the local Claude Code settings and
# the lockfile.
{
  find . -path ./.git -prune -o -type d -path '*/.husky/_' -prune -print
  find . -path ./.git -prune -o -type f -name '.env*' -print |
    while IFS= read -r env_file; do
      git ls-files --error-unmatch -- "$env_file" >/dev/null 2>&1 || printf '%s\n' "$env_file"
    done
  for local_file in ./.claude/settings.local.json ./pnpm-lock.yaml; do
    [ ! -e "$local_file" ] || printf '%s\n' "$local_file"
  done
} >"$removable"
drift_remove_paths "$removable"

drift_print_success "Workspace reset: installing from scratch."
rm -f "$removable"
trap - EXIT
drift_flush_results
exec pnpm install
