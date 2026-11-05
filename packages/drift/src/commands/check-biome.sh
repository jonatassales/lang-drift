#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift check biome" 2
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# Biome carries part of the rules in .claude/rules: imports never climb with
# ../ (01-filesystem) and imports are grouped as packages, aliases, then
# relative paths (.agents/context.guide.md). This check runs fixtures against
# a copy of biome.json and fails when the configuration, or a Biome upgrade,
# stops enforcing them. The fixtures are written here rather than committed: biome
# check . would lint them, and their imports are unsorted on purpose.

biome=$DRIFT_PROJECT_ROOT/node_modules/.bin/biome
[ -x "$biome" ] || drift_die "Biome is not installed. Run pnpm install first."

sandbox=$(mktemp -d)
findings=$sandbox/findings
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$sandbox"; drift_flush_results' EXIT

cp biome.json "$sandbox/biome.json"

# run_biome <command> [argument...]: Biome inside the sandbox. The sandbox is
# not a Git checkout, so the VCS integration that biome.json enables is off.
run_biome() {
  subcommand=$1
  shift
  (cd "$sandbox" && "$biome" "$subcommand" --vcs-enabled=false --config-path "$sandbox" "$@")
}

report() {
  printf '%s\n' "$1" >>"$findings"
}

# write_fixture <path>: the source on stdin becomes the file in the sandbox.
write_fixture() {
  mkdir -p "$sandbox/${1%/*}"
  cat >"$sandbox/$1"
}

# expect <label> <expected> <actual>
expect() {
  [ "$2" = "$3" ] && return 0
  report "$1: expected $2, got ${3:-nothing}"
}

# ---- Import restriction ----------------------------------------------------

fixture=apps/website/src/fixture.ts

# restricted_count <path>: the import restrictions Biome reports in the fixture.
restricted_count() {
  run_biome lint --only=style/noRestrictedImports --reporter=json "$1" 2>/dev/null |
    jq '[.diagnostics[] | select(.category == "lint/style/noRestrictedImports")] | length'
}

# module_source <form> <module>: one way to name a module.
module_source() {
  case "$1" in
    import) printf "import { item } from '%s'\nexport { item }\n" "$2" ;;
    import-type) printf "import type { Item } from '%s'\nexport type { Item }\n" "$2" ;;
    reexport) printf "export { item } from '%s'\n" "$2" ;;
    reexport-all) printf "export * from '%s'\n" "$2" ;;
    side-effect) printf "import '%s'\n" "$2" ;;
    dynamic) printf "export const modulePromise = import('%s')\n" "$2" ;;
    require) printf "import loaded = require('%s')\nexport { loaded }\n" "$2" ;;
  esac
}

forms='import import-type reexport reexport-all side-effect dynamic require'

# expect_gap <label>: a gap the owner knows about. The fixture on disk must stay
# unreported; when Biome or biome.json starts catching it, the check fails so
# the case moves to the enforced ones on purpose.
expect_gap() {
  [ "$(restricted_count "$fixture")" = 0 ] && return 0
  report "$1 is now rejected: move it from the known gaps to the enforced cases"
}

check_parent_traversal() {
  for module in .. ../parent ../../parent; do
    for form in $forms; do
      module_source "$form" "$module" | write_fixture "$fixture"
      expect "$form of '$module'" 1 "$(restricted_count "$fixture")"
    done
  done
}

# The pattern group in biome.json is ".." and "../**", which does not match a
# ".." segment inside a longer path. Adding "**/.." and "**/../**" does.
check_embedded_traversal_gap() {
  for module in ./nested/../peer ./nested/.. @components/../peer; do
    module_source import "$module" | write_fixture "$fixture"
    expect_gap "import of '$module'"
  done
}

# What Biome's rule cannot see, whatever the pattern: it reads import and
# export statements and import() calls with a literal string.
check_biome_limits() {
  write_fixture "$fixture" <<'EOF'
export const loaded = require('../parent')
EOF
  expect_gap "require('../parent')"

  write_fixture "$fixture" <<'EOF'
const segment = '../parent'
export const loaded = import(segment)
EOF
  expect_gap "import() of a variable"

  write_fixture "$fixture" <<'EOF'
export type Item = import('../parent').Item
EOF
  expect_gap "an import type"

  write_fixture "$fixture" <<'EOF'
export const loaded = import(`../parent`)
EOF
  expect_gap "import() of a template literal"

  write_fixture "$fixture" <<'EOF'
export { item } from '..\\parent'
EOF
  expect_gap "a backslash path"
}

