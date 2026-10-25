#!/bin/sh
# shellcheck disable=SC2034 # Its variables are read by the scripts that source it.

: "${DRIFT_LOGS:=false}"

drift_stdout_is_tty() {
  command -v tty >/dev/null 2>&1 && tty -s <&1 2>/dev/null
}

if [ "${NO_COLOR+x}" != x ] && { [ "${DRIFT_FORCE_COLOR:-0}" = 1 ] || drift_stdout_is_tty; }; then
  DRIFT_COLOR_BRAND=$(printf '\033[38;5;56m')
  DRIFT_COLOR_YELLOW=$(printf '\033[33m')
  DRIFT_COLOR_GREEN=$(printf '\033[32m')
  DRIFT_COLOR_RED=$(printf '\033[31m')
  DRIFT_COLOR_CYAN=$(printf '\033[36m')
  DRIFT_COLOR_DIM=$(printf '\033[2m')
  DRIFT_COLOR_RESET=$(printf '\033[0m')
else
  DRIFT_COLOR_BRAND=
  DRIFT_COLOR_YELLOW=
  DRIFT_COLOR_GREEN=
  DRIFT_COLOR_RED=
  DRIFT_COLOR_CYAN=
  DRIFT_COLOR_DIM=
  DRIFT_COLOR_RESET=
fi

DRIFT_ICON_SUCCESS='✅'
DRIFT_ICON_WARNING='⚠️'
DRIFT_ICON_ERROR='❌'
DRIFT_ICON_INFO='ℹ️'
DRIFT_ICON_LOG='🔎'
DRIFT_ICON_TADA='🎉'
DRIFT_ICON_DOCTOR='🧙'
DRIFT_ICON_SETUP='⏳'
DRIFT_ICON_CLEANUP='🧹'
DRIFT_ICON_ENV='⚙️'
DRIFT_DOCTOR_PHRASE='let me check your spells...'
DRIFT_SETUP_PHRASE='exploring...'
DRIFT_CLEANUP_PHRASE='cleaning the house...'
DRIFT_RESET_PHRASE='starting over...'
DRIFT_ENV_PHRASE='organizing your vars...'

drift_result_file=
drift_results_flushed=false
drift_loader_pid=
drift_loader_delay_pid=
drift_loader_dir=

# The loader spins on stderr while a command works. It waits one second
# before appearing, so quick commands never flicker.
drift_loader_active() {
  [ -z "${CI:-}" ] && { [ "${DRIFT_FORCE_LOADER:-0}" = 1 ] || drift_stdout_is_tty; }
}

drift_loader_frame() {
  case $1 in
    0) printf '|' ;;
    1) printf '/' ;;
    2) printf '-' ;;
    3) printf '%s' "\\" ;;
  esac
}

drift_loader_loop() {
  dir=$1
  i=0
  while [ -f "$dir/alive" ]; do
    if [ -f "$dir/ready" ]; then
      if [ ! -f "$dir/hidden" ]; then
        printf '\033[?25l' >&2
        : > "$dir/hidden"
      fi
      : > "$dir/shown"
      printf '\r\033[K%s%s%s%s' "$(cat "$dir/phrase")" "$DRIFT_COLOR_BRAND" "$(drift_loader_frame "$i")" "$DRIFT_COLOR_RESET" >&2
    fi
    i=$(((i + 1) % 4))
    sleep 0.08 2>/dev/null || sleep 1
  done
}

# The one-second delay runs in its own background job. Stop it before its
# directory goes away, or a command that ends within the second leaves it
# writing into a deleted directory after the prompt returns.
drift_loader_cancel_delay() {
  [ -n "$drift_loader_delay_pid" ] || return 0
  kill "$drift_loader_delay_pid" 2>/dev/null || true
  wait "$drift_loader_delay_pid" 2>/dev/null || true
  drift_loader_delay_pid=
}

drift_loader_start() {
  drift_loader_active || return 0
  if [ -z "$drift_loader_dir" ]; then
    drift_loader_dir=$(mktemp -d "${TMPDIR:-/tmp}/drift-loader.XXXXXX")
  fi
  if [ -f "$drift_loader_dir/shown" ]; then
    printf '\r\033[K%s\n' "$(cat "$drift_loader_dir/phrase")" >&2
    rm -f "$drift_loader_dir/shown"
  fi
  drift_loader_cancel_delay
  rm -f "$drift_loader_dir/ready"
  printf '%s ' "$1" > "$drift_loader_dir/phrase"
  : > "$drift_loader_dir/alive"
  (sleep 1 && : > "$drift_loader_dir/ready") &
  drift_loader_delay_pid=$!
  if [ -z "$drift_loader_pid" ] || ! kill -0 "$drift_loader_pid" 2>/dev/null; then
    drift_loader_loop "$drift_loader_dir" &
    drift_loader_pid=$!
  fi
}

