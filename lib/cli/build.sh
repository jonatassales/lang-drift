#!/bin/sh
set -eu

if [ "$#" -gt 1 ]; then
  printf '%s\n' 'Usage: sh build.sh [--release]' >&2
  exit 2
fi

case "${1:-}" in
  '') _dft_mode=host ;;
  --release) _dft_mode=release ;;
  *)
    printf '%s\n' 'Usage: sh build.sh [--release]' >&2
    exit 2
    ;;
esac

# Frontend-only deployments may install the workspace without using the CLI.
# Keep this explicit: an ordinary source installation still requires Go.
if [ "${LANGDRIFT_SKIP_DFT_BUILD:-0}" = "1" ]; then
  printf '%s\n' 'Skipping the LangDrift CLI build (LANGDRIFT_SKIP_DFT_BUILD=1). Build the CLI separately before using dft.' >&2
  exit 0
fi

# Build only this package using an existing toolchain; never install Go or tools.
_dft_package=$(CDPATH='' cd -P -- "$(dirname -- "$0")" && pwd)
if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'LangDrift source installation requires Go 1.26+. Install Go separately, then rerun build.sh in this CLI package.' >&2
  exit 1
fi

_dft_version=$(sed -n 's/.*"version":[[:space:]]*"\([^"]*\)".*/\1/p' "$_dft_package/package.json" | head -n 1)
if [ -z "$_dft_version" ]; then
  printf '%s\n' 'Failed to read version from package.json.' >&2
  exit 1
fi

mkdir -p "$_dft_package/dist"
_dft_lock=$_dft_package/dist/.build.lock
if ! mkdir "$_dft_lock" 2>/dev/null; then
  printf '%s\n' 'Another LangDrift build is active, or dist/.build.lock remains from an interrupted build. Wait for it, or remove the stale lock after confirming no build is running.' >&2
  exit 1
fi
_dft_output=$_dft_package/dist/.dft-build-$$
trap 'rm -f "$_dft_output"; rmdir "$_dft_lock"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP

install_output() {
  chmod 755 "$_dft_output"
  mv -f "$_dft_output" "$1"
}

build_target() {
  _target_os=$1
  _target_arch=$2
  _target_dest=$3
  GOTOOLCHAIN=local CGO_ENABLED=0 GOOS="$_target_os" GOARCH="$_target_arch" \
    go -C "$_dft_package" build -ldflags="-s -w -X main.version=$_dft_version" -o "$_dft_output" ./src
  install_output "$_target_dest"
}

# GOTOOLCHAIN=local forbids automatic toolchain installation. Go module fetching
# may still access dependency registries; source installation is not offline.
_dft_os=$(GOTOOLCHAIN=local go env GOHOSTOS)
_dft_arch=$(GOTOOLCHAIN=local go env GOHOSTARCH)
case $_dft_os in darwin|linux) ;; *) printf '%s\n' 'LangDrift currently supports macOS and Linux.' >&2; exit 1 ;; esac

if [ "$_dft_mode" = "release" ]; then
  build_target darwin arm64 "$_dft_package/dist/dft-darwin-arm64"
  build_target darwin amd64 "$_dft_package/dist/dft-darwin-amd64"
  build_target linux amd64 "$_dft_package/dist/dft-linux-amd64"
  build_target linux arm64 "$_dft_package/dist/dft-linux-arm64"
  cp -f "$_dft_package/dist/dft-$_dft_os-$_dft_arch" "$_dft_output"
  install_output "$_dft_package/dist/dft"
  exit 0
fi

build_target "$_dft_os" "$_dft_arch" "$_dft_package/dist/dft"
