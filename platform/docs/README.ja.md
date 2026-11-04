<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: エージェント型 FinOps プラットフォーム" width="820">
  </picture>
</p>

<p align="center">
  <strong>チームの予算に責任を持つ、エージェント支援のソフトウェア配信。</strong><br>
  実行前のコストガバナンスと、AI コーディングエージェントのための境界付きランタイム。人間とエージェントのために設計。
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <strong>日本語</strong> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@langdrift/ex"><img alt="npm @langdrift/ex" src="https://img.shields.io/npm/v/@langdrift/ex?logo=npm&amp;logoColor=white&amp;label=%40langdrift%2Fex&amp;color=CB3837&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/ja/">Webサイト</a> ·
  <a href="#主な機能">主な機能</a> ·
  <a href="#使い方">使い方</a> ·
  <a href="#設定ファイル">設定ファイル</a> ·
  <a href="#参考文献">参考文献</a> ·
  <a href="#謝辞">謝辞</a> ·
  <a href="#ライセンス">ライセンス</a>
</p>

## 主な機能

- **実行前の予算事前予約:** すべてのモデル呼び出しやクラウドインフラ利用は日付付き定価で試算され、実行開始前にチームの予算枠から予約されます。試算が上限を超える場合、費用が発生する前に実行が一時停止します。
- **境界付きマイクロランタイム (Microruntimes):** 既存のソフトウェアに合わせて設計された隔離環境 (Docker、kind、Terraform、Kubernetes) をローカルで計画・実行。独自のベンダーSDKをアプリコードに埋め込む必要はありません。
- **好みのコーディングエージェントを活用:** 構造化されたワークフローステージ (計画、実装、レビュー) を通じて、Claude Code などのエージェントドライバーを活用し、決定論的な人間の確認ポイントを設けることができます。
- **決定論的な監査証跡:** トークンを含まないエビデンス、検証可能な元帳、次回以降の計画に活かせるエピソード記憶を記録します。

## 使い方

### 前提条件

- **macOS** または **Linux**
- **Node.js 24+** (`npx` 経由で `@langdrift/ex` を実行)
- **Docker**、**kind**、**kubectl**、**Terraform** (境界付きローカルランタイムの実行用。いつでも `npx @langdrift/ex doctor` で確認可能)
- 環境変数 `ANTHROPIC_API_KEY` の設定 (モデルプランナーおよび評価用)

### 1. プロジェクトの初期化

リポジトリのルートディレクトリで初期化コマンドを実行します:

```sh
npx @langdrift/ex init
```

これにより、`.langdrift/config.json` にプロジェクト設定が作成され、`.langdrift/microruntimes/AGENTS.md` にローカル実行ルールが配置されます。既存のファイルは変更されません。`/.langdrift/` は `ex` が書き込むローカル状態なので、`.gitignore` に追加してください。

### 2. チケットから実行計画を作成

`ex` にチケットや仕様書を指定します。リンクされたドキュメントを読み込み、予算枠内でCPUやメモリのサイズを決定し、最適なモデルと作業強度を提案して境界付き実行計画を出力します:

```sh
npx @langdrift/ex agent plan <path-to-ticket.md> --component name=console,path=apps/console
```

> **ヒント:** `--no-judge` を付けると LLM 呼び出しを行わずにヒューリスティックのみで計画を作成できるため、費用はかかりません。

<p align="center">
  <img src="/platform/docs/assets/images/ex-agent-plan.webp" alt="Terminal running ex agent plan on ticket 0070. It reads seven linked documents, applies the team's 2000 millicore, 3072 MiB envelope and $25.00 per-run budget, suggests claude-sonnet-5-5 at medium effort, sizes apps/console at 500m and 1024 MiB, and ends with a proposed plan, revision 1." width="820">
</p>

計画は `.langdrift/microruntimes/runtimes/<runtime-id>/runtime.md` に出力されます。`langdrift` ブロックに実行可能な設定が含まれます:

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

