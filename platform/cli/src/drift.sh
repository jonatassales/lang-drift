#!/bin/sh
set -eu

DRIFT_CLI_DIR=$(CDPATH='' cd -P "$(dirname "$0")" && pwd)
DRIFT_PROJECT_ROOT=$(CDPATH='' cd -P "$DRIFT_CLI_DIR/../../.." && pwd)
export DRIFT_CLI_DIR DRIFT_PROJECT_ROOT

. "$DRIFT_CLI_DIR/core/common.sh"

drift_usage_error() {
  drift_print_error "drift: $1"
  drift_print_error "Run 'drift --help' for usage."
  exit 2
}

command=${1:---help}
[ "$#" -eq 0 ] || shift

if [ "$command" = --logs ] || [ "${1:-}" = --logs ]; then
  DRIFT_LOGS=true
  export DRIFT_LOGS
  if [ "$command" = --logs ]; then
    command=${1:---help}
    [ "$#" -eq 0 ] || shift
  else
    shift
  fi
fi

drift_log "command=$command"
commands="$DRIFT_CLI_DIR/commands"

case "$command" in
  --help | -h) exec "$commands/help.sh" "$@" ;;
  setup) exec "$commands/setup.sh" "$@" ;;
  doctor) exec "$commands/doctor.sh" "$@" ;;
  cleanup) exec "$commands/cleanup.sh" "$@" ;;
  reset) exec "$commands/reset.sh" "$@" ;;
  check)
    subcommand=${1:-}
    [ "$#" -eq 0 ] || shift
    case "$subcommand" in
      dead-code) exec "$commands/check-dead-code.sh" "$@" ;;
      localization) exec "$commands/check-localization.sh" "$@" ;;
      workflows) exec "$commands/check-workflows.sh" "$@" ;;
      *) drift_usage_error "Usage: drift check <dead-code|localization|workflows>" ;;
    esac
    ;;
  ci)
    subcommand=${1:-}
    [ "$#" -eq 0 ] || shift
    case "$subcommand" in
      summary) exec "$commands/ci-summary.sh" "$@" ;;
      record) exec "$commands/ci-record.sh" "$@" ;;
      ship) exec "$commands/ci-ship.sh" "$@" ;;
      *) drift_usage_error "Usage: drift ci <summary|record|ship>" ;;
    esac
    ;;
  env)
    subcommand=${1:-}
    [ "$#" -eq 0 ] || shift
    case "$subcommand" in
      setup) exec "$commands/env-setup.sh" "$@" ;;
      validate) exec "$commands/env-validate.sh" "$@" ;;
      *) drift_usage_error "Usage: drift env <setup|validate>" ;;
    esac
    ;;
  git)
    subcommand=${1:-}
    [ "$#" -eq 0 ] || shift
    case "$subcommand" in
      setup) exec "$commands/git-setup.sh" "$@" ;;
      doctor) exec "$commands/git-doctor.sh" "$@" ;;
      pre-commit) exec "$commands/git-pre-commit.sh" "$@" ;;
      commit-msg) exec "$commands/git-commit-msg.sh" "$@" ;;
      lint-history) exec "$commands/git-lint-history.sh" "$@" ;;
      *) drift_usage_error "Usage: drift git <setup|doctor|pre-commit|commit-msg|lint-history>" ;;
    esac
    ;;
  --*) drift_usage_error "Unknown option: $command" ;;
  *) drift_usage_error "Unknown command: $command" ;;
esac
