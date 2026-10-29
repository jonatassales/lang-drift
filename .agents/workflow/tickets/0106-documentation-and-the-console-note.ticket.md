# 0106: The folder, the commands and the vocabulary are documented

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the CLI README describes the `.langdrift/` folder it writes (the table the sample's README carried), the commands with their exit codes, the inputs and mirrors, the session's lines and the JSON events, the judge's credential and the driver configuration, and the upgrade note about old private state; the context file's glossary gains episode, judge, stage driver, mirror and input record, and its commands table names the new commands; the dead-code allow list drops anything the removed scaffold kept; and a follow-up spec is drafted for the console to read real records through a connector, mapping the plan fence and the run record to the console's ticket phases.

**Blocked by:** 0093 to 0105.

**Status:** ready-for-agent

- [ ] The README's command table lists every command the binary registers and no other, and its boundary section matches what the code writes.
- [ ] The glossary defines the new terms with their avoid-lists, in the file's format.
- [ ] `pnpm turbo run check`, `drift check localization` and `drift check dead-code` pass.
- [ ] The follow-up spec exists with status needs-triage and names the fields it maps.

### Verification

- `pnpm turbo run check`; `drift check dead-code`; `./dist/ex --help` compared with the README table.
