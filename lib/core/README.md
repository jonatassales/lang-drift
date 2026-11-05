# @langdrift/core

The core of LangDrift Microruntimes, published for the people who run the platform's engine outside this repository. It is organized by surface, one folder per language, so that a contract has one definition per language and the languages sit side by side:

| Surface | Contract |
| --- | --- |
| [`contracts/go`](https://github.com/jonatassales/lang-drift/blob/main/lib/core/contracts/go/README.md) | The Go module the LangDrift CLI imports: the `runtime.md` fence version 3, the judge's three answers, the five stage reports and the hand-off records the console receives, beside the older shapes that leave in tickets 0323, 0324 and 0325. |
| `@langdrift/core/configuration/*` | The `runtime.md` fence's JSON Schema and the simulation's fence as a sample (`contracts/go/configuration/schema/`). |
| `@langdrift/core/judgment/*` | The JSON Schemas and samples of the judge's three answers, `criteria`, `episodes` and `suggestion` (`contracts/go/judgment/schema/`). |
| `@langdrift/core/stages/*` | The JSON Schemas and samples of the stage reports, `planning`, `spec`, `tickets`, `implement` and `review` (`contracts/go/stages/schema/`), so a skill's author can validate a report before a run. |
| `@langdrift/core/judge/*` | The older judge contract's JSON Schemas and samples (`contracts/go/judge/schema/`), which the console's judge route tests against, until ticket 0325. |

Each schema is the same shape its Go package declares, and the Go tests check each schema and its Go check against the same samples. The hand-off records are Go types only; no schema is published for them.

Each surface is published in the form its language expects:

- **npm** carries the JSON Schemas and samples, and nothing else:

  ```bash
  pnpm add @langdrift/core
  ```

  ```ts
  import reportSchema from '@langdrift/core/stages/implement.report.schema.json' with {
  	type: 'json',
  }
  ```

- **Go** gets the module by git tag, as [its README](https://github.com/jonatassales/lang-drift/blob/main/lib/core/contracts/go/README.md) explains:

  ```bash
  go get github.com/jonatassales/tao/lib/core/contracts/go@latest
  ```

The manifest's scripts run the Go module's tests and `vet` through Turborepo like any package's.
