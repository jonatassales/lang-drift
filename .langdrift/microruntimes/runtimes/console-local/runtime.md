# LangDrift console with extub

Configuration is interpreted by the connected project agent. Prose does not execute commands.

**Ticket:** 0001-TICKET-console-microruntime.ticket.md · run `0001-console-microruntime`
**Context:** sha256 8eab41020062c36b9e88e094b6b818e312e532974cb12b7cba2ffae5ac667e17
**Authorized by:** the owner, 2026-10-28 (illustrative)

The plan took the component from the spec's `langdrift-plan` fence and the envelope from `harness/envelope.json`, because no earlier run measured this runtime. The model and effort configure the console's recommendation; the agent makes no model call.

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
