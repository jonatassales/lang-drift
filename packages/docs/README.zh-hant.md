<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/packages/docs/assets/images/langdrift-banner-dark.png">
    <img src="/packages/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: 代理型 FinOps 平台" width="820">
  </picture>
</p>

<p align="center">
  <strong>代理協助交付，對團隊預算負責。</strong><br>
  在執行前落實成本治理，為 AI 程式開發代理提供邊界 runtime。為人類與代理打造。
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/packages/docs/README.pt-br.md">Português (Brasil)</a> ·
  <strong>繁體中文</strong> ·
  <a href="/packages/docs/README.ja.md">日本語</a> ·
  <a href="/packages/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@langdrift/cli"><img alt="npm @langdrift/cli" src="https://img.shields.io/npm/v/@langdrift/cli?logo=npm&amp;logoColor=white&amp;label=%40langdrift%2Fcli&amp;color=CB3837&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/zh-hant/">網站</a> ·
  <a href="#功能特色">功能特色</a> ·
  <a href="#如何使用">如何使用</a> ·
  <a href="#設定檔">設定檔</a> ·
  <a href="#參考資料">參考資料</a> ·
  <a href="#致謝">致謝</a> ·
  <a href="#授權條款">授權條款</a>
</p>

## 功能特色

- **執行前預算保留：** 每次模型呼叫與雲端資源都會依具日期的標價預估，並在執行開始前自團隊預算額度中保留。若預估超出上限，執行會在產生任何花費前暫停。
- **具邊界的微執行環境 (Microruntimes)：** 為你既有的軟體規劃並執行隔離的本機環境 (Docker、kind、Terraform、Kubernetes)，無需在應用程式程式碼中引入專屬廠商 SDK。
- **自選程式開發代理：** 透過你偏好的代理驅動器 (如 Claude Code) 進行交付，涵蓋結構化的工作流程階段 (規劃、實作、檢視)，並具備確定性的人工作業確認點。
- **確定性稽核紀錄：** 留下不含 token 的證據鏈、可驗證的執行帳本，以及作為後續規劃依據的記憶紀錄。

## 如何使用

### 先決條件

- **macOS** 或 **Linux**
- **Node.js 24+** (透過 `npx` 執行 `@langdrift/cli`)
- **Docker**、**kind**、**kubectl** 與 **Terraform** (用於執行本機邊界 runtime；隨時使用 `npx @langdrift/cli doctor` 檢查環境)
- 在環境中設定 `ANTHROPIC_API_KEY` (供模型規劃器與評判使用)

### 1. 初始化專案

在儲存庫根目錄執行初始化指令：

```sh
npx @langdrift/cli init
```

這會在 `.langdrift/config.json` 建立專案識別設定，並在 `.langdrift/microruntimes/AGENTS.md` 建立本機執行規則。既有檔案保持不變。請將 `/.langdrift/` 加入 `.gitignore`：這是 `ex` 所產生的本機狀態。

### 2. 從工單進行規劃

將 `ex` 指向工單或規格文件。它會讀取關聯文件、計算預算額度內的 CPU 與記憶體邊界、建議最佳模型與推理強度，並寫入邊界執行計畫：

```sh
npx @langdrift/cli agent plan <path-to-ticket.md> --component name=console,path=apps/console
```

> **提示：** 使用 `--no-judge` 可純粹透過啟發式規則產生計畫，無需呼叫 LLM，零成本。

<p align="center">
  <img src="/packages/docs/assets/images/ex-agent-plan.webp" alt="Terminal running ex agent plan on ticket 0070. It reads seven linked documents, applies the team's 2000 millicore, 3072 MiB envelope and $25.00 per-run budget, suggests claude-sonnet-5-5 at medium effort, sizes apps/console at 500m and 1024 MiB, and ends with a proposed plan, revision 1." width="820">
</p>

計畫會儲存於 `.langdrift/microruntimes/runtimes/<runtime-id>/runtime.md`。其中 `langdrift` 區塊包含可執行的設定：

```json
{
  "llmBudgetUsd": 25,
  "version": 2,
  "runtimeId": "console-microruntime",
  "revision": 1,
  "cpuMillis": 2000,
  "memoryMiB": 3072,
  "model": "claude-sonnet-5-5",
  "effort": "medium",
  "components": [
    {
      "name": "console",
      "path": "apps/console",
      "dockerfile": "Dockerfile",
      "port": 3001
    }
  ]
}
```

元件指向檔案：它的名稱與資料夾，資料夾中有一個 Dockerfile。`ex` 從專案根目錄以該 Dockerfile 執行 `docker build`，將映像檔標記為 `langdrift-<runtime-id>-<name>:r<revision>`，並從第一個 `EXPOSE` 取得連接埠、從 Dockerfile 自身的 `HEALTHCHECK` 取得就緒判斷。`--component` 可另外指定 `dockerfile=`、`port=`、`readiness=tcp` 或 `readiness=http:/path`、不含機密且已提交的 `env-file=`，以及 `ex` 啟動時從自身環境讀取的變數 `pass-env=A;B`；`ex` 會將兩者寫入每個元件各一個 Kubernetes Secret。

