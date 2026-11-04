# Security Policy

## Reporting a vulnerability

Use either channel, whichever you prefer:

- **GitHub private vulnerability reporting**, from the
  [Security tab](https://github.com/jonatassales/lang-drift/security/advisories/new).
  This is the fastest route: it opens a private thread attached to the
  repository, so the report reaches the engineers who own the code with no
  inbox routing in between.
- **Email security@langdrift.md**, if you would rather not use GitHub or want
  to report something outside this repository.

Please do not open a public GitHub issue, pull request, or public discussion for a
security problem. A public report puts every LangDrift user at risk during the
window before a fix ships.

A useful report includes:

- the affected component and version, or the commit you tested against
- what an attacker can do, and what access they need to start
- steps to reproduce, ideally with the exact requests or a short script
- anything you already know about impact or blast radius

You do not need a finished exploit. A clear description of the weakness is
enough to open a case.

If you have a PGP key preference or need an encrypted channel, say so in your
first email and we will arrange one.

## What we commit to

| Stage | Target |
| --- | --- |
| Acknowledgement that a human has the report | 2 business days |
| Initial assessment, including whether we reproduced it | 5 business days |
| Fix or documented mitigation for critical and high severity | 30 days |
| Fix or documented mitigation for medium and low severity | 90 days |

We will keep you updated as the work progresses rather than going quiet, and we
will tell you if a date is going to slip.

When a fix ships we publish a GitHub Security Advisory for the affected package
or component, request a CVE where one is warranted, and mark affected released
versions so that dependency scanners pick the fix up.

## Coordinated disclosure

We ask for 90 days from your report before public disclosure, or until a fix is
released, whichever comes first. If a fix is going to take longer than that we
will come back to you and explain why rather than expect silence by default.

We are glad to credit you in the advisory and the release notes. Tell us how you
would like to be named and whether to include an affiliation. If you would
rather stay anonymous, that is fine too.

## Scope

### In scope

- **FinOps governance & budget boundaries**: flaws that allow an agent driver,
  crafted prompt, or malicious payload to bypass pre-execution budget reservations,
  trigger unauthorized model spend, manipulate list price tables, or evade
  deterministic human-in-the-loop approval gates (*"Recommendation is not authorization"*).
- **Microruntime isolation & execution bounds**: vulnerabilities in the `ex` CLI
  or Microruntimes harness allowing container escapes (Docker, kind, Terraform),
  arbitrary host filesystem traversal, command injection outside intended task bounds,
  or unauthorized network egress.
- **Secret & evidence integrity**: leakage of project credentials (from
  `.langdrift/config.json`, console secrets, or environment variables) into logs,
  token-free evidence trails, run ledgers, or episodic memory.
- **Platform assets**: `https://langdrift.md` (the public site), the console
  (`apps/console`), and stub APIs (`apps/extub`).
- **Published packages**: `@langdrift/ex` (including platform binaries and launchers)
  and `@langdrift/contracts`.

### Out of scope

- findings that only affect a local development or self-hosted deployment where
  the operator deliberately disabled sandboxing or ran containers in privileged host mode.
- theoretical issues requiring pre-existing root or administrative access on the
  developer's local workstation.
- vulnerabilities in third-party model providers (Anthropic, OpenAI, etc.) or cloud
  vendors; these should be reported directly to that vendor.
- reports produced purely by automated vulnerability scanners with no demonstrated
  attack path or impact.
- missing hardening headers, cookie flags, or TLS configuration without an exploitable
  threat vector.
- rate limiting on local development loopback stubs (`apps/extub`) with no security
  consequence.
- social engineering or physical attacks against LangDrift maintainers or users.
- outdated dependencies with no reachable or exploitable path in LangDrift.

## Safe harbour

If you make a good-faith effort to follow this policy, we will treat your
research as authorized. We will not pursue or support legal action against you,
and if a third party brings action over research that followed this policy, we
will make it known that you were authorized.

Good faith means:

- you only access data belonging to accounts or environments you own or have
  permission to test.
- you stop as soon as you have confirmed a vulnerability, rather than pivoting
  deeper into the system or accessing unauthorized host files.
- you do not degrade service for other users, exfiltrate data, or destroy or
  modify data that is not yours.
- if you encounter sensitive data or credentials, you stop, do not save them,
  and report them immediately.
- you give us a reasonable window to fix the issue before going public.

If you are unsure whether something is in scope or whether a test crosses a
line, email us first at security@langdrift.md. We would rather answer a question
than have you guess.

## Reporting a security problem in someone else's LangDrift deployment

If you discover an exposed public LangDrift console or stub instance that is
unauthenticated and not yours, please alert us at security@langdrift.md and we
will attempt to reach the operator. Do not test against or exploit it.
