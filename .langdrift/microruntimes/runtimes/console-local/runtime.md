# LangDrift console with extub

Configuration is interpreted by the connected project agent. Prose does not execute commands.

**Ticket:** 0070-console-microruntime.ticket.md · run `0001-console-microruntime`
**Context:** sha256 0d8ce74c47569bfd35cc820e302898ecde272d93715728f8de68e77ffe8f0dc7
**Envelope and budget:** console input, record `runtime-console-local-r1`, set by the owner on 2026-10-28 (illustrative)

The planner inferred the component from the ticket and the files it links, and the owner confirmed it in the console. Nothing collected set a limit: the envelope and the budget are console input. The model and effort configure the console's recommendation; the agent makes no model call.

```langdrift
{
  "llmBudgetUsd": 25,
  "version": 1,
  "runtimeUid": "701706be7e70d681d608826d3ad548b9",
  "runtimeId": "console-local",
  "projectId": "a759267ef752f98d6c25c4114ad65d0c",
  "revision": 1,
  "name": "LangDrift console with extub",
  "cpuMillis": 2000,
  "memoryMiB": 3072,
  "model": "claude-sonnet-5-5",
  "effort": "medium",
  "components": [
    {
      "path": ".",
      "image": "langdrift-console:local",
      "port": 3001,
      "readinessPath": "/icon.svg",
      "build": {
        "executable": "platform/microruntime/console/build-image.sh",
        "args": []
      }
    }
  ]
}
```
