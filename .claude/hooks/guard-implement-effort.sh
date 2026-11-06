#!/bin/sh
# PreToolUse guardrail for the Skill tool.
# Blocks the "implement" skill unless an LLM-effort analysis has been
# recorded for the current session.
#
# Needs: jq.
#
# Marker convention: the llm-effort-analyzer skill
# (.claude/skills/llm-effort-analyzer/SKILL.md) ends by creating
# .agents/.effort-log/<session_id>.md. This hook checks that file
# exists before letting "implement" run.
#
# Unverified assumption, flagged on purpose: tool_input.name is the field
# holding the invoked skill's name. This is inferred from the documented
# permission-rule syntax "Skill(name)" (same pattern as Bash(command) ->
# tool_input.command), not from a literal tool_input example in the docs.
# Verify once a skill literally named "implement" exists, by temporarily
# logging $INPUT to a file from this hook.

set -eu

if ! command -v jq >/dev/null 2>&1; then
  echo '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"guard-implement-effort.sh needs jq, which is not installed. Install jq or remove this hook."}}'
  exit 0
fi

INPUT="$(cat)"

TOOL_NAME="$(printf '%s' "$INPUT" | jq -r '.tool_name // empty')"
if [ "$TOOL_NAME" != "Skill" ]; then
  exit 0
fi

SKILL_NAME="$(printf '%s' "$INPUT" | jq -r '.tool_input.name // .tool_input.skill // .tool_input.command // empty')"
case "$SKILL_NAME" in
  implement|*:implement) : ;;
  *) exit 0 ;;
esac

SESSION_ID="$(printf '%s' "$INPUT" | jq -r '.session_id // empty')"
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
MARKER="$PROJECT_DIR/.agents/.effort-log/$SESSION_ID.md"

if [ -f "$MARKER" ]; then
  exit 0
fi

REASON="No LLM-effort analysis is recorded for this session. Run the llm-effort-analyzer skill, which records its result at .agents/.effort-log/$SESSION_ID.md, before running implement."
jq -n --arg reason "$REASON" '{hookSpecificOutput:{hookEventName:"PreToolUse",permissionDecision:"deny",permissionDecisionReason:$reason}}'
exit 0
