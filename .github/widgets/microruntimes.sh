#!/bin/sh
set -eu

: "${DRIFT_PROJECT_ROOT:=$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
: "${DRIFT_CLI_DIR:=$DRIFT_PROJECT_ROOT/platform/cli/src}"
. "$DRIFT_CLI_DIR/core/common.sh"

base_sha=${1:-${BASE_SHA:-}}
head_sha=${2:-${HEAD_SHA:-HEAD}}

# Ex package version
ex_pkg_version=$(grep -m 1 '"version":' "$DRIFT_PROJECT_ROOT/packages/ex/package.json" 2>/dev/null | sed -E 's/.*"version":[[:space:]]*"([^"]+)".*/\1/' || echo "0.1.0")

# Ex binary version
ex_bin_version="ex 0.0.0-dev"
if [ -f "$DRIFT_PROJECT_ROOT/packages/ex/bin/ex" ]; then
  raw_bin_ver=$(sh "$DRIFT_PROJECT_ROOT/packages/ex/bin/ex" version 2>/dev/null || true)
  [ -z "$raw_bin_ver" ] || ex_bin_version="$raw_bin_ver"
fi

# Check if packages/ex was updated in commit range
ex_updates="No (0 files changed)"
ex_status="🟢 Stable"
contracts_updates="No drift"
contracts_status="🟢 In Sync"

if [ -n "$base_sha" ]; then
  ex_diff=$(git diff --name-only "$base_sha" "$head_sha" -- packages/ex 2>/dev/null || true)
  if [ -n "$ex_diff" ]; then
    count=$(printf '%s\n' "$ex_diff" | grep -c '.' || echo 1)
    ex_updates="Yes (${count} files updated)"
    ex_status="⚡ Updated"
  fi

  contracts_diff=$(git diff --name-only "$base_sha" "$head_sha" -- packages/contracts 2>/dev/null || true)
  if [ -n "$contracts_diff" ]; then
    count=$(printf '%s\n' "$contracts_diff" | grep -c '.' || echo 1)
    contracts_updates="Yes (${count} contracts modified)"
    contracts_status="🟡 Modified"
  fi
fi

cat <<EOF
### <img src="https://langdrift.md/metadata/icon-192.png" width="16" height="16" align="middle" alt="" /> Microruntimes Gate

| Component | Property | Current Value | PR Status |
| :--- | :--- | :--- | :---: |
| **LangDrift-X (CLI)** | Package Version | \`@langdrift/ex@${ex_pkg_version}\` | 📦 Pinned |
| **LangDrift-X (Binary)** | Engine Build | \`${ex_bin_version}\` | 🛠️ Development |
| **LangDrift-X (Updates)**| Code Changes | ${ex_updates} | ${ex_status} |
| **Contracts** | State Protocol | \`@langdrift/contracts\` | ${contracts_status} |
| **Bounded Context** | Resource Limits | 2000m CPU · 2048MiB Mem | 🛡️ Bounded |

> **Engine Contract:** Microruntimes plan, run, and record bounded execution around the team's software. Every plan is tied to the FinOps envelope.
EOF
