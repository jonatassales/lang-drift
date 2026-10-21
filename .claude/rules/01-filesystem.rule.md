# Filesystem rules

Naming and path conventions for files and directories in this repository. Loaded in every session.

1. **Names use kebab-case.** Files and directories: `early-access`, not `earlyAccess` or `early_access`.

2. **Markdown documents carry a type suffix before the extension**, naming what kind of document they are. A rule is `<name>.rule.md`.

3. **Paths never climb with `../`.** A relative path (`./name`) only reaches a sibling file; anything else uses an absolute path. This holds in every language and framework. In JavaScript and TypeScript, Biome reports `../` imports as errors.
