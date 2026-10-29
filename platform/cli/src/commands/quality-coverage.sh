#!/bin/sh
set -eu
. "$DRIFT_CLI_DIR/core/common.sh"
drift_need jq
cd "$DRIFT_PROJECT_ROOT"

# The test coverage of every workspace package that has a test script, measured
# by the package's own runner and printed as a Markdown table. The runners count
# differently: Vitest counts every source file under src/, node --test only the
# files its tests load, and Go every package of the module, so the Files column
# says what each figure covers. Arguments narrow the run to some packages:
# drift quality coverage apps/console packages/ex

workdir=$(mktemp -d)
# drift's own EXIT trap prints the queued results, so this one chains to it.
trap 'rm -rf "$workdir"; drift_flush_results' EXIT

packages=$*
[ -n "$packages" ] || packages=$(find apps packages -mindepth 2 -maxdepth 2 -name package.json -exec dirname {} \; | sort)

status=0

printf '| Package | Runner | Lines | Branches | Files measured |\n'
printf '| --- | --- | --- | --- | --- |\n'

for package in $packages; do
  manifest=$package/package.json
  [ -f "$manifest" ] || drift_die "$package has no package.json." 2

  test_script=$(jq -r '.scripts.test // empty' "$manifest")
  [ -n "$test_script" ] || continue

  case "$test_script" in
    vitest*) runner=vitest ;;
    'node --test'*) runner=node ;;
    'go test'*) runner=go ;;
    *)
      printf '| %s | unknown | | | no coverage for `%s` |\n' "$package" "$test_script"
      continue
      ;;
  esac

  report=$workdir/${package##*/}
  tests_passed=yes

  case "$runner" in
    vitest)
      drift_need pnpm
      (cd "$package" && pnpm exec vitest run --coverage.enabled --coverage.provider=v8 \
        --coverage.reporter=lcov --coverage.reportsDirectory="$report" \
        --coverage.include='src/**/*.{ts,tsx}' --coverage.exclude='**/*.{unit,test}.{ts,tsx}') \
        >"$report.log" 2>&1 || tests_passed=no
      ;;
    node)
      # node --test reads every flag before its patterns, so the patterns come
      # from the package's test script, the same command pnpm test runs.
      eval "set -- ${test_script#node --test}"
      (cd "$package" && node --test --experimental-test-coverage --test-reporter=lcov \
        --test-reporter-destination="$report.lcov" "$@") >"$report.log" 2>&1 || tests_passed=no
      ;;
    go)
      drift_need go
      (cd "$package" && go test -coverpkg=./... -coverprofile="$report.out" ./...) \
        >"$report.log" 2>&1 || tests_passed=no
      ;;
  esac

  if [ "$tests_passed" = no ]; then
    printf '| %s | %s | tests failed | | |\n' "$package" "$runner"
    drift_print_error "$package: the tests failed; the end of their output follows."
    tail -n 20 "$report.log" >&2
    status=1
    continue
  fi

  if [ "$runner" = go ]; then
    sources=$(git ls-files "$package" | grep '\.go$' | grep -v '_test\.go$' | grep -cv '/testdata/' || true)
    measured=$(awk -F: 'NR > 1 { print $1 }' "$report.out" | sort -u | wc -l | tr -d ' ')
    statements=$(go tool cover -func="$report.out" | awk '/^total:/ { print $NF }')
    printf '| %s | go | %s of statements | n/a | %s of %s |\n' "$package" "$statements" "$measured" "$sources"
    continue
  fi

  lcov=$report.lcov
  [ "$runner" = node ] || lcov=$report/lcov.info
  [ -f "$lcov" ] || : >"$lcov"

  sources=$(git ls-files "$package/src" | grep -E '\.(ts|tsx)$' | grep -cvE '\.(unit|test)\.(ts|tsx)$' || true)
  awk -F: -v package="$package" -v runner="$runner" -v sources="$sources" '
    /^SF:/ { files++ }
    /^LF:/ { found_lines += $2 }
    /^LH:/ { hit_lines += $2 }
    /^BRF:/ { found_branches += $2 }
    /^BRH:/ { hit_branches += $2 }
    END {
      lines = found_lines ? sprintf("%.1f%%", hit_lines * 100 / found_lines) : "no data"
      branches = found_branches ? sprintf("%.1f%%", hit_branches * 100 / found_branches) : "n/a"
      printf "| %s | %s | %s | %s | %d of %d |\n", package, runner, lines, branches, files, sources
    }' "$lcov"
done

exit "$status"
