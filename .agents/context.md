# Context

Project knowledge for agents working in this repository. This file grows as decisions are made; it does not get restructured into subfolders unless the owner asks.

## Glossary

### Apps

**App**:
One LangDrift offering with its own URL, such as the website. Apps share account, plan and shell.
_Avoid_: Product, module

### Design tokens

**Shared token**:
A design token with the same value in every app.
_Avoid_: Global token, base token

**App theme**:
The tokens through which one app expresses its brand. It holds only the primary color and its hover state; today only the website has one.
_Avoid_: Theme (alone), product theme, brand theme, skin

**Foundation**:
A raw value on a scale, such as gray 50 to 950, that roles draw from. Components never use a foundation directly.
_Avoid_: Primitive, palette color

**Role**:
A token named for its purpose, such as background, foreground or primary, rather than for its value. Components use roles only.
_Avoid_: Semantic color, alias

**Status color**:
A role that signals success, warning, danger or info. It is the same in every app and never takes the brand color.
_Avoid_: Feedback color, alert color

**Classification color**:
The color of one Product Vision classification: intentional, aligned, unexplained or review. It is the same in every app and color scheme, and never takes the brand color.
_Avoid_: Category color, label color

**Classification text tone**:
The classification color at a contrast fit for small text: darker in light, lighter in dark. It exists for intentional, review and unexplained.
_Avoid_: Classification ink, label tone

**Drift color**:
The color of a drift's sign: negative, positive or zero. It is the same in every app and never takes the brand color.
_Avoid_: Delta color, trend color

**Color scheme**:
Light or dark. Dark is the default in every app.
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
