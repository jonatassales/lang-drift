# The events one GitHub Actions run sends to Better Stack, one JSON object per
# line: the run, each job, each step and each Turborepo task.
#
# Inputs, passed with --slurpfile to jq -n:
#   $run        the run attempt, from the GitHub API
#   $jobs       its jobs with their steps, from the GitHub API
#   $summaries  the Turborepo run summaries the run uploaded as artifacts
#
# The API reports what GitHub observed, but a fork's pull request runs its own
# workflow file, so names and branches are capped, and the summaries, which the
# run wrote itself, are read as untrusted data: every field is coerced.

def text: tostring | .[0:200];
def number: if type == "number" then . else null end;
def seconds_between($from; $to):
  if $from and $to then ($to | fromdate) - ($from | fromdate) else null end;

$run[0] as $run
| {
    repository: $run.repository.full_name,
    workflow: ($run.name | text),
    event: $run.event,
    run_id: $run.id,
    run_attempt: $run.run_attempt,
    run_url: $run.html_url,
    sha: $run.head_sha,
    branch: ($run.head_branch // "" | text),
    pull_request: ($run.pull_requests[0].number // null),
    fork: ($run.head_repository.full_name != $run.repository.full_name),
    actor: ($run.triggering_actor.login // $run.actor.login)
  } as $common

| ($common + {
    dt: $run.run_started_at,
    kind: "run",
    conclusion: $run.conclusion,
    duration_seconds: seconds_between($run.run_started_at; $run.updated_at)
  }
  | .message = "\(.workflow) \(.conclusion) in \(.duration_seconds) s"),

  ($jobs[] as $job
  | ($common + {
      dt: ($job.started_at // $job.created_at),
      kind: "job",
      job: ($job.name | text),
      conclusion: $job.conclusion,
      queued_seconds: seconds_between($job.created_at; $job.started_at),
      duration_seconds: seconds_between($job.started_at; $job.completed_at),
      runner: ($job.labels[0] // null)
    }
    | .message = "\(.job) \(.conclusion) in \(.duration_seconds) s"),

    ($job.steps[]? | $common + {
      dt: (.started_at // $job.started_at),
      kind: "step",
      job: ($job.name | text),
      step: (.name | text),
      conclusion: .conclusion,
      duration_seconds: seconds_between(.started_at; .completed_at)
    }
    | .message = "\(.job) › \(.step) \(.conclusion) in \(.duration_seconds) s")),

  ($summaries[]
  | select(type == "object")
  | .tasks
  | if type == "array" then .[:500][] else empty end
  | select(type == "object")
  | (.execution.startTime? | number) as $start
  | (.execution.endTime? | number) as $end
  | $common + {
      dt: (if $start then $start / 1000 | floor | todate else $run.run_started_at end),
      kind: "task",
      task_id: (.taskId | text),
      task: (.task | text),
      package: (.package | text),
      cache: (.cache.status? | if . == "HIT" or . == "MISS" then . else "UNKNOWN" end),
      exit_code: (.execution.exitCode? | number),
      duration_ms: (if $start and $end then $end - $start else null end)
    }
  | .message = "\(.task_id) \(.cache) in \(.duration_ms) ms, exit \(.exit_code)")
