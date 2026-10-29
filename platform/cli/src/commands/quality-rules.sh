#!/bin/sh
set -eu
# Globs travel to rg and grep as text; the shell never expands them.
set -f
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift quality rules" 2
drift_need rg
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# What the tree does against the rules in .claude/rules, as a Markdown report:
# one row per kind of violation with its count, then where each one is. The
# detectors (platform/cli/src/quality/rule-detectors.tsv) are patterns, so the
# report is a list to review, not a verdict; drift check rules stays the gate.

detectors=$DRIFT_CLI_DIR/quality/rule-detectors.tsv
kebab_case=$DRIFT_CLI_DIR/quality/kebab-case.awk

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

tracked=$workdir/tracked
git ls-files >"$tracked"

summary=$workdir/summary
details=$workdir/details
matches=$workdir/matches
found=$workdir/found
: >"$summary"
: >"$details"

while IFS='	' read -r rule violation scope mode roots globs pattern; do
  case "$rule" in
    '' | '#'*) continue ;;
  esac

  if [ "$mode" = unmeasured ]; then
    printf '%s\t%s\tnot measured\n' "$rule" "$violation" >>"$summary"
    continue
  fi

  set --
  for glob in $(printf '%s' "$globs" | tr ',' ' '); do
    [ "$glob" = - ] || set -- "$@" -g "$glob"
  done

  [ "$scope" != code ] || set -- "$@" -g '!**/*_test.go' -g '!**/*.unit.*' -g '!**/*.test.*' -g '!**/testdata/**'

  for root in $(printf '%s' "$roots" | tr ',' ' '); do
    [ "$root" = - ] || set -- "$@" "$root"
  done

  search_status=0

  case "$mode" in
    lines) rg --json -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    multiline) rg --json -U -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    per-file) rg -c -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    files) rg -l -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    tracked) grep -E -e "$pattern" "$tracked" >"$matches" || search_status=$? ;;
    names) grep -v '/testdata/' "$tracked" | awk -F/ -f "$kebab_case" >"$matches" ;;
    markdown)
      grep '\.md$' "$tracked" | grep -v '/testdata/' |
        grep -Ev '(^|/)(README|AGENTS|CLAUDE|SKILL)(\.[a-z-]+)?\.md$' |
        grep -Ev '\.(rule|spec|ticket|adr|research|policy|prompt)\.md$' >"$matches" || search_status=$?
      ;;
    *) drift_die "Unknown detector mode $mode: $rule · $violation" ;;
  esac

  [ "$search_status" -le 1 ] || drift_die "The detector failed: $rule · $violation"

  case "$mode" in
    lines | multiline) jq -r 'select(.type == "match") | "\(.data.path.text):\(.data.line_number)"' "$matches" >"$found" ;;
    per-file) awk -F: '$NF > 1 { print $1 " (" $NF ")" }' "$matches" >"$found" ;;
    *) sort -u "$matches" >"$found" ;;
  esac

  count=$(wc -l <"$found" | tr -d ' ')
  printf '%s\t%s\t%s\n' "$rule" "$violation" "$count" >>"$summary"

  [ "$count" -gt 0 ] || continue

  printf '\n#### %s · %s (%s)\n\n' "$rule" "$violation" "$count" >>"$details"
  head -n 10 "$found" | sed 's/^/- `/; s/$/`/' >>"$details"
  [ "$count" -le 10 ] || printf -- '- and %s more\n' "$((count - 10))" >>"$details"
done <"$detectors"

total=$(awk -F'\t' '$3 ~ /^[0-9]+$/ { total += $3 } END { print total + 0 }' "$summary")

printf '## Rule violations\n\n'
printf 'What the tree does against `.claude/rules`, found by patterns: a list to review, not a verdict. Test files and testdata are skipped, except by the testing rules.\n\n'
printf '| Rule | Violation | Count |\n'
printf '| --- | --- | ---: |\n'
awk -F'\t' '{ printf "| %s | %s | %s |\n", $1, $2, $3 }' "$summary"
printf '| | **Total** | **%s** |\n' "$total"

[ -s "$details" ] || exit 0

printf '\n### Where\n'
cat "$details"
