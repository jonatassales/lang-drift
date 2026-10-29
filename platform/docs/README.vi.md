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
  <a href="#giấy-phép">Giấy phép</a>
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

## Giấy phép

*Nội dung giấy phép được giữ nguyên bằng tiếng Anh.*

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
