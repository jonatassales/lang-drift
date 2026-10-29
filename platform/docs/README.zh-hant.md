<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift：Agentic FinOps Platform。全隊共用一份預算額度，每次付費呼叫都先保留再執行。" width="820">
  </picture>
</p>

<p align="center">
  <strong>代理協助的交付，對團隊預算負責。</strong>
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <strong>繁體中文</strong> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="#專案狀態"><img alt="狀態：開發中的候選版本" src="https://img.shields.io/badge/status-development%20candidate-FBBF24?labelColor=18181B" height="20"></a>
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Pull request 檢查" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 或更新版本" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="https://langdrift.md/zh-hant/"><img alt="語言：en、pt-BR、zh-Hant、ja、vi" src="https://img.shields.io/badge/locales-en%20%C2%B7%20pt--BR%20%C2%B7%20zh--Hant%20%C2%B7%20ja%20%C2%B7%20vi-6D0AD6?labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/zh-hant/">網站</a> ·
  <a href="#產品">產品</a> ·
  <a href="#預算額度">預算額度</a> ·
  <a href="#開始使用-ex">開始使用 <code>ex</code></a> ·
  <a href="#專案狀態">專案狀態</a>
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

## 預算額度

授權的預算屬於團隊，不屬於某位開發者或某張工單。工單在其中得到分配，而每次執行都有上限。

在一圈開始之前，`ex` 會以牌價估算這次執行，並加上團隊設定的保留邊際。落在單次執行上限、工單分配與本期剩餘額度之內的保留會開啟這一圈；超出的則會暫停，而且什麼都還沒付。這一圈進行時，每一次付費呼叫前都會經過閘門。結束時，未用到的保留會退回預算。

| 閘門 | 檢查什麼 | 未通過時 |
| --- | --- | --- |
| **FG-01** | 有一份有日期的價格來源涵蓋付費工作。 | 阻擋這一圈。 |
| **FG-02** | 工單有一筆分配。 | 阻擋這一圈。 |
| **FG-03** | 估算加上邊際，落在單次執行上限、分配與本期剩餘額度之內。 | 在付任何錢之前暫停。 |
| **FG-04** | 這一圈觀測到的花費尚未達到它的保留。 | 暫停付費呼叫。 |
| **FG-05** | 負責人已確認先前某一圈的超支。 | 需要負責人以 `--acknowledge-overspend` 確認。 |
| **FG-06** | 供應商已為這次執行開立帳單。 | 數字維持為暫定。 |

**未知永遠不是零。** 沒有回報用量的呼叫、沒有人能確認的清理、尚未送達的帳單，都會維持未知並附上原因，它的保留也會持續保留，直到對帳完成。

## 你的技術棧，你的決定

LangDrift 負責邊界與機制；你的團隊負責應用程式與其領域上的選擇。

- **很小的足跡。** 在你的專案中，`ex` 只寫入 `.langdrift/`。它讀取工單連結的產出物，絕不讀取專案之外的內容，並拒絕秘密與符號連結。
- **你的交付，由它觀察。** `ex` 從不 push、不開 pull request、也不合併。交付屬於你的 harness；`ex` 透過 `git` 與 `gh` 觀察分支、pull request 與它的檢查。
- **建議不等於授權。** 規劃器會建議一個模型與推理強度；真正授權計畫的，是執行這個 runtime。模型的建議永遠不會擴大範圍、資金或存取權。
- **憑證留在環境裡。** 規劃器在呼叫當下讀取 `ANTHROPIC_API_KEY`，程式開發代理只拿到最小的環境，加上你列出的變數。`ex` 從不以命令列引數接收憑證。

每一項可能阻擋或放行一次執行的輸入，都會連同設定者、時間與途徑一起記錄。`.langdrift/` 下的 JSON 檔是這些紀錄的唯讀鏡像；編輯鏡像不會改變任何閘門讀取的內容。

## 專案狀態

