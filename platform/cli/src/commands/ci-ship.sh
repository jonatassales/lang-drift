#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 1 ] || drift_die "Usage: drift ci ship <events.ndjson>" 2
[ -f "$1" ] || drift_die "No events file: $1"
drift_need curl

# Sends NDJSON events to a Better Stack HTTP source (Telemetry, Sources, HTTP),
# which answers 202 when it accepts them.
[ -n "${BETTERSTACK_SOURCE_TOKEN:-}" ] || drift_die "BETTERSTACK_SOURCE_TOKEN is not set."
[ -n "${BETTERSTACK_INGESTING_HOST:-}" ] || drift_die "BETTERSTACK_INGESTING_HOST is not set."

# A bare hostname only, so the token never reaches another path or scheme.
case "$BETTERSTACK_INGESTING_HOST" in
  .* | *. | *[!a-z0-9.-]*)
    drift_die "BETTERSTACK_INGESTING_HOST must be the source's ingesting host, a bare hostname."
    ;;
esac

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

# Better Stack takes up to 10 MiB per request; 1000 events stay far below it.
split -l 1000 "$1" "$workdir/batch-"

for batch in "$workdir"/batch-*; do
  # curl reads the token from stdin, so it never appears in argv or a log.
  status=$(printf 'header = "Authorization: Bearer %s"\n' "$BETTERSTACK_SOURCE_TOKEN" |
    curl --config - --silent --show-error --max-time 30 --retry 2 \
      --header 'Content-Type: application/x-ndjson' \
      --data-binary "@$batch" \
      --output /dev/null --write-out '%{http_code}' \
      "https://$BETTERSTACK_INGESTING_HOST") ||
    drift_die "Better Stack could not be reached at $BETTERSTACK_INGESTING_HOST."

  case "$status" in
    202) ;;
    402) drift_die "Better Stack refused the events: the plan's quota or spending limit is reached." ;;
    403) drift_die "Better Stack refused the events: BETTERSTACK_SOURCE_TOKEN is not this source's token." ;;
    406) drift_die "Better Stack refused the events: the batch is not valid NDJSON." ;;
    413) drift_die "Better Stack refused the events: the batch is over 10 MiB." ;;
    *) drift_die "Better Stack answered HTTP $status." ;;
  esac
done

events=$(wc -l <"$1" | tr -d ' ')
drift_print_success "Shipped $events events to Better Stack."
