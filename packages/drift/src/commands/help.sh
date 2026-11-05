#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Help does not accept arguments." 2

drift_print_opening_logo

command_line() {
  printf '  %-54s %s\n' "$1" "$2"
}

printf '\n%s\n' 'Commands:'
command_line setup 'Link the drift launcher into your bin directory'
command_line 'doctor [--ci]' 'Check the toolchain and repository layout'
command_line cleanup 'Remove generated workspace state'
command_line reset 'Remove generated state, local config and the lockfile, then install'
command_line 'check [biome|localization|workflows]' 'Check imports, catalogs and workflows'
command_line 'env [setup|validate]' 'Create and validate environment files'
command_line 'git [setup|doctor|pre-commit|commit-msg|lint-history]' 'Install, check and run Git hooks and guardrails'
command_line 'quality [check|coverage|dead-code]' 'Rule coverage, non-compliant paths and dead code'
command_line 'widget [<name>|check-run|sticky-comment]' 'Print PR widgets and update check runs or comments'
printf '\n%s\n' 'Add --logs to any command to print debug logs.'