drift_loader_stop() {
  [ -n "$drift_loader_dir" ] || return 0
  drift_loader_cancel_delay
  rm -f "$drift_loader_dir/alive"
  if [ -n "$drift_loader_pid" ]; then
    kill "$drift_loader_pid" 2>/dev/null || true
    wait "$drift_loader_pid" 2>/dev/null || true
    drift_loader_pid=
  fi
  if [ -f "$drift_loader_dir/hidden" ]; then
    printf '\r\033[K\033[?25h' >&2
  fi
  rm -rf "$drift_loader_dir"
  drift_loader_dir=
}

drift_print_logo() {
  printf '%s\n' \
    "${DRIFT_COLOR_BRAND}██████╗ ██████╗ ██╗███████╗████████╗" \
    "██╔══██╗██╔══██╗██║██╔════╝╚══██╔══╝" \
    "██║  ██║██████╔╝██║█████╗     ██║" \
    "██║  ██║██╔══██╗██║██╔══╝     ██║" \
    "██████╔╝██║  ██║██║██║        ██║" \
    "╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝        ╚═╝${DRIFT_COLOR_RESET}"
}

drift_print_opening_logo() {
  [ -z "${CI:-}" ] || return 0
  drift_print_logo
}

# Results are queued while the loader runs and printed when the command ends.
drift_queue_result() {
  kind=$1
  shift
  if [ -z "$drift_result_file" ]; then
    drift_result_file=$(mktemp "${TMPDIR:-/tmp}/drift-result.XXXXXX")
  fi
  printf '%s\t%s\n' "$kind" "$*" >> "$drift_result_file"
}

drift_emit_result() {
  kind=$1
  suffix=$2
  shift 2
  case "$kind" in
    success) printf '%s%s %s%s%s\n' "$DRIFT_COLOR_GREEN" "$DRIFT_ICON_SUCCESS" "$*" "$DRIFT_COLOR_RESET" "$suffix" ;;
    warning) printf '%s%s %s%s%s\n' "$DRIFT_COLOR_YELLOW" "$DRIFT_ICON_WARNING" "$*" "$DRIFT_COLOR_RESET" "$suffix" >&2 ;;
    error) printf '%s%s %s%s%s\n' "$DRIFT_COLOR_RED" "$DRIFT_ICON_ERROR" "$*" "$DRIFT_COLOR_RESET" "$suffix" >&2 ;;
    info) printf '%s%s %s%s%s\n' "$DRIFT_COLOR_CYAN" "$DRIFT_ICON_INFO" "$*" "$DRIFT_COLOR_RESET" "$suffix" ;;
    banner) printf '\n%s%s\n\n' "$*" "$suffix" ;;
  esac
}

drift_flush_results() {
  [ "$drift_results_flushed" = false ] || return 0
  drift_results_flushed=true
  drift_loader_stop
  [ -n "$drift_result_file" ] || return 0

  # The last line gets a celebration when the command reported any success.
  total=$(wc -l < "$drift_result_file" | tr -d ' ')
  had_success=false
  ! grep -q "^success	" "$drift_result_file" || had_success=true
  index=0
  while IFS='	' read -r kind text; do
    index=$((index + 1))
    suffix=
    if [ "$index" -eq "$total" ] && [ "$had_success" = true ] && [ -z "${CI:-}" ]; then
      suffix=" ${DRIFT_COLOR_BRAND}${DRIFT_ICON_TADA}${DRIFT_COLOR_RESET}"
    fi
    drift_emit_result "$kind" "$suffix" "$text"
  done < "$drift_result_file"
  rm -f "$drift_result_file"
}

drift_queue_scene() {
  [ -z "${CI:-}" ] || return 0
  drift_queue_result banner "$1 $2"
}

drift_print_success() { drift_queue_result success "$*"; }
drift_print_warning() { drift_queue_result warning "$*"; }
drift_print_error() { drift_queue_result error "$*"; }
drift_print_info() { drift_queue_result info "$*"; }
drift_log() {
  [ "$DRIFT_LOGS" = true ] || return 0
  printf '%s%s %s%s\n' "$DRIFT_COLOR_DIM" "$DRIFT_ICON_LOG" "$*" "$DRIFT_COLOR_RESET" >&2
}

if [ "${drift_exit_trap_set:-0}" != 1 ]; then
  drift_exit_trap_set=1
  trap drift_flush_results EXIT
  trap 'exit 130' INT
  trap 'exit 143' TERM
fi
