#!/bin/sh
set -eu
# Globs travel to rg and grep as text; the shell never expands them.
set -f
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift quality coverage" 2
drift_need rg
cd "$DRIFT_PROJECT_ROOT"

# How much of the tree applies each rule file in .claude/rules, as a Markdown
# table: of the files a rule's detectors look at, the share with no finding.
# The detectors (platform/cli/src/quality/rule-detectors.tsv) are patterns, so
# the percentage is an estimate; drift quality check stays the gate.

detectors=$DRIFT_CLI_DIR/quality/rule-detectors.tsv
kebab_case=$DRIFT_CLI_DIR/quality/kebab-case.awk

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

tracked=$workdir/tracked
git ls-files >"$tracked"

matches=$workdir/matches
mkdir "$workdir/scope" "$workdir/findings"

while IFS='	' read -r rule violation scope mode roots globs pattern; do
  case "$rule" in
    '' | '#'*) continue ;;
  esac

  [ "$mode" != unmeasured ] || continue

  set --
  for glob in $(printf '%s' "$globs" | tr ',' ' '); do
    [ "$glob" = - ] || set -- "$@" -g "$glob"
  done

  [ "$scope" != code ] || set -- "$@" -g '!**/*_test.go' -g '!**/*.unit.*' -g '!**/*.test.*' -g '!**/testdata/**'

  for root in $(printf '%s' "$roots" | tr ',' ' '); do
    [ "$root" = - ] || set -- "$@" "$root"
  done

  # The files the detector looks at, then the files it finds something in.
  scope_status=0

  case "$mode" in
    tracked) cat "$tracked" ;;
    names) grep -v '/testdata/' "$tracked" ;;
    markdown) grep '\.md$' "$tracked" | grep -v '/testdata/' ;;
    *) rg --files "$@" ;;
  esac >>"$workdir/scope/$rule" || scope_status=$?

  [ "$scope_status" -le 1 ] || drift_die "The detector could not list its files: $rule · $violation"

  search_status=0

  case "$mode" in
    lines | files) rg -l -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    multiline) rg -l -U -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    per-file)
      rg -c -e "$pattern" "$@" >"$workdir/counts" || search_status=$?
      awk -F: '$NF > 1 { print $1 }' "$workdir/counts" >"$matches"
      ;;
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

  cat "$matches" >>"$workdir/findings/$rule"
done <"$detectors"

printf '## Rules applied\n\n'
printf '| Rule | Applied |\n'
printf '| --- | ---: |\n'

# One row per rule file, in the detectors' order.
rules=$(awk -F'\t' '$1 != "" && $1 !~ /^#/ && !seen[$1]++ { print $1 }' "$detectors")

for rule in $rules; do
  if [ ! -f "$workdir/scope/$rule" ]; then
    printf '| %s | not measured |\n' "$rule"
    continue
  fi

  checked=$(sort -u "$workdir/scope/$rule" | wc -l | tr -d ' ')
  flagged=$(sort -u "$workdir/findings/$rule" | wc -l | tr -d ' ')
  applied=$(awk -v checked="$checked" -v flagged="$flagged" 'BEGIN { if (checked == 0) print "n/a"; else printf "%.1f%%", (checked - flagged) * 100 / checked }')
  printf '| %s | %s |\n' "$rule" "$applied"
done
