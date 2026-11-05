# 0228: release `@langdrift/cli` 0.5.0 and deploy the judge (owner)

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** The owner publishes the release that carries the judge as a platform service and the run that brings the runtime up first: `@langdrift/cli` 0.5.0 on npm, the tag `lib/cli/v0.5.0` (and `lib/core/contracts/go/v0.2.0` for the judge surface) and the GitHub Release; on the console's Vercel project, `RUNTIME_SECRET` and the platform's provider key set for production and preview, so `/api/judge` answers; then this repository connects its own project and plans and runs its microruntime with the published CLI, closing ticket 0198's run criteria.

**Blocked by:** 0221, 0223, 0227, 0313, 0314, 0315.

**Status:** ready-for-human

- [ ] `@langdrift/cli@0.5.0` is on npm; `pnpm dlx @langdrift/cli@0.5.0 --version` prints `dft 0.5.0`; the tags and the GitHub Release exist.
- [ ] `RUNTIME_SECRET` and the platform's provider key exist on `langdrift-console` for production and preview, and `POST /api/judge` on the deployed console answers 401 without the secret and 200 with a sample question.
- [ ] In this repository, with `RUNTIME_SECRET` exported, `pnpm exec dft agent connect` records the connection, `sh platform/microruntime/plan.sh` writes a judged plan, and `pnpm exec dft agent run book-a-demo-microruntime --timeout 30m` shows four Ready pods and waits at `planning`.
- [ ] Ticket 0198's two open criteria are ticked, and the root `package.json` pins `@langdrift/cli` 0.5.0.

## Dev part

**Override:** No part qualifies: a release and deployment configuration (policy, Categories: never configuration). 2026-11-04.

## Verification

- `pnpm dlx @langdrift/cli@0.5.0 --version`
- `curl -i -X POST https://console.langdrift.md/api/judge` (401), then with the secret and a sample question (200).
- `sh platform/microruntime/plan.sh` and `pnpm exec dft agent run book-a-demo-microruntime --timeout 30m`, then `kubectl --kubeconfig .langdrift/agent/state/cluster.kubeconfig get pods -n book-a-demo-microruntime`.
