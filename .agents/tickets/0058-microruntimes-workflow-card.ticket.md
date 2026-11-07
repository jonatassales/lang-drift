# 0058: The Microruntimes hero closes with a Workflow card that plays the run

**Spec:** .agents/specs/0007-microruntimes-hero.spec.md (revised by the owner: the Workflow card shows US dollars for model calls; the runtime's CPU and memory stay in units)

**What to build:** a fourth card, Workflow, opens when the plan session ends. It summarizes the session as money: the run's expected spend for a real Claude model and effort against the $25.00 per-run budget the design doc sets. It lists the delivery stages with their skills and a Human in the loop checkbox each, and a Play workflow button. Play shows a one-second loader, then the terminal clears and plays `ex agent run` through the console's gate, the stages, the microruntime, the push, `#19`, the checks and ready for review. A reservation over the budget pauses the run at the gate instead.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Plan session

- [ ] The guardrail reads FIN-02 · FIN-03 · FIN-04, and the judging line glosses each as the guardrail file does.
- [ ] The price line and the Plan card's price basis speak of the runtime only.
- [ ] The suggestion is `claude-sonnet-5-5` in the ticket plan, the runtime.md editor and the Analysis card, which shows the display name with the id beneath; the design-doc fact carries the budget.
- [ ] The hero summary and the terminal label describe four cards and a simulated ex session.

### Workflow

- [ ] Claude Opus 5.5, Sonnet 5.5 and Haiku 5.5 at Anthropic list prices (Oct 2026), and the efforts low, medium and high. Money is computed in integer cents; the reservation is 20% over the estimate.
- [ ] Stages Planning (human), Spec, Tickets, Implement, Review (human, can be unchecked); Merge is always the owner's and has no checkbox.
- [ ] Choices are locked while a run launches or plays; the run reads a snapshot of them. Replay keeps them.
- [ ] Play is unavailable until the plan has ended; the terminal scrolls into view when the run starts.
- [ ] No gate line names a FIN criterion; observed spend is a dash until the run completes and "nothing paid" at a pause.
- [ ] One status region announces only the run's end.

### Verification

- `pnpm check`, `pnpm turbo run typecheck test build --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/microruntimes/` at 1280, 1024 and 375: defaults end at `#19 · ready for review`; Opus · high pauses at the gate; unchecking Review ends "to merge"; Replay restarts the plan with the choices kept.
