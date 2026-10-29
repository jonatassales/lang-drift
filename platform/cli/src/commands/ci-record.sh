#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 2 ] || drift_die "Usage: drift ci record <run-id> <attempt>" 2
drift_need gh
drift_need jq

for number in "$1" "$2"; do
  case "$number" in
    '' | *[!0-9]*) drift_die "The run id and the attempt are numbers." 2 ;;
  esac
done

# The record of one workflow run attempt as NDJSON on stdout, the events that
# drift ci ship sends to Better Stack. The run and its jobs come from the GitHub
# API; the Turborepo summaries come from the run's turbo-summary-* artifacts and
# are only parsed, never executed (platform/observability/ci-events.jq).

# Actions sets GITHUB_REPOSITORY; elsewhere it is the checkout's repository.
repository=${GITHUB_REPOSITORY:-$(gh repo view --json nameWithOwner --jq .nameWithOwner)}
attempt="repos/$repository/actions/runs/$1/attempts/$2"

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

gh api "$attempt" >"$workdir/run.json"
gh api --paginate "$attempt/jobs" --jq '.jobs[]' >"$workdir/jobs.ndjson"

summaries=$(gh api --paginate "repos/$repository/actions/runs/$1/artifacts" \
  --jq '.artifacts[] | select(.expired | not) | .name | select(startswith("turbo-summary-"))')
if [ -n "$summaries" ]; then
  gh run download "$1" --repo "$repository" --pattern 'turbo-summary-*' --dir "$workdir/artifacts"
fi

# A summary over 5 MiB (10240 blocks of 512 bytes) is not one turbo wrote. The
# list is read from a file rather than a pipe, so the warnings reach this shell.
: >"$workdir/summaries.ndjson"
find "$workdir" -path "$workdir/artifacts/*" -type f -name '*.json' -size -10240 | sort >"$workdir/summary-files"
while IFS= read -r summary; do
  jq -c . "$summary" >>"$workdir/summaries.ndjson" 2>/dev/null ||
    drift_warn "Skipped a Turborepo summary that is not JSON: ${summary#"$workdir"/artifacts/}"
done <"$workdir/summary-files"

jq -n -c \
  --slurpfile run "$workdir/run.json" \
  --slurpfile jobs "$workdir/jobs.ndjson" \
  --slurpfile summaries "$workdir/summaries.ndjson" \
  -f "$DRIFT_PROJECT_ROOT/platform/observability/ci-events.jq"
