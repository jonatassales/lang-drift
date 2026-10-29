#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -gt 0 ] || drift_die "Usage: drift ci summary <run-summary.json>..." 2
drift_need jq

# Each Turborepo run summary as a Markdown table, for $GITHUB_STEP_SUMMARY. A job
# that failed before turbo ran has none, so an unmatched glob prints nothing.
for summary in "$@"; do
  [ -f "$summary" ] || continue
  jq -r -f "$DRIFT_PROJECT_ROOT/platform/observability/turbo-table.jq" "$summary"
done
