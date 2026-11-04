# @langdrift/core

The core of LangDrift Microruntimes, published for the people who run the platform's engine outside this repository. It is organized by surface, one folder per language, so that a contract has one definition per language and the languages sit side by side:

| Surface | Contract |
| --- | --- |
| [`contracts/go`](contracts/go/README.md) | The Go module the `ex` CLI imports: the `runtime.md` fence, the component env file and image tag, the lexical rules, and the connected-agent protocol. |

The Go surface is published as a Go module by git tag, `lib/core/contracts/go/vX.Y.Z`, the form Go expects of a module in a subdirectory. This manifest stays private until a surface that npm can carry exists; the scripts run the Go module's tests and `vet` through Turborepo like any package's.