コンポーネントはファイルを指します: 名前と、Dockerfile を含むフォルダです。`ex` はプロジェクトのルートからその Dockerfile で `docker build` を実行し、イメージに `langdrift-<runtime-id>-<name>:r<revision>` というタグを付け、最初の `EXPOSE` からポートを、Dockerfile 自身の `HEALTHCHECK` から準備完了の判定を取得します。`--component` には `dockerfile=`、`port=`、`readiness=tcp` または `readiness=http:/path`、シークレットを含まないコミット済みの `env-file=`、起動時に `ex` が自身の環境から読み取る変数の `pass-env=A;B` を追加できます。`ex` はこれらをコンポーネントごとに 1 つの Kubernetes Secret に書き込みます。

### 3. 計画を実行

ランタイムを実行することで計画が承認され、配信サイクルが開始されます:

```sh
npx @langdrift/ex agent run <runtime-id>
```

実行中、`ex` は配信ライフサイクル全体を制御します:

1. **ゲート検証:** 定価が実行内容をカバーしているか確認し、チケットの割り当てと期間予算枠に対して試算コストを予約します。
2. **ワークフローステージ:** 設定されたコーディングエージェントを通じて各ステージを進めます。人間の承認が必要な場合、ターミナルまたは `npx @langdrift/ex agent confirm <runtime-id> <stage>` で確認されるまで停止します。
3. **ローカルランタイム:** 各コンポーネントを Dockerfile からビルドし、コンポーネント名を付けてローカル kind クラスタ上で起動し、準備完了プローブを待ちます。
4. **配信の監視:** Git コミット、ブランチ、PR チェックを読み取り専用で監視します。
5. **精算:** ローカルランタイムを停止し、未使用の予約額を予算に戻し、監査元帳に実行結果を記録します。

### 4. 状態の確認と診断

実行中のセッションを確認し、ツールチェーンを検証します:

```sh
# 現在の実行ステータスとラップ状態を表示
npx @langdrift/ex agent status <runtime-id>

# 過去の実行からエピソード記憶インデックスを再構築
npx @langdrift/ex agent reindex

# 前提ツールと環境の健全性をチェック
npx @langdrift/ex doctor
```

`status` に `--json` を付けると、ダッシュボードやコンソールとの連携に適した NDJSON イベントを出力できます。

## 設定ファイル

リポジトリ内の JSON ファイルで、ワークフローステージ、エージェントドライバー、コストパラメータをカスタマイズできます:

<details>
<summary><code>workflow.json</code>: 配信ステージと人間の確認ポイント</summary>

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
<summary><code>driver.json</code>: コーディングエージェントドライバーの定義</summary>

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
<summary><code>prices.json</code>: トークン価格表</summary>

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
<summary><code>forecasts.json</code>: 作業強度ごとの予測トークン数</summary>

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

## 参考文献

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/)、Jonathan Morley、FinOps Foundation、2026 年 3 月 27 日。実務者たちが、支出を予算と照合し、無駄を見つけ、開発ライフサイクルのより早い段階でコストポリシーを適用するエージェントについて語っています。最後のテーマ、つまり FinOps をデリバリーの中へ持ち込むことこそ、LangDrift が取り組む領域です。
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html)、Birgitta Böckeler、martinfowler.com、2026 年 4 月 2 日。チームがコーディングエージェントの周りに築く外側のハーネス、つまりエージェントが動く前に方向づけるガイドと、自己修正を可能にするセンサーについて述べています。`ex harness` はそのハーネスをプロジェクトごとに記録し、FinOps のゲートはそれが使える金額を制限します。

## 謝辞

[Matt Pocock](https://x.com/mattpocockuk) に感謝します。このリポジトリで使っているエージェントのスキル、grilling から to-spec、to-tickets、implement、tdd、code-review までは、MIT ライセンスで公開されている [mattpocock/skills](https://github.com/mattpocock/skills) を改変したもので、上のワークフロー例でもそのいくつかを使っています。

## ライセンス

LangDrift は [Apache 2.0 ライセンス](/LICENSE) のもとで公開されているオープンソースソフトウェアです。
