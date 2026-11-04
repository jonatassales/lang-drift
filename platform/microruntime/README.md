# This repository's microruntime

The Book a Demo services run as a microruntime of the platform's own engine: `@langdrift/cli`, pinned at the repository root to a published release, plans and runs them the way it plans and runs a consumer's software (spec 0032, ADR 0004). This folder holds the committed inputs and the script that records them and plans the runtime; `.langdrift/`, where `ex` keeps the project's state, stays machine-local and ignored.

## Prerequisites

- Docker with a reachable engine, kind, kubectl and Terraform; `pnpm exec ex doctor` lists what it finds.
- `pnpm install` at the root, which installs the pinned `@langdrift/cli`.
- For a run: the secrets the components read from the environment at start, exported in the shell that runs `ex agent run`: `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` for `storage`; `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY` and `DEMO_REQUESTS_TOKEN` for `demo-requests`; `RESEND_API_KEY` and `DEMO_LEAD_EMAIL` for `mailer`. The archive's access and secret key are the MinIO root credentials.

## Commands

| Command | Does |
| --- | --- |
| `pnpm exec ex doctor` | Checks the tools and reports the project cluster, if one exists. |
| `sh platform/microruntime/plan.sh` | Initializes the project when `.langdrift/config.json` is absent, records the harness and FinOps inputs from the files in this folder, and plans the runtime from ticket 0198. Idempotent: an initialized project is kept, the inputs are recorded again, the prices once, and a plan after the first raises the runtime's revision. The runtime id is `book-a-demo-microruntime`. |
| `pnpm exec ex agent run book-a-demo-microruntime --timeout 30m` | Reserves against the budget and plays the workflow. Every stage waits for a person, so the run pauses at `planning`; `ex agent confirm book-a-demo-microruntime planning` lets it play on. The images are built and the pods come up when the stages have been played. |
| `pnpm exec ex agent status book-a-demo-microruntime` | Where the run stands. |
| `kubectl --kubeconfig .langdrift/agent/state/cluster.kubeconfig get pods -n book-a-demo-microruntime` | The pods, once the runtime is up. |

A runtime that is waiting with a lap open cannot be planned again until the lap ends; `ex` refuses it with "the runtime is running". There is no command to abandon a waiting run yet.

## Files

| File | Holds |
| --- | --- |
| `plan.sh` | The POSIX script above: `ex init`, the `ex harness` and `ex finops` inputs, then `ex agent plan` with the four components, without the judge, so planning spends nothing and needs no provider key. |
| `workflow.json` | The delivery stages, every one with `"human": "required"` (`merge` is `owner-only`), so a run pauses at the first stage and nothing is played without a person. |
| `driver.json` | The coding agent a confirmed stage would run: `claude -p`, reading `ANTHROPIC_API_KEY`, pausing when a call's usage is not measured. |
| `prices.json` | Illustrative list prices per million tokens for the models the driver may run; the owner replaces them with the provider's current list. |
| `forecasts.json` | Illustrative token forecasts by effort, and for planning, until this repository's own runs produce observed figures. |

The envelope (2,000 millicores and 3,072 MiB per execution, under a 4,000 and 6,144 ceiling), the model choice (`claude-sonnet-5-5`, effort `medium`) and the month's budget ($100 for the period, $25 per run, $25 allocated to ticket 0198, a 20% reserve margin) are recorded by `plan.sh` from its own flags; they are illustrative figures the owner adjusts there.

## Components

| Component | Folder | Port | Env |
| --- | --- | --- | --- |
| `messaging` | `infrastructure/messaging` | 6379 | none |
| `storage` | `infrastructure/storage` | 9000 | pass-through `MINIO_ROOT_USER`, `MINIO_ROOT_PASSWORD` |
| `demo-requests` | `apps/microservices/demo-requests` | 3402 | `runtime.env` (port, `redis://messaging:6379`, `storage:9000`, bucket) and pass-through `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `DEMO_REQUESTS_TOKEN` |
| `mailer` | `apps/microservices/mailer` | 3403 | `runtime.env` (port, `redis://messaging:6379`) and pass-through `RESEND_API_KEY`, `DEMO_LEAD_EMAIL` |

Each component is its folder and the `Dockerfile` in it; `ex` derives the image tag, the build (`docker build --file <folder>/Dockerfile .` from the repository root, bounded by the root `.dockerignore`), the port from `EXPOSE` and the readiness probe from `HEALTHCHECK`. Components address each other by name: a Service per component. A `runtime.env` holds no secret; the secrets arrive through the pass-through names, which `ex` reads from its own environment and writes into one Kubernetes Secret per component.
