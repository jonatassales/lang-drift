#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift check dead-code" 2
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# Code that only served something a refactor removed: an alias nothing imports
# through, a module or export nothing uses, a dependency nothing imports, a
# turbo.json variable nothing reads, a CSS class no markup carries. A reference
# in a comment counts, so code parked in comments keeps what it needs. Folders
# held by a .gitkeep are never examined. platform/quality/dead-code.allow lists
# what the owner keeps on purpose.

allowlist=platform/quality/dead-code.allow
findings=$(mktemp)
sources=$(mktemp)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -f "$findings" "$sources"; drift_flush_results' EXIT

# allowed <kind> <key>
allowed() {
  [ -f "$allowlist" ] || return 1
  while read -r allowed_kind allowed_pattern _; do
    [ "$allowed_kind" = "$1" ] || continue
    # shellcheck disable=SC2254 # The allowlist pattern is a glob on purpose.
    case "$2" in $allowed_pattern) return 0 ;; esac
  done <"$allowlist"
  return 1
}

# report <kind> <key> <message>
report() {
  allowed "$1" "$2" && return 0
  printf '%s\n' "$3" >>"$findings"
}

escape_regex() {
  printf '%s' "$1" | sed 's/[.[\*^$()+?{}|]/\\&/g'
}

# Lines of a file outside `export { … } from` statements, which only forward a
# name and never use it.
without_reexports() {
  awk '
    /^export (type )?\{/ && !/\}/ { forwarding = 1; next }
    forwarding && /\} from/ { forwarding = 0; next }
    forwarding { next }
    /^export (type )?\{.*\} from/ { next }
    { print }
  ' "$1"
}

is_framework_file() {
  case "${1##*/}" in
    page.tsx | layout.tsx | template.tsx | loading.tsx | error.tsx | global-error.tsx | \
      not-found.tsx | default.tsx | route.ts | proxy.ts | middleware.ts | instrumentation.ts | \
      _meta.ts)
      return 0
      ;;
  esac
  return 1
}

workspaces=$(for workspace in apps/* packages/*; do
  [ -f "$workspace/package.json" ] && [ ! -f "$workspace/go.mod" ] && printf '%s\n' "$workspace"
done)

# Every non-test source file of the JavaScript workspaces, one corpus for the
# CSS class search.
for workspace in $workspaces; do
  git ls-files -- "$workspace/src" | grep -E '\.(ts|tsx|js|jsx)$' | grep -v -E '\.(unit|test)\.' || true
done | while IFS= read -r source; do cat "$source"; done >"$sources"

check_aliases() {
  for paths in $(git ls-files -- 'apps/*/tsconfig.paths.json' 'packages/*/tsconfig.paths.json'); do
    workspace=${paths%/*}
    jq -r '.compilerOptions.paths // {} | keys[]' "$paths" | while IFS= read -r alias; do
      prefix=${alias%\*}
      [ "$prefix" = "$alias" ] && prefix="$alias'"
      if ! git grep -q -F -e "'$prefix" -e "\"$prefix" -- "$workspace" ":!$paths"; then
        report alias "$paths:$alias" "$paths: nothing imports through $alias"
      fi
    done
  done
}

check_modules() {
  for workspace in $workspaces; do
    git ls-files -- "$workspace/src" | grep -E '\.(ts|tsx|js|jsx|json|css)$' |
      grep -v -E '\.(unit|test)\.|/testdata/' | while IFS= read -r module; do
      is_framework_file "$module" && continue
      base=${module##*/}
      case "$base" in
        *.json | *.css) pattern="[/'\"]$(escape_regex "$base")['\"]" ;;
        index.*)
          directory=${module%/*}
          directory=${directory##*/}
          pattern="[/'\"@]$(escape_regex "$directory")['\"]|/$(escape_regex "$directory")/index['\"]"
          ;;
        *) pattern="[/'\"@]$(escape_regex "${base%.*}")(\\.(ts|tsx|js|jsx))?['\"]" ;;
      esac
      # shellcheck disable=SC2086 # One pathspec per workspace.
      if ! git grep -q -E -e "$pattern" -- $workspaces ":!$module" ':!*.unit.*' ':!*.test.*'; then
        report module "$module" "$module: nothing imports it"
      fi
    done
  done
}

