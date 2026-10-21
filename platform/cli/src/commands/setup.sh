#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"

# pnpm runs this on postinstall; CI machines need no personal launcher.
if [ -n "${CI:-}" ]; then
  drift_print_info "CI detected: the drift launcher is not linked."
  exit 0
fi

bin_dir=
while [ "$#" -gt 0 ]; do
  case "$1" in
    --bin-dir)
      shift
      [ "$#" -gt 0 ] || drift_die "--bin-dir requires a directory." 2
      bin_dir=$1
      ;;
    *) drift_die "Unknown setup option: $1" 2 ;;
  esac
  shift
done

drift_print_opening_logo
drift_loader_start "Linking launcher"
drift_queue_scene "$DRIFT_ICON_SETUP" "$DRIFT_SETUP_PHRASE"

[ -n "$bin_dir" ] || bin_dir=$(drift_default_bin_dir) || drift_die "Unable to determine a user bin directory."
mkdir -p "$bin_dir"
ln -sfn "$DRIFT_PROJECT_ROOT/platform/cli/drift" "$bin_dir/drift"
drift_log "bin=$bin_dir"
drift_print_success "drift launcher linked at $bin_dir/drift"
drift_check_launcher "$bin_dir"
