# Vidro editor preview package

- Package: `vidropro@0.1.0-preview.3`, MIT, a local build not yet published to npm.
- Upstream repository: `https://github.com/jonatassales/vidro`, `packages/editor`.
- Source revision: `ba9966c0d0f4daaee68a3b71181639cb0c3a2158` (main, PR 12 merged).
- Archive: `vidropro-0.1.0-preview.3.tgz`, SHA-256 `501ed9c62a3146ff4eb1db8afaa0c14239e390b5b4f84441c55d22349f647e1f`.

The interface locales are `en`, `pt-BR`, `zh` (Simplified), `ja` and `vi`.
There is no Traditional Chinese catalog, so a `zh-Hant` page sets the editor's
`uiLocale` to `en`. Simplified Chinese never stands in for it.

## Remove after publication

This directory is temporary. Once `vidropro` is published, replace the
`file:vendor/vidropro-0.1.0-preview.3.tgz` dependency in `packages/web` with the
published version, regenerate `pnpm-lock.yaml` and delete this directory.
