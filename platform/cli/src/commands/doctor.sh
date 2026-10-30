#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"

ci=false
if [ "${1:-}" = --ci ]; then
  ci=true
  shift
fi
[ "$#" -eq 0 ] || drift_die "Usage: drift doctor [--ci]" 2
if [ "$ci" = false ]; then
  drift_loader_start "Checking environment"
  drift_queue_scene "$DRIFT_ICON_DOCTOR" "$DRIFT_DOCTOR_PHRASE"
fi

status=0
fail() {
  drift_print_error "$1"
  status=1
}

for tool in git node pnpm; do
  if drift_has "$tool"; then drift_print_success "$tool available"; else fail "$tool missing"; fi
done

# require_major <tool> <minimum> <version>
require_major() {
  major=${3#v}
  major=${major%%.*}
  case "$major" in
    '' | *[!0-9]*) fail "Unable to read the $1 version: $3" ;;
    *) [ "$major" -ge "$2" ] || fail "$1 $2+ required; found $3" ;;
  esac
}
if drift_has node; then require_major Node 24 "$(node --version)"; fi
if drift_has pnpm; then require_major pnpm 12 "$(pnpm --version)"; fi

required_paths="AGENTS.md apps packages platform/cli/drift $DRIFT_HOOKS_DIR .agents .claude"
for path in $required_paths; do
  if [ -e "$DRIFT_PROJECT_ROOT/$path" ]; then drift_print_success "$path"; else fail "Missing $path"; fi
done

if [ "$ci" = false ]; then
  if drift_git_checkout; then
    branch=$(git -C "$DRIFT_PROJECT_ROOT" branch --show-current 2>/dev/null || true)
    drift_print_info "Git branch: ${branch:-detached}"
  else
    drift_print_warning "Not a Git checkout yet: run git init so the hooks can install."
  fi
fi

exit "$status"
