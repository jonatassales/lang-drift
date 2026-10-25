#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
. "$DRIFT_CLI_DIR/core/commit-rules.sh"
message=${1:-}
if [ -z "$message" ] || [ "$#" -ne 1 ]; then
  drift_die "Usage: drift git commit-msg <message-file>" 2
fi
cd "$DRIFT_PROJECT_ROOT"

# 1. No AI-generated signature, trailer or footer.
if drift_has_ai_signature < "$message"; then
  drift_die "$DRIFT_AI_SIGNATURE_ERROR"
fi

# 2. Changes under apps/ or packages/ name the ticket they implement.
if git diff --cached --name-only | drift_touches_workspaces && ! drift_names_ticket < "$message"; then
  drift_die "$DRIFT_TICKET_ERROR"
fi

# 3. Conventional Commits.
[ -x node_modules/.bin/commitlint ] || drift_die "commitlint is missing. Run pnpm install."
exec pnpm exec commitlint --config "$DRIFT_COMMITLINT_CONFIG" --edit "$message"
