<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>Agent-assisted delivery, accountable to your team's budget.</strong><br>
  Built for humans and AI Agents.
</p>

<p align="center">
  <strong>English</strong> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Pull request checks" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 or newer" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md">Website</a> ·
  <a href="#the-product">The product</a> ·
  <a href="#license">License</a>
</p>

## The product

**LangDrift FinOps** makes agent-assisted delivery accountable to a team budget. It connects what a ticket asks for, the work an agent does, the evidence that work leaves and what it cost, in one chain a team can inspect.

**Microruntimes** is its engine. Its CLI, `ex`, reads a ticket and the documents it links, plans a bounded local runtime around the software you already have, and drives the delivery through your own coding agent, with every paid call reserved against the team's envelope before it runs.

A model call that times out, a retry that costs more than planned, a pull request whose checks never report: delivery leaves the plan in ordinary ways. LangDrift's aim is that none of them becomes spend nobody approved or a number nobody can explain.

### What a team gets

| Capability | Value for the team |
| --- | --- |
| **Plans from tickets** | `ex agent plan` reads the ticket and only the documents it links, retrieves comparable earlier runs and writes a runtime plan sized inside the team's envelope. |
| **Reservations before spend** | Every paid call, the planner's and the coding agent's, is estimated at dated list prices and reserved against the per-run cap, the ticket's allocation and what is left of the period. |
| **Delivery through your harness** | Your workflow stages and skills run through the coding agent you configure. People confirm where the team says they take part, and the merge always stays the owner's. |
| **Runtimes around existing software** | Your components come up on a bounded local kind, Terraform and Kubernetes lifecycle, with no LangDrift dependency in your code. |
| **Memory with evidence** | Every run leaves token-free evidence, a ledger and a memory that the next plan retrieves. |

LangDrift has its own word for a pass of a run: a microruntime is the car, a run is the race, and a **lap** is one pass of it, either an attempt paused at a gate or the execution that finished.

## License

MIT License

Copyright (c) 2026 Jonatas Sales

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
