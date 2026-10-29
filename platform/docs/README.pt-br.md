<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>Entrega assistida por agentes, prestando contas ao orçamento do seu time.</strong>
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
  <a href="#o-envelope">O envelope</a> ·
  <a href="#comece-com-ex">Comece com <code>ex</code></a> ·
  <a href="#status-do-projeto">Status do projeto</a>
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

## O envelope

O orçamento autorizado pertence ao time, não a um desenvolvedor nem a um ticket. Os tickets recebem alocações dentro dele, e cada execução tem um teto.

Antes de abrir uma volta, o `ex` estima a execução a preços de tabela e soma a margem de reserva do time. Uma reserva que cabe no teto por execução, na alocação do ticket e no saldo do período abre a volta; uma que não cabe a pausa sem nada pago. Enquanto a volta roda, um gate é verificado antes de cada chamada paga. Quando ela termina, a reserva não usada volta ao orçamento.

| Gate | O que verifica | Quando não passa |
| --- | --- | --- |
| **FG-01** | Uma fonte de preços datada cobre o trabalho pago. | Bloqueia a volta. |
| **FG-02** | O ticket tem uma alocação. | Bloqueia a volta. |
| **FG-03** | A estimativa mais a margem cabe no teto por execução, na alocação e no saldo do período. | Pausa antes de qualquer pagamento. |
| **FG-04** | O gasto observado da volta não atingiu a sua reserva. | Pausa as chamadas pagas. |
| **FG-05** | O owner reconheceu o gasto excedente de uma volta anterior. | Exige o owner, com `--acknowledge-overspend`. |
| **FG-06** | O provedor já faturou a execução. | Mantém os valores como provisórios. |

**Desconhecido nunca é zero.** Uma chamada que não informa uso, uma limpeza que ninguém conseguiu confirmar ou uma fatura que ainda não chegou continua desconhecida, com o motivo, e mantém a sua reserva retida até a reconciliação.

## Sua stack, suas decisões

O LangDrift é dono das fronteiras e dos mecanismos; o seu time é dono da aplicação e das escolhas do seu domínio.

- **Uma pegada pequena.** No seu projeto, o `ex` escreve apenas em `.langdrift/`. Ele lê os artefatos que um ticket referencia, nunca fora do projeto, e recusa segredos e links simbólicos.
- **Sua entrega, observada.** O `ex` nunca faz push, nunca abre pull request e nunca faz merge. A entrega pertence ao seu harness; o `ex` observa a branch, o pull request e as verificações por meio de `git` e `gh`.
- **Recomendação não é autorização.** O planejador sugere um modelo e um esforço; é rodar o runtime que autoriza o plano. A sugestão de um modelo nunca amplia escopo, verba ou acesso.
- **Credenciais ficam no ambiente.** O planejador lê `ANTHROPIC_API_KEY` no momento da chamada, e o agente de código recebe um ambiente mínimo mais apenas as variáveis que você listar. O `ex` nunca recebe uma credencial como argumento de comando.

Toda entrada que pode bloquear ou liberar uma execução é registrada com quem a definiu, quando e por qual meio. Os arquivos JSON em `.langdrift/` espelham esses registros para leitura; editar um espelho não muda nada do que um gate lê.

## Status do projeto

