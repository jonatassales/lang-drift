#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Help does not accept arguments." 2

drift_print_opening_logo

command_line() {
  printf '  %-30s %s\n' "$1" "$2"
}

printf '\n%s\n' 'Commands:'
command_line setup 'Link the drift launcher into your bin directory'
command_line 'doctor [--ci]' 'Check the toolchain and repository layout'
command_line cleanup 'Remove generated workspace state'
command_line reset 'Remove generated state, local config and the lockfile, then install'
command_line 'check biome' 'Check that Biome still enforces the import rules'
command_line 'check dead-code' 'Find code that only served something a refactor removed'
command_line 'check localization' 'Check that translation catalogs are JSON and loaders hold no copy'
command_line 'check rules' 'Enforce the repository rules in .claude/rules/*.rule.md'
command_line 'check workflows' 'Lint the GitHub Actions workflows and actions'
command_line 'ci summary <run.json>...' 'Print Turborepo run summaries as Markdown tables'
command_line 'ci record <run-id> <attempt>' 'Print a workflow run as Better Stack events'
command_line 'ci ship <events.ndjson>' 'Send events to Better Stack'
command_line 'widget <name>' 'Print pull request widgets (workflow, finops, microruntimes, code-quality, agentic, dashboard)'
command_line 'widget check-run <name> ...' 'Publish a check run (<conclusion> <title>) with a summary from stdin'
command_line 'widget sticky-comment <key>' 'Post or update the pull request comment for <key> from stdin'
command_line 'env setup' 'Create each .env from its .env.example'
command_line 'env validate' 'Check that each .env.example has its .env'
command_line 'git setup' 'Install the Git hooks'
command_line 'git doctor' 'Check the Git hooks'
command_line 'git pre-commit' 'Run the pre-commit checks'
command_line 'git commit-msg <file>' 'Run the commit message guardrails'
command_line 'git lint-history <from> <to>' 'Run the commit guardrails on a range'
printf '\n%s\n' 'Add --logs to any command to print debug logs.'
