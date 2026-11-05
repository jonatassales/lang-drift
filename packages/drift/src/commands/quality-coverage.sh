#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -le 1 ] || drift_die "Usage: drift quality coverage [rule]" 2
cd "$DRIFT_PROJECT_ROOT"

# How much of the tree applies each rule file in .claude/rules, as a Markdown
# table: of the files a rule's detectors look at, the share with no finding.
# A rule (05-react.rule.md, react.rule.md or react) narrows it to that row.
# The detectors are patterns, so the percentage is an estimate.

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

sh "$DRIFT_CLI_DIR/quality/run-detectors.sh" "$workdir" "$@" || exit $?

printf '## Rules applied\n\n'
printf '| Rule | Applied |\n'
printf '| --- | ---: |\n'

while IFS= read -r rule; do
  if [ ! -f "$workdir/scope/$rule" ]; then
    printf '| %s | not measured |\n' "$rule"
    continue
  fi

  checked=$(sort -u "$workdir/scope/$rule" | wc -l | tr -d ' ')
  flagged=$(cut -f1 "$workdir/findings/$rule" | sort -u | wc -l | tr -d ' ')
  applied=$(awk -v checked="$checked" -v flagged="$flagged" 'BEGIN { if (checked == 0) print "n/a"; else printf "%.1f%%", (checked - flagged) * 100 / checked }')
  printf '| %s | %s |\n' "$rule" "$applied"
done <"$workdir/rules"
