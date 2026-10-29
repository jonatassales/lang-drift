# One Turborepo run summary (.turbo/runs/<id>.json) as Markdown for the GitHub
# job summary: the run's totals, then each task with its cache status,
# duration and exit code.

def seconds: . / 100 | round / 10;
def duration: if .startTime and .endTime then "\(.endTime - .startTime | seconds) s" else "–" end;

"### `\(.execution.command)`",
"",
"\(.execution.attempted) tasks, \(.execution.cached) cached, \(.execution.failed) failed, in \(.execution | duration)",
"",
"| Task | Cache | Duration | Exit |",
"| --- | --- | ---: | ---: |",
(.tasks | sort_by(.taskId)[] |
  "| `\(.taskId)` | \(.cache.status) | \(.execution // {} | duration) | \(.execution.exitCode // "–") |"),
""
