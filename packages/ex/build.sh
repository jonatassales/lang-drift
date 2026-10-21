#!/bin/sh
set -eu

# Frontend-only deployments may install the workspace without using the CLI.
# Keep this explicit: an ordinary source installation still requires Go.
if [ "${EXPHERE_SKIP_CLI_BUILD:-0}" = "1" ]; then
  printf '%s\n' 'Skipping the Exphere CLI build (EXPHERE_SKIP_CLI_BUILD=1). Build the CLI separately before using ex.' >&2
  exit 0
fi

# Build only this package using an existing toolchain; never install Go or tools.
_ex_package=$(CDPATH= cd -P -- "$(dirname -- "$0")" && pwd)
if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'Exphere source installation requires Go 1.23+. Install Go separately, then rerun build.sh in this CLI package.' >&2
  exit 1
fi
mkdir -p "$_ex_package/dist"
_ex_lock=$_ex_package/dist/.build.lock
if ! mkdir "$_ex_lock" 2>/dev/null; then
  printf '%s\n' 'Another Exphere build is active, or dist/.build.lock remains from an interrupted build. Wait for it, or remove the stale lock after confirming no build is running.' >&2
  exit 1
fi
_ex_output=$_ex_package/dist/.ex-build-$$
trap 'rm -f "$_ex_output"; rmdir "$_ex_lock"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
# GOTOOLCHAIN=local forbids automatic toolchain installation. Go module fetching
# may still access dependency registries; source installation is not offline.
_ex_os=$(GOTOOLCHAIN=local go env GOHOSTOS)
_ex_arch=$(GOTOOLCHAIN=local go env GOHOSTARCH)
case $_ex_os in darwin|linux) ;; *) printf '%s\n' 'Exphere currently supports macOS and Linux.' >&2; exit 1 ;; esac
GOTOOLCHAIN=local GOOS="$_ex_os" GOARCH="$_ex_arch" go -C "$_ex_package" build -o "$_ex_output" ./src
chmod 755 "$_ex_output"
mv -f "$_ex_output" "$_ex_package/dist/ex"
