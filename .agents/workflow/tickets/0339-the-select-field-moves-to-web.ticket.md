# 0339: the select field moves to web

**Spec:** .agents/workflow/specs/0068-the-session-panel-primitives-move-to-web.spec.md

**What to build:** `packages/web` gains `SelectField`, under its own `web/ui/select-field` entry: a native select with a label, a caret, an optional note and a locked state, generic over its option values. The Workflow card's model and effort choices compose it in place of the website's local workflow choice, which is deleted. The choices reprice the run and lock while it plays, exactly as before.

**Blocked by:** 0337. Both tickets change the Workflow card.

**Status:** ready-for-agent

**Alignment:** change, in code only: no simulated line or card field changes. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The Workflow card's "Model" select, with its price note, and its "Effort" select are drawn by `web/ui/select-field`.
- [ ] Choosing a model or an effort reprices the run and warns over the budget. While a run plays, both selects are `aria-disabled` and the choice does not change. The session cards tests prove both with no changed assertion.
- [ ] The select's text is 16px on phones and 12px from `sm`.
- [ ] The website no longer has its workflow choice component.
- [ ] `pnpm turbo run check`, `typecheck` and `test` pass.

## Dev part

**Category:** `hard-part`. This ticket's difficulty is a native select that stays type-safe over its option values: the DOM hands back a string, and the piece must turn it back into the option's typed value and never into a raw string. Locked, the select keeps its focus and changes nothing on its own: refusing the change stays with the consumer.

**Paths:**

- `packages/web/src/ui/select-field/select-field.tsx`
- `packages/web/src/ui/select-field/select-field.unit.tsx`

**Boundary:** the agent writes these, outside the paths:

- `packages/web/src/ui/select-field/select-field-option.type.ts`, which exports `interface SelectFieldOption<Value extends string> { label: string; value: Value }`;
- the barrel, which re-exports `SelectField`, `SelectFieldProps` and `SelectFieldOption`, and the package's `./ui/select-field` export entry.

The Workflow card calls `<SelectField<Value>>` with `SelectFieldProps<Value>`:

- `label: string`;
- `locked: boolean`;
- `name: string`;
- `note?: string`;
- `onChange: (value: Value) => void`;
- `options: readonly SelectFieldOption<Value>[]`;
- `value: Value`.

**Goal:** write the piece and its test, with the classes and ARIA of the website's workflow choice.

- [ ] The label comes through `Eyebrow`, the caret is decorative, and the note shows only when given.
- [ ] A change calls `onChange` with the chosen option's typed value, and with nothing when no option matches.
- [ ] Locked, the select is `aria-disabled` and keeps focus.
- [ ] The unit test, described as `web/ui/select-field/select-field`, covers the options, the value, the typed change, the locked state and the note.
