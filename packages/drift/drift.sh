#!/bin/sh
set -eu

# Resolve symlinks so `drift` also works when linked into a bin directory.
source_path=$0
while [ -L "$source_path" ]; do
  source_dir=$(CDPATH='' cd -P "$(dirname "$source_path")" && pwd)
  link=$(readlink "$source_path")
  case "$link" in
    /*) source_path=$link ;;
    *) source_path=$source_dir/$link ;;
  esac
done

cli_dir=$(CDPATH='' cd -P "$(dirname "$source_path")" && pwd)
exec "$cli_dir/src/drift.sh" "$@"
