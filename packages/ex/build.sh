#!/bin/sh
set -eu

_ex_mode=host
if [ "$#" -gt 0 ]; then
  case "$1" in
    --release)
      _ex_mode=release
      shift
      ;;
    *)
      printf '%s\n' 'Usage: sh build.sh [--release]' >&2
      exit 2
      ;;
  esac
fi
if [ "$#" -gt 0 ]; then
  printf '%s\n' 'Usage: sh build.sh [--release]' >&2
  exit 2
fi

# Frontend-only deployments may install the workspace without using the CLI.
# Keep this explicit: an ordinary source installation still requires Go.
if [ "${LANGDRIFT_SKIP_EX_BUILD:-0}" = "1" ]; then
  printf '%s\n' 'Skipping the LangDrift CLI build (LANGDRIFT_SKIP_EX_BUILD=1). Build the CLI separately before using ex.' >&2
  exit 0
fi

# Build only this package using an existing toolchain; never install Go or tools.
_ex_package=$(CDPATH='' cd -P -- "$(dirname -- "$0")" && pwd)
if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'LangDrift source installation requires Go 1.26+. Install Go separately, then rerun build.sh in this CLI package.' >&2
  exit 1
fi

_ex_version=$(sed -n 's/.*"version":[[:space:]]*"\([^"]*\)".*/\1/p' "$_ex_package/package.json" | head -n 1)
if [ -z "$_ex_version" ]; then
  _ex_version=0.0.0-dev
fi

mkdir -p "$_ex_package/dist"
_ex_lock=$_ex_package/dist/.build.lock
if ! mkdir "$_ex_lock" 2>/dev/null; then
  printf '%s\n' 'Another LangDrift build is active, or dist/.build.lock remains from an interrupted build. Wait for it, or remove the stale lock after confirming no build is running.' >&2
  exit 1
fi
_ex_output=$_ex_package/dist/.ex-build-$$
trap 'rm -f "$_ex_output"; rmdir "$_ex_lock"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP

build_target() {
  _target_os=$1
  _target_arch=$2
  _target_dest=$3
  GOTOOLCHAIN=local CGO_ENABLED=0 GOOS="$_target_os" GOARCH="$_target_arch" \
    go -C "$_ex_package" build -ldflags="-s -w -X main.version=$_ex_version" -o "$_ex_output" ./src
  chmod 755 "$_ex_output"
  mv -f "$_ex_output" "$_target_dest"
}

# GOTOOLCHAIN=local forbids automatic toolchain installation. Go module fetching
# may still access dependency registries; source installation is not offline.
_ex_os=$(GOTOOLCHAIN=local go env GOHOSTOS)
_ex_arch=$(GOTOOLCHAIN=local go env GOHOSTARCH)
case $_ex_os in darwin|linux) ;; *) printf '%s\n' 'LangDrift currently supports macOS and Linux.' >&2; exit 1 ;; esac

if [ "$_ex_mode" = "release" ]; then
  build_target darwin arm64 "$_ex_package/dist/ex-darwin-arm64"
  build_target darwin amd64 "$_ex_package/dist/ex-darwin-amd64"
  build_target linux amd64 "$_ex_package/dist/ex-linux-amd64"
  build_target linux arm64 "$_ex_package/dist/ex-linux-arm64"
  cp -f "$_ex_package/dist/ex-$_ex_os-$_ex_arch" "$_ex_output"
  chmod 755 "$_ex_output"
  mv -f "$_ex_output" "$_ex_package/dist/ex"
  exit 0
fi

build_target "$_ex_os" "$_ex_arch" "$_ex_package/dist/ex"
