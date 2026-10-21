# Connected project agent

Run `ex agent connect <console-url>` with `RUNTIME_SECRET` supplied through the environment, then `ex agent run`. The agent is part of the standalone Go executable; application dependencies are unnecessary. It uses an MCP client, explicit repository build recipes, and the private project kind/Terraform/Kubernetes lifecycle. It does not install tools or commit local changes.

The console owns runtime composition and requested intent. `runtime.md` contains one strictly parsed `exphere` JSON fence; free prose is retained as evidence and never executed. Editing managed files locally produces a conflict rather than silently overwriting them. Secrets are stored only in the ignored runtime `.env`; connection state and pending events remain in `.exphere/state`.

Loss of the 15-second communication lease cancels active work and requests bounded cleanup. Cleanup failure stays unconfirmed. Reconnection synchronizes events before consulting current console intent. A user stop, deletion, or revoked credential prevents resume. The shared cluster and project sources are preserved.

## Consumer-owned context

The bundled [SDD workflow](../../workflow/README.md), [templates](../../templates), and [skills](../../workflow/skills/README.md) are reference material from Exphere's installed snapshot. They do not establish organizational adoption, grant permissions, install hooks, or run an agent framework. Source-project references unavailable in this bundle are labeled consumer-owned; adapt those references to your project. Existing consumer files and customizations always take precedence. No project root AGENTS.md is created or replaced.

Pocock material retains its [upstream MIT license](../../workflow/skills/UPSTREAM-LICENSE) and attribution. The 26 installed skills are bundled; `writing-for-agents` is excluded as in the source catalog. No automatic skill discovery is claimed.
