#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift env setup" 2
drift_loader_start "Preparing env"
drift_queue_scene "$DRIFT_ICON_ENV" "$DRIFT_ENV_PHRASE"

# Create each missing .env empty and private; values are filled in by hand.
drift_find_env_templates | while IFS= read -r template; do
  target=${template%.development}
  if [ ! -e "$target" ]; then
    (umask 077 && : > "$target")
    drift_print_success "Created ${target#"$DRIFT_PROJECT_ROOT"/}"
  fi
done