**LangDrift 仍在開發中。** 網站已在 [langdrift.md](https://langdrift.md/zh-hant/) 上線，提供五種語言。主控台讀取 `extub`，這是一個植入單一專案 LangDrift Studio 的 stub API，並把它的 FinOps 數據標示為示意。`ex` 是本機開發的候選版本，需要從原始碼建置；目前沒有發佈任何二進位檔。

`ex agent plan` 與 `ex agent run` 可以在專案內於本機運作，金錢閘門在 `ex` 內針對專案自己的帳本執行；透過伺服器共享的團隊帳本尚未建置。連線代理（`ex agent connect`、`project`、`supervise` 與 `migrate`）需要一個本儲存庫還沒有的伺服器。測試套件以腳本化的模型、程式開發代理與生命週期驅動 CLI；它們並不證明真實叢集能啟動、真實的花費或已可用於正式環境。網站的註冊目前暫停。

| 路徑 | 內容 |
| --- | --- |
| [`apps/website`](/apps/website/README.md) | langdrift.md 上的平台網站：首頁、FinOps、Microruntimes、價格與搶先體驗。 |
| [`apps/console`](/apps/console/README.md) | 團隊管理 FinOps 與其 microruntimes 的主控台。 |
| [`apps/extub`](/apps/extub/README.md) | 主控台讀取的本機 Go stub API。 |
| [`packages/ex`](/packages/ex/README.md) | 以 Go 撰寫的 `ex` CLI。 |
| [`packages/contracts`](/packages/contracts/README.md) | Microruntimes 的契約：`runtime.md` 區塊與代理協定。 |
| [`packages/tokens`](/packages/tokens) · [`packages/web`](/packages/web) | 設計 token 與共用 UI。 |

## 開始使用 `ex`

你需要 macOS 或 Linux，以及 Go 1.26 或更新版本；這樣就足以規劃。規劃器的模型評審會讀取 `ANTHROPIC_API_KEY`。執行一圈還需要 Docker、kind、kubectl 與 Terraform（`ex doctor` 會檢查它們，且不安裝任何東西）、你設定的程式開發代理（例如 Claude Code），以及用來觀察 pull request 的 `gh`。

### 1. 建置 CLI

```sh
git clone https://github.com/jonatassales/lang-drift.git
cd lang-drift/packages/ex
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex version
```

`version` 會印出 `ex 0.0.0-dev`。你的系統可能已經有一個 `ex`，也就是 Vi 的行編輯器，所以請用路徑呼叫這一個。以下步驟以 `ex` 代表 `packages/ex/dist/ex`；只在這個 shell 中執行 `alias ex="$PWD/dist/ex"` 即可。

### 2. 初始化專案

想在現成的專案上試用，就複製測試套件所規劃的那個專案：一張工單、它連結的 spec、研究、ADR 與 guardrail，以及兩個應用程式的 manifest。

```sh
cp -R src/testdata/langdrift/project /tmp/ex-demo
cd /tmp/ex-demo
ex init
```

`ex init` 會在 `.langdrift/config.json` 建立專案識別，並在 `.langdrift/microruntimes/AGENTS.md` 寫入執行規則，既有檔案保持不變。在你自己的儲存庫裡，請把 `/.langdrift/` 加進 `.gitignore`：那是 `ex` 寫入的本機狀態，團隊撰寫的內容都不會放在那裡。

### 3. 記錄 harness

把[範例輸入檔](#範例輸入檔)存到專案中，然後記錄交付 harness：硬體預算額度、工作流程階段、預設模型與推理強度，以及執行各階段的程式開發代理。

```sh
ex harness envelope set --per-execution-cpu-millis 4000 --per-execution-memory-mib 6144 \
  --default-cpu-millis 2000 --default-memory-mib 3072 --not-runtime-field gpu
ex harness workflow set --file workflow.json
ex harness model set claude-sonnet-5-5 --effort medium
ex harness driver set --file driver.json
```

### 4. 記錄金錢規則

這份預算給團隊 2026 年 10 月共 500.00 美元，每次執行上限 25.00 美元，在每次估算上保留 20% 的邊際，並分配 25.00 美元給工單 0070。金額以美分為單位。

```sh
ex finops budget set --period 2026-10 --period-cents 50000 --per-run-cents 2500 \
  --reserve-margin-percent 20 --allocate 0070-console-microruntime=2500
ex finops gates set --default
ex finops prices set --file prices.json --slug list
ex finops forecasts set --file forecasts.json
```

每個命令都會印出它寫入的紀錄，而每個輸入命令都接受 `--by`，用來記下設定者。

### 5. 從工單規劃

指定 runtime 要建置並執行的元件，然後讓 `ex` 讀取工單。`--no-judge` 會在不呼叫模型的情況下規劃，所以不會產生任何費用；不加這個旗標並設定 `ANTHROPIC_API_KEY` 時，規劃器會在與一圈相同的閘門下透過模型評斷各個來源，費用記在帳本的規劃那一行。

```sh
ex agent plan .agents/workflow/tickets/0070-console-microruntime.ticket.md --no-judge \
  --component "path=apps/console,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=docker,args=build -t langdrift-console:local ."
```

CLI 的輸出是英文：

```text
ex agent plan 0070-console-microruntime.ticket.md
# plan only · nothing is built or started · the planner's calls are paid and gated
agent › tool call: read 0070-console-microruntime.ticket.md
agent › needs-triage · links 0010-exphere-and-the-langdrift-folder.spec.md · ex-microruntimes-capability.research.md · 0002-console-input-is-untrusted.adr.md · 0001-finops.gr.md · package.json · go.mod
judge      unavailable · --no-judge; verdicts, facts and the suggestion are the inputs' own
# analysis · reading the linked documents, then the project's episodes
agent › tool call: read 0010-exphere-and-the-langdrift-folder.spec.md · read
agent › tool call: read ex-microruntimes-capability.research.md · read
agent › tool call: read 0002-console-input-is-untrusted.adr.md · read
agent › tool call: read 0001-finops.gr.md · read
agent › tool call: read package.json · read
agent › tool call: read go.mod · read
✓ read       7 documents · nothing outside the ticket's links
reads      .langdrift/harness/guardrails/0001-finops.gr.md · FIN-01 · FIN-02 · FIN-03 · FIN-04 · FIN-05 · FIN-06 · FIN-07 · FIN-08
bounded    cpuMillis 2000 · memoryMiB 3072 · set by owner in the CLI · record envelope-b228
budget     $25.00 per run for model calls · set by owner in the CLI
agent › RAG over project episodes · a model and effort for this ticket
agent › embedding the ticket intent · 1 query vector
agent › vector search · project episodes · top 6 by cosine similarity
agent › 0 candidates · no episode indexed yet
agent › dropped 0 below 0.55 · kept 0
agent › 0 episodes cited · none with measured resources
agent › reasoning · no cited episode measured cpu or memory · sizing follows the inputs, not history
✓ suggests   claude-sonnet-5-5 · effort medium
# plan · resolving the components, then sizing them inside the envelope
agent › resolving 1 components named in the CLI · paths, images, ports
agent › tool call: inspect apps/console · package.json · :3001
✓ apps/console          Node.js+Next.js · web 500m · 1024 MiB
agent › cross-checking · 1 ports distinct · 1 images resolved · no collisions
✓ size       500m of 2000m · 1024 of 3072 MiB · 1 one-replica Deployments · proposed shares
limits     gpu and disk are not runtime.md fields · none requested
estimate   runtime has no price source · cpu and memory stay in millicores and MiB
agent › tool call: validate runtime.md draft against the schema · 0 errors
✓ plan       runtime.md draft, revision 1 · proposed · authorize with ex agent run
# next: once the team authorizes, ex agent run applies the recorded intent
```

計畫會寫入 `.langdrift/microruntimes/runtimes/console-microruntime/runtime.md`。其中的 `langdrift` 區塊是 runtime 會執行的內容；第二個區塊 `langdrift-plan` 保存分析。檔案中的文字只作為脈絡保留，永遠不會被執行。

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

這段摘錄省略了區塊的版本、名稱，以及由 `ex` 管理的 runtime 與專案識別。在 runtime 停止時再次規劃同一張工單，會提高修訂版本。

### 6. 執行計畫

執行 runtime 就代表授權它的計畫。一圈會建置映像檔、啟動專案的本機叢集，並透過你的程式開發代理進行付費呼叫，所以請先執行 `ex doctor`。

```sh
ex agent run console-microruntime
```

1. **閘門。** 執行會印出意圖、以牌價計算的估算與單次執行預算，然後進行保留。落在所有界限之內的保留會開啟這一圈；超出的則暫停，什麼都還沒付。
2. **階段。** 各階段依序透過程式開發代理執行，並各印出一行摘要。需要人員參與的階段會等待負責人：在終端機回答、使用 `--confirm`，或之後執行 `ex agent confirm console-microruntime planning`。`owner-only` 階段永遠不會執行。
3. **Runtime。** 元件依區塊在專案的 kind 叢集上啟動，就緒時執行會告訴你。
4. **交付。** 執行會觀察分支、pull request 與它的檢查；`--checks-wait` 限制追蹤的時間。
5. **結算。** runtime 關閉，未用到的保留退回預算，寫入記憶 `run.md`，並重建 episode 索引。執行會說明它是否已可審查或合併，以及接下來要做什麼。

每一項變更都會記錄在這一圈的檔案、`run.json` 與帳本中。

### 7. 追蹤執行

```sh
ex agent status console-microruntime
ex agent reindex
ex doctor
```

`status` 顯示執行紀錄與進行中的那一圈，`reindex` 從每個 runtime 的記憶重建 episode 索引，`doctor` 檢查原生工具。加上 `--json` 時，每一行工作階段輸出都是一個 NDJSON 事件 `{at, session, kind, id, text, card, field, data}`，讓主控台或終端機介面能呈現同一個工作階段。

| 結束代碼 | 意義 |
| --- | --- |
| `0` | 成功，或執行正在等待人員。 |
| `1` | 操作失敗：某個 skill、hook、建置或啟動。 |
| `2` | 用法無效。 |
| `3` | 被防護或政策拒絕，例如金錢閘門。 |
| `4` | 缺少先決條件：工具、輸入、叢集或 runtime 紀錄。 |
| `130` | 已中斷。 |

完整的命令參考、`.langdrift/` 的結構與復原規則，請見英文的 [`ex` README](/packages/ex/README.md)。

### 範例輸入檔

<details>
<summary><code>workflow.json</code>：各階段、它們的 skills，以及人員參與的位置</summary>

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

階段也可以列出 `hooks`：它會先在專案內透過 `sh` 執行的、已提交的腳本。

</details>

<details>
<summary><code>driver.json</code>：執行各階段的程式開發代理</summary>

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

執行檔從專案根目錄執行，絕不經過 shell。`usage` 指出輸出在哪裡回報 token；沒有回報的呼叫會讓花費維持未知，而 `"unmeasured": "pause"` 會讓 FG-04 暫停下一次付費呼叫。

</details>

<details>
<summary><code>prices.json</code>：每次估算所依據的有日期牌價</summary>

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

這些數字僅為示意，用來演示格式，並不是經過驗證的價目表。

</details>

<details>
<summary><code>forecasts.json</code>：每種推理強度預期使用的 token</summary>

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

`planning` 涵蓋規劃器自己的呼叫。當工單先前的執行都沒有量測過 token 時，就會使用這些預測。

</details>

### 探索網站與主控台

使用 Node.js 24 與 pnpm 12，在儲存庫根目錄執行：

```sh
pnpm install
pnpm website
```

網站在 `http://localhost:3000` 執行。`pnpm console` 會在 `http://127.0.0.1:3001` 啟動主控台，並在旁邊的 `127.0.0.1:3401` 啟動 `extub`，這需要 Go。主控台顯示的每個數字都來自 stub。

## LangDrift 背後的工作

LangDrift 由 [Jonatas Sales](https://github.com/jonatassales) 打造，結合了 FinOps 實務、代理協助的交付與平台工程。它源自 Exphere 專案，以 pnpm 與 Turborepo 的 monorepo 建構：網站與主控台以 TypeScript 搭配 Next.js、React 與 Tailwind CSS 撰寫，`ex` 與 Microruntimes 契約則以 Go 撰寫，規劃器使用官方的 Anthropic Go SDK。

這個儲存庫以 spec 先行的方式建構：每一項變更都從 [`.agents/workflow/`](/.agents/workflow) 下的 spec 與工單開始，也就是 `ex` 所規劃的同一種工單。歡迎透過 [GitHub issues](https://github.com/jonatassales/lang-drift/issues) 提供回饋。

## 參考資料

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/)，Jonathan Morley，FinOps Foundation，2026 年 3 月 27 日。實務工作者描述了代理如何將支出與預算對帳、找出浪費，並在開發生命週期更早的階段落實成本政策；最後這個把 FinOps 帶進交付的主題，正是 LangDrift 投入的地方。
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html)，Birgitta Böckeler，martinfowler.com，2026 年 4 月 2 日。文章描述團隊圍繞程式開發代理所打造的外層 harness：在代理行動前引導它的指引，以及讓它自我修正的感測器。`ex harness` 為專案記錄這套 harness，而 FinOps 閘門則限制它能花多少錢。

## 致謝

感謝 [Matt Pocock](https://x.com/mattpocockuk)。本儲存庫所使用的代理 skills，從 grilling 到 to-spec、to-tickets、implement、tdd 與 code-review，都改編自以 MIT 授權釋出的 [mattpocock/skills](https://github.com/mattpocock/skills)，上方的範例工作流程也使用了其中幾個。
