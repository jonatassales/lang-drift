#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift check localization" 2
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# Translation catalogs are JSON under each app's messages directory, the Next
# tree's src/app/template/i18n/messages or a Nest app's src/i18n/messages, and
# the loaders beside them hold no translated value
# (.claude/rules/08-localization.rule.md, .claude/rules/11-nestjs.rule.md).

status=0
fail() {
  drift_print_error "$1"
  status=1
}

# A line of a loader that assigns a string literal to a key carries copy. A JSON
# import attribute (`with { type: 'json' }`) describes module loading, not a value.
loader_carries_copy() {
  grep -Ev "^import .* from ['\"][^'\"]+\\.json['\"] with \\{ type: ['\"]json['\"] \\}\$" "$1" |
    grep -Eq "^[[:space:]]*[^:]+:[[:space:]]*['\"]"
}

directories=$(find apps -type d \( -path '*/src/app/template/i18n/messages' -o -path '*/src/i18n/messages' \) -not -path '*/node_modules/*')
[ -n "$directories" ] || drift_die "No messages directory under apps/*/src/app/template/i18n/ or apps/microservices/*/src/i18n/."

for directory in $directories; do
  catalogs=0
  for file in $(find "$directory" -type f -name '*.json' | sort); do
    catalogs=$((catalogs + 1))
    jq empty "$file" 2>/dev/null || fail "Invalid JSON: $file"
  done

  for file in $(find "$directory" -type f \( -name '*.ts' -o -name '*.tsx' -o -name '*.js' -o -name '*.jsx' \) | sort); do
    case "$file" in
      "$directory"/index.ts | "$directory"/messages.types.ts)
        if loader_carries_copy "$file"; then
          fail "Translated values are forbidden in catalog loaders: $file"
        fi
        ;;
      *) fail "Translation catalogs must be JSON: $file" ;;
    esac
  done

  [ "$status" -eq 0 ] && drift_print_success "$directory: $catalogs catalogs are valid JSON and the loaders hold no copy"
done

exit "$status"
