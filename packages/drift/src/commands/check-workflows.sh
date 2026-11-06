#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift check workflows" 2
drift_need go
cd "$DRIFT_PROJECT_ROOT"

# The GitHub Actions conventions in .agents/context.md (Continuous integration).
# actionlint checks syntax, expressions and, when shellcheck is installed, the
# shell in run scripts. The checks below it are this repository's own.

status=0
fail() {
  drift_print_error "$1"
  status=1
}

# fail_lines <file> <message> <line numbers>: one finding per line.
fail_lines() {
  [ -n "$3" ] || return 0
  while IFS= read -r line; do
    fail "$1:$line: $2"
  done <<EOF
$3
EOF
}

is_kebab_case() {
  case "$1" in
    '' | -* | *- | *--* | *[!a-z0-9-]*) return 1 ;;
  esac
}

# actionlint v1.7.12, pinned to the tag's commit like every downloaded tool.
actionlint=github.com/rhysd/actionlint/cmd/actionlint@914e7df21a07ef503a81201c76d2b11c789d3fca
GOWORK=off go run "$actionlint" || fail "actionlint found problems in .github/workflows"

# Layout: workflows sit directly in .github/workflows as kebab-case .yml files;
# each composite action is .github/actions/<name>/action.yml and nothing else,
# because scripts belong to drift. The pull request widgets are the one
# exception: kebab-case .sh files directly in .github/widgets.
workflows=$(find .github/workflows -type f | sort)
actions=$(find .github/actions -type f | sort)
widgets=
[ ! -d .github/widgets ] || widgets=$(find .github/widgets -type f | sort)

for file in $workflows; do
  name=${file#.github/workflows/}
  case "$name" in
    */*) fail "$file: workflows sit directly in .github/workflows" ;;
    *.yml) is_kebab_case "${name%.yml}" || fail "$file: workflow names are kebab-case" ;;
    *) fail "$file: workflows are .yml files" ;;
  esac
  grep -q '^permissions:' "$file" || fail "$file: declare the token's permissions at the top level"
done

for file in $actions; do
  name=${file#.github/actions/}
  case "$name" in
    */*/*) fail "$file: an action is .github/actions/<name>/action.yml" ;;
    */action.yml) is_kebab_case "${name%/action.yml}" || fail "$file: action names are kebab-case" ;;
    *) fail "$file: an action folder holds only its action.yml; scripts belong to drift" ;;
  esac
done

for file in $widgets; do
  name=${file#.github/widgets/}
  case "$name" in
    */*) fail "$file: widgets sit directly in .github/widgets" ;;
    *.sh) is_kebab_case "${name%.sh}" || fail "$file: widget names are kebab-case" ;;
    *) fail "$file: widgets are .sh files" ;;
  esac
done

for file in $workflows $actions; do
  # Every action outside this repository runs at a full commit SHA, with its
  # version in a comment; a container image runs at its digest.
  unpinned=$(grep -n -E '^[[:space:]]*(-[[:space:]]+)?uses:' "$file" |
    grep -v -E 'uses:[[:space:]]*\./' |
    grep -v -E 'uses:[[:space:]]*[^@[:space:]]+@[0-9a-f]{40}[[:space:]]+#[[:space:]]*v[0-9]' |
    grep -v -E 'uses:[[:space:]]*docker://[^@[:space:]]+@sha256:[0-9a-f]{64}' |
    cut -d: -f1 || true)
  fail_lines "$file" "pin the action to a full commit SHA with its version in a comment" "$unpinned"

  # An expression inside a run script is spliced into the shell program before
  # it runs, so values reach a script through env: and are read as variables.
  spliced=$(awk '
    in_script {
      if ($0 ~ /^[[:space:]]*$/) next
      match($0, /^ */)
      if (RLENGTH > script_indent) {
        if (index($0, "${{")) print FNR
        next
      }
      in_script = 0
    }
    /^[[:space:]]*(-[[:space:]]+)?run:/ {
      if (index($0, "${{")) print FNR
      match($0, /^[[:space:]]*(-[[:space:]]+)?/)
      script_indent = RLENGTH
      in_script = ($0 ~ /run:[[:space:]]*[|>][-+0-9]*[[:space:]]*$/)
    }
  ' "$file")
  fail_lines "$file" "pass the expression to the run script through env:" "$spliced"

  # Every uploaded artifact states how long it is kept. A step starts at its
  # "- " and runs to the next one.
  unbounded=$(awk '
    function close_step() {
      if (uploads && !retention) print step_line
      uploads = 0
      retention = 0
    }
    /^[[:space:]]*- / { close_step(); step_line = FNR }
    /uses:[[:space:]]*actions\/upload-artifact@/ { uploads = 1 }
    /^[[:space:]]*retention-days:/ { retention = 1 }
    END { close_step() }
  ' "$file")
  fail_lines "$file" "set retention-days on the upload-artifact step" "$unbounded"
done

[ "$status" -ne 0 ] || drift_print_success "The workflows and actions follow the CI conventions"
exit "$status"
