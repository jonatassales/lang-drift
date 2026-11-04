#!/bin/sh
# Plans this repository's Book a Demo microruntime with the platform's own
# CLI, pinned at the root as @langdrift/cli. Run it from anywhere inside the
# checkout; it is idempotent: a project already initialized is kept, every
# input is recorded again from the committed files beside this script, and
# a plan after the first raises the runtime's revision.
#
#   sh platform/microruntime/plan.sh
set -eu

root=$(git rev-parse --show-toplevel)
inputs=$root/platform/microruntime
ticket=.agents/workflow/tickets/0198-book-a-demo-microruntime.ticket.md
allocation=0198-book-a-demo-microruntime

cd "$root"

# The project: .langdrift/ is machine-local and ignored, so a fresh checkout
# initializes it. No CI gates workflow is generated into this repository.
if [ ! -f .langdrift/config.json ]; then
  pnpm exec ex init --ask=false
fi

# The harness: the envelope and the model the plan needs, the workflow with
# every stage waiting for a person, and the driver that would play a stage.
pnpm exec ex harness envelope set \
  --per-execution-cpu-millis 4000 \
  --per-execution-memory-mib 6144 \
  --default-cpu-millis 2000 \
  --default-memory-mib 3072
pnpm exec ex harness model set claude-sonnet-5-5 --effort medium
pnpm exec ex harness workflow set --file "$inputs/workflow.json"
pnpm exec ex harness driver set --file "$inputs/driver.json"

# FinOps: the month's budget with this ticket's allocation, the forecasts,
# and the prices once, since each prices record is a new numbered file.
pnpm exec ex finops budget set \
  --period "$(date -u +%Y-%m)" \
  --period-cents 10000 \
  --per-run-cents 2500 \
  --reserve-margin-percent 20 \
  --allocate "$allocation=2500"
pnpm exec ex finops forecasts set --file "$inputs/forecasts.json"
prices=.langdrift/finops/prices
if [ ! -d "$prices" ] || [ -z "$(find "$prices" -name '*-illustrative.prices.json')" ]; then
  pnpm exec ex finops prices set --file "$inputs/prices.json" --slug illustrative
fi

# The runtime: four components named after their folders, the services with
# their committed env files and the secrets ex reads from its environment at
# start. No judge: this plan spends nothing and needs no provider key.
pnpm exec ex agent plan "$ticket" --no-judge \
  --component "name=messaging,path=infrastructure/messaging,tech=redis:Redis" \
  --component "name=storage,path=infrastructure/storage,tech=minio:MinIO,pass-env=MINIO_ROOT_USER;MINIO_ROOT_PASSWORD" \
  --component "name=demo-requests,path=apps/microservices/demo-requests,env-file=apps/microservices/demo-requests/runtime.env,pass-env=MINIO_ACCESS_KEY;MINIO_SECRET_KEY;DEMO_REQUESTS_TOKEN" \
  --component "name=mailer,path=apps/microservices/mailer,env-file=apps/microservices/mailer/runtime.env,pass-env=RESEND_API_KEY;DEMO_LEAD_EMAIL"
