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
  <a href="/packages/ex/README.md"><img alt="Go 1.26 以降" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/ja/">ウェブサイト</a> ·
  <a href="#プロダクト">プロダクト</a> ·
  <a href="#ライセンス">ライセンス</a>
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

## ライセンス

*ライセンス本文は英語のままです。*

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
