#!/bin/sh
# shellcheck disable=SC2034 # Its variables are read by the scripts that source it.

# Commit guardrails shared by `drift git commit-msg` (the commit being
# written) and `drift git lint-history` (every commit in a range).
# Each check reads a commit message, or a list of changed paths, on stdin.

DRIFT_AI_SIGNATURE_ERROR='Remove the AI-generated signature, trailer or footer.'
DRIFT_TICKET_ERROR="Changes under apps/ or packages/ need a 'Ticket: <NNNN>' trailer naming a file in .agents/workflow/tickets/."

# True when the message carries an AI-generated signature, trailer or footer.
drift_has_ai_signature() {
  ai_tools='claude|anthropic|copilot|openai|chatgpt|codex|cursor|gemini'
  grep -Eiq "^co-authored-by:.*($ai_tools)|^claude-session:|generated with.*($ai_tools)"
}

# True when a changed path lies under apps/ or packages/.
drift_touches_workspaces() {
  grep -E '^(apps|packages)/' | grep -Evq '/\.gitkeep$'
}

# True when the message's Ticket trailer names a file in the tickets folder.
drift_names_ticket() {
  ticket=$(git interpret-trailers --parse | awk -F': *' 'tolower($1) == "ticket" { print $2; exit }')
  case "$ticket" in
    '' | *[!0-9]*) return 1 ;;
  esac
  for candidate in "$DRIFT_PROJECT_ROOT"/.agents/workflow/tickets/"$ticket"-*.ticket.md; do
    [ -f "$candidate" ] && return 0
  done
  return 1
}
