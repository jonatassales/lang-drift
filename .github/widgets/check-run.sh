#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 3 ] || drift_die "Usage: drift widget check-run <name> <conclusion> <title> <summary.md" 2
drift_need gh
drift_need jq

# A completed check run on the commit a pull request shows: its name, a
# conclusion, a title and the Markdown summary read from stdin. The job that
# publishes it needs the checks: write permission.

name=$1
conclusion=$2
title=$3

[ -n "$name" ] || drift_die "The check run needs a name." 2
[ -n "$title" ] || drift_die "The check run needs a title." 2

case "$conclusion" in
  success | failure | neutral | cancelled | skipped | timed_out | action_required) ;;
  *) drift_die "The conclusion is success, failure, neutral, cancelled, skipped, timed_out or action_required." 2 ;;
esac

case "${GITHUB_REPOSITORY:-}" in
  */*/* | /* | */ | *[!A-Za-z0-9._/-]*) drift_die "GITHUB_REPOSITORY must name the repository as owner/name." 2 ;;
  */*) ;;
  *) drift_die "GITHUB_REPOSITORY must name the repository as owner/name." 2 ;;
esac

[ -n "${WIDGET_SHA:-}" ] || drift_die "WIDGET_SHA must be the full SHA of the commit the check run belongs to." 2

case "$WIDGET_SHA" in
  *[!0-9a-f]*) drift_die "WIDGET_SHA is a full commit SHA." 2 ;;
esac
[ "${#WIDGET_SHA}" -eq 40 ] || drift_die "WIDGET_SHA is a full commit SHA." 2

summary_file=$(mktemp)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -f "$summary_file"; drift_flush_results' EXIT

[ ! -t 0 ] || drift_die "Pipe the Markdown summary into stdin." 2
cat >"$summary_file"
[ -s "$summary_file" ] || drift_die "The summary on stdin is empty." 2

# GitHub rejects a summary over 65535 characters, so a longer one is cut at
# 65400, which leaves room for the note that says so.
payload=$(jq -Rs \
  --arg name "$name" \
  --arg sha "$WIDGET_SHA" \
  --arg conclusion "$conclusion" \
  --arg title "$title" \
  '{
    name: $name,
    head_sha: $sha,
    status: "completed",
    conclusion: $conclusion,
    output: {
      title: $title,
      summary: (if length > 65535 then .[0:65400] + "\n\n_Truncated: the summary passed the GitHub limit._" else . end)
    }
  }' <"$summary_file")

printf '%s' "$payload" | gh api --method POST "repos/$GITHUB_REPOSITORY/check-runs" --input - --jq .html_url ||
  drift_die "GitHub refused the check run $name."
