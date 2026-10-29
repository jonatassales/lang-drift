<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>Bàn giao có tác tử hỗ trợ, chịu trách nhiệm trước ngân sách của đội.</strong><br>
  Được xây dựng cho con người và tác tử AI.
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <strong>Tiếng Việt</strong>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Kiểm tra pull request" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 trở lên" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/vi/">Trang web</a> ·
  <a href="#sản-phẩm">Sản phẩm</a> ·
  <a href="#hạn-mức">Hạn mức</a> ·
  <a href="#bắt-đầu-với-ex">Bắt đầu với <code>ex</code></a> ·
  <a href="#trạng-thái-dự-án">Trạng thái dự án</a>
</p>

## Sản phẩm

**LangDrift FinOps** buộc việc bàn giao có tác tử hỗ trợ phải chịu trách nhiệm trước ngân sách của đội. Nó nối điều một ticket yêu cầu, công việc một tác tử thực hiện, bằng chứng mà công việc đó để lại và chi phí của nó thành một chuỗi duy nhất mà cả đội có thể kiểm tra.

**Microruntimes** là động cơ của nó. CLI của nó, `ex`, đọc một ticket cùng các tài liệu mà ticket liên kết, lập kế hoạch cho một runtime cục bộ có giới hạn quanh phần mềm bạn đã có, và đưa việc bàn giao đi qua chính tác tử lập trình của bạn, với mọi lần gọi trả phí đều được giữ chỗ trong hạn mức của đội trước khi chạy.

Một lần gọi mô hình bị quá thời gian, một lần thử lại tốn hơn dự kiến, một pull request mà các bước kiểm tra không bao giờ báo kết quả: việc bàn giao lệch khỏi kế hoạch theo những cách rất bình thường. Mục tiêu của LangDrift là không điều nào trong số đó biến thành khoản chi không ai phê duyệt hay con số không ai giải thích được.

### Đội của bạn nhận được gì

| Năng lực | Giá trị cho đội |
| --- | --- |
| **Kế hoạch từ ticket** | `ex agent plan` chỉ đọc ticket và các tài liệu mà nó liên kết, truy xuất các lượt chạy trước có thể so sánh và viết ra một kế hoạch runtime được định cỡ trong hạn mức của đội. |
| **Giữ chỗ trước khi chi** | Mọi lần gọi trả phí, của bộ lập kế hoạch lẫn của tác tử lập trình, đều được ước tính theo bảng giá niêm yết có ghi ngày và được giữ chỗ trong mức trần mỗi lượt chạy, phần phân bổ của ticket và phần còn lại của kỳ. |
| **Bàn giao qua harness của bạn** | Các giai đoạn và skill trong quy trình của bạn chạy qua tác tử lập trình mà bạn cấu hình. Con người xác nhận ở những chỗ đội đã quy định, và việc merge luôn thuộc về chủ sở hữu. |
| **Runtime quanh phần mềm sẵn có** | Các thành phần của bạn khởi chạy trên một vòng đời cục bộ có giới hạn với kind, Terraform và Kubernetes, mà mã của bạn không cần phụ thuộc gì vào LangDrift. |
| **Bộ nhớ có bằng chứng** | Mỗi lượt chạy để lại bằng chứng không chứa token, một sổ cái và một bộ nhớ mà kế hoạch tiếp theo sẽ truy xuất. |

LangDrift có từ riêng cho mỗi lượt đi qua của một lần chạy: microruntime là chiếc xe, một lượt chạy (run) là cuộc đua, và một **vòng** (lap) là một lần đi qua, có thể là một lần thử bị tạm dừng ở cổng kiểm soát hoặc là lần thực thi đã về đích.

## Hạn mức

Ngân sách được phê duyệt thuộc về đội, không thuộc về một lập trình viên hay một ticket. Các ticket nhận phần phân bổ bên trong nó, và mỗi lượt chạy có một mức trần.

Trước khi mở một vòng, `ex` ước tính lượt chạy theo giá niêm yết và cộng thêm biên độ giữ chỗ của đội. Một khoản giữ chỗ nằm trong mức trần mỗi lượt chạy, phần phân bổ của ticket và phần còn lại của kỳ sẽ mở vòng; khoản nào vượt quá sẽ tạm dừng vòng mà chưa trả gì. Trong khi vòng chạy, cổng kiểm soát được kiểm tra trước mỗi lần gọi trả phí. Khi vòng kết thúc, phần giữ chỗ chưa dùng được trả lại ngân sách.

