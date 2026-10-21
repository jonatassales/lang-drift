---
name: code-review
description: "Review the changes since a fixed point (commit, branch, tag, or merge-base) along two axes: Standards (does the code follow this repo's documented coding standards?) and Spec (does the code match what the originating issue/spec asked for?). Runs both reviews in parallel sub-agents and reports them side by side. Use when the user wants to review a branch, a PR, work-in-progress changes, or asks to \"review since X\"."
---

Two-axis review of a scoped change against a fixed point, including relevant uncommitted changes when reviewing work in progress:

- **Standards**: does the code conform to this repo's documented coding standards?
- **Spec**: does the code faithfully implement the originating issue / spec?

Both axes run as **parallel sub-agents** so they don't pollute each other's context, then this skill aggregates their findings.

Read `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md` for the configured tracker, artifact paths and states. If it is missing, report that configuration gap; do not infer a remote tracker.

## Process

### 1. Pin the scope and comparison

Use the paths and task scope the user authorized. Inspect repository status to identify committed, staged, unstaged and untracked changes, then select the relevant paths; unrelated edits, generated artifacts and vendored files are outside the review unless explicitly in scope. State exclusions. If ownership is ambiguous, resolve that ambiguity before reviewing the affected paths.

Use the fixed point the user supplied (commit, branch, tag or merge-base), or the recorded starting commit for this task. Validate explicit refs before proceeding. For a work-in-progress review without an earlier base, use the current `HEAD` and state that the review covers local changes only. For a branch review whose intended base cannot be inferred, ask for the base rather than silently choosing one.

Capture one consistent review snapshot and provide it to both reviewers:

- **Committed branch changes:** `git diff <fixed-point>...HEAD -- <scoped-paths>` and `git log <fixed-point>..HEAD --oneline`.
- **Work in progress:** include `git diff --cached -- <scoped-paths>` and `git diff -- <scoped-paths>`, plus final file contents where staged and unstaged edits overlap. When an earlier base is part of the request, also include the committed changes from that base. An empty committed diff does not mean there is no work to review.
- **Relevant untracked files:** enumerate with `git ls-files --others --exclude-standard`, select only task-owned paths, and read those files as additions. Avoid ingesting secrets or unrelated files. No staging or committing is needed for review.
- **Repository with no commit yet:** compare staged content with the empty tree (obtain its ID using `git hash-object -t tree /dev/null`), include unstaged changes and scoped untracked additions. State that the baseline is an empty repository, not a previous version.
- **Directory without Git metadata:** review only the explicit task file manifest or user-specified paths as initial additions. State that no historical diff is available; do not initialize Git merely to review.

Record base, selected paths and the snapshot's contents or hashes so both axes examine the same state. If the working files change before completion, refresh the affected snapshot or report the review as stale. If no scoped changes exist across all relevant sources, report that fact and stop; never expand to the entire repository just to obtain a diff.

### 2. Identify the spec source

Look for the originating spec, in this order:

1. Issue references in the commit messages (`#123`, `Closes #45`, GitLab `!67`, etc.), fetched via the workflow in the tracker doc.
2. A path the user passed as an argument.
3. The matching existing spec under `.exphere/artifacts/specs/`, following `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md`; use other locations only when the project explicitly configures them.
4. If nothing is found, ask the user where the spec is. If they say there isn't one, the **Spec** sub-agent will skip and report "no spec available".

### 3. Identify the standards sources

Search the repo for every file that documents how code should be written. When `CODING_STANDARDS.md` or `CONTRIBUTING.md` exists, it must be on the list.

On top of whatever the repo documents, the Standards axis always carries the **smell baseline** below: a fixed set of Fowler code smells (_Refactoring_, ch.3) that applies even when a repo documents nothing. Two rules bind it:

- **The repo overrides.** A documented repo standard always wins; where it endorses something the baseline would flag, suppress the smell.
- **Always a judgment call.** Each smell is a labeled heuristic ("possible Feature Envy"), never a hard violation. Like any standard here, skip anything tooling already enforces.

Each smell reads *what it is* → *how to fix*; match it against the diff:

- **Mysterious Name**: a function, variable, or type whose name doesn't reveal what it does or holds. → rename it; if no honest name comes, the design's murky.
- **Duplicated Code**: the same logic shape appears in more than one hunk or file in the change. → extract the shared shape, call it from both.
- **Feature Envy**: a method that reaches into another object's data more than its own. → move the method onto the data it envies.
- **Data Clumps**: the same few fields or params keep travelling together (a type wanting to be born). → bundle them into one type, pass that.
- **Primitive Obsession**: a primitive or string standing in for a domain concept that deserves its own type. → give the concept its own small type.
- **Repeated Switches**: the same `switch`/`if`-cascade on the same type recurs across the change. → replace with polymorphism, or one map both sites share.
- **Shotgun Surgery**: one logical change forces scattered edits across many files in the diff. → gather what changes together into one module.
- **Divergent Change**: one file or module is edited for several unrelated reasons. → split so each module changes for one reason.
- **Speculative Generality**: abstraction, parameters, or hooks added for needs the spec doesn't have. → delete it; inline back until a real need shows.
- **Message Chains**: long `a.b().c().d()` navigation the caller shouldn't depend on. → hide the walk behind one method on the first object.
- **Middle Man**: a class or function that mostly just delegates onward. → cut it, call the real target direct.
- **Refused Bequest**: a subclass or implementer that ignores or overrides most of what it inherits. → drop the inheritance, use composition.

### 4. Spawn both sub-agents in parallel

Issue both sub-agent calls together, in the foreground, and aggregate the reports they return.

**Standards sub-agent prompt** should include:

- The scoped review snapshot, selected paths, comparison commands, base and commit list (or the stated absence of history).
- The list of standards-source files you found in step 3, **plus the smell baseline from step 3** pasted in full (the sub-agent has no other access to it).
- The brief: "Report, per file/hunk where relevant, (a) every place the diff violates a documented standard: cite the standard (file + the rule); and (b) any baseline smell you spot: name it and quote the hunk. Distinguish hard violations from judgment calls: documented-standard breaches can be hard, but baseline smells are always judgment calls, and a documented repo standard overrides the baseline. Skip anything tooling enforces. Under 400 words."

**Spec sub-agent prompt** should include:

- The same scoped review snapshot, selected paths, comparison commands, base and commit list (or the stated absence of history).
- The path or fetched contents of the spec.
- The brief: "Report: (a) requirements the spec asked for that are missing or partial; (b) behavior in the diff that wasn't asked for (scope creep); (c) requirements that look implemented but where the implementation looks wrong. Quote the spec line for each finding. Under 400 words."

If the spec is missing, skip the Spec sub-agent and note this in the final report.

### 5. Aggregate

Present the two reports under `## Standards` and `## Spec` headings, verbatim or lightly cleaned. Do **not** merge or rerank findings, because the two axes are deliberately separate (see _Why two axes_).

End with a one-line summary: total findings per axis, and the worst issue _within each axis_ (if any). Don't pick a single winner across axes: that's the reranking the separation exists to prevent.

## Why two axes

A change can pass one axis and fail the other:

- Code that follows every standard but implements the wrong thing → **Standards pass, Spec fail.**
- Code that does exactly what the issue asked but breaks the project's conventions → **Spec pass, Standards fail.**

Reporting them separately stops one axis from masking the other.
