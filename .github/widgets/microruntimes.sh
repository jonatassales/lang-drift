#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/packages/drift/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

base_sha=${1:-${BASE_SHA:-}}
head_sha=${2:-${HEAD_SHA:-HEAD}}

# Ex package version
dft_pkg_version=$(grep -m 1 '"version":' "$DRIFT_PROJECT_ROOT/lib/cli/package.json" 2>/dev/null | sed -E 's/.*"version":[[:space:]]*"([^"]+)".*/\1/' || echo "0.1.0")

# Ex binary version
dft_bin_version="dft 0.0.0-dev"
if [ -f "$DRIFT_PROJECT_ROOT/lib/cli/bin/dft.sh" ]; then
  raw_bin_ver=$(sh "$DRIFT_PROJECT_ROOT/lib/cli/bin/dft.sh" version 2>/dev/null || true)
  [ -z "$raw_bin_ver" ] || dft_bin_version="$raw_bin_ver"
fi

# Check if lib/cli was updated in commit range
dft_updates="No (0 files changed)"
dft_status="Stable"
contracts_updates="No drift"
contracts_status="In Sync"

if [ -n "$base_sha" ]; then
  dft_diff=$(git diff --name-only "$base_sha" "$head_sha" -- lib/cli 2>/dev/null || true)
  if [ -n "$dft_diff" ]; then
    count=$(printf '%s\n' "$dft_diff" | grep -c '.' || echo 1)
    dft_updates="Yes (${count} files updated)"
    dft_status="Updated"
  fi

  contracts_diff=$(git diff --name-only "$base_sha" "$head_sha" -- lib/core/contracts/go 2>/dev/null || true)
  if [ -n "$contracts_diff" ]; then
    count=$(printf '%s\n' "$contracts_diff" | grep -c '.' || echo 1)
    contracts_updates="Yes (${count} contracts modified)"
    contracts_status="Modified"
  fi
fi

cat <<EOF
### Microruntimes Gate

| Component | Property | Current Value | PR Status |
| :--- | :--- | :--- | :---: |
| **LangDrift CLI** | Package Version | \`@langdrift/cli@${dft_pkg_version}\` | [Pinned](https://console.langdrift.md/sign-in) |
| **LangDrift CLI binary** | Engine Build | \`${dft_bin_version}\` | [Development](https://console.langdrift.md/sign-in) |
| **LangDrift CLI updates** | Code Changes | ${dft_updates} | [${dft_status}](https://console.langdrift.md/sign-in) |
| **Contracts** | State Protocol | \`@langdrift/core\` | [${contracts_status}](https://console.langdrift.md/sign-in) |
| **Bounded Context** | Resource Limits | 2000m CPU · 2048MiB Mem | [Bounded](https://console.langdrift.md/sign-in) |

> **Engine Contract:** Microruntimes plan, run, and record bounded execution around the team's software. Every plan is tied to the FinOps envelope.
EOF