| Cổng | Kiểm tra điều gì | Khi không đạt |
| --- | --- | --- |
| **FG-01** | Có một nguồn giá ghi ngày bao phủ phần việc trả phí. | Chặn vòng. |
| **FG-02** | Ticket có phần phân bổ. | Chặn vòng. |
| **FG-03** | Ước tính cộng biên độ nằm trong mức trần mỗi lượt chạy, phần phân bổ và phần còn lại của kỳ. | Tạm dừng trước khi trả bất cứ khoản nào. |
| **FG-04** | Chi tiêu quan sát được của vòng chưa chạm tới khoản giữ chỗ. | Tạm dừng các lần gọi trả phí. |
| **FG-05** | Chủ sở hữu đã xác nhận khoản vượt chi của một vòng trước. | Cần chủ sở hữu xác nhận bằng `--acknowledge-overspend`. |
| **FG-06** | Nhà cung cấp đã lập hóa đơn cho lượt chạy. | Giữ các con số ở trạng thái tạm tính. |

**Chưa rõ không bao giờ là không.** Một lần gọi không báo mức sử dụng, một lần dọn dẹp không ai xác nhận được hay một hóa đơn chưa đến đều giữ nguyên là chưa rõ, kèm lý do, và khoản giữ chỗ của nó vẫn được giữ cho đến khi đối soát xong.

## Công nghệ của bạn, quyết định của bạn

LangDrift nắm các ranh giới và cơ chế; đội của bạn nắm ứng dụng và các lựa chọn thuộc lĩnh vực của mình.

- **Dấu chân nhỏ.** Trong dự án của bạn, `ex` chỉ ghi vào `.langdrift/`. Nó đọc các tài liệu mà ticket liên kết, không bao giờ đọc bên ngoài dự án, và từ chối bí mật cùng liên kết tượng trưng.
- **Việc bàn giao của bạn, chỉ được quan sát.** `ex` không bao giờ push, mở pull request hay merge. Việc bàn giao thuộc về harness của bạn; `ex` quan sát nhánh, pull request và các bước kiểm tra của nó qua `git` và `gh`.
- **Đề xuất không phải là phê duyệt.** Bộ lập kế hoạch đề xuất một mô hình và một mức nỗ lực; chính việc chạy runtime mới phê duyệt kế hoạch. Đề xuất của mô hình không bao giờ mở rộng phạm vi, nguồn vốn hay quyền truy cập.
- **Thông tin xác thực nằm trong môi trường.** Bộ lập kế hoạch đọc `ANTHROPIC_API_KEY` ngay lúc gọi, còn tác tử lập trình chỉ nhận một môi trường tối thiểu cộng với các biến bạn liệt kê. `ex` không bao giờ nhận thông tin xác thực qua tham số dòng lệnh.

Mọi đầu vào có thể chặn hoặc cho phép một lượt chạy đều được ghi lại cùng người đặt, thời điểm và kênh đặt. Các tệp JSON trong `.langdrift/` là bản phản chiếu của những bản ghi đó để đọc; sửa một bản phản chiếu không thay đổi bất cứ điều gì mà cổng kiểm soát đọc.

## Trạng thái dự án

