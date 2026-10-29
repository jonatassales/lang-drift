<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>Entrega assistida por agentes, prestando contas ao orçamento do seu time.</strong><br>
  Feito para pessoas e agentes de IA.
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <strong>Português (Brasil)</strong> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Verificações de pull request" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 ou mais recente" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/pt-br/">Site</a> ·
  <a href="#o-produto">O produto</a> ·
  <a href="#licença">Licença</a>
</p>

## O produto

O **LangDrift FinOps** faz a entrega assistida por agentes prestar contas ao orçamento do time. Ele conecta o que um ticket pede, o trabalho que um agente faz, as evidências que esse trabalho deixa e quanto ele custou, em uma única cadeia que o time pode inspecionar.

**Microruntimes** é o seu motor. A CLI dele, `ex`, lê um ticket e os documentos que ele referencia, planeja um runtime local e delimitado em volta do software que você já tem e conduz a entrega pelo seu próprio agente de código, com cada chamada paga reservada contra o envelope do time antes de rodar.

Uma chamada ao modelo que estoura o tempo, uma retentativa que custa mais que o planejado, um pull request cujas verificações nunca respondem: a entrega sai do plano de jeitos comuns. O objetivo do LangDrift é que nenhum deles vire um gasto que ninguém aprovou ou um número que ninguém consegue explicar.

### O que o time ganha

| Capacidade | Valor para o time |
| --- | --- |
| **Planos a partir de tickets** | `ex agent plan` lê o ticket e apenas os documentos que ele referencia, recupera execuções anteriores comparáveis e escreve um plano de runtime dimensionado dentro do envelope do time. |
| **Reserva antes do gasto** | Toda chamada paga, do planejador e do agente de código, é estimada a preços de tabela datados e reservada contra o teto por execução, a alocação do ticket e o que resta do período. |
| **Entrega pelo seu harness** | Os estágios e as skills do seu workflow rodam pelo agente de código que você configurar. As pessoas confirmam onde o time definiu que elas participam, e o merge continua sempre com o owner. |
| **Runtimes em volta do software existente** | Seus componentes sobem em um ciclo de vida local e delimitado com kind, Terraform e Kubernetes, sem nenhuma dependência do LangDrift no seu código. |
| **Memória com evidências** | Cada execução deixa evidências sem tokens, um livro-razão e uma memória que o próximo plano recupera. |

O LangDrift tem uma palavra própria para cada passagem de uma execução: o microruntime é o carro, a execução (run) é a corrida, e uma **volta** (lap) é uma passagem dela, seja uma tentativa pausada em um gate, seja a execução que terminou.

## Licença

*O texto da licença é mantido em inglês.*

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
