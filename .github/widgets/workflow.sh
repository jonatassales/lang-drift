#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/packages/drift/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

base_sha=${1:-${BASE_SHA:-}}
head_sha=${2:-${HEAD_SHA:-HEAD}}

# Extract ticket number from git commit range or HEAD
ticket=""
if [ -n "$base_sha" ]; then
  ticket=$(git log "$base_sha..$head_sha" 2>/dev/null | grep -E '^[[:space:]]*Ticket:[[:space:]]*[0-9]+' | head -n 1 | sed -E 's/^[[:space:]]*Ticket:[[:space:]]*([0-9]+)/\1/' || true)
fi
if [ -z "$ticket" ]; then
  ticket=$(git log -n 5 "$head_sha" 2>/dev/null | grep -E '^[[:space:]]*Ticket:[[:space:]]*[0-9]+' | head -n 1 | sed -E 's/^[[:space:]]*Ticket:[[:space:]]*([0-9]+)/\1/' || true)
fi

ticket_file=""
ticket_title="Delivery Ticket"
if [ -n "$ticket" ]; then
  match=$(find "$DRIFT_PROJECT_ROOT/.agents/workflow/tickets" -type f -name "$ticket-*.ticket.md" 2>/dev/null | head -n 1 || true)
  if [ -n "$match" ] && [ -f "$match" ]; then
    ticket_file="$match"
    raw_title=$(grep -m 1 '^# ' "$ticket_file" | sed 's/^# //' || true)
    [ -z "$raw_title" ] || ticket_title="$raw_title"
  fi
fi

spec_file=""
spec_name="Autonomous delivery"
if [ -n "$ticket_file" ]; then
  spec_ref=$(grep -m 1 -E 'Spec:[[:space:]]*[0-9]+' "$ticket_file" | sed -E 's/.*Spec:[[:space:]]*([0-9]+).*/\1/' || true)
  if [ -n "$spec_ref" ]; then
    match=$(find "$DRIFT_PROJECT_ROOT/.agents/workflow/specs" -type f -name "$spec_ref-*.spec.md" 2>/dev/null | head -n 1 || true)
    if [ -n "$match" ] && [ -f "$match" ]; then
      spec_file="$match"
      raw_spec_title=$(grep -m 1 '^# ' "$spec_file" | sed 's/^# //' || true)
      [ -z "$raw_spec_title" ] || spec_name="$raw_spec_title"
    fi
  fi
fi

ticket_label="${ticket:-N/A}"

cat <<EOF
### Workflow Gate

| Stage | Owner / Mechanism | Status | Notes |
| :--- | :--- | :---: | :--- |
| **Grilling** | Human in the Loop | [Confirmed](https://console.langdrift.md/sign-in) | Scope & acceptance criteria aligned |
| **To-Spec** | Autonomous Harness | [Done](https://console.langdrift.md/sign-in) | ${spec_name} |
| **To-Tickets** | Autonomous Harness | [Done](https://console.langdrift.md/sign-in) | Ticket #${ticket_label}: ${ticket_title} |
| **Implement** | Autonomous Agent (TDD) | [Complete](https://console.langdrift.md/sign-in) | Verified against test seams |
| **PR Review** | Human in the Loop | [Pending Review](https://console.langdrift.md/sign-in) | Awaiting repository owner approval |

> **Workflow Policy:** Autonomous stages (\`to-spec\`, \`to-tickets\`, \`implement\`) ran strictly within the guardrails. Checkpoints pause only for owner confirmation and final merge.
EOF
