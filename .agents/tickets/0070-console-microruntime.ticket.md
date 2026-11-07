# 0070: Run the console with console-stubs as a microruntime

**Spec:** .agents/specs/0010-exphere-and-the-langdrift-folder.spec.md

**What to build:** `apps/console` and `apps/microservices/console-stubs` run together as one microruntime: one image with two processes, console-stubs on `127.0.0.1:3401` and the console on `0.0.0.0:3001`, because both apps accept only loopback hosts and a microruntime runs one container per component without injected variables. Readiness uses `/icon.svg`, a static route the kubelet can probe without a loopback Host. The image build is pre-warmed, because each agent operation is bounded to 120 seconds. The team expects the run to fit in about 2,000 millicores and 3,072 MiB; the limits themselves are set in the console.

**Blocked by:** the Exphere server, not yet ticketed.

**Status:** wontfix. The owner closed it on 2026-11-05: it will not be implemented.

- [ ] The runtime reaches ready, and `/en/home` answers through a port-forward.
- [ ] CPU and memory are observed against the envelope, with their timestamps.
- [ ] Stopping removes every owned resource and keeps the evidence.

### Context

- [What the real `ex` reproduces of the Microruntimes page](/.agents/research/ex-microruntimes-capability.research.md)
- [The console's build recipes are allowlisted](/.agents/adr/0002-console-input-is-untrusted.adr.md)
- [FinOps guardrails](/.langdrift/harness/guardrails/0001-finops.gr.md)
- [The console's manifest](/apps/console/package.json)
- [the console-stubs module](/apps/microservices/console-stubs/go.mod)
