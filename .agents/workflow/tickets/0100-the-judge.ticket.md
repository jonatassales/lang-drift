# 0100: The planner judges, extracts and suggests through a model

**Spec:** .agents/workflow/specs/0015-ex-performs-the-microruntimes-session.spec.md

**What to build:** the plan's judged lines print for real: a one-line fact per document read, the guardrail criteria that apply, a reason per candidate episode, the weighing of effort and model, and the suggestion, bounded to the price list and the three efforts. The judge is a port with an adapter on the official Anthropic Go SDK, using structured JSON outputs validated against the plan's schemas and the credential the SDK resolves at the call; each call is paid work, gated on the planning forecast before the first call, watched during, and written to the ledger's planning line with its tokens and cents. Without a credential or with `--no-judge` the plan prints one line saying so and keeps the judged fields `null`. The planner's model is the input recorded by `ex harness model set`.

**Blocked by:** 0098, 0099.

**Status:** ready-for-agent

- [ ] With a fake judge returning the sample's facts, verdicts, reasons and suggestion, the whole plan session prints the 43 lines and the plan fence equals the expected fixture, including the planning line of the ledger.
- [ ] The adapter is tested against a local test server: a schema-valid answer is accepted, an invalid one is refused without a retry loop, usage is recorded per call, and a refusal stop reason leaves the fields `null`.
- [ ] A planning reservation above the per-run cap pauses before the first call, with nothing paid and exit code 3.
- [ ] No test reaches the network.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; a real plan on this repository with a credential prints judged lines and the ledger's planning line shows the cents paid.
