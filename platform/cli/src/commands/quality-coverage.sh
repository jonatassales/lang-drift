#!/bin/sh
set -eu
# Test patterns travel to the runners as text; the shell never expands them.
set -f
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

tracked=$workdir/tracked
git ls-files >"$tracked"

packages=$*

if [ -z "$packages" ]; then
  find apps packages -mindepth 2 -maxdepth 2 -name package.json >"$workdir/manifests"
  packages=$(sed 's#/package\.json$##' "$workdir/manifests" | sort)
fi

for package in $packages; do
  [ -f "$package/package.json" ] || drift_die "$package has no package.json." 2
done

status=0

printf '| Package | Runner | Lines | Branches | Files measured |\n'
printf '| --- | --- | --- | --- | --- |\n'

for package in $packages; do
  test_script=$(jq -r '.scripts.test // empty' "$package/package.json")
  [ -n "$test_script" ] || continue

  case "$test_script" in
    vitest*) runner=vitest ;;
    'node --test'*) runner=node ;;
    'go test'*) runner=go ;;
    *)
      printf '| %s | unsupported | | | no coverage for `%s` |\n' "$package" "$test_script"
      continue
      ;;
  esac

  # The files each runner can measure: Vitest is told to count src/, Go counts
  # its module, and node --test can only meet the files its tests load.
  sources=$(awk -v prefix="$package/" -v runner="$runner" '
    index($0, prefix) != 1 || /\/testdata\// || /(_test\.go|\.(unit|test)\.[a-z]+)$/ { next }
    runner == "vitest" && index($0, prefix "src/") == 1 && /\.(ts|tsx)$/ { count++ }
    runner == "node" && /\.(ts|tsx|js|mjs|cjs)$/ { count++ }
    runner == "go" && /\.go$/ { count++ }
    END { print count + 0 }' "$tracked")

  report=$workdir/$(printf '%s' "$package" | tr '/' '-')
  log=$report.log
  lcov=
  tests_passed=yes

  case "$runner" in
    vitest)
      drift_need pnpm
      lcov=$report/lcov.info
      (cd "$package" && pnpm exec vitest run --coverage.enabled --coverage.provider=v8 \
        --coverage.reporter=lcov --coverage.reportsDirectory="$report" \
        --coverage.include='src/**/*.{ts,tsx}' --coverage.exclude='**/*.{unit,test}.{ts,tsx}') \
        >"$log" 2>&1 || tests_passed=no
      ;;
    node)
      lcov=$report.lcov
      patterns=${test_script#node --test}

      case "$patterns" in
        *[\;\&\|\$\`\<\>]*)
          printf '| %s | node | | | unsupported test script `%s` |\n' "$package" "$test_script"
          continue
          ;;
      esac

      # node --test reads every flag before its patterns, so the patterns come
      # from the package's test script, split as the shell would split them.
      : >"$report.patterns"
      [ -z "$patterns" ] || printf '%s\n' "$patterns" | xargs -n 1 printf '%s\n' >"$report.patterns"

      set --
      while IFS= read -r pattern; do
        set -- "$@" "$pattern"
      done <"$report.patterns"

      (cd "$package" && node --test --experimental-test-coverage --test-reporter=lcov \
        --test-reporter-destination="$lcov" "$@") >"$log" 2>&1 || tests_passed=no
      ;;
    go)
      drift_need go
      # Without the test script's -race, which coverage does not need.
      (cd "$package" && go test -coverpkg=./... -coverprofile="$report.out" ./...) \
        >"$log" 2>&1 || tests_passed=no
      ;;
  esac

  if [ "$tests_passed" = no ]; then
    printf '| %s | %s | tests failed | | |\n' "$package" "$runner"
    printf 'drift: %s: the tests failed; the end of their output follows.\n' "$package" >&2
    tail -n 20 "$log" >&2
    status=1
    continue
  fi

  if [ -z "$lcov" ]; then
    go tool cover -func="$report.out" >"$report.func"
    statements=$(awk '/^total:/ { print $NF }' "$report.func")
    measured=$(awk -F: 'NR > 1 && !seen[$1]++ { count++ } END { print count + 0 }' "$report.out")
    printf '| %s | go | %s of statements | n/a | %s of %s |\n' "$package" "$statements" "$measured" "$sources"
    continue
  fi

  if [ ! -s "$lcov" ]; then
    printf '| %s | %s | no coverage report | n/a | 0 of %s |\n' "$package" "$runner" "$sources"
    continue
  fi

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

printf '\nVitest counts every source file under src/, node --test only the files its tests load and Go every package of its module, so percentages from different runners are not comparable.\n'

exit "$status"