check_exports() {
  for workspace in $workspaces; do
    git ls-files -- "$workspace/src" | grep -E '\.(ts|tsx)$' | grep -v -E '\.(unit|test)\.' |
      while IFS= read -r module; do
        sed -n -E 's/^export (async )?(function|const|let|class|enum) ([A-Za-z0-9_]+).*/\3/p' "$module" |
          while IFS= read -r name; do
            case "$name" in
              generateMetadata | generateStaticParams | generateViewport | metadata | viewport | \
                dynamic | revalidate | runtime | config | proxy | middleware | \
                GET | HEAD | POST | PUT | PATCH | DELETE | OPTIONS)
                continue
                ;;
            esac
            [ "$(grep -c -w -e "$name" "$module")" -gt 1 ] && continue
            used=false
            # shellcheck disable=SC2086 # One pathspec per workspace.
            for consumer in $(git grep -l -w -e "$name" -- $workspaces ":!$module" ':!*.unit.*' ':!*.test.*'); do
              if without_reexports "$consumer" | grep -q -w -e "$name"; then
                used=true
                break
              fi
            done
            [ "$used" = true ] || report export "$module#$name" "$module: nothing uses the export $name"
          done
      done
  done
}

check_dependencies() {
  for workspace in $workspaces; do
    jq -r '.dependencies // {} | keys[]' "$workspace/package.json" | while IFS= read -r dependency; do
      pattern="['\"]$(escape_regex "$dependency")['\"/]"
      if ! git grep -q -E -e "$pattern" -- "$workspace" ":!$workspace/package.json"; then
        report dependency "$workspace:$dependency" "$workspace/package.json: nothing imports $dependency"
      fi
    done
  done
}

check_environment() {
  [ -f turbo.json ] || return 0
  jq -r '
    [.globalEnv, .globalPassThroughEnv, (.tasks[] | .env, .passThroughEnv)]
    | map(. // []) | add | unique[] | select(startswith("!") | not)
  ' turbo.json | while IFS= read -r variable; do
    pattern="(^|[^A-Z0-9_])$(escape_regex "$variable" | sed 's/\\\*/[A-Z0-9_]*/g')([^A-Z0-9_]|\$)"
    if ! git grep -q -E -e "$pattern" -- . ':!turbo.json' ':!*.md' ':!.agents' ':!.claude' ':!pnpm-lock.yaml'; then
      report env "$variable" "turbo.json: nothing reads $variable"
    fi
  done
}

check_classes() {
  for workspace in $workspaces; do
    git ls-files -- "$workspace/src" | grep -E '\.css$' | while IFS= read -r stylesheet; do
      sed -e 's#/\*.*\*/##g' -e '/^@import/d' -e "s/'[^']*'//g" -e 's/"[^"]*"//g' -e 's/url([^)]*)//g' "$stylesheet" |
        grep -o -E '\.[A-Za-z_-][A-Za-z0-9_-]*' | sort -u | while IFS= read -r selector; do
        class=${selector#.}
        if ! grep -q -E -e "(^|[^A-Za-z0-9_-])$(escape_regex "$class")([^A-Za-z0-9_-]|\$)" "$sources"; then
          report class "$class" "$stylesheet: no markup carries .$class"
        fi
      done
    done
  done
}

check_aliases
check_modules
check_exports
check_dependencies
check_environment
check_classes

if [ -s "$findings" ]; then
  while IFS= read -r finding; do drift_print_error "$finding"; done <"$findings"
  drift_print_error "Remove the dead code, or list what the owner keeps on purpose in $allowlist."
  exit 1
fi

drift_print_success "No dead code: every alias, module, export, dependency, turbo.json variable and CSS class is in use"
