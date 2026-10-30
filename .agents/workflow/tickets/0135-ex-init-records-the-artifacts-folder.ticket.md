# 0135: ex init records the artifacts folder

**Spec:** `.agents/workflow/specs/0024-ex-init-asks-finops-and-generates-ci-gates.spec.md`

**What to build:** `ex init` asks for the folder that holds the team's tickets and specs, or takes `--artifacts`, and stores it in `config.json` as a slash-separated path relative to the project root. An empty answer skips it. A second run does not ask again once a folder is stored. The folder is only stored here: `ex agent plan` does not read it yet.

**Blocked by:** 0134 (ex init asks for the budget and the default gates).

**Status:** ready-for-agent

- [ ] An existing relative folder is stored as `artifacts` in `config.json`, and the result lists `artifacts` as recorded.
- [ ] An absolute path, a path that climbs with `..` or a missing folder is refused with the reason, and nothing is recorded.
- [ ] Without an artifacts folder, `config.json` keeps today's exact layout.
- [ ] A second `init` does not ask for the folder once one is stored.
- [ ] The README's `.langdrift` table describes the new field.
