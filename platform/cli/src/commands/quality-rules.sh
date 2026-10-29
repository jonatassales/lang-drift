#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift quality rules" 2
drift_need rg
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# What the tree does against the rules in .claude/rules, as a Markdown report:
# one row per kind of violation with its count, then where each one is. The
# detectors are patterns, so the report is a list to review, not a verdict;
# drift check rules stays the gate. Test files and testdata are skipped, except
# by the testing rules, which only look at tests.

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

summary=$workdir/summary
details=$workdir/details
: >"$summary"
: >"$details"

# One detector per line, fields separated by a tab: the rule file, the
# violation, the scope (code skips tests, tests looks only at them), the mode,
# the roots, the file globs and the rg pattern. The mode says what is counted:
# lines (each matching line), multiline (each match across lines), per-file
# (each file with more than one match), files (each file with a match) and
# paths (each file a glob names; no pattern).
while IFS='	' read -r rule violation scope mode roots globs pattern; do
  [ -n "$rule" ] || continue

  set --
  for glob in $(printf '%s' "$globs" | tr ',' ' '); do
    set -- "$@" -g "$glob"
  done

  if [ "$scope" = code ]; then
    set -- "$@" -g '!**/*_test.go' -g '!**/*.unit.*' -g '!**/*.test.*' -g '!**/testdata/**'
  fi

  for root in $(printf '%s' "$roots" | tr ',' ' '); do
    set -- "$@" "$root"
  done

  matches=$workdir/matches
  found=$workdir/found
  code=0

  case "$mode" in
    lines) rg --json -e "$pattern" "$@" >"$matches" || code=$? ;;
    multiline) rg --json -U -e "$pattern" "$@" >"$matches" || code=$? ;;
    per-file) rg -c -e "$pattern" "$@" >"$matches" || code=$? ;;
    files) rg -l -e "$pattern" "$@" >"$matches" || code=$? ;;
    paths) rg --files "$@" >"$matches" || code=$? ;;
    *) drift_die "Unknown detector mode $mode for: $violation" ;;
  esac

  [ "$code" -le 1 ] || drift_die "rg failed on the detector for: $violation"

  case "$mode" in
    lines | multiline)
      jq -r 'select(.type == "match") | "\(.data.path.text):\(.data.line_number)"' "$matches" >"$found"
      ;;
    per-file) awk -F: '$NF > 1 { print $1 " (" $NF ")" }' "$matches" >"$found" ;;
    files | paths) sort "$matches" >"$found" ;;
  esac

  count=$(wc -l <"$found" | tr -d ' ')
  printf '%s\t%s\t%s\n' "$rule" "$violation" "$count" >>"$summary"

  [ "$count" -gt 0 ] || continue

  printf '\n#### %s · %s (%s)\n\n' "$rule" "$violation" "$count" >>"$details"
  head -n 10 "$found" | sed 's/^/- `/; s/$/`/' >>"$details"
  [ "$count" -le 10 ] || printf -- '- and %s more\n' "$((count - 10))" >>"$details"
