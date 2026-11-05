# Testing rules

Test placement, browser runners, local evidence and test names, in every language. Loaded in every session.

1. **Tests stay with their owner.** Unit and component tests and native service integration tests stay with their owning application or package. End-to-end tests belong to the microruntime `ex` generates for the runtime, under `.langdrift/`, never under `apps/`. A native service integration test is not classified as E2E just because it exercises HTTP.
2. **No browser runner.** The repository does not install or run Playwright, Puppeteer or another browser runner: no application E2E directories, browser-runner configurations, `test:e2e` scripts or direct Playwright, Puppeteer or axe browser dependencies. Inert upstream development or optional-peer metadata inside reviewed dependencies is not installed test infrastructure. Removing application E2E does not authorize moving those suites or introducing a replacement runner.
3. **Local evidence stays local.** `.audits/` is ignored local experiment and diagnostic output. Its files stay out of commits and pull request links; the pull request summarizes the validation. Local outputs may reference other local outputs.
4. **Test descriptions name the file.** A TypeScript test description starts with the application or package name, follows the module path and ends with the file name, never the exported symbol.

```ts
// Do:
describe('website/early-access/form.validation', () => {})
describe('web/ui/share-bar', () => {})

// Don't:
describe('ShareBar', () => {})
describe('EarlyAccessForm validation', () => {})
```
