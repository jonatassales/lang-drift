<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: 代理型 FinOps 平台" width="820">
  </picture>
</p>

<p align="center">
  <strong>代理協助交付，對團隊預算負責。</strong><br>
  在執行前落實成本治理，為 AI 程式開發代理提供邊界 runtime。為人類與代理打造。
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <strong>繁體中文</strong> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@langdrift/ex"><img alt="npm @langdrift/ex" src="https://img.shields.io/npm/v/@langdrift/ex?logo=npm&amp;logoColor=white&amp;label=%40langdrift%2Fex&amp;color=CB3837&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/zh-hant/">網站</a> ·
  <a href="#功能特色">功能特色</a> ·
  <a href="#如何使用">如何使用</a> ·
  <a href="#設定檔">設定檔</a> ·
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
- **Node.js 24+** (透過 `npx` 執行 `@langdrift/ex`)
- **Docker**、**kind**、**kubectl** 與 **Terraform** (用於執行本機邊界 runtime；隨時使用 `npx @langdrift/ex doctor` 檢查環境)
- 在環境中設定 `ANTHROPIC_API_KEY` (供模型規劃器與評判使用)

### 1. 初始化專案

在儲存庫根目錄執行初始化指令：

```sh
npx @langdrift/ex init
```

這會在 `.langdrift/config.json` 建立專案識別設定，並在 `.langdrift/microruntimes/AGENTS.md` 建立本機執行規則。既有檔案保持不變。請將 `/.langdrift/` 加入 `.gitignore`：這是 `ex` 所產生的本機狀態。

### 2. 從工單進行規劃

將 `ex` 指向工單或規格文件。它會讀取關聯文件、計算預算額度內的 CPU 與記憶體邊界、建議最佳模型與推理強度，並寫入邊界執行計畫：

```sh
npx @langdrift/ex agent plan <path-to-ticket.md>
```

> **提示：** 使用 `--no-judge` 可純粹透過啟發式規則產生計畫，無需呼叫 LLM，零成本。

<p align="center">
  <img src="/platform/docs/assets/images/ex-agent-plan.webp" alt="Terminal running ex agent plan on ticket 0070. It reads seven linked documents, applies the team's 2000 millicore, 3072 MiB envelope and $25.00 per-run budget, suggests claude-sonnet-5-5 at medium effort, sizes apps/console at 500m and 1024 MiB, and ends with a proposed plan, revision 1." width="820">
</p>

計畫會儲存於 `.langdrift/microruntimes/runtimes/<runtime-id>/runtime.md`。其中 `langdrift` 區塊包含可執行的設定：

```json
{
  "llmBudgetUsd": 25,
  "runtimeId": "console-microruntime",
  "revision": 1,
  "cpuMillis": 2000,
  "memoryMiB": 3072,
  "model": "claude-sonnet-5-5",
  "effort": "medium",
  "components": [
    {
      "path": "apps/console",
      "image": "langdrift-console:local",
      "port": 3001,
      "readinessPath": "/icon.svg",
      "build": { "executable": "docker", "args": ["build", "-t", "langdrift-console:local", "."] }
    }
  ]
}
```

### 3. 執行計畫

執行 runtime 即代表授權該計畫並啟動交付週期：

```sh
npx @langdrift/ex agent run <runtime-id>
```

執行期間，`ex` 會協調整個交付生命週期：

1. **閘門檢查：** 確認標價涵蓋執行內容，並針對工單配額與期間預算保留預估成本。
2. **工作流程階段：** 驅動程式開發代理執行工作流程中所定義的階段。若需要人工確認，執行會暫停並等待終端機輸入或 `npx @langdrift/ex agent confirm <runtime-id> <stage>`。
3. **本機 runtime：** 建置映像檔並在本機 kind 叢集中啟動元件，持續輪詢就緒檢查端點直到狀況正常。
4. **交付觀察：** 以唯讀方式觀察 Git 提交、分支與 PR 檢查。
5. **結算清算：** 銷毀本機 runtime 資源、釋出未使用的預算保留額，並將執行紀錄寫入帳本。

### 4. 檢查狀態與診斷

監控作用中的執行、檢視工作階段事件，並驗證本機工具鏈：

```sh
# 顯示目前執行狀態與圈次進度
npx @langdrift/ex agent status <runtime-id>

# 從歷次執行重新建構情節記憶索引
npx @langdrift/ex agent reindex

# 檢查先決條件與工具鏈健全狀況
npx @langdrift/ex doctor
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

## 授權條款

LangDrift 為依據 [MIT 授權條款](/LICENSE) 發布的開源軟體。
