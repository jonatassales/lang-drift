<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>代理協助的交付，對團隊預算負責。</strong><br>
  為人類與 AI 代理打造。
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <strong>繁體中文</strong> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Pull request 檢查" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 或更新版本" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/zh-hant/">網站</a> ·
  <a href="#產品">產品</a> ·
  <a href="#授權條款">授權條款</a>
</p>

## 產品

**LangDrift FinOps** 讓代理協助的交付對團隊預算負責。它把工單要求的事、代理所做的工作、那份工作留下的證據，以及它花了多少錢，串成一條團隊可以檢視的鏈。

**Microruntimes** 是它的引擎。它的 CLI `ex` 會讀取一張工單與它連結的文件，在你既有的軟體周圍規劃一個有界的本機 runtime，並透過你自己的程式開發代理推動交付；每一次付費呼叫在執行前，都會先對團隊的預算額度保留。

模型呼叫逾時、重試的花費超出計畫、pull request 的檢查遲遲沒有回報：交付偏離計畫的方式都很平常。LangDrift 的目標，是讓這些情況都不會變成沒有人核准的支出，或沒有人說得清楚的數字。

### 團隊得到什麼

| 能力 | 對團隊的價值 |
| --- | --- |
| **從工單產生計畫** | `ex agent plan` 只讀取工單與它連結的文件，取回可比較的過往執行，並寫出一份在團隊預算額度內配置好資源的 runtime 計畫。 |
| **花費前先保留** | 每一次付費呼叫，無論來自規劃器或程式開發代理，都會以有日期的牌價估算，並對單次執行上限、工單的分配與本期剩餘額度先行保留。 |
| **透過你的 harness 交付** | 你的工作流程階段與 skills 透過你設定的程式開發代理執行。人員在團隊指定參與的地方確認，而合併永遠由負責人決定。 |
| **圍繞既有軟體的 runtime** | 你的元件在有界的本機 kind、Terraform 與 Kubernetes 生命週期上啟動，你的程式碼不需要任何 LangDrift 相依套件。 |
| **有證據的記憶** | 每次執行都會留下不含 token 的證據、一本帳本，以及下一次規劃會取回的記憶。 |

LangDrift 為一次執行中的每一輪取了自己的名字：microruntime 是賽車，一次執行（run）是一場比賽，而一**圈**（lap）是其中的一輪，可能是在閘門前暫停的一次嘗試，也可能是跑完的那次執行。

## 授權條款

*授權條款的文字維持英文。*

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
