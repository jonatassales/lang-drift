#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
. "$DRIFT_CLI_DIR/core/commit-rules.sh"
[ "$#" -eq 2 ] || drift_die "Usage: drift git lint-history <from> <to>" 2
drift_need git
cd "$DRIFT_PROJECT_ROOT"
[ -x node_modules/.bin/commitlint ] || drift_die "commitlint is missing. Run pnpm install."

# The commit-msg guardrails, applied to every commit in the range, so a
# commit made with --no-verify cannot skip them.
status=0
for commit in $(git rev-list --reverse "$1..$2"); do
  commit_title=$(git log -1 --format='%h %s' "$commit")

  if git log -1 --format=%B "$commit" | drift_has_ai_signature; then
    drift_print_error "$commit_title: $DRIFT_AI_SIGNATURE_ERROR"
    status=1
  fi

  if git diff-tree --root --no-commit-id --name-only -r "$commit" | drift_touches_workspaces &&
    ! git log -1 --format=%B "$commit" | drift_names_ticket; then
    drift_print_error "$commit_title: $DRIFT_TICKET_ERROR"
    status=1
  fi
done

pnpm exec commitlint --from "$1" --to "$2" || status=1
exit "$status"
