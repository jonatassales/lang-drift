# 0095: `ex init` writes the project identity and the execution rules only

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** `ex init` creates `.langdrift/config.json`, `.langdrift/microruntimes/AGENTS.md` (the Microruntimes execution rules, byte-identical to the sample's, embedded in the binary, added when missing and never replaced) and the private state directory with its own ignore file, and nothing else. The embedded scaffold of skills, templates, rules, guardrails and the workflow proposal goes, with its bootstrap and its tests. `config.json` is written with the sample's tab layout and a readable mode, since the folder is the consumer's to read. A second `ex init` changes nothing and reports the existing project.

**Blocked by:** 0093 (the fixture the rules are compared with), 0094 (the moved paths).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `ex init` in an empty directory leaves exactly `config.json`, `microruntimes/AGENTS.md` and `agent/state/.gitignore` under `.langdrift/`.
- [ ] The written `AGENTS.md` equals the fixture byte for byte and `config.json` equals the sample's layout with a fresh id.
- [ ] No scaffold asset remains in the source tree and `drift check dead-code` reports nothing new.
- [ ] The README's command table and boundary section describe what `ex init` writes now.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; `./dist/ex init` in a scratch directory; `drift check dead-code`.
