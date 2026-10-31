# 0156: Every website file declares one function, named after it

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** Each website source file with more than one function becomes single-function files grouped in a folder named after their concern: the pages and layouts with their `generateMetadata` and `generateStaticParams`, the microruntimes cards, the content and workflow data, the metadata helpers, the services, the factory, the proxy's locale logic and the route handlers. The early-access and book-a-demo form logic is left to 0160. The messages barrel holds only named re-exports, its logic in files of its own. Files are named after their one export, in kebab-case, with the suffix of the TypeScript rules' table where one names their responsibility. Relative imports that reach into a subfolder use an alias. Test helpers stay in their test.

**Blocked by:** 0155

**Status:** ready-for-agent

- [ ] `drift quality check general` reports no "More than one function" finding in `apps/website`.
- [ ] Every route file re-exports the `generateMetadata` or `generateStaticParams` it no longer declares.
- [ ] `template/i18n/messages/index.ts` holds named re-exports only.
- [ ] Every file is named after its export, with the suffix of its responsibility where the table has one.
- [ ] No relative import reaches beyond a sibling file.
- [ ] Behavior is unchanged: the test suite, the build and the captures still pass.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check general` and `platform/cli/drift quality check filesystem`
