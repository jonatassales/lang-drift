#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
[ "$#" -eq 0 ] || drift_die "Usage: drift quality check" 2
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# Enforces repository rules defined in .claude/rules/*.rule.md.
# Checked in pre-commit and CI.

status=0
fail() {
  drift_print_error "$1"
  status=1
}

# ------------------------------------------------------------------------------
# 00-general.rule.md
# ------------------------------------------------------------------------------
# Platform and tooling scripts use POSIX shell (.sh), never .js, .mjs, .ts or .mts.
non_shell_platform_scripts=$(find platform -type f \( -name '*.js' -o -name '*.mjs' -o -name '*.ts' -o -name '*.mts' \) -not -path '*/node_modules/*' -not -path '*/dist/*')
if [ -n "$non_shell_platform_scripts" ]; then
  for file in $non_shell_platform_scripts; do
    fail "00-general: Platform and tooling scripts must be POSIX shell (.sh): $file"
  done
fi

# pnpm is the only package manager; no alien lockfiles.
alien_locks=$(find . -maxdepth 3 \( -name 'package-lock.json' -o -name 'yarn.lock' -o -name 'bun.lockb' \) -not -path '*/node_modules/*' -not -path '*/.turbo/*')
if [ -n "$alien_locks" ]; then
  for file in $alien_locks; do
    fail "00-general: Only pnpm is allowed; unexpected lockfile: $file"
  done
fi

# ------------------------------------------------------------------------------
# 01-filesystem.rule.md
# ------------------------------------------------------------------------------
# Markdown documents under .agents/workflow/ must carry an approved type suffix.
workflow_markdown=$(find .agents/workflow -type f -name '*.md' -not -name 'README.md' -not -path '*/.effort-log/*')
for file in $workflow_markdown; do
  case "$file" in
    *.spec.md | *.ticket.md | *.adr.md | *.research.md) ;;
    *) fail "01-filesystem: Workflow markdown missing approved type suffix (.spec.md, .ticket.md, .adr.md, .research.md): $file" ;;
  esac
done

rules_markdown=$(find .claude/rules -type f -name '*.md')
for file in $rules_markdown; do
  case "$file" in
    *.rule.md) ;;
    *) fail "01-filesystem: Rule document missing .rule.md suffix: $file" ;;
  esac
done

# ------------------------------------------------------------------------------
# 02-frontend.rule.md & 03-tailwind.rule.md
# ------------------------------------------------------------------------------
# Plain CSS is forbidden. Only the approved exception list (ADR 0003) and entry points are allowed.
allowed_css_files="
apps/website/src/app/[locale]/(pages)/early-access-sent/early-access-sent.css
apps/website/src/app/globals.css
apps/website/src/app/hero.css
apps/website/src/app/lib/components/conversion.css
apps/website/src/app/lib/components/early-access/early-access.css
apps/website/src/app/lib/components/website-footer.css
apps/website/src/app/lib/components/website-header.css
apps/website/src/app/template/tailwind/app.css
apps/website/src/app/template/tailwind/base.css
apps/website/src/app/template/tailwind/utilities.css
apps/console/src/app/template/tailwind/app.css
packages/web/src/styles.css
packages/web/src/tokens.css
packages/web/src/vendors/smoothui/animated-tabs.css
packages/web/src/vendors/smoothui/basic-accordion.css
packages/web/src/vendors/smoothui/smoothui.css
"

tracked_css_files=$(git ls-files '*.css')
for file in $tracked_css_files; do
  is_allowed=0
  for allowed in $allowed_css_files; do
    if [ "$file" = "$allowed" ]; then
      is_allowed=1
      break
    fi
  done
  if [ "$is_allowed" -eq 0 ]; then
    fail "02-frontend: Plain CSS stylesheets forbidden. File not in approved exception list: $file"
  fi
done

# Ratchet ceiling: allowed legacy CSS files may only shrink (Rule 5).
check_ratchet() {
  css_path=$1
  max_lines=$2
  if [ -f "$css_path" ]; then
    actual_lines=$(wc -l <"$css_path" | tr -d ' ')
    if [ "$actual_lines" -gt "$max_lines" ]; then
      fail "02-frontend: Legacy CSS file grew from $max_lines to $actual_lines lines (Rule 5 ratchet violated): $css_path"
    fi
  fi
}

