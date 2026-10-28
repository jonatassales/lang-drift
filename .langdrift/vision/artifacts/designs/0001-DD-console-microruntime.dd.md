# 0001-DD: Console microruntime

**Status:** accepted
**PRD:** 0001-PRD-console-preview.prd.md

## Envelope and budget

The team sets the run inside the project's default envelope, 2,000 millicores and 3,072 MiB, and caps model calls at $25.00 per run. The binding values live in `harness/envelope.json` and `finops/budgets.json`; this document records why they were chosen.

## Shape

One image with two processes, because the console and extub accept only loopback hosts and a microruntime runs one container per component without injected variables. The console listens on `0.0.0.0:3001`; extub stays on `127.0.0.1:3401`, reachable only from inside the container.

## Risks

- The cold image build exceeds the agent's 120-second operation bound. The build is pre-warmed.
- The kubelet probes with the pod address as Host, which every data route refuses. Readiness uses `/icon.svg`.
