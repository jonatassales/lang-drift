# @langdrift/core

The core of LangDrift Microruntimes, published for the people who run the platform's engine outside this repository. It is organized by surface, one folder per language, so that a contract has one definition per language and the languages sit side by side:

| Surface | Contract |
| --- | --- |
| [`contracts/go`](https://github.com/jonatassales/lang-drift/blob/main/lib/core/contracts/go/README.md) | The Go module the LangDrift CLI imports: the `runtime.md` fence, the component env file and image tag, the lexical rules, the connected-agent protocol, and the judge's contract. |
| `@langdrift/core/judge/*` | The judge's JSON Schemas and samples (`contracts/go/judge/schema/`), published on npm for the console and anyone who runs the engine: the same shapes the Go `judge` package declares. Both sides test against these files, so a change on one side is a change on both. |

Each surface is published in the form its language expects:

- **npm** carries the judge's JSON Schemas and samples, and nothing else:

  ```bash
  pnpm add @langdrift/core
  ```

  ```ts
  import requestSchema from '@langdrift/core/judge/request.schema.json' with {
  	type: 'json',
  }
  ```

- **Go** gets the module by git tag, as [its README](https://github.com/jonatassales/lang-drift/blob/main/lib/core/contracts/go/README.md) explains:

  ```bash
  go get github.com/jonatassales/tao/lib/core/contracts/go@latest
  ```

The manifest's scripts run the Go module's tests and `vet` through Turborepo like any package's.
