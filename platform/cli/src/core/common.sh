#!/bin/sh
# shellcheck disable=SC2034 # Its variables are read by the scripts that source it.

: "${DRIFT_PROJECT_ROOT:?DRIFT_PROJECT_ROOT must be set by platform/cli/src/drift.sh}"
: "${DRIFT_CLI_DIR:?DRIFT_CLI_DIR must be set by platform/cli/src/drift.sh}"

. "$DRIFT_CLI_DIR/core/output.sh"

# Tool configuration lives under platform/, grouped by concern.
DRIFT_HOOKS_DIR=platform/versioning/.husky
DRIFT_COMMITLINT_CONFIG=platform/versioning/commitlint.config.json
DRIFT_LINT_STAGED_CONFIG=platform/quality/lint-staged.config.json

drift_die() {
  drift_print_error "drift: $1"
  exit "${2:-1}"
}

drift_warn() { drift_print_warning "drift: $*"; }
drift_has() { command -v "$1" >/dev/null 2>&1; }
drift_need() { drift_has "$1" || drift_die "Required command not found: $1" 127; }

drift_git_checkout() {
  git -C "$DRIFT_PROJECT_ROOT" rev-parse --is-inside-work-tree >/dev/null 2>&1
}

# Every .env.example that Git tracks or would track describes one .env file.
drift_find_env_templates() {
  git -C "$DRIFT_PROJECT_ROOT" ls-files --cached --others --exclude-standard |
    while IFS= read -r path; do
      case "$path" in
        .agents/* | .claude/*) continue ;;
        */.env.example | .env.example) printf '%s/%s\n' "$DRIFT_PROJECT_ROOT" "$path" ;;
      esac
    done
}

# What the toolchain regenerates, matched by name at any depth: dependencies,
# caches, build output and the files tsc, Next.js, macOS and pnpm write.
drift_find_generated_state() {
  find . -path ./.git -prune -o \
    \( -type d \( -name node_modules -o -name .pnpm-store -o -name .turbo -o -name .next -o -name dist -o -name coverage \) \
    -o -type f \( -name '*.tsbuildinfo' -o -name next-env.d.ts -o -name .DS_Store -o -name '*.log' \) \) \
    -prune -print
}

# A .gitkeep only holds an otherwise empty folder; once the folder has anything
# else, it is redundant.
drift_find_redundant_gitkeeps() {
  find . -path ./.git -prune -o -type f -name .gitkeep -print |
    while IFS= read -r gitkeep; do
      if ls -A "${gitkeep%/*}" | grep -q -v -x .gitkeep; then
        printf '%s\n' "$gitkeep"
      fi
    done
}

# drift_remove_paths <file>: removes each path the file lists, one per line,
# and reports it. Read from a file rather than a pipe, so the loop runs in this
# shell and drift_removed_count survives it.
drift_removed_count=0
drift_remove_paths() {
  while IFS= read -r removed_path; do
    [ -n "$removed_path" ] || continue
    rm -rf "$removed_path"
    drift_removed_count=$((drift_removed_count + 1))
    drift_print_success "Removed ${removed_path#./}"
  done <"$1"
}

drift_path_has() {
  case ":${PATH:-}:" in
    *:"$1":*) return 0 ;;
    *) return 1 ;;
  esac
}

drift_default_bin_dir() {
  if [ -n "${DRIFT_BIN_DIR:-}" ]; then
    printf '%s\n' "$DRIFT_BIN_DIR"
  elif [ -n "${PNPM_HOME:-}" ]; then
    if drift_path_has "$PNPM_HOME" && ! drift_path_has "$PNPM_HOME/bin"; then
      printf '%s\n' "$PNPM_HOME"
    else
      printf '%s\n' "$PNPM_HOME/bin"
    fi
  elif [ -n "${XDG_BIN_HOME:-}" ]; then
    printf '%s\n' "$XDG_BIN_HOME"
  elif [ -n "${HOME:-}" ]; then
    printf '%s\n' "$HOME/.local/bin"
  else
    return 1
  fi
}

# Another program named `drift` may come first on PATH (for example a link
# from an older checkout).
# Report which one the shell would run.
drift_check_launcher() {
  launcher="$1/drift"
  resolved=$(command -v drift 2>/dev/null || true)
  if [ "$resolved" = "$launcher" ]; then
    drift_print_info "You can run: drift --help"
  elif [ -n "$resolved" ]; then
    drift_warn "'drift' runs $resolved, not $launcher. Put $1 earlier on PATH."
  else
    drift_warn "$1 is not on PATH."
  fi
}
