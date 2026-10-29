<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>エージェントによるデリバリーを、チームの予算に対して説明できるものに。</strong><br>
  人間と AI エージェントのために作られています。
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <strong>日本語</strong> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="プルリクエストのチェック" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="https://www.npmjs.com/package/@langdrift/ex"><img alt="npm @langdrift/ex" src="https://img.shields.io/npm/v/@langdrift/ex?logo=npm&amp;logoColor=white&amp;label=%40langdrift%2Fex&amp;color=CB3837&amp;labelColor=18181B" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 以降" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/ja/">ウェブサイト</a> ·
  <a href="#プロダクト">プロダクト</a> ·
  <a href="#予算枠">予算枠</a> ·
  <a href="#ex-をはじめる"><code>ex</code> をはじめる</a> ·
  <a href="#プロジェクトの状況">プロジェクトの状況</a>
</p>

## プロダクト

**LangDrift FinOps** は、エージェントによるデリバリーをチームの予算に対して説明できるものにします。チケットが求めること、エージェントが行う作業、その作業が残す証拠、そしてそれにかかった費用を、チームが確認できるひとつのつながりとして結びます。

**Microruntimes** はそのエンジンです。CLI の `ex` は、チケットとそこからリンクされたドキュメントを読み、既存のソフトウェアの周りに範囲の限られたローカル runtime を計画し、あなた自身のコーディングエージェントを通してデリバリーを進めます。有料呼び出しはすべて、実行前にチームの予算枠に対して予約されます。

タイムアウトしたモデル呼び出し、計画より高くついた再試行、チェック結果が返ってこないプルリクエスト。デリバリーはごく普通のことで計画から外れます。LangDrift が目指すのは、そのどれもが誰も承認していない支出や、誰も説明できない数字にならないことです。

### チームが得るもの

| 機能 | チームにとっての価値 |
| --- | --- |
| **チケットから計画** | `ex agent plan` はチケットとそこからリンクされたドキュメントだけを読み、比較可能な過去の実行を取り出し、チームの予算枠に収まるよう見積もった runtime の計画を書き出します。 |
| **支出の前に予約** | プランナーのものもコーディングエージェントのものも、有料呼び出しはすべて日付付きの定価で見積もられ、実行ごとの上限、チケットへの配分、期間の残額に対して予約されます。 |
| **あなたのハーネスでデリバリー** | ワークフローのステージとスキルは、あなたが設定したコーディングエージェントで実行されます。チームが決めた場面では人が確認し、マージは常にオーナーのものです。 |
| **既存ソフトウェアの周りに runtime** | コンポーネントは kind、Terraform、Kubernetes による範囲の限られたローカルのライフサイクルで起動し、あなたのコードに LangDrift への依存は入りません。 |
| **証拠のある記憶** | 実行のたびに、トークンを含まない証拠、台帳、そして次の計画が取り出す記憶が残ります。 |

LangDrift は実行の一周ごとに独自の言葉を使います。microruntime が車、ラン（run）がレース、そして**ラップ**（lap）がその一周です。ゲートで一時停止した試行の場合も、完走した実行の場合もあります。

## 予算枠

承認された予算は、開発者やチケットではなくチームのものです。チケットはその中で配分を受け、実行にはそれぞれ上限があります。

ラップを始める前に、`ex` は実行を定価で見積もり、チームの予約マージンを上乗せします。実行ごとの上限、チケットへの配分、期間の残額に収まる予約ならラップが始まり、収まらなければ何も支払わずに一時停止します。ラップの進行中は、有料呼び出しのたびにその前でゲートが確認します。終了すると、使われなかった予約は予算に戻ります。

| ゲート | 確認すること | 通らなかったとき |
| --- | --- | --- |
| **FG-01** | 日付付きの価格情報が有料の作業をカバーしている。 | ラップをブロックします。 |
| **FG-02** | チケットに配分がある。 | ラップをブロックします。 |
| **FG-03** | 見積もりにマージンを加えた額が、実行ごとの上限、配分、期間の残額に収まる。 | 何も支払う前に一時停止します。 |
| **FG-04** | ラップで観測された支出が予約額に達していない。 | 有料呼び出しを一時停止します。 |
| **FG-05** | 以前のラップの超過支出をオーナーが確認済みである。 | `--acknowledge-overspend` によるオーナーの確認が必要です。 |
| **FG-06** | プロバイダーが実行分を請求済みである。 | 数字を暫定のままにします。 |

**不明はゼロではありません。** 使用量を報告しない呼び出し、誰も確認できなかったクリーンアップ、まだ届いていない請求は、理由とともに不明のまま残り、照合されるまで予約も保持されます。

## あなたのスタック、あなたの判断

LangDrift が持つのは境界と仕組みです。アプリケーションとそのドメインの選択はあなたのチームのものです。

