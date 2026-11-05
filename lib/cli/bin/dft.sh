#!/bin/sh
set -eu

unsupported_platform() {
  printf '%s\n' 'LangDrift CLI supports macOS and Linux (x64 and arm64).' >&2
  exit 4
}

# Resolve symlinks so `dft` works through `npx` and `node_modules/.bin/dft`.
_dft_source=$0
while [ -L "$_dft_source" ]; do
  _dft_dir=$(CDPATH='' cd -P -- "$(dirname -- "$_dft_source")" && pwd)
  _dft_link=$(readlink "$_dft_source")
  case "$_dft_link" in
    /*) _dft_source=$_dft_link ;;
    *) _dft_source=$_dft_dir/$_dft_link ;;
  esac
done

_dft_bin_dir=$(CDPATH='' cd -P -- "$(dirname -- "$_dft_source")" && pwd)
_dft_package_dir=$(dirname -- "$_dft_bin_dir")

case "$(uname -s 2>/dev/null || true)" in
  Darwin) _dft_os=darwin ;;
  Linux) _dft_os=linux ;;
  *) unsupported_platform ;;
esac

case "$(uname -m 2>/dev/null || true)" in
  arm64 | aarch64) _dft_arch=arm64 ;;
  x86_64 | amd64) _dft_arch=amd64 ;;
  *) unsupported_platform ;;
esac

_dft_binary="$_dft_package_dir/dist/dft-$_dft_os-$_dft_arch"
if [ ! -f "$_dft_binary" ]; then
  _dft_binary="$_dft_package_dir/dist/dft"
fi

if [ -f "$_dft_binary" ] && [ ! -x "$_dft_binary" ]; then
  chmod 755 "$_dft_binary" 2>/dev/null || true
fi

if [ ! -x "$_dft_binary" ]; then
  printf '%s\n' "LangDrift CLI binary is missing for $_dft_os/$_dft_arch. Run build.sh or reinstall @langdrift/cli." >&2
  exit 4
fi

exec "$_dft_binary" "$@"