**LangDrift đang được phát triển.** Trang web đã hoạt động tại [langdrift.md](https://langdrift.md/vi/) với năm ngôn ngữ. Bảng điều khiển đọc `extub`, một API giả lập khởi tạo sẵn một dự án là LangDrift Studio và ghi rõ các số liệu FinOps của nó là minh họa. CLI được phát hành trên npm dưới tên `@langdrift/ex` (`npx @langdrift/ex`).

`ex agent plan` và `ex agent run` hoạt động cục bộ trong một dự án, và các cổng kiểm soát chi tiêu chạy bên trong `ex` trên sổ cái riêng của dự án; một sổ cái chung của đội qua máy chủ vẫn chưa được xây dựng. Tác tử kết nối (`ex agent connect`, `project`, `supervise` và `migrate`) cần một máy chủ mà kho mã này chưa có. Các bộ kiểm thử điều khiển CLI bằng mô hình, tác tử lập trình và vòng đời được lập kịch bản; chúng không chứng minh việc khởi chạy cụm thật, chi tiêu thật hay sự sẵn sàng cho môi trường production. Việc đăng ký trên trang web đang tạm dừng.

| Đường dẫn | Nội dung |
| --- | --- |
| [`apps/website`](/apps/website/README.md) | Trang web của nền tảng tại langdrift.md: trang chủ, FinOps, Microruntimes, bảng giá và truy cập sớm. |
| [`apps/console`](/apps/console/README.md) | Bảng điều khiển nơi đội quản lý FinOps và các microruntime của mình. |
| [`apps/extub`](/apps/extub/README.md) | API giả lập cục bộ viết bằng Go mà bảng điều khiển đọc. |
| [`packages/ex`](/packages/ex/README.md) | CLI `ex`, viết bằng Go. |
| [`packages/contracts`](/packages/contracts/README.md) | Các hợp đồng của Microruntimes: khối `runtime.md` và giao thức của tác tử. |
| [`packages/tokens`](/packages/tokens) · [`packages/web`](/packages/web) | Design token và giao diện dùng chung. |

## Bắt đầu với `ex`

Bạn cần macOS hoặc Linux và Node.js để chạy `@langdrift/ex` bằng `npx`. Như vậy là đủ để lập kế hoạch. Phần đánh giá bằng mô hình của bộ lập kế hoạch đọc `ANTHROPIC_API_KEY`. Để chạy một vòng, bạn còn cần Docker, kind, kubectl và Terraform (`npx @langdrift/ex doctor` kiểm tra chúng và không cài đặt gì), tác tử lập trình mà bạn cấu hình, chẳng hạn Claude Code, và `gh` để quan sát pull request.

### 1. Khởi tạo dự án

Từ thư mục gốc của kho mã dự án của bạn:

```sh
npx @langdrift/ex init
```

`npx @langdrift/ex init` tạo định danh dự án trong `.langdrift/config.json` và các quy tắc thực thi trong `.langdrift/microruntimes/AGENTS.md`, đồng thời giữ nguyên các tệp đã có. Hãy thêm `/.langdrift/` vào `.gitignore`: đó là trạng thái cục bộ của máy do `ex` ghi, và không có gì do đội bạn soạn nằm trong đó.

### 2. Lập kế hoạch từ một ticket

Khai báo các thành phần mà runtime cần build và chạy, rồi để `ex` đọc ticket. `--no-judge` lập kế hoạch mà không gọi mô hình, nên không phát sinh chi phí; nếu bỏ cờ này và đặt `ANTHROPIC_API_KEY`, bộ lập kế hoạch sẽ đánh giá các nguồn qua một mô hình dưới cùng các cổng kiểm soát như một vòng, và chi phí được ghi vào dòng lập kế hoạch của sổ cái.

Đầu ra của CLI bằng tiếng Anh:

<p align="center">
  <img src="/platform/docs/assets/images/ex-agent-plan.webp" alt="Terminal chạy ex agent plan cho ticket 0070. Nó đọc bảy tài liệu được liên kết, áp dụng hạn mức 2000 millicore và 3072 MiB của đội cùng ngân sách 25,00 USD mỗi lượt chạy, đề xuất claude-sonnet-5-5 với mức nỗ lực trung bình, định cỡ apps/console ở 500m và 1024 MiB, rồi kết thúc bằng một kế hoạch được đề xuất, bản sửa đổi 1." width="820">
</p>

Kế hoạch được ghi vào `.langdrift/microruntimes/runtimes/console-microruntime/runtime.md`. Khối `langdrift` là thứ runtime sẽ thực thi; khối thứ hai, `langdrift-plan`, chứa phần phân tích. Văn bản trong tệp chỉ được giữ làm ngữ cảnh và không bao giờ được thực thi.

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

Đoạn trích này lược bỏ phiên bản, tên và các định danh runtime và dự án của khối, những thứ do `ex` quản lý. Lập kế hoạch lại cùng ticket khi runtime đang dừng sẽ tăng số phiên bản sửa đổi.

### 3. Chạy kế hoạch

Chạy runtime chính là phê duyệt kế hoạch của nó. Một vòng sẽ build image, khởi động cụm cục bộ của dự án và thực hiện các lần gọi trả phí qua tác tử lập trình của bạn, vì vậy hãy chạy `npx @langdrift/ex doctor` trước.

```sh
npx @langdrift/ex agent run console-microruntime
```

1. **Cổng kiểm soát.** Lượt chạy in ra ý định, ước tính theo giá niêm yết và ngân sách mỗi lượt chạy, rồi giữ chỗ. Một khoản giữ chỗ nằm trong mọi giới hạn sẽ mở vòng; khoản vượt giới hạn sẽ tạm dừng mà chưa trả gì.
2. **Các giai đoạn.** Mỗi giai đoạn chạy lần lượt qua tác tử lập trình và in ra một dòng tóm tắt. Giai đoạn có con người tham gia sẽ chờ chủ sở hữu: trả lời ngay trên terminal hoặc bằng `--confirm`, hoặc sau đó bằng `npx @langdrift/ex agent confirm console-microruntime planning`. Giai đoạn `owner-only` không bao giờ chạy.
3. **Runtime.** Các thành phần khởi chạy từ khối trên cụm kind của dự án, và lượt chạy báo khi chúng sẵn sàng.
4. **Bàn giao.** Lượt chạy quan sát nhánh, pull request và các bước kiểm tra của nó; `--checks-wait` giới hạn thời gian theo dõi.
5. **Quyết toán.** Runtime dừng lại, phần giữ chỗ chưa dùng trả về ngân sách, bộ nhớ `run.md` được ghi và chỉ mục episode được dựng lại. Lượt chạy cho biết nó đã sẵn sàng để review hay để merge, và bước tiếp theo là gì.

Mọi thay đổi đều được ghi vào tệp của vòng, `run.json` và sổ cái.

### 4. Theo dõi lượt chạy

```sh
npx @langdrift/ex agent status console-microruntime
npx @langdrift/ex agent reindex
npx @langdrift/ex doctor
```

`status` hiển thị bản ghi lượt chạy và vòng đang mở, `reindex` dựng lại chỉ mục episode từ bộ nhớ của mọi runtime, và `doctor` kiểm tra các công cụ gốc. Với `--json`, mỗi dòng của phiên là một sự kiện NDJSON, `{at, session, kind, id, text, card, field, data}`, để bảng điều khiển hoặc giao diện terminal có thể hiển thị cùng một phiên.

| Mã thoát | Ý nghĩa |
| --- | --- |
| `0` | Thành công, hoặc một lượt chạy đang chờ con người. |
| `1` | Thao tác thất bại: một skill, một hook, một lần build hoặc một lần khởi chạy. |
| `2` | Cách dùng không hợp lệ. |
| `3` | Bị từ chối bởi một cơ chế bảo vệ hoặc chính sách, chẳng hạn cổng kiểm soát chi tiêu. |
| `4` | Thiếu điều kiện tiên quyết: một công cụ, một đầu vào, một cụm hoặc bản ghi runtime. |
| `130` | Bị ngắt. |

Tài liệu tham khảo đầy đủ về các lệnh, cấu trúc của `.langdrift/` và các quy tắc khôi phục nằm trong [README của `ex`](/packages/ex/README.md), bằng tiếng Anh.

### Tệp đầu vào mẫu

<details>
<summary><code>workflow.json</code>: các giai đoạn, skill của chúng và nơi con người tham gia</summary>

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

Một giai đoạn cũng có thể liệt kê `hooks`: các script đã commit mà nó chạy trước tiên, qua `sh`, bên trong dự án.

</details>

<details>
<summary><code>driver.json</code>: tác tử lập trình chạy các giai đoạn</summary>

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

Tệp thực thi chạy từ thư mục gốc của dự án, không bao giờ qua shell. `usage` chỉ ra chỗ đầu ra báo số token; lần gọi nào không báo sẽ để chi tiêu của nó ở trạng thái chưa rõ, và `"unmeasured": "pause"` khiến FG-04 tạm dừng lần gọi trả phí tiếp theo.

</details>

<details>
<summary><code>prices.json</code>: bảng giá niêm yết có ghi ngày đứng sau mọi ước tính</summary>

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

Các con số này chỉ mang tính minh họa để thể hiện định dạng; chúng không phải là bảng giá đã được kiểm chứng.

</details>

<details>
<summary><code>forecasts.json</code>: số token dự kiến cho mỗi mức nỗ lực</summary>

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

`planning` bao gồm các lần gọi của chính bộ lập kế hoạch. Các dự báo được dùng khi chưa có lượt chạy trước nào của ticket đo số token của nó.

</details>

## Tài liệu tham khảo

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/), Jonathan Morley, FinOps Foundation, ngày 27 tháng 3 năm 2026. Những người thực hành mô tả các tác tử đối soát chi tiêu với ngân sách, tìm ra lãng phí và áp dụng chính sách chi phí sớm hơn trong vòng đời phát triển; chủ đề cuối cùng đó, đưa FinOps vào ngay trong quá trình bàn giao, chính là nơi LangDrift hoạt động.
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html), Birgitta Böckeler, martinfowler.com, ngày 2 tháng 4 năm 2026. Bài viết mô tả lớp harness bên ngoài mà một đội xây quanh tác tử lập trình: các chỉ dẫn định hướng tác tử trước khi nó hành động và các cảm biến giúp nó tự sửa lỗi. `ex harness` ghi lại harness đó cho từng dự án, còn các cổng kiểm soát FinOps giới hạn số tiền nó được chi.

## Lời cảm ơn

Xin cảm ơn [Matt Pocock](https://x.com/mattpocockuk). Các skill tác tử mà kho mã này sử dụng, từ grilling đến to-spec, to-tickets, implement, tdd và code-review, được chuyển thể từ [mattpocock/skills](https://github.com/mattpocock/skills) theo giấy phép MIT, và quy trình mẫu ở trên dùng một vài skill trong số đó.