check_ratchet "apps/website/src/app/[locale]/(pages)/early-access-sent/early-access-sent.css" 37
check_ratchet "apps/website/src/app/globals.css" 50
check_ratchet "apps/website/src/app/hero.css" 21
check_ratchet "apps/website/src/app/lib/components/conversion.css" 191
check_ratchet "apps/website/src/app/lib/components/early-access/early-access.css" 154
check_ratchet "apps/website/src/app/lib/components/website-footer.css" 104
check_ratchet "apps/website/src/app/lib/components/website-header.css" 174
check_ratchet "apps/website/src/app/template/tailwind/base.css" 9
check_ratchet "apps/website/src/app/template/tailwind/utilities.css" 0
check_ratchet "packages/web/src/styles.css" 5
check_ratchet "packages/web/src/vendors/smoothui/animated-tabs.css" 100
check_ratchet "packages/web/src/vendors/smoothui/basic-accordion.css" 67
check_ratchet "packages/web/src/vendors/smoothui/smoothui.css" 42

# ------------------------------------------------------------------------------
# 04-code.rule.md
# ------------------------------------------------------------------------------
# No env module, registry, getEnv, requireEnv or Zod env object.
forbidden_env_files=$(find apps packages -type f \( -name 'env.ts' -o -name 'getEnv.ts' -o -name 'env.config.ts' \) -not -path '*/node_modules/*' -not -path '*/dist/*')
if [ -n "$forbidden_env_files" ]; then
  for file in $forbidden_env_files; do
    fail "04-code: Centralized env module forbidden; read external input where used: $file"
  done
fi

# ------------------------------------------------------------------------------
# 05-react.rule.md
# ------------------------------------------------------------------------------
# Barrels export symbols by name, never with export *.
if grep -rn "export \*" packages/web/src/ 2>/dev/null; then
  fail "05-react: Shared UI barrel in packages/web must export by name, never export *"
fi

# ------------------------------------------------------------------------------
# 06-nextjs.rule.md
# ------------------------------------------------------------------------------
# No utils/, helpers/, common/, or misc/ directories in apps.
forbidden_dirs=$(find apps -type d \( -name 'utils' -o -name 'helpers' -o -name 'common' -o -name 'misc' \) -not -path '*/node_modules/*' -not -path '*/.next/*')
if [ -n "$forbidden_dirs" ]; then
  for dir in $forbidden_dirs; do
    fail "06-nextjs: Forbidden utility directory: $dir (no utils/, helpers/, common/ or misc/)"
  done
fi

# App Router under src/app; no Pages Router.
pages_dirs=$(find apps -type d -name 'pages' -not -path '*/(pages)*' -not -path '*/.next/*' -not -path '*/node_modules/*')
if [ -n "$pages_dirs" ]; then
  for dir in $pages_dirs; do
    fail "06-nextjs: Pages Router forbidden; use App Router under src/app: $dir"
  done
fi

# Private app manifests carry no version.
for manifest in apps/*/package.json; do
  version=$(jq -r '.version // empty' "$manifest")
  if [ -n "$version" ]; then
    fail "06-nextjs: Private app manifest carries no version: $manifest"
  fi
done

# ------------------------------------------------------------------------------
# 07-jotai.rule.md
# ------------------------------------------------------------------------------
# Shared atoms in apps/*/src/app/atoms must be named *.atom.ts and never have 'use client'.
atom_dirs=$(find apps -type d -path '*/src/app/atoms' -not -path '*/node_modules/*')
for dir in $atom_dirs; do
  for file in $(find "$dir" -type f); do
    case "$file" in
      *.atom.ts | *.test.ts | *.types.ts | */index.ts) ;;
      *) fail "07-jotai: Atom file must end in .atom.ts: $file" ;;
    esac
    if grep -q "['\"]use client['\"]" "$file"; then
      fail "07-jotai: 'use client' forbidden in atom module: $file"
    fi
  done
done

[ "$status" -ne 0 ] || drift_print_success "All repository rule checks passed (.claude/rules/*.rule.md)"
exit "$status"
