# Run console-microruntime: Run the console with extub as a microruntime

**Ticket:** 0070-console-microruntime.ticket.md
**Kind:** ticket
**Outcome:** completed in 2 laps, ready for review
**Proposed ticket status:** in-review, pending the owner

Lap 0001 stopped in the pit: the owner's choice over the suggestion, Claude Opus 5.5 at high effort asked for a reservation above the per-run cap, and nothing was paid. Lap 0002 followed the suggestion, Claude Sonnet 5.5 at medium effort and finished with $9.80 observed of $11.40 reserved.

## Worth remembering

- At this ticket's size, Claude Opus 5.5 at high effort needs a reservation above the $25.00 cap: $40.80 was asked for.
- 690 MiB of 3,072 MiB and 410m of 2000m at peak: the next plan can propose with this headroom in mind.
- Claude Sonnet 5.5 at medium effort used 3.7M input and 240K output tokens; the next estimate for this ticket can start from these instead of the effort table.
- `HEALTHCHECK of Dockerfile` passed readiness for `.`.

```langdrift-run
{
  "run": "console-microruntime",
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
    "claude-opus-5-5-high-exceeds-per-run-cap",
    "memory-peak-690-of-3072-mib",
    "cpu-peak-410-of-2000m",
    "claude-sonnet-5-5-medium-used-3.7m-input-240k-output",
    "readiness-healthcheck-of-dockerfile",
    "checks-4-of-4-passed"
  ]
}
```
