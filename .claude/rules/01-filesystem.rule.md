# Filesystem rules

Naming and path conventions for every file and directory in this repository, in every language. Loaded in every session.

1. **Names use kebab-case.** Files and directories: `early-access`, not `earlyAccess` or `early_access`. Where a language toolchain requires another form, the toolchain wins: Go keeps the `_test.go` suffix and build-constraint suffixes such as `_linux.go`, and a Go package directory is one lowercase word. Conventional files keep their names: `README.md`, `AGENTS.md`, `LICENSE`, `NOTICE`, `SECURITY.md`, `SKILL.md`.
2. **Markdown documents carry a type suffix** before the extension, naming what kind of document they are: `.rule.md`, `.prd.md`, `.spec.md`, `.tdd.md`, `.ticket.md`, `.adr.md`, `.research.md`, `.policy.md`, `.prompt.md`, `.guide.md`. A rule is `<name>.rule.md`.
3. **Paths never climb with `../`.** This holds in every language and framework, including CSS imports and shell `source`, and it includes traversal embedded in a longer path. A relative path (`./name`) reaches only a sibling file; anything else uses an absolute path or a configured alias. In JavaScript and TypeScript, Biome reports `../` imports as errors.
4. **No bucket directories.** No `utils/`, `helpers/`, `common/` or `misc/`. Organize by domain, feature or infrastructure concern.

```text
Do:     early-access/  submit-form.handler.ts  circuit-breaker.ts  website-header.tsx
Don't:  earlyAccess/   submitForm.handler.ts   CircuitBreaker.ts   website_header.tsx
```