check_allowed_imports() {
  for module in ./peer ./nested/child @components/child @/domain/product react \
    @testing-library/react web node:path; do
    module_source reexport "$module" | write_fixture "$fixture"
    expect "import of '$module'" 0 "$(restricted_count "$fixture")"
  done
}

# biome.json has no overrides, so every app and package is held to the rule.
check_restriction_scope() {
  for path in apps/website/src/fixture.ts apps/console/src/fixture.ts \
    packages/web/src/fixture.ts packages/tokens/fixture.ts; do
    module_source reexport ../parent | write_fixture "$path"
    expect "../parent in $path" 1 "$(restricted_count "$path")"
  done
}

# ---- Import groups ---------------------------------------------------------

# organize <path> [argument...]: the import organizer alone, without the linter
# and the formatter.
organize() {
  path=$1
  shift
  run_biome check --linter-enabled=false --formatter-enabled=false "$@" "$path" >/dev/null 2>&1
}

# expect_organized <label>: the expected source on stdin. Organizing the fixture
# must produce exactly it. Redirect the fixture itself to expect no change.
expect_organized() {
  cat >"$sandbox/expected.ts"
  organize "$fixture" --write || report "$1: Biome could not organize the imports"
  cmp -s "$sandbox/$fixture" "$sandbox/expected.ts" || report "$1: the imports came out differently"
}

# One import per entry of the alias group in biome.json, so a dropped matcher
# sends its alias back to the packages.
check_import_groups() {
  write_fixture "$fixture" <<'EOF'
import local from './local'
import state from '@state'
import views from '@views/page'
import React from 'react'
import ui from '@ui/button'
import services from '@services/mail'
import scoped from '@testing-library/react'
import web from 'web'
import factories from '@factories/cookie-service'
import app from '@app/page'
import alias from '@/domain/product'
import components from '@components/button'
import request from '@i18n/request'
import proxy from '@proxy/locale'
import header from '@layout/website-header/website-header'
import home from '@routes/home/finops-section'
import reply from '@api/book-a-demo/book-demo-reply.map'
import type { Product } from '@data/products'
EOF
  if organize "$fixture"; then
    report "unsorted imports: Biome accepted them"
  fi

  expect_organized "grouping packages, aliases and relative paths" <<'EOF'
import scoped from '@testing-library/react'
import React from 'react'

import reply from '@api/book-a-demo/book-demo-reply.map'
import app from '@app/page'
import components from '@components/button'
import type { Product } from '@data/products'
import factories from '@factories/cookie-service'
import request from '@i18n/request'
import header from '@layout/website-header/website-header'
import proxy from '@proxy/locale'
import home from '@routes/home/finops-section'
import services from '@services/mail'
import state from '@state'
import ui from '@ui/button'
import views from '@views/page'
import web from 'web'
import alias from '@/domain/product'

import local from './local'
EOF
  organize "$fixture" || report "grouped imports: Biome rejected them"
}

check_side_effect_imports() {
  write_fixture "$fixture" <<'EOF'
import 'z-setup'
import './globals.css'
import 'a-setup'
EOF
  expect_organized "side-effect imports keep their order" <"$sandbox/$fixture"

  write_fixture "$fixture" <<'EOF'
import before from './before'
import React from 'react'
import './setup'
import after from './after'
import state from '@state'
EOF
  expect_organized "a side-effect import is a barrier" <<'EOF'
import React from 'react'

import before from './before'
import './setup'

import state from '@state'

import after from './after'
EOF
}

check_comments() {
  write_fixture "$fixture" <<'EOF'
import local from './local'
// Package-specific intent.
import React from 'react'
import state from '@state'
EOF
  expect_organized "an attached comment moves with its import" <<'EOF'
// Package-specific intent.
import React from 'react'

import state from '@state'

import local from './local'
EOF

  write_fixture "$fixture" <<'EOF'
import local from './local'

// Separate execution group.

import React from 'react'
import state from '@state'
EOF
  expect_organized "a detached comment keeps its group boundary" <<'EOF'
import local from './local'

// Separate execution group.

import React from 'react'

import state from '@state'
EOF
}

check_parent_traversal
check_embedded_traversal_gap
check_biome_limits
check_allowed_imports
check_restriction_scope
check_import_groups
check_side_effect_imports
check_comments

if [ -s "$findings" ]; then
  while IFS= read -r finding; do drift_print_error "$finding"; done <"$findings"
  drift_print_error "biome.json no longer enforces the import rules in .claude/rules/01-filesystem.rule.md and .agents/context.guide.md."
  exit 1
fi

version=$("$biome" --version)
drift_print_success "${version#Version: } rejects ../ imports and groups imports as packages, aliases and relative paths"
