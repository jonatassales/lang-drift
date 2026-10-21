# Context

Project knowledge for agents working in this repository. This file grows as decisions are made; it does not get restructured into subfolders unless the owner asks.

## Glossary

### Products

**Product**:
One LangDrift offering with its own URL, such as Vision or FinOps. Products share account, plan and shell.
_Avoid_: App, module

### Design tokens

**Shared token**:
A design token with the same value in every product.
_Avoid_: Global token, base token

**Product theme**:
The tokens through which one product expresses its brand. Today it holds only the primary color and its hover state.
_Avoid_: Theme (alone), brand theme, skin

**Foundation**:
A raw value on a scale, such as gray 50 to 950, that roles draw from. Components never use a foundation directly.
_Avoid_: Primitive, palette color

**Role**:
A token named for its purpose, such as background, foreground or primary, rather than for its value. Components use roles only.
_Avoid_: Semantic color, alias

**Status color**:
A role that signals success, warning, danger or info. It is the same in every product and never takes the brand color.
_Avoid_: Feedback color, alert color

**Color scheme**:
Light or dark. Dark is the default in every product.
_Avoid_: Theme, mode

### Ex CLI

**Microruntime**:
One console-defined workload unit that ex runs on the project's local cluster: its configuration, execution and cleanup.
_Avoid_: Runtime (alone, in new writing), service

**Connected agent**:
The long-running `ex agent` process that supervises a project's microruntimes from authenticated console intent.
_Avoid_: Daemon, worker

**Build recipe**:
The console-supplied command that builds a microruntime's image. Only a committed project script or an allowlisted tool may run.
_Avoid_: Build command, build step

**Execution ledger**:
The on-disk record of what an operation created, written before side effects and settled afterwards. Cleanup acts only on what the ledger owns.
_Avoid_: State file, journal

**Evidence**:
The token-free record ex keeps of every lifecycle operation and inspection, for later audit.
_Avoid_: Log, report
