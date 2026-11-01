<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Plataforma FinOps Agêntica" width="820">
  </picture>
</p>

<p align="center">
  <strong>Entrega assistida por agentes, prestando contas ao orçamento do seu time.</strong><br>
  Governança de custos pré-execução e microruntimes delimitados para agentes de IA. Feito para pessoas e agentes.
</p>

<p align="center">
  <a href="/README.md">English</a> ·
  <strong>Português (Brasil)</strong> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@langdrift/ex"><img alt="npm @langdrift/ex" src="https://img.shields.io/npm/v/@langdrift/ex?logo=npm&amp;logoColor=white&amp;label=%40langdrift%2Fex&amp;color=CB3837&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md/pt-br/">Site</a> ·
  <a href="#recursos">Recursos</a> ·
  <a href="#como-usar">Como usar</a> ·
  <a href="#configuracao">Configuração</a> ·
  <a href="#referencias">Referências</a> ·
  <a href="#agradecimentos">Agradecimentos</a> ·
  <a href="#licenca">Licença</a>
</p>

## Recursos

- **Reserva de orçamento antes da execução:** Cada chamada de modelo e recurso de nuvem é estimada a preços de tabela datados e reservada contra o envelope do time antes de qualquer execução. Se uma estimativa ultrapassar o teto, a execução pausa antes que qualquer custo ocorra.
- **Microruntimes delimitados:** Planeje e execute ambientes locais isolados (Docker, kind, Terraform, Kubernetes) adaptados ao software que você já possui, sem introduzir dependências ou SDKs proprietários no seu código.
- **Traga seu próprio agente de código:** Conduza a entrega pelo driver de agente de sua preferência (como Claude Code) através de estágios estruturados de workflow (planejamento, implementação, revisão), com pontos claros de confirmação humana.
- **Trilhas de auditoria determinísticas:** Deixe evidências sem tokens, livros-razão verificáveis e memórias de execuções anteriores que refinam os próximos planos.

## Como usar

### Pré-requisitos

- **macOS** ou **Linux**
- **Node.js 24+** (para executar `@langdrift/ex` via `npx`)
- **Docker**, **kind**, **kubectl** e **Terraform** (para executar microruntimes locais delimitados; verifique seu ambiente a qualquer momento com `npx @langdrift/ex doctor`)
- `ANTHROPIC_API_KEY` configurada no ambiente (para o planejador de modelos)

### 1. Inicializar o projeto

Execute o inicializador na raiz do seu repositório:

```sh
npx @langdrift/ex init
```

Isso cria a configuração de identidade do projeto em `.langdrift/config.json` e as regras locais de execução em `.langdrift/microruntimes/AGENTS.md`. Os arquivos existentes permanecem inalterados. Adicione `/.langdrift/` ao seu `.gitignore`: essa pasta armazena o estado gerado pelo `ex`.

### 2. Planejar a partir de uma issue ou ticket

Aponte o `ex` para um ticket ou especificação. Ele lê os documentos vinculados, dimensiona limites de CPU e memória dentro do envelope do seu time, escolhe o modelo ideal e escreve um plano de execução delimitado:

```sh
npx @langdrift/ex agent plan <caminho-do-ticket.md>
```

> **Dica:** Use `--no-judge` para gerar o plano por heurística sem realizar chamadas de LLM, sem nenhum custo.

<p align="center">
  <img src="/platform/docs/assets/images/ex-agent-plan.webp" alt="Terminal executando ex agent plan no ticket 0070. Ele lê sete documentos vinculados, aplica o envelope de 2000 milicores e 3072 MiB do time e o orçamento de $25.00 por execução, sugere claude-sonnet-5-5 com esforço médio, dimensiona apps/console para 500m e 1024 MiB e encerra com uma proposta de plano, revisão 1." width="820">
</p>

O plano é salvo em `.langdrift/microruntimes/runtimes/<runtime-id>/runtime.md`. O bloco `langdrift` contém a configuração executável:

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

### 3. Executar o plano

Executar o runtime autoriza o plano e inicia o ciclo de entrega:

```sh
npx @langdrift/ex agent run <runtime-id>
```

Durante a execução, o `ex` coordena todo o ciclo de entrega:

1. **Checagem de portão:** Valida se os preços cobrem a execução e reserva o custo estimado contra a alocação do ticket e o envelope do período.
2. **Estágios do workflow:** Guia o agente de código pelos estágios definidos no workflow. Quando a aprovação humana é exigida, a execução pausa até ser confirmada no terminal ou com `npx @langdrift/ex agent confirm <runtime-id> <estagio>`.
3. **Runtime local:** Constrói imagens e inicia os componentes em um cluster kind local, verificando os endpoints de prontidão.
4. **Observação da entrega:** Monitora commits, branches e verificações de PR em modo somente leitura.
5. **Liquidação:** Desativa recursos do runtime local, devolve reservas de orçamento não utilizadas e registra a execução no livro-razão.

### 4. Inspecionar status e diagnósticos

Acompanhe execuções ativas, inspecione eventos de sessão e verifique ferramentas locais:

```sh
# Exibir status da execução atual e estado da volta
npx @langdrift/ex agent status <runtime-id>

# Reconstruir o índice de memória a partir de execuções anteriores
npx @langdrift/ex agent reindex

# Verificar pré-requisitos e integridade das ferramentas
npx @langdrift/ex doctor
```

Adicione `--json` ao `status` para receber eventos NDJSON compatíveis com consoles e dashboards.

## Configuração

Você pode personalizar estágios de workflow, drivers de agentes e parâmetros de custos usando arquivos de configuração JSON no seu repositório:

<details>
<summary><code>workflow.json</code>: estágios de entrega e pontos de aprovação humana</summary>

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
<summary><code>driver.json</code>: definição do driver do agente de código</summary>

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
<summary><code>prices.json</code>: tabela de preços de tokens</summary>

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
<summary><code>forecasts.json</code>: estimativas de tokens por nível de esforço</summary>

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

## Referências

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/), Jonathan Morley, FinOps Foundation, 27 de março de 2026. Praticantes descrevem agentes que reconciliam gastos com orçamentos, encontram desperdício e aplicam políticas de custo mais cedo no ciclo de desenvolvimento; esse último tema, levar o FinOps para dentro da entrega, é onde o LangDrift trabalha.
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html), Birgitta Böckeler, martinfowler.com, 2 de abril de 2026. O artigo descreve o harness externo que um time constrói em volta de um agente de código: guias que orientam o agente antes de ele agir e sensores que permitem que ele se corrija. O `ex harness` registra esse harness para um projeto, e os gates de FinOps limitam o que ele pode gastar.

## Agradecimentos

Obrigado a [Matt Pocock](https://x.com/mattpocockuk). As skills de agente com que este repositório trabalha, de grilling a to-spec, to-tickets, implement, tdd e code-review, são adaptadas de [mattpocock/skills](https://github.com/mattpocock/skills), sob a licença MIT, e o workflow de exemplo acima usa várias delas.

## Licença

O LangDrift é um software de código aberto licenciado sob a [Licença MIT](/LICENSE).
