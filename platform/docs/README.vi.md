<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform. Một hạn mức cho cả đội. Mọi lần gọi trả phí đều được giữ chỗ trước khi chạy." width="820">
  </picture>
</p>

<p align="center">
  <strong>Bàn giao có tác tử hỗ trợ, chịu trách nhiệm trước ngân sách của đội.</strong>
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <strong>Tiếng Việt</strong>
</p>

<p align="center">
  <a href="#trạng-thái-dự-án"><img alt="Trạng thái: bản ứng viên đang phát triển" src="https://img.shields.io/badge/status-development%20candidate-FBBF24?labelColor=18181B" height="20"></a>
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Kiểm tra pull request" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 trở lên" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="https://langdrift.md/vi/"><img alt="Ngôn ngữ: en, pt-BR, zh-Hant, ja, vi" src="https://img.shields.io/badge/locales-en%20%C2%B7%20pt--BR%20%C2%B7%20zh--Hant%20%C2%B7%20ja%20%C2%B7%20vi-6D0AD6?labelColor=18181B" height="20"></a>
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

**LangDrift đang được phát triển.** Trang web đã hoạt động tại [langdrift.md](https://langdrift.md/vi/) với năm ngôn ngữ. Bảng điều khiển đọc `extub`, một API giả lập khởi tạo sẵn một dự án là LangDrift Studio và ghi rõ các số liệu FinOps của nó là minh họa. `ex` là bản ứng viên phát triển cục bộ mà bạn tự build từ mã nguồn; chưa có bản nhị phân nào được phát hành.

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

Bạn cần macOS hoặc Linux và Go 1.26 trở lên; như vậy là đủ để lập kế hoạch. Phần đánh giá bằng mô hình của bộ lập kế hoạch đọc `ANTHROPIC_API_KEY`. Để chạy một vòng, bạn còn cần Docker, kind, kubectl và Terraform (`ex doctor` kiểm tra chúng và không cài đặt gì), tác tử lập trình mà bạn cấu hình, chẳng hạn Claude Code, và `gh` để quan sát pull request.

### 1. Build CLI

```sh
git clone https://github.com/jonatassales/lang-drift.git
cd lang-drift/packages/ex
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex version
```

`version` in ra `ex 0.0.0-dev`. Hệ thống của bạn có thể đã có một `ex` khác, trình soạn thảo dòng của Vi, nên hãy gọi bản này bằng đường dẫn. Các bước dưới đây viết `ex` thay cho `packages/ex/dist/ex`; chỉ trong shell này, `alias ex="$PWD/dist/ex"` sẽ làm điều đó.

### 2. Khởi tạo dự án

Để thử trên một dự án có sẵn, hãy sao chép dự án mà bộ kiểm thử lập kế hoạch: một ticket, bản spec, nghiên cứu, ADR và guardrail mà nó liên kết, cùng hai manifest ứng dụng.

```sh
cp -R src/testdata/langdrift/project /tmp/ex-demo
cd /tmp/ex-demo
ex init
```

`ex init` tạo định danh dự án trong `.langdrift/config.json` và các quy tắc thực thi trong `.langdrift/microruntimes/AGENTS.md`, đồng thời giữ nguyên các tệp đã có. Trong kho mã của bạn, hãy thêm `/.langdrift/` vào `.gitignore`: đó là trạng thái cục bộ của máy do `ex` ghi, và không có gì do đội bạn soạn nằm trong đó.

### 3. Ghi lại harness

Lưu các [tệp đầu vào mẫu](#tệp-đầu-vào-mẫu) vào dự án, rồi ghi lại harness bàn giao: hạn mức phần cứng, các giai đoạn của quy trình, mô hình và mức nỗ lực mặc định, và tác tử lập trình chạy các giai đoạn.

```sh
ex harness envelope set --per-execution-cpu-millis 4000 --per-execution-memory-mib 6144 \
  --default-cpu-millis 2000 --default-memory-mib 3072 --not-runtime-field gpu
ex harness workflow set --file workflow.json
ex harness model set claude-sonnet-5-5 --effort medium
ex harness driver set --file driver.json
```

### 4. Ghi lại quy tắc chi tiêu

Ngân sách này cấp cho đội 500,00 USD cho tháng 10 năm 2026, giới hạn mỗi lượt chạy ở 25,00 USD, giữ thêm biên độ 20% trên mỗi ước tính và phân bổ 25,00 USD cho ticket 0070. Số tiền tính bằng cent.

```sh
ex finops budget set --period 2026-10 --period-cents 50000 --per-run-cents 2500 \
  --reserve-margin-percent 20 --allocate 0070-console-microruntime=2500
ex finops gates set --default
ex finops prices set --file prices.json --slug list
ex finops forecasts set --file forecasts.json
```

Mỗi lệnh in ra bản ghi mà nó đã viết, và mọi lệnh đầu vào đều nhận `--by` để ghi tên người đặt.

### 5. Lập kế hoạch từ một ticket

Khai báo các thành phần mà runtime cần build và chạy, rồi để `ex` đọc ticket. `--no-judge` lập kế hoạch mà không gọi mô hình, nên không phát sinh chi phí; nếu bỏ cờ này và đặt `ANTHROPIC_API_KEY`, bộ lập kế hoạch sẽ đánh giá các nguồn qua một mô hình dưới cùng các cổng kiểm soát như một vòng, và chi phí được ghi vào dòng lập kế hoạch của sổ cái.

```sh
ex agent plan .agents/workflow/tickets/0070-console-microruntime.ticket.md --no-judge \
  --component "path=apps/console,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=docker,args=build -t langdrift-console:local ."
```

Đầu ra của CLI bằng tiếng Anh:

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

### 6. Chạy kế hoạch

Chạy runtime chính là phê duyệt kế hoạch của nó. Một vòng sẽ build image, khởi động cụm cục bộ của dự án và thực hiện các lần gọi trả phí qua tác tử lập trình của bạn, vì vậy hãy chạy `ex doctor` trước.

```sh
ex agent run console-microruntime
```

1. **Cổng kiểm soát.** Lượt chạy in ra ý định, ước tính theo giá niêm yết và ngân sách mỗi lượt chạy, rồi giữ chỗ. Một khoản giữ chỗ nằm trong mọi giới hạn sẽ mở vòng; khoản vượt giới hạn sẽ tạm dừng mà chưa trả gì.
2. **Các giai đoạn.** Mỗi giai đoạn chạy lần lượt qua tác tử lập trình và in ra một dòng tóm tắt. Giai đoạn có con người tham gia sẽ chờ chủ sở hữu: trả lời ngay trên terminal hoặc bằng `--confirm`, hoặc sau đó bằng `ex agent confirm console-microruntime planning`. Giai đoạn `owner-only` không bao giờ chạy.
3. **Runtime.** Các thành phần khởi chạy từ khối trên cụm kind của dự án, và lượt chạy báo khi chúng sẵn sàng.
4. **Bàn giao.** Lượt chạy quan sát nhánh, pull request và các bước kiểm tra của nó; `--checks-wait` giới hạn thời gian theo dõi.
5. **Quyết toán.** Runtime dừng lại, phần giữ chỗ chưa dùng trả về ngân sách, bộ nhớ `run.md` được ghi và chỉ mục episode được dựng lại. Lượt chạy cho biết nó đã sẵn sàng để review hay để merge, và bước tiếp theo là gì.

Mọi thay đổi đều được ghi vào tệp của vòng, `run.json` và sổ cái.

### 7. Theo dõi lượt chạy

```sh
ex agent status console-microruntime
ex agent reindex
ex doctor
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

### Khám phá trang web và bảng điều khiển

Với Node.js 24 và pnpm 12, từ thư mục gốc của kho mã:

```sh
pnpm install
pnpm website
```

Trang web chạy tại `http://localhost:3000`. `pnpm console` khởi động bảng điều khiển tại `http://127.0.0.1:3001` cùng `extub` bên cạnh tại `127.0.0.1:3401`, vốn cần Go. Mọi con số mà bảng điều khiển hiển thị đều đến từ bản giả lập.

## Công việc đằng sau LangDrift

Được tạo bởi [Jonatas Sales](https://github.com/jonatassales), LangDrift kết hợp thực hành FinOps, bàn giao có tác tử hỗ trợ và kỹ thuật nền tảng. Nó ra đời từ dự án Exphere và được xây dựng như một monorepo với pnpm và Turborepo: trang web và bảng điều khiển viết bằng TypeScript với Next.js, React và Tailwind CSS, còn `ex` và các hợp đồng của Microruntimes viết bằng Go, với SDK Go chính thức của Anthropic cho bộ lập kế hoạch.

Kho mã được xây dựng theo hướng spec trước: mọi thay đổi bắt đầu từ một spec và các ticket của nó trong [`.agents/workflow/`](/.agents/workflow), cũng chính là loại ticket mà `ex` lập kế hoạch. Mọi góp ý xin gửi qua [GitHub issues](https://github.com/jonatassales/lang-drift/issues).

## Tài liệu tham khảo

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/), Jonathan Morley, FinOps Foundation, ngày 27 tháng 3 năm 2026. Những người thực hành mô tả các tác tử đối soát chi tiêu với ngân sách, tìm ra lãng phí và áp dụng chính sách chi phí sớm hơn trong vòng đời phát triển; chủ đề cuối cùng đó, đưa FinOps vào ngay trong quá trình bàn giao, chính là nơi LangDrift hoạt động.
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html), Birgitta Böckeler, martinfowler.com, ngày 2 tháng 4 năm 2026. Bài viết mô tả lớp harness bên ngoài mà một đội xây quanh tác tử lập trình: các chỉ dẫn định hướng tác tử trước khi nó hành động và các cảm biến giúp nó tự sửa lỗi. `ex harness` ghi lại harness đó cho từng dự án, còn các cổng kiểm soát FinOps giới hạn số tiền nó được chi.

## Lời cảm ơn

Xin cảm ơn [Matt Pocock](https://x.com/mattpocockuk). Các skill tác tử mà kho mã này sử dụng, từ grilling đến to-spec, to-tickets, implement, tdd và code-review, được chuyển thể từ [mattpocock/skills](https://github.com/mattpocock/skills) theo giấy phép MIT, và quy trình mẫu ở trên dùng một vài skill trong số đó.
