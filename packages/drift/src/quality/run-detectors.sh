#!/bin/sh
set -eu
# Globs travel to rg and grep as text; the shell never expands them.
set -f
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -ge 1 ] && [ "$#" -le 2 ] || drift_die "Usage: run-detectors.sh <workdir> [rule]" 2
drift_need rg
cd "$DRIFT_PROJECT_ROOT"

# Runs the detectors of packages/drift/src/quality/rule-detectors.tsv, the one
# source drift quality check and drift quality coverage read. Into <workdir>
# it writes rules (the rule files measured, in the detectors' order), and per
# rule file scope/<rule> (the paths its detectors look at) and
# findings/<rule> (each path with a finding, a tab, then the violation). A
# rule names one file of .claude/rules: 05-react.rule.md, react.rule.md or
# react.

workdir=$1
rule_argument=${2:-}
filter=${rule_argument%.rule.md}

quality=$DRIFT_CLI_DIR/quality
detectors=$quality/rule-detectors.tsv

rules=$(awk -F'\t' -v filter="$filter" '
  $1 == "" || $1 ~ /^#/ || seen[$1]++ { next }
  filter == "" || $1 == filter || $1 ~ ("^[0-9]+-" filter "$") { print $1 }' "$detectors")

[ -n "$rules" ] || drift_die "No rule file in .claude/rules matches $rule_argument." 2

printf '%s\n' "$rules" >"$workdir/rules"
mkdir -p "$workdir/scope" "$workdir/findings"

tracked=$workdir/tracked
matches=$workdir/matches
git ls-files >"$tracked"

while IFS='	' read -r rule violation scope mode roots globs pattern; do
  case "$rule" in
    '' | '#'*) continue ;;
  esac

  grep -qx -e "$rule" "$workdir/rules" || continue
  [ "$mode" != unmeasured ] || continue

  set --
  for glob in $(printf '%s' "$globs" | tr ',' ' '); do
    [ "$glob" = - ] || set -- "$@" -g "$glob"
  done

  [ "$scope" != code ] || set -- "$@" -g '!**/*_test.go' -g '!**/*.unit.*' -g '!**/*.test.*' -g '!**/testdata/**'

  for root in $(printf '%s' "$roots" | tr ',' ' '); do
    [ "$root" = - ] || set -- "$@" "$root"
  done

  # The paths the detector looks at.
  scope_status=0

  case "$mode" in
    tracked | paths) cat "$tracked" ;;
    names) grep -v '/testdata/' "$tracked" ;;
    markdown) grep '\.md$' "$tracked" | grep -v '/testdata/' ;;
    *) rg --files "$@" ;;
  esac >"$workdir/listed" || scope_status=$?

  [ "$scope_status" -le 1 ] || drift_die "The detector could not list its files: $rule · $violation"

  cat "$workdir/listed" >>"$workdir/scope/$rule"

  # The paths it finds something in.
  search_status=0

  case "$mode" in
    lines | files) rg -l -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    multiline) rg -l -U -e "$pattern" "$@" >"$matches" || search_status=$? ;;
    per-file)
      rg -c -e "$pattern" "$@" >"$workdir/counts" || search_status=$?
      awk -F: '$NF > 1 { print $1 }' "$workdir/counts" >"$matches"
      ;;
    awk)
      : >"$matches"
      [ ! -s "$workdir/listed" ] ||
        tr '\n' '\0' <"$workdir/listed" | xargs -0 awk -f "$quality/$pattern" >"$matches" || search_status=$?
      ;;
    paths) awk -f "$quality/$pattern" "$tracked" >"$matches" || search_status=$? ;;
    biome-imports)
      : >"$workdir/biome"
      biome_status=0
      [ ! -s "$workdir/listed" ] ||
        tr '\n' '\0' <"$workdir/listed" |
        xargs -0 pnpm exec biome check --linter-enabled=false --formatter-enabled=false \
          --reporter=github --max-diagnostics=none >"$workdir/biome" 2>"$workdir/biome-errors" || biome_status=$?

      # Biome exits 1 when it reports something, which xargs turns into 123;
      # 123 without a report means Biome itself failed.
      biome_reported=false
      ! grep -q '^::' "$workdir/biome" || biome_reported=true
      biome_failed=false
      [ "$biome_status" -eq 0 ] || { [ "$biome_status" -eq 123 ] && [ "$biome_reported" = true ]; } || biome_failed=true

      if [ "$biome_failed" = true ]; then
        tail -n 5 "$workdir/biome-errors" >&2
        search_status=$biome_status
      fi

      sed -n 's/^::error title=assist\/source\/organizeImports,file=\([^,]*\),.*/\1/p' "$workdir/biome" |
        sed "s#^$DRIFT_PROJECT_ROOT/##" >"$matches"
      ;;
    tracked) grep -E -e "$pattern" "$tracked" >"$matches" || search_status=$? ;;
    names)
      grep -v '/testdata/' "$tracked" >"$workdir/names"
      awk -F/ -f "$quality/kebab-case.awk" "$workdir/names" >"$matches" || search_status=$?
      ;;
    markdown)
      grep '\.md$' "$tracked" | grep -v '/testdata/' |
        grep -Ev '(^|/)(README|AGENTS|CLAUDE|SKILL|SECURITY|context)(\.[a-z-]+)?\.md$' |
        grep -Ev '\.(rule|spec|ticket|adr|research|policy|prompt|guide|prd|tdd)\.md$' >"$matches" || search_status=$?
      ;;
    *) drift_die "Unknown detector mode $mode: $rule · $violation" ;;
  esac

  [ "$search_status" -le 1 ] || drift_die "The detector failed: $rule · $violation"

  awk -v violation="$violation" '{ print $0 "\t" violation }' "$matches" >>"$workdir/findings/$rule"
done <"$detectors"

# A measured rule file always has its two files, even when nothing matched.
while IFS= read -r rule; do
  [ ! -f "$workdir/scope/$rule" ] || : >>"$workdir/findings/$rule"
done <"$workdir/rules"
