#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
message=${1:-}
if [ -z "$message" ] || [ "$#" -ne 1 ]; then
  drift_die "Usage: drift git commit-msg <message-file>" 2
fi
cd "$DRIFT_PROJECT_ROOT"

# 1. No AI-generated signature, trailer or footer.
ai_tools='claude|anthropic|copilot|openai|chatgpt|codex|cursor|gemini'
if grep -Eiq "^co-authored-by:.*($ai_tools)|^claude-session:|generated with.*($ai_tools)" "$message"; then
  drift_die "Remove the AI-generated signature, trailer or footer."
fi

# 2. Changes under apps/ or packages/ name the ticket they implement.
if git diff --cached --name-only | grep -E '^(apps|packages)/' | grep -Evq '/\.gitkeep$'; then
  ticket=$(git interpret-trailers --parse "$message" | awk -F': *' 'tolower($1) == "ticket" { print $2; exit }')
  ticket_file=
  case "$ticket" in
    '' | *[!0-9]*) ;;
    *)
      for candidate in .agents/workflow/tickets/"$ticket"-*.ticket.md; do
        [ -f "$candidate" ] && ticket_file=$candidate && break
      done
      ;;
  esac
  [ -n "$ticket_file" ] || drift_die "Changes under apps/ or packages/ need a 'Ticket: <NNNN>' trailer naming a file in .agents/workflow/tickets/."
fi

# 3. Conventional Commits.
[ -x node_modules/.bin/commitlint ] || drift_die "commitlint is missing. Run pnpm install."
exec pnpm exec commitlint --config "$DRIFT_COMMITLINT_CONFIG" --edit "$message"
