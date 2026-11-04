# Run the console with extub as a microruntime

Configuration is interpreted by the connected project agent. Prose does not execute commands.

**Ticket:** 0070-console-microruntime.ticket.md · run `console-microruntime`
**Context:** sha256 0d8ce74c47569bfd35cc820e302898ecde272d93715728f8de68e77ffe8f0dc7
**Envelope and budget:** CLI input, record `envelope-7c2a`, set by owner on 2026-10-28

The planner read the ticket and the files it links, and the owner confirmed the components through the CLI. Nothing collected set a limit: the envelope and the budget are CLI input. The model and effort configure the suggestion.

```langdrift
{
  "llmBudgetUsd": 25,
  "version": 2,
  "runtimeUid": "701706be7e70d681d608826d3ad548b9",
  "runtimeId": "console-microruntime",
  "projectId": "a759267ef752f98d6c25c4114ad65d0c",
  "revision": 1,
  "name": "Run the console with extub as a microruntime",
  "cpuMillis": 2000,
  "memoryMiB": 3072,
  "model": "claude-sonnet-5-5",
  "effort": "medium",
  "components": [
    {
      "name": "console",
      "path": ".",
      "dockerfile": "Dockerfile",
      "port": 3001
    }
  ]
}
```

```langdrift-plan
{
  "version": 1,
  "run": "console-microruntime",
  "state": "authorized",
  "ticket": {
    "path": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
    "number": "0070",
    "slug": "console-microruntime",
    "title": "Run the console with extub as a microruntime",
    "status": "needs-triage",
    "namedBy": { "person": "owner", "via": "cli" }
  },
  "collection": {
    "agent": {
      "model": "claude-sonnet-5-5",
      "promptVersion": "judge-1",
      "paid": true,
      "ledgerLine": "planning"
    },
    "limits": {
      "maxDepth": 3,
      "maxFiles": 64,
      "maxBytes": 4194304,
      "secrets": "refused",
      "outsideProject": "refused",
      "symlinks": "refused"
    },
    "reached": { "files": 7, "depth": 1, "bytes": 38125 },
    "context": {
      "sha256": "0d8ce74c47569bfd35cc820e302898ecde272d93715728f8de68e77ffe8f0dc7",
      "rule": "sha256 over the sorted path<TAB>sha256 lines of the sources read; byte hashes until canonical hashing lands"
    }
  },
  "sources": [
    {
      "path": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "mediaType": "text/markdown",
      "bytes": 1459,
      "sha256": "6d7935ac7e4eb3b4036b5a9d6ec8fe1d22fd1d0ec05fc0a000d30305ac8faae7",
      "depth": 0,
      "linkedFrom": null,
      "linkSyntax": "named in the CLI",
      "state": "read",
      "verdict": "kept",
      "reason": "the ticket itself",
      "confidence": 1,
      "inferredRole": "ticket"
    },
    {
      "path": ".agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md",
      "mediaType": "text/markdown",
      "bytes": 10570,
      "sha256": "3dde504fc77f063715c765dae8d51fce6535de5631068ad4b215545b8981c98a",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "field line",
      "state": "read",
      "verdict": "kept",
      "reason": "states that guardrails come only from console input and that Microruntimes ships with FinOps",
      "confidence": 0.9,
      "inferredRole": "spec"
    },
    {
      "path": ".agents/workflow/research/ex-microruntimes-capability.research.md",
      "mediaType": "text/markdown",
      "bytes": 17725,
      "sha256": "7f187a177ccf87f94bae151a50b6692503b61875c98dbbb504047e55058e5e25",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "markdown link",
      "state": "read",
      "verdict": "discarded",
      "reason": "audits the website simulation and the CLI; adds no input for sizing, budget or effort beyond what the ticket states",
      "confidence": 0.71,
      "inferredRole": "research"
    },
    {
      "path": ".agents/workflow/adrs/0002-console-input-is-untrusted.adr.md",
      "mediaType": "text/markdown",
      "bytes": 1314,
      "sha256": "94cd83fdbd581c5be2d2b6068f498631a074ea0d236995b0b206a9502bae90f4",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "markdown link",
      "state": "read",
      "verdict": "kept",
      "reason": "bounds the build recipe the component may use",
      "confidence": 0.88,
      "inferredRole": "decision"
    },
    {
      "path": ".langdrift/harness/guardrails/0001-finops.gr.md",
      "mediaType": "text/markdown",
      "bytes": 5159,
      "sha256": "52dec7dc62975afb1a8b46734a6297e0593a77109aff98c37a224f5c1baea446",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "markdown link",
      "state": "read",
      "verdict": "kept",
      "reason": "names the FinOps criteria the run is reviewed against",
      "confidence": 0.93,
      "inferredRole": "guardrail"
    },
    {
      "path": "apps/console/package.json",
      "mediaType": "application/json",
      "bytes": 854,
      "sha256": "c53acfb6db82dbb7e60a2ca80181137343109de285b96136e9b080cf978437f7",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "markdown link",
      "state": "read",
      "verdict": "kept",
      "reason": "the console's stack and port",
      "confidence": 0.95,
      "inferredRole": "component manifest"
    },
    {
      "path": "apps/extub/go.mod",
      "mediaType": "text/plain",
      "bytes": 1044,
      "sha256": "aafa98e9c16d359d76ab8fce921f64232162e7808132050b02ab0da8e5fc1f3f",
      "depth": 1,
      "linkedFrom": ".agents/workflow/tickets/0070-console-microruntime.ticket.md",
      "linkSyntax": "markdown link",
      "state": "read",
      "verdict": "kept",
      "reason": "extub's Go toolchain, needed to build the image",
      "confidence": 0.9,
      "inferredRole": "component manifest"
    }
  ],
  "facts": [
    {
      "topic": "resource",
      "text": "expected to fit in about 2,000 millicores and 3,072 MiB",
      "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L5"
    },
    {
      "topic": "component",
      "text": "one image, two processes: extub on 127.0.0.1:3401, the console on 0.0.0.0:3001",
      "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L5"
    },
    {
      "topic": "component",
      "text": "readiness on /icon.svg",
      "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L5"
    },
    {
      "topic": "resource",
      "text": "each agent operation is bounded to 120 seconds; pre-warm the build",
      "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L5"
    },
    {
      "topic": "acceptance",
      "text": "3 acceptance criteria",
      "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L11-L13"
    },
    {
      "topic": "component",
      "text": "Next.js 16.4.0; scripts bind port 3001",
      "source": "apps/console/package.json#L6-L17"
    },
    {
      "topic": "component",
      "text": "Go 1.26.0, toolchain go1.26.10",
      "source": "apps/extub/go.mod#L3-L5"
    },
    {
      "topic": "component",
      "text": "a build recipe is a committed script or docker build, go, npm, pnpm, make",
      "source": ".agents/workflow/adrs/0002-console-input-is-untrusted.adr.md#L7"
    },
    {
      "topic": "budget",
      "text": "FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load · review rules, no automatic enforcement",
      "source": ".langdrift/harness/guardrails/0001-finops.gr.md#L16-L18"
    }
  ],
  "inputs": {
    "cpuMillis": {
      "value": 2000,
      "record": "envelope-7c2a",
      "setBy": "owner",
      "via": "cli",
      "setAt": "2026-10-28T13:46:50Z"
    },
    "memoryMiB": {
      "value": 3072,
      "record": "envelope-7c2a",
      "setBy": "owner",
      "via": "cli",
      "setAt": "2026-10-28T13:46:50Z"
    },
    "modelBudgetCents": {
      "value": 2500,
      "record": "budget-2026-10-9d4e",
      "setBy": "owner",
      "via": "cli",
      "setAt": "2026-10-28T13:46:50Z"
    },
    "ceiling": {
      "cpuMillis": 4000,
      "memoryMiB": 6144,
      "record": "envelope-7c2a",
      "setBy": "owner",
      "via": "cli",
      "setAt": "2026-10-28T13:46:50Z"
    },
    "notRuntimeFields": ["gpu", "disk"]
  },
  "discrepancies": [
    {
      "topic": "resource",
      "fact": "about 2,000 millicores and 3,072 MiB",
      "input": "2,000 millicores and 3,072 MiB",
      "result": "consistent"
    }
  ],
  "retrieval": {
    "method": "vector+structured",
    "criteria": [
      "same-component-path",
      "same-stack",
      "measured-resources",
      "measured-tokens"
    ],
    "keepAtMatches": 2,
    "considered": [],
    "note": "first run of this project; no earlier run to cite"
  },
  "sizing": {
    "basis": "input",
    "measuredRunsCited": [],
    "note": "no cited episode measured cpu or memory · sizing follows the inputs, not history"
  },
  "suggestion": {
    "by": "recommendation",
    "model": "claude-sonnet-5-5",
    "modelName": "Claude Sonnet 5.5",
    "effort": "medium",
    "reasons": [
      "one component, no broker wiring: high effort is not warranted",
      "Opus 5.5 at high effort reserves more than the per-run budget"
    ],
    "authorizedBy": "owner"
  },
  "components": [
    {
      "path": ".",
      "name": "console",
      "kind": "web",
      "technologies": [
        { "id": "nextjs", "label": "Next.js" },
        { "id": "go", "label": "Go" }
      ],
      "dockerfile": "Dockerfile",
      "image": "langdrift-console-microruntime-console:r1",
      "port": 3001,
      "readiness": {
        "kind": "exec",
        "from": "HEALTHCHECK of Dockerfile",
        "path": null,
        "command": [
          "/bin/sh",
          "-c",
          "wget -q --spider http://127.0.0.1:3001/icon.svg"
        ],
        "periodSeconds": 5,
        "timeoutSeconds": 2,
        "failureThreshold": 12,
        "initialDelaySeconds": 0
      },
      "envFile": null,
      "passEnv": [],
      "inferredFrom": ["Dockerfile"],
      "confirmedBy": { "person": "owner", "via": "cli" },
      "proposed": { "cpuMillis": 500, "memoryMiB": 1024 },
      "executed": {
        "cpuMillis": 2000,
        "memoryMiB": 3072,
        "rule": "equal split of the runtime envelope; runtime.md v2 carries no per-component share"
      }
    }
  ],
  "totals": {
    "proposed": { "cpuMillis": 500, "memoryMiB": 1024 },
    "bounded": { "cpuMillis": 2000, "memoryMiB": 3072 },
    "headroom": { "cpuMillis": 1500, "memoryMiB": 2048 }
  },
  "materialization": { "deployments": 1, "replicasEach": 1 },
  "price": {
    "record": "prices-0001",
    "mirror": ".langdrift/finops/prices/0001-illustrative.prices.json",
    "readOn": "2026-10-06",
    "runtimePriced": false,
    "note": "the runtime has no price source; cpu and memory stay in units"
  },
  "checks": [
    { "id": "schema", "result": "passed", "errors": 0 },
    { "id": "ports-distinct", "result": "passed", "detail": "3001" },
    { "id": "dockerfiles-exist", "result": "passed", "detail": "Dockerfile" },
    {
      "id": "readiness-known",
      "result": "passed",
      "detail": "HEALTHCHECK of Dockerfile"
    },
    {
      "id": "env-file-exists",
      "result": "passed",
      "detail": "no env file named"
    },
    {
      "id": "within-ceiling",
      "result": "passed",
      "detail": "2,000m of 4,000m · 3,072 MiB of 6,144 MiB"
    }
  ],
  "runtime": {
    "path": ".langdrift/microruntimes/runtimes/console-microruntime/runtime.md",
    "revision": 1,
    "sha256": "25c3c8b0e381a60d60e3c88cffe71e21a20ea82635a2deb0da8dbeea0fbe3459",
    "state": "authorized"
  },
  "contractGaps": [
    "per-component cpu and memory are not runtime.md v2 fields, so the proposal and the execution can differ",
    "technologies are not runtime.md v2 fields; the plan keeps them for display"
  ],
  "authorization": { "by": "owner", "via": "cli", "at": "2026-10-28T13:46:50Z" }
}
```