**O LangDrift está em desenvolvimento.** O site está no ar em [langdrift.md](https://langdrift.md/pt-br/), em cinco idiomas. O console lê o `extub`, uma API stub que semeia um projeto, o LangDrift Studio, e marca as suas leituras de FinOps como ilustrativas. Você compila o `ex` a partir do código-fonte; nenhum binário é publicado.

`ex agent plan` e `ex agent run` funcionam localmente dentro de um projeto, e os gates de dinheiro rodam dentro do `ex` contra o livro-razão do próprio projeto; um livro-razão do time compartilhado por um servidor ainda não foi construído. O agente conectado (`ex agent connect`, `project`, `supervise` e `migrate`) precisa de um servidor que este repositório ainda não tem. As suítes de teste conduzem a CLI com modelos, agentes de código e ciclos de vida roteirizados; elas não comprovam o início de um cluster real, gasto real nem prontidão para produção. O cadastro no site está pausado.

| Caminho | O que contém |
| --- | --- |
| [`apps/website`](/apps/website/README.md) | O site da plataforma em langdrift.md: início, FinOps, Microruntimes, preços e acesso antecipado. |
| [`apps/console`](/apps/console/README.md) | O console em que um time gerencia o FinOps e os seus microruntimes. |
| [`apps/extub`](/apps/extub/README.md) | A API stub local em Go que o console lê. |
| [`packages/ex`](/packages/ex/README.md) | A CLI `ex`, em Go. |
| [`packages/contracts`](/packages/contracts/README.md) | Os contratos do Microruntimes: o bloco do `runtime.md` e o protocolo do agente. |
| [`packages/tokens`](/packages/tokens) · [`packages/web`](/packages/web) | Os design tokens e a UI compartilhada. |

## Comece com `ex`

Você precisa de macOS ou Linux e de Go 1.26 ou mais recente; isso basta para planejar. O juiz do planejador lê `ANTHROPIC_API_KEY`. Rodar uma volta também exige Docker, kind, kubectl e Terraform (`ex doctor` verifica e não instala nada), o agente de código que você configurar, como o Claude Code, e `gh` para observar pull requests.

### 1. Compile a CLI

```sh
git clone https://github.com/jonatassales/lang-drift.git
cd lang-drift/packages/ex
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex version
```

`version` imprime `ex 0.0.0-dev`. O seu sistema pode já ter um `ex`, o editor de linha do Vi, então chame este pelo caminho. Os passos abaixo escrevem `ex` no lugar de `packages/ex/dist/ex`; apenas neste shell, `alias ex="$PWD/dist/ex"` faz isso.

### 2. Inicialize um projeto

Para experimentar em um projeto pronto, copie o que a suíte de testes planeja: um ticket, a spec, a pesquisa, o ADR e o guardrail que ele referencia, e dois manifestos de apps.

```sh
cp -R src/testdata/langdrift/project /tmp/ex-demo
cd /tmp/ex-demo
ex init
```

`ex init` cria a identidade do projeto em `.langdrift/config.json` e as regras de execução em `.langdrift/microruntimes/AGENTS.md`, e não altera arquivos existentes. No seu repositório, adicione `/.langdrift/` ao `.gitignore`: é estado local da máquina que o `ex` escreve, e nada que o seu time escreve fica ali.

### 3. Registre o harness

Salve os [arquivos de entrada de exemplo](#arquivos-de-entrada-de-exemplo) no projeto e registre o harness de entrega: o envelope de hardware, os estágios do workflow, o modelo e o esforço padrão, e o agente de código que executa os estágios.

```sh
ex harness envelope set --per-execution-cpu-millis 4000 --per-execution-memory-mib 6144 \
  --default-cpu-millis 2000 --default-memory-mib 3072 --not-runtime-field gpu
ex harness workflow set --file workflow.json
ex harness model set claude-sonnet-5-5 --effort medium
ex harness driver set --file driver.json
```

### 4. Registre as regras de dinheiro

Este orçamento dá ao time US$ 500,00 para outubro de 2026, limita cada execução a US$ 25,00, reserva uma margem de 20% sobre cada estimativa e aloca US$ 25,00 ao ticket 0070. Os valores são em centavos.

```sh
ex finops budget set --period 2026-10 --period-cents 50000 --per-run-cents 2500 \
  --reserve-margin-percent 20 --allocate 0070-console-microruntime=2500
ex finops gates set --default
ex finops prices set --file prices.json --slug list
ex finops forecasts set --file forecasts.json
```

Cada comando imprime o registro que escreveu, e todo comando de entrada aceita `--by` para dizer quem definiu o valor.

### 5. Planeje a partir de um ticket

Informe os componentes que o runtime deve compilar e rodar e deixe o `ex` ler o ticket. `--no-judge` planeja sem chamada a modelo, então nada é pago; sem ele, e com `ANTHROPIC_API_KEY` definida, o planejador julga as fontes por um modelo sob os mesmos gates de uma volta, e o custo entra na linha de planejamento do livro-razão.

```sh
ex agent plan .agents/workflow/tickets/0070-console-microruntime.ticket.md --no-judge \
  --component "path=apps/console,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=docker,args=build -t langdrift-console:local ."
```

A saída da CLI é em inglês:

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

O plano fica em `.langdrift/microruntimes/runtimes/console-microruntime/runtime.md`. O bloco `langdrift` é o que o runtime executa; um segundo bloco, `langdrift-plan`, guarda a análise. A prosa do arquivo é mantida como contexto e nunca é executada.

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

Este trecho omite a versão, o nome e as identidades do runtime e do projeto, que o `ex` gerencia. Planejar o mesmo ticket de novo, com o runtime parado, aumenta a revisão.

### 6. Rode o plano

Rodar o runtime autoriza o seu plano. Uma volta compila imagens, inicia o cluster local do projeto e faz chamadas pagas pelo seu agente de código, então rode `ex doctor` antes.

```sh
ex agent run console-microruntime
```

1. **O gate.** A execução imprime a intenção, a estimativa a preços de tabela e o orçamento por execução, e então reserva. Uma reserva dentro de todos os limites abre a volta; uma acima deles pausa sem nada pago.
2. **Os estágios.** Cada estágio roda em ordem pelo agente de código e imprime um resumo de uma linha. Um estágio em que uma pessoa participa espera pelo owner: responda no terminal ou com `--confirm`, ou depois com `ex agent confirm console-microruntime planning`. O estágio `owner-only` nunca roda.
3. **O runtime.** Os componentes sobem a partir do bloco no cluster kind do projeto, e a execução avisa quando estão prontos.
4. **A entrega.** A execução observa a branch, o pull request e as suas verificações; `--checks-wait` limita por quanto tempo ela os acompanha.
5. **O acerto.** O runtime desce, a reserva não usada volta ao orçamento, a memória `run.md` é escrita e o índice de episódios é reconstruído. A execução diz se está pronta para revisão ou para merge, e o que vem a seguir.

Toda mudança fica registrada no arquivo da volta, no `run.json` e no livro-razão.

### 7. Acompanhe a execução

```sh
ex agent status console-microruntime
ex agent reindex
ex doctor
```

`status` mostra o registro da execução e a volta aberta, `reindex` reconstrói o índice de episódios a partir da memória de cada runtime, e `doctor` verifica as ferramentas nativas. Com `--json`, cada linha da sessão é um evento NDJSON, `{at, session, kind, id, text, card, field, data}`, para que um console ou uma interface de terminal renderize a mesma sessão.

| Código de saída | Significado |
| --- | --- |
| `0` | Sucesso, ou uma execução esperando por uma pessoa. |
| `1` | A operação falhou: uma skill, um hook, uma compilação ou um início. |
| `2` | Uso inválido. |
| `3` | Recusado por uma proteção ou política, como um gate de dinheiro. |
| `4` | Falta um pré-requisito: uma ferramenta, uma entrada, um cluster ou o registro de um runtime. |
| `130` | Interrompido. |

A referência completa de comandos, o layout de `.langdrift/` e as regras de recuperação estão no [README do `ex`](/packages/ex/README.md), em inglês.

### Arquivos de entrada de exemplo

<details>
<summary><code>workflow.json</code>: os estágios, as suas skills e onde uma pessoa participa</summary>

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

Um estágio também pode listar `hooks`: scripts versionados que ele roda primeiro, via `sh`, dentro do projeto.

</details>

<details>
<summary><code>driver.json</code>: o agente de código que executa os estágios</summary>

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

O executável roda a partir da raiz do projeto, nunca por um shell. `usage` indica onde a saída informa os tokens; uma chamada que não informa nenhum deixa o seu gasto desconhecido, e `"unmeasured": "pause"` faz o FG-04 pausar a próxima chamada paga.

</details>

<details>
<summary><code>prices.json</code>: os preços de tabela datados por trás de cada estimativa</summary>

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

Estes valores são ilustrativos e servem para exercitar o formato; não são uma tabela de preços verificada.

</details>

<details>
<summary><code>forecasts.json</code>: os tokens que cada esforço deve consumir</summary>

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

`planning` cobre as chamadas do próprio planejador. As previsões são usadas quando nenhuma execução anterior do ticket mediu os seus tokens.

</details>

### Explore o site e o console

Com Node.js 24 e pnpm 12, a partir da raiz do repositório:

```sh
pnpm install
pnpm website
```

O site roda em `http://localhost:3000`. `pnpm console` inicia o console em `http://127.0.0.1:3001` com o `extub` ao lado em `127.0.0.1:3401`, o que exige Go. Todo número que o console mostra vem do stub.

## O trabalho por trás do LangDrift

Criado por [Jonatas Sales](https://github.com/jonatassales), o LangDrift reúne a prática de FinOps, a entrega com agentes e a engenharia de plataforma. Ele nasceu do projeto Exphere e é construído como um monorepo com pnpm e Turborepo: o site e o console em TypeScript com Next.js, React e Tailwind CSS, e o `ex` e os contratos do Microruntimes em Go, com o SDK oficial da Anthropic para Go no planejador.

O repositório é construído a partir de specs: toda mudança começa por uma spec e pelos seus tickets em [`.agents/workflow/`](/.agents/workflow), o mesmo tipo de ticket que o `ex` planeja. Feedback é bem-vindo pelas [issues do GitHub](https://github.com/jonatassales/lang-drift/issues).

## Referências

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/), Jonathan Morley, FinOps Foundation, 27 de março de 2026. Praticantes descrevem agentes que reconciliam gastos com orçamentos, encontram desperdício e aplicam políticas de custo mais cedo no ciclo de desenvolvimento; esse último tema, levar o FinOps para dentro da entrega, é onde o LangDrift trabalha.
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html), Birgitta Böckeler, martinfowler.com, 2 de abril de 2026. O artigo descreve o harness externo que um time constrói em volta de um agente de código: guias que orientam o agente antes de ele agir e sensores que permitem que ele se corrija. O `ex harness` registra esse harness para um projeto, e os gates de FinOps limitam o que ele pode gastar.

## Agradecimentos

Obrigado a [Matt Pocock](https://x.com/mattpocockuk). As skills de agente com que este repositório trabalha, de grilling a to-spec, to-tickets, implement, tdd e code-review, são adaptadas de [mattpocock/skills](https://github.com/mattpocock/skills), sob a licença MIT, e o workflow de exemplo acima usa várias delas.
