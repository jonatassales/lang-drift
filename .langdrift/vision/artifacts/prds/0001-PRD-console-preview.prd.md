# 0001-PRD: A console preview anyone on the team can run

**Status:** accepted

## Problem

The LangDrift console only runs on a developer's machine, as two processes started by hand: the Next.js console and the extub API it reads. Nobody can bring up the same console for a review, a demo or a test without repeating that setup, and nothing records how much it consumed.

## Outcome

One request brings the console up, bounded, observed and cleaned up afterwards, and leaves a record of what it used.

## Scope

- In: the console with its stub data, on the project's local cluster.
- Out: real accounts, a hosted environment, any paid model call.
- Open: whether the preview runs per pull request.

## Use cases

1. As a reviewer, I want to open the console from a ticket's run, so that I review the change without setting it up.
2. As the owner, I want each run's CPU and memory recorded against the envelope, so that the next plan sizes from evidence instead of a guess.
