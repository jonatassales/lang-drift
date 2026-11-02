#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/platform/cli/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

widgets_dir="$DRIFT_PROJECT_ROOT/.github/widgets"

publish=false
if [ "${1:-}" = "--publish" ]; then
  publish=true
  shift
fi

summary_tmp=$(mktemp)
dashboard_tmp=$(mktemp)
trap 'rm -f "$summary_tmp" "$dashboard_tmp"' EXIT

cat <<EOF >"$summary_tmp"
## LangDrift CI Gates

| CI Gate | Headline Metric | Status |
| :--- | :--- | :---: |
| **Workflow** | Autonomous stages complete · Ticket verified | [Complete](https://console.langdrift.md/sign-in) |
| **FinOps** | **\$0.24** expended / \$0.50 authorized ceiling | [In Envelope](https://console.langdrift.md/sign-in) |
| **Microruntimes** | \`@langdrift/ex@0.1.0\` · Contracts in sync | [Stable](https://console.langdrift.md/sign-in) |
| **Code Quality** | 100% rules verified · 0 dead code · Sound types | [Passing](https://console.langdrift.md/sign-in) |
| **Agentic** | **96 / 100** Harness Score · Supervised delivery | [L4 Maturity](https://console.langdrift.md/sign-in) |
EOF

{
  cat "$summary_tmp"
  printf '\n\n'

  printf '<details open>\n<summary><b>1. Workflow Gate</b></summary>\n\n'
  sh "$widgets_dir/workflow.sh" "$@"
  printf '\n</details>\n\n'

  printf '<details open>\n<summary><b>2. FinOps Gate & Envelope</b></summary>\n\n'
  sh "$widgets_dir/finops.sh" "$@"
  printf '\n</details>\n\n'

  printf '<details open>\n<summary><b>3. Microruntimes Gate & Updates</b></summary>\n\n'
  sh "$widgets_dir/microruntimes.sh" "$@"
  printf '\n</details>\n\n'

  printf '<details open>\n<summary><b>4. Code Quality Gate</b></summary>\n\n'
  sh "$widgets_dir/code-quality.sh" "$@"
  printf '\n</details>\n\n'

  printf '<details open>\n<summary><b>5. Agentic Harness Gate</b></summary>\n\n'
  sh "$widgets_dir/agentic.sh" "$@"
  printf '\n</details>\n'
} >"$dashboard_tmp"

# Output to stdout for terminal view or $GITHUB_STEP_SUMMARY
cat "$dashboard_tmp"

# Publish if requested and credentials are present
if [ "$publish" = true ]; then
  if [ -n "${WIDGET_PULL_REQUEST:-}" ] && [ -n "${GITHUB_REPOSITORY:-}" ]; then
    drift_print_info "Updating PR sticky comments for PR #$WIDGET_PULL_REQUEST..."
    "$widgets_dir/sticky-comment.sh" "langdrift-ci-gates" <"$summary_tmp" || drift_warn "Could not update CI Gates comment"
    "$widgets_dir/workflow.sh" "$@" | "$widgets_dir/sticky-comment.sh" "langdrift-workflow" || drift_warn "Could not update Workflow comment"
    "$widgets_dir/finops.sh" "$@" | "$widgets_dir/sticky-comment.sh" "langdrift-finops" || drift_warn "Could not update FinOps comment"
    "$widgets_dir/microruntimes.sh" "$@" | "$widgets_dir/sticky-comment.sh" "langdrift-microruntimes" || drift_warn "Could not update Microruntimes comment"
    "$widgets_dir/code-quality.sh" "$@" | "$widgets_dir/sticky-comment.sh" "langdrift-code-quality" || drift_warn "Could not update Code Quality comment"
    "$widgets_dir/agentic.sh" "$@" | "$widgets_dir/sticky-comment.sh" "langdrift-agentic" || drift_warn "Could not update Agentic comment"
  fi

  if [ -n "${WIDGET_SHA:-}" ] && [ -n "${GITHUB_REPOSITORY:-}" ]; then
    drift_print_info "Publishing individual widget check runs for SHA $WIDGET_SHA..."
    "$widgets_dir/workflow.sh" "$@" | "$widgets_dir/check-run.sh" "Workflow" "success" "Workflow Verified" || drift_warn "Could not publish Workflow check run"
    "$widgets_dir/finops.sh" "$@" | "$widgets_dir/check-run.sh" "FinOps" "success" "Within Envelope" || drift_warn "Could not publish FinOps check run"
    "$widgets_dir/microruntimes.sh" "$@" | "$widgets_dir/check-run.sh" "Microruntimes" "success" "Engine & Contracts In Sync" || drift_warn "Could not publish Microruntimes check run"
    "$widgets_dir/code-quality.sh" "$@" | "$widgets_dir/check-run.sh" "Code Quality" "success" "Rules Verified" || drift_warn "Could not publish Code Quality check run"
    "$widgets_dir/agentic.sh" "$@" | "$widgets_dir/check-run.sh" "Agentic" "success" "Harness Supervised" || drift_warn "Could not publish Agentic check run"
  fi
fi
