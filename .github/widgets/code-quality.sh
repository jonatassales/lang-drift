#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/packages/drift/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

# How much of the tree applies each rule file in .claude/rules, as drift
# quality coverage measures it. Its detectors are patterns, and the Code Quality
# check run is an active, required CI gate.

coverage=$("$DRIFT_PROJECT_ROOT/packages/drift/drift.sh" quality coverage) ||
  drift_die "drift quality coverage could not measure the rules."

cat <<EOF_WIDGET
### Code Quality Gate

$(printf '%s\n' "$coverage" | sed '1,2d')
EOF_WIDGET
