# Run 0001: the console preview as a microruntime

> **Illustrative sample.** No run happened. Hashes match the files in this tree; times, ids and figures show the format. The branch, commit and pull request in `run.json` are placeholders.

**Ticket:** 0070-console-microruntime.ticket.md
**Kind:** ticket
**Outcome:** completed in two laps
**Proposed ticket status:** in-review, pending the owner

The first lap stopped in the pit: the owner chose Claude Opus 5.5 at high effort over the suggestion, and its reservation exceeded the per-run cap. The second lap followed the suggestion, Claude Sonnet 5.5 at medium effort, and finished inside the envelope and the budget.

## Worth remembering

- At this ticket's size, Opus at high effort needs a reservation above the $25.00 cap. Plan Sonnet at medium, or raise the cap first.
- A cold image build of the console takes longer than the agent's 120-second operation bound; build it once before the start.
- `/icon.svg` passes readiness without a loopback Host; data routes answer 403 to the kubelet.
- Inside the pod, extub's WebRTC offers only loopback candidates, so the console falls back to HTTP polling.
- 690 MiB of 3,072 MiB at peak: the next run can propose 1,024 MiB with headroom.
- Sonnet at medium used 3.7M input and 240K output tokens; the next estimate for this ticket can start from these instead of the effort table.

```langdrift-run
{
  "run": "0001-console-microruntime",
  "sample": true,
  "ticket": "0070-console-microruntime",
  "kind": "ticket",
  "outcome": "completed",
  "laps": 2,
  "componentPaths": ["."],
  "stacks": ["nextjs", "go"],
  "measuredResources": true,
  "measuredTokens": true,
  "observed": {
    "cpuMillisPeak": 410,
    "memoryMiBPeak": 690,
    "tokensThousands": { "input": 3700, "output": 240 }
  },
  "facts": [
    "opus-high-exceeds-per-run-cap",
    "prewarm-build-over-120s",
    "readiness-icon-svg",
    "webrtc-loopback-fallback"
  ]
}
```