done <<'DETECTORS'
00-general	More than one function in a TypeScript file	code	per-file	apps,packages	*.ts,*.tsx	^(export (default )?)?(async )?function [A-Za-z_$]|^(export )?const [A-Za-z_$][\w$]* = (async )?(\([^)]*\)|[A-Za-z_$][\w$]*) =>
00-general	More than one function or type in a Go file	code	per-file	apps,packages	*.go	^(func [A-Za-z_]|type [A-Za-z_]\w* (struct|interface))
00-general	One-letter name in TypeScript	code	lines	apps,packages	*.ts,*.tsx	\b(const|let|var)\s+[a-z]\b\s*[=:]|\(\s*[a-z](\s*,\s*[a-z])*\s*\)\s*=>|(^|[^\w$.])[a-z]\s*=>
00-general	One-letter name in Go	code	lines	apps,packages	*.go	\b[a-z]\s*:=|\bfor\s+[a-z]\s*(,|:=)|^func \([a-z] |\bvar [a-z]
00-general	One-line literal with several fields in TypeScript	code	lines	apps,packages	*.ts,*.tsx	\{\s*[A-Za-z_$][\w$]*\s*:\s*[^,{}\n]+,\s*[A-Za-z_$][\w$]*\s*:[^\n]*\}
00-general	One-line literal with several fields in Go	code	lines	apps,packages	*.go	\b[A-Za-z_]\w*\{[A-Za-z_]\w*:\s*[^{}\n]*,\s*[A-Za-z_]\w*:|\bmap\[[^\]\n]+\][^{\n]*\{[^}\n]*:[^}\n]*,[^}\n]*:
00-general	else after an exit	code	multiline	apps,packages,platform	*.ts,*.tsx,*.go	\b(return|throw|continue|break)\b[^\n]*\n\s*\}\s*else\b
00-general	else after an exit in shell	code	multiline	platform,.github	*.sh	\b(exit|return)\b[^\n]*\n\s*else\b
00-general	Fallback for an environment variable	code	lines	apps,packages	*.ts,*.tsx	process\.env\.[A-Z_][A-Z0-9_]*\s*(\?\?|\|\|)
00-general	Empty catch	code	multiline	apps,packages	*.ts,*.tsx	catch\s*(\([^)]*\))?\s*\{\s*\}
00-general	Error discarded with _	code	lines	apps,packages	*.go	\b_\s*=\s*err\b
01-filesystem	../ in a TypeScript or JavaScript import	code	lines	apps,packages	*.ts,*.tsx,*.js,*.mjs	(\bfrom\s+|\bimport\(|\brequire\()['"]\.\./
01-filesystem	../ in a CSS import	code	lines	apps,packages	*.css	@import\s+['"]\.\./
01-filesystem	../ in a shell path	code	lines	platform,.github	*.sh	^[^#\n]*["'/]\.\./
02-tooling	Automation script that is not shell	code	paths	platform	*.js,*.mjs,*.cjs,*.ts,*.mts	-
02-tooling	npm, npx, yarn or corepack in tooling	code	lines	.github,platform,package.json	*.yml,*.yaml,*.sh,package.json	\b(npx|yarn|corepack)\b|\bnpm (run|install|i|ci|exec)\b
03-tailwind	Stylesheet with hand-written rules	code	files	apps,packages	*.css	^\s*[^@/\s}][^{]*\{
03-tailwind	Palette color class	code	lines	apps,packages	*.ts,*.tsx,*.css	\b(text|bg|border|ring|fill|stroke|from|to|via|outline|decoration|divide|accent|caret|shadow|placeholder)-(red|blue|green|gray|grey|slate|zinc|neutral|stone|amber|yellow|orange|lime|emerald|teal|cyan|sky|indigo|violet|purple|fuchsia|pink|rose)-[0-9]{2,3}\b
03-tailwind	Arbitrary color value	code	lines	apps,packages	*.ts,*.tsx	-\[(#|rgb|hsl|oklch|oklab|color:)
03-tailwind	Theme switched through data-theme	code	lines	apps,packages	*.ts,*.tsx,*.css	data-theme
03-tailwind	@layer components, CSS Module or CSS-in-JS	code	lines	apps,packages	*.ts,*.tsx,*.css	@layer components|\.module\.css|from ['"](styled-components|@emotion/)
04-typescript	Translator named t	code	lines	apps,packages	*.ts,*.tsx	\bconst t = (await )?(use|get)Translations\b
04-typescript	export * in a barrel	code	lines	apps,packages	*.ts,*.tsx	^export \* from
04-typescript	await and ?. in one statement	code	lines	apps,packages	*.ts,*.tsx	\bawait\b[^;\n]*\?\.
04-typescript	Dynamic process.env access	code	lines	apps,packages	*.ts,*.tsx	process\.env\[
05-react	Arrow-function component	code	lines	apps,packages	*.tsx	^(export )?const [A-Z][A-Za-z0-9]* = (async )?(\(|memo\(|forwardRef\()
05-react	Destructured props parameter	code	lines	apps,packages	*.tsx	^(export )?(default )?function [A-Z][A-Za-z0-9]*\(\{
05-react	Inline event handler	code	lines	apps,packages	*.tsx	\bon[A-Z][A-Za-z]*=\{\s*(async\s*)?\(?[\w\s,]*\)?\s*=>
05-react	useEffect call that needs a written reason	code	lines	apps,packages	*.ts,*.tsx	\buseEffect\(
05-react	Loading flag in useState	code	lines	apps,packages	*.ts,*.tsx	const \[(is)?([Ll]oading|[Pp]ending|[Ss]ubmitting|[Ss]aving)\b
05-react	React Context for shared state	code	lines	apps,packages	*.ts,*.tsx	\bcreateContext\(
06-nextjs	Internal link written as <a>	code	lines	apps	*.tsx	<a\s[^>]*href=["']/
06-nextjs	Metadata generated by TypeScript	code	paths	apps	**/app/**/sitemap.ts,**/app/**/robots.ts,**/app/**/manifest.ts,**/app/**/opengraph-image.tsx,**/app/**/twitter-image.tsx,**/app/**/icon.tsx,**/app/**/apple-icon.tsx	-
07-jotai	'use client' in an atom module	code	lines	apps	*.atom.ts	^['"]use client['"]
10-testing	Test description without the file path	tests	lines	apps,packages	*.unit.ts,*.unit.tsx,*.test.ts,*.test.tsx	\bdescribe\(\s*['"][^'"/]*['"]
10-testing	Browser runner dependency	code	lines	package.json,apps,packages	package.json	"(playwright|@playwright/test|puppeteer|axe-core|@axe-core/[^"]+)"\s*:
DETECTORS

# Naming checks that read the tracked paths rather than file contents.
names=$workdir/names
git ls-files | grep -v '/testdata/' | awk -F/ '
  {
    path = ""
    for (part = 1; part <= NF; part++) {
      name = $part
      path = path (part > 1 ? "/" : "") name
      if (name ~ /^(README|AGENTS|CLAUDE|SKILL)(\.[a-z-]+)?\.md$/ || name == "LICENSE") continue
      sub(/^\./, "", name)
      gsub(/[\[\]()@]|\.\.\./, "", name)
      sub(/_test\.go$/, ".go", name)
      sub(/_(linux|darwin|windows|amd64|arm64)\.go$/, ".go", name)
      if (name !~ /^([a-z0-9]+(-[a-z0-9]+)*)?(\.[a-z0-9]+(-[a-z0-9]+)*)*$/ && !(path in seen)) {
        seen[path] = 1
        print path
      }
    }
  }' >"$names"

markdown=$workdir/markdown
git ls-files '*.md' | grep -v '/testdata/' |
  grep -Ev '(^|/)(README|AGENTS|CLAUDE|SKILL)(\.[a-z-]+)?\.md$' |
  grep -Ev '\.(rule|spec|ticket|adr|research|policy|prompt)\.md$' >"$markdown" || true

buckets=$workdir/buckets
git ls-files | awk -F/ '
  {
    path = ""
    for (part = 1; part < NF; part++) {
      path = path (part > 1 ? "/" : "") $part
      if ($part ~ /^(utils|helpers|common|misc)$/ && !(path in seen)) {
        seen[path] = 1
        print path "/"
      }
    }
  }' >"$buckets"

for check in "names|Name not in kebab-case" "markdown|Markdown without a type suffix" "buckets|Bucket directory"; do
  found=$workdir/${check%%|*}
  violation=${check#*|}
  count=$(wc -l <"$found" | tr -d ' ')
  printf '01-filesystem\t%s\t%s\n' "$violation" "$count" >>"$summary"

  [ "$count" -gt 0 ] || continue

  printf '\n#### 01-filesystem · %s (%s)\n\n' "$violation" "$count" >>"$details"
  head -n 10 "$found" | sed 's/^/- `/; s/$/`/' >>"$details"
  [ "$count" -le 10 ] || printf -- '- and %s more\n' "$((count - 10))" >>"$details"
done

total=$(awk -F'\t' '{ total += $3 } END { print total + 0 }' "$summary")

printf '## Rule violations\n\n'
printf 'What the tree does against `.claude/rules`, found by patterns: a list to review, not a verdict. Test files and testdata are skipped, except by the testing rules.\n\n'
printf '| Rule | Violation | Count |\n'
printf '| --- | --- | ---: |\n'
sort -t '	' -k1,1 -s "$summary" | awk -F'\t' '{ printf "| %s | %s | %s |\n", $1, $2, $3 }'
printf '| | **Total** | **%s** |\n' "$total"

[ -s "$details" ] || exit 0

printf '\n### Where\n'
cat "$details"