- **小さな足跡。** プロジェクトの中で `ex` が書き込むのは `.langdrift/` だけです。チケットがリンクする成果物を読み、プロジェクトの外は決して読まず、シークレットとシンボリックリンクは拒否します。
- **デリバリーは観測するだけ。** `ex` は push も、プルリクエストの作成も、マージもしません。デリバリーはあなたのハーネスのものであり、`ex` は `git` と `gh` を通してブランチ、プルリクエスト、そのチェックを観測します。
- **提案は承認ではありません。** プランナーはモデルと推論強度を提案しますが、計画を承認するのは runtime を実行することです。モデルの提案がスコープ、資金、アクセスを広げることはありません。
- **認証情報は環境の中に。** プランナーは呼び出しの時点で `ANTHROPIC_API_KEY` を読み、コーディングエージェントには最小限の環境と、あなたが指定した変数だけが渡されます。`ex` が認証情報をコマンド引数で受け取ることはありません。

実行をブロックまたは許可しうる入力はすべて、誰が、いつ、何を通して設定したかとともに記録されます。`.langdrift/` の JSON ファイルはその記録を読むためのミラーであり、ミラーを編集してもゲートが読む内容は何も変わりません。

## プロジェクトの状況

**LangDrift は開発中です。** ウェブサイトは [langdrift.md](https://langdrift.md/ja/) で五つの言語で公開されています。コンソールは `extub` を読み込みます。これは LangDrift Studio というひとつのプロジェクトを用意するスタブ API で、その FinOps の数値はイメージであると明示されています。CLI は npm に `@langdrift/ex`（`npx @langdrift/ex`）として公開されています。

`ex agent plan` と `ex agent run` はプロジェクト内でローカルに動作し、お金のゲートは `ex` の中でプロジェクト自身の台帳に対して実行されます。サーバーで共有するチームの台帳はまだ作られていません。接続エージェント（`ex agent connect`、`project`、`supervise`、`migrate`）には、このリポジトリにまだないサーバーが必要です。テストスイートは、スクリプト化したモデル、コーディングエージェント、ライフサイクルで CLI を動かすもので、実際のクラスターの起動、実際の支出、本番環境での準備が整っていることを示すものではありません。ウェブサイトでの登録は一時停止中です。

| パス | 内容 |
| --- | --- |
| [`apps/website`](/apps/website/README.md) | langdrift.md のプラットフォームサイト：ホーム、FinOps、Microruntimes、料金、早期アクセス。 |
| [`apps/console`](/apps/console/README.md) | チームが FinOps と microruntimes を管理するコンソール。 |
| [`apps/extub`](/apps/extub/README.md) | コンソールが読み込むローカルの Go スタブ API。 |
| [`packages/ex`](/packages/ex/README.md) | Go で書かれた `ex` CLI。 |
| [`packages/contracts`](/packages/contracts/README.md) | Microruntimes の契約：`runtime.md` のブロックとエージェントのプロトコル。 |
| [`packages/tokens`](/packages/tokens) · [`packages/web`](/packages/web) | デザイントークンと共有 UI。 |

## `ex` をはじめる

macOS または Linux と、`npx` で `@langdrift/ex` を実行するための Node.js が必要です。計画するだけならこれで十分です。プランナーのモデルによる判定は `ANTHROPIC_API_KEY` を読みます。ラップを実行するには、さらに Docker、kind、kubectl、Terraform（`npx @langdrift/ex doctor` が確認し、何もインストールしません）、Claude Code などあなたが設定するコーディングエージェント、そしてプルリクエストを観測するための `gh` が必要です。

### 1. プロジェクトを初期化する

プロジェクトのリポジトリのルートで次を実行します。

```sh
npx @langdrift/ex init
```

`npx @langdrift/ex init` は `.langdrift/config.json` にプロジェクトの識別情報を、`.langdrift/microruntimes/AGENTS.md` に実行ルールを作成し、既存のファイルには手を加えません。`/.langdrift/` を `.gitignore` に追加してください。そこは `ex` が書き込むマシンローカルの状態であり、チームが書くものは何も置かれません。

### 2. チケットから計画する

runtime がビルドして実行するコンポーネントを指定し、`ex` にチケットを読ませます。`--no-judge` はモデルを呼び出さずに計画するので、何も支払いは発生しません。これを付けずに `ANTHROPIC_API_KEY` を設定すると、プランナーはラップと同じゲートのもとでモデルを通して情報源を判定し、その費用は台帳の計画の行に記録されます。

CLI の出力は英語です。

<p align="center">
  <img src="/platform/docs/assets/images/ex-agent-plan.webp" alt="チケット 0070 に対して ex agent plan を実行しているターミナル。リンクされた 7 つのドキュメントを読み、チームの 2000 ミリコア・3072 MiB の予算枠と実行ごとの 25.00 ドルの予算を適用し、claude-sonnet-5-5 を中程度の推論強度で提案し、apps/console を 500m と 1024 MiB で見積もって、提案済みの計画（リビジョン 1）で終わります。" width="820">
</p>

計画は `.langdrift/microruntimes/runtimes/console-microruntime/runtime.md` に書き出されます。`langdrift` ブロックが runtime の実行する内容で、二つ目の `langdrift-plan` ブロックに分析が入ります。ファイル内の文章は文脈として残るだけで、実行されることはありません。

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

この抜粋では、ブロックのバージョン、名前、そして `ex` が管理する runtime とプロジェクトの識別情報を省いています。runtime が停止している間に同じチケットを再び計画すると、リビジョンが上がります。

### 3. 計画を実行する

runtime を実行することが、その計画の承認になります。ラップはイメージをビルドし、プロジェクトのローカルクラスターを起動し、コーディングエージェントを通して有料呼び出しを行うので、先に `npx @langdrift/ex doctor` を実行してください。

```sh
npx @langdrift/ex agent run console-microruntime
```

1. **ゲート。** 実行は意図、定価での見積もり、実行ごとの予算を表示してから予約します。すべての上限に収まる予約ならラップが始まり、超えていれば何も支払わずに一時停止します。
2. **ステージ。** 各ステージはコーディングエージェントを通して順に実行され、一行の要約を表示します。人が関わるステージはオーナーを待ちます。ターミナルで答えるか、`--confirm` を使うか、あとで `npx @langdrift/ex agent confirm console-microruntime planning` を実行します。`owner-only` のステージは決して実行されません。
3. **runtime。** コンポーネントはブロックに従ってプロジェクトの kind クラスター上で起動し、準備が整うと実行がそれを知らせます。
4. **デリバリー。** 実行はブランチ、プルリクエスト、そのチェックを観測します。`--checks-wait` で追いかける時間を制限できます。
5. **精算。** runtime が停止し、使われなかった予約が予算に戻り、記憶 `run.md` が書き込まれ、エピソードのインデックスが再構築されます。実行は、レビューやマージの準備ができているかどうかと、次に何をすべきかを伝えます。

すべての変更は、ラップのファイル、`run.json`、台帳に記録されます。

### 4. 実行を追う

```sh
npx @langdrift/ex agent status console-microruntime
npx @langdrift/ex agent reindex
npx @langdrift/ex doctor
```

`status` は実行の記録と進行中のラップを、`reindex` はすべての runtime の記憶からエピソードのインデックスを再構築し、`doctor` はネイティブツールを確認します。`--json` を付けると、セッションの各行が NDJSON のイベント `{at, session, kind, id, text, card, field, data}` になるので、コンソールやターミナル UI で同じセッションを表示できます。

| 終了コード | 意味 |
| --- | --- |
| `0` | 成功、または人を待っている実行。 |
| `1` | 操作が失敗した：スキル、フック、ビルド、起動のいずれか。 |
| `2` | 使い方が正しくない。 |
| `3` | お金のゲートなど、ガードやポリシーによって拒否された。 |
| `4` | 前提条件が足りない：ツール、入力、クラスター、runtime の記録のいずれか。 |
| `130` | 中断された。 |

コマンドの完全なリファレンス、`.langdrift/` の構成、復旧のルールは、英語の [`ex` README](/packages/ex/README.md) にあります。

### 入力ファイルの例

<details>
<summary><code>workflow.json</code>：ステージ、そのスキル、人が関わる場面</summary>

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

ステージには `hooks` も指定できます。プロジェクトの中で、`sh` を通して最初に実行される、コミット済みのスクリプトです。

</details>

<details>
<summary><code>driver.json</code>：ステージを実行するコーディングエージェント</summary>

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

実行ファイルはプロジェクトのルートから、シェルを介さずに実行されます。`usage` は出力のどこにトークン数があるかを示します。報告のない呼び出しは支出が不明のままになり、`"unmeasured": "pause"` によって FG-04 が次の有料呼び出しを一時停止します。

</details>

<details>
<summary><code>prices.json</code>：すべての見積もりの根拠となる日付付きの定価</summary>

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

これらの数字は形式を示すためのイメージであり、検証済みの価格表ではありません。

</details>

<details>
<summary><code>forecasts.json</code>：推論強度ごとに見込まれるトークン</summary>

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

`planning` はプランナー自身の呼び出しを対象とします。チケットの過去の実行でトークンが計測されていないときに、この予測が使われます。

</details>

## 参考文献

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/)、Jonathan Morley、FinOps Foundation、2026 年 3 月 27 日。実務者たちが、支出を予算と照合し、無駄を見つけ、開発ライフサイクルのより早い段階でコストポリシーを適用するエージェントについて語っています。最後のテーマ、つまり FinOps をデリバリーの中へ持ち込むことこそ、LangDrift が取り組む領域です。
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html)、Birgitta Böckeler、martinfowler.com、2026 年 4 月 2 日。チームがコーディングエージェントの周りに築く外側のハーネス、つまりエージェントが動く前に方向づけるガイドと、自己修正を可能にするセンサーについて述べています。`ex harness` はそのハーネスをプロジェクトごとに記録し、FinOps のゲートはそれが使える金額を制限します。

## 謝辞

[Matt Pocock](https://x.com/mattpocockuk) に感謝します。このリポジトリで使っているエージェントのスキル、grilling から to-spec、to-tickets、implement、tdd、code-review までは、MIT ライセンスで公開されている [mattpocock/skills](https://github.com/mattpocock/skills) を改変したもので、上のワークフロー例でもそのいくつかを使っています。
