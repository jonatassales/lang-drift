# 0159: The website's components follow the React rules

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** Every component with props declares `interface <Name>Props` right above it and destructures `props` on its first line, followed by a blank line; no props interface carries another component's name. Component scopes follow the rule's order. Event handlers are named functions in the component scope, with no inline function in JSX and no arrow function in component scope. `'use client'` stays only where a component needs the browser. Effects that derive state move to render or to the event handler with the same behavior, and every remaining effect carries a comment with its reason.

**Blocked by:** 0158

**Status:** ready-for-agent

- [ ] Every props interface sits right above its component and matches its name.
- [ ] Every component scope follows the rule's order.
- [ ] No inline event handler or arrow function remains in component scope.
- [ ] `'use client'` marks only components that need the browser.
- [ ] No effect derives state, and every effect says why it exists.
- [ ] Behavior is unchanged: the test suite, the build and the captures still pass.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check react`
- Built-in browser captures compared with the baseline.