### 3. 執行計畫

執行 runtime 即代表授權該計畫並啟動交付週期：

```sh
npx @langdrift/cli agent run <runtime-id>
```

執行期間，`ex` 會協調整個交付生命週期：

1. **閘門檢查：** 確認標價涵蓋執行內容，並針對工單配額與期間預算保留預估成本。
2. **工作流程階段：** 驅動程式開發代理執行工作流程中所定義的階段。若需要人工確認，執行會暫停並等待終端機輸入或 `npx @langdrift/cli agent confirm <runtime-id> <stage>`。
3. **本機 runtime：** 依各元件的 Dockerfile 建置，並以元件名稱在本機 kind 叢集中啟動，等待就緒探針直到狀況正常。
4. **交付觀察：** 以唯讀方式觀察 Git 提交、分支與 PR 檢查。
5. **結算清算：** 銷毀本機 runtime 資源、釋出未使用的預算保留額，並將執行紀錄寫入帳本。

### 4. 檢查狀態與診斷

監控作用中的執行、檢視工作階段事件，並驗證本機工具鏈：

```sh
# 顯示目前執行狀態與圈次進度
npx @langdrift/cli agent status <runtime-id>

# 從歷次執行重新建構情節記憶索引
npx @langdrift/cli agent reindex

# 檢查先決條件與工具鏈健全狀況
npx @langdrift/cli doctor
```

在 `status` 後加上 `--json` 可輸出適合主控台與儀表板整合的 NDJSON 事件串流。

## 設定檔

你可以透過儲存庫中的 JSON 設定檔自訂工作流程階段、代理驅動器與成本參數：

<details>
<summary><code>workflow.json</code>：交付階段與人工確認點</summary>

```json
{
  "version": 1,
  "stages": [
    { "id": "planning", "label": "Planning", "skills": ["grilling"], "human": "required" },
    { "id": "spec", "label": "Spec", "skills": ["to-spec"], "human": "optional" },
    { "id": "implement", "label": "Implement", "skills": ["implement", "tdd", "code-review"], "human": "optional" },
    { "id": "review", "label": "Review", "skills": [], "human": "optional" },
    { "id": "merge", "label": "Merge", "skills": [], "human": "owner-only" }
  ],
  "defaults": { "humanAt": ["planning", "review", "merge"] },
  "efforts": ["low", "medium", "high"]
}
```

</details>

<details>
<summary><code>driver.json</code>：程式開發代理驅動器定義</summary>

```json
{
  "executable": "claude",
  "args": ["-p", "{{prompt}}", "--output-format", "json"],
  "passEnv": ["ANTHROPIC_API_KEY"],
  "timeoutSeconds": 1800,
  "usage": { "inputTokens": "usage.input_tokens", "outputTokens": "usage.output_tokens" },
  "summary": "result",
  "unmeasured": "pause"
}
```

</details>

<details>
<summary><code>prices.json</code>：Token 標價清單</summary>

```json
{
  "version": 1,
  "illustrative": true,
  "note": "Replace with your provider's current list prices.",
  "readOn": "2026-10-06",
  "currency": "USD",
  "unit": "cents-per-million-tokens",
  "cache": "uncached",
  "models": [
    { "id": "claude-sonnet-5-5", "name": "Claude Sonnet 5.5", "inputCents": 200, "outputCents": 1000 }
  ]
}
```

</details>

<details>
<summary><code>forecasts.json</code>：依推理強度劃分的 Token 預估量</summary>

```json
{
  "version": 1,
  "illustrative": true,
  "unit": "thousands-of-tokens",
  "byEffort": {
    "low": { "input": 2000, "output": 120 },
    "medium": { "input": 3500, "output": 250 },
    "high": { "input": 6000, "output": 500 }
  },
  "planning": { "input": 35, "output": 5 }
}
```

</details>

## 參考資料

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/)，Jonathan Morley，FinOps Foundation，2026 年 3 月 27 日。實務工作者描述了代理如何將支出與預算對帳、找出浪費，並在開發生命週期更早的階段落實成本政策；最後這個把 FinOps 帶進交付的主題，正是 LangDrift 投入的地方。
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html)，Birgitta Böckeler，martinfowler.com，2026 年 4 月 2 日。文章描述團隊圍繞程式開發代理所打造的外層 harness：在代理行動前引導它的指引，以及讓它自我修正的感測器。`ex harness` 為專案記錄這套 harness，而 FinOps 閘門則限制它能花多少錢。

## 致謝

感謝 [Matt Pocock](https://x.com/mattpocockuk)。本儲存庫所使用的代理 skills，從 grilling 到 to-spec、to-tickets、implement、tdd 與 code-review，都改編自以 MIT 授權釋出的 [mattpocock/skills](https://github.com/mattpocock/skills)，上方的範例工作流程也使用了其中幾個。

## 授權條款

LangDrift 為依據 [Apache 2.0 授權條款](/LICENSE) 發布的開源軟體。
