# General rules

Rules for every task, file and language in this repository: TypeScript, Go, shell and anything else. Loaded in every session. The stack files add to these rules and never weaken them.

## Authority

1. **Optimize for human readers.** Code must be elegant and obvious for the next person who reads it, not merely convenient for the tool or agent that produced it.
2. **Precedence.** When two sources disagree, the higher one wins:
   1. The latest explicit owner instruction.
   2. The rules in `.claude/rules/`.
   3. An accepted ADR or spec.
   4. The current framework or runtime constraint.
   5. The current tree.
   6. Guardrails, harnesses, audits and tests.
3. **No exception unless the owner states it.** No exception to these rules applies unless the owner states it in the task.
4. **Rules define; checks verify.** Rules state what code must look like and how it must behave. Guardrails, audits, tests and harnesses verify that; they do not define the architecture, and a passing check does not make an invalid implementation acceptable. When a check conflicts with a rule, follow the rule and report the conflict. Tool-specific rule files only activate these rules.
5. **Existing code does not legalize a prohibited pattern.** Adopting a rule does not certify that the current tree complies. New code follows the rule; an existing violation is recorded for remediation, never copied, and a mid-refactor tree's transitional violations never reach new code.
6. **Direct over indirect.** Use the direct dependency, the platform or framework API, external input read at its use site, guard clauses and the framework's file conventions. Add a wrapper, helper, compatibility layer, generated file, fallback or unusual placement only when the current implementation needs that boundary.
7. **Mandatory two-axis code review.** Before every commit or delivery, run the `code-review` skill against the spec and against every `.claude/rules/*.rule.md`, and report its findings under `## Standards` and `## Spec`. A breach of a rule or of the spec is fixed before delivery; no commit is made with unaddressed findings. Several rules have no tool behind them, among them statement spacing, multi-field literals and single responsibility, so the review checks every rule on every change.

## Single responsibility

1. **One responsibility.** Every file, type, component and function has one reason to change. A unit that starts serving a second concern is split.
2. **One function per file.** A source file declares exactly one top-level function, or exactly one type with its methods: a class, a Go type and its methods, a React component, a hook. Functions nested inside it, such as a component's event handlers or a closure, belong to it. Types and constants that exist only for that function may stay in its file. Any other function, including a helper the main function uses, lives in its own file. In Next.js route handlers (`route.ts`), each exported HTTP method handler (`GET`, `POST`, etc.) is an entrypoint.
3. **Group with folders, not with files.** Related single-function files sit together in a folder named after their concern; in Go, that folder is the package. A folder groups them; a multi-function file never does.
4. **The file names its one export.** A file is named after the function or type it declares: `to-initials.fmt.ts` declares `toInitials`, `cookie.service.ts` declares `CookieService`.

## Control flow

1. **Guard clauses and early returns.** Leave a function, loop or branch as soon as the result is known.
2. **No `else` after an exit.** After `return`, `throw`, `panic`, `continue` or `break`, the rest of the function is the next statement; it is never wrapped in `else`.
3. **Name compound predicates.** When a condition combines several business or runtime rules, give it a name.
4. **At most two conditional levels.** Beyond that, use a guard, a named function, a `switch`, an explicit state or a lookup map. Never flatten the logic into a boolean expression that is harder to read than the branches.

```ts
// TypeScript
if (!user) return
if (!user.active) return

if (!shouldPersistLocale(request, preference, locale)) return response

return persist(user)
```

```go
// Go
func activeProject(user *User) (*Project, error) {
	if user == nil {
		return nil, errMissingUser
	}

	if !user.Active {
		return nil, errInactiveUser
	}

	return user.Project, nil
}
```

## Statement spacing

1. **Before `return`.** A `return` that follows another statement in its block has a blank line before it. A `return` that is the first statement of its block, such as the body of a guard, needs none.
2. **Between phases.** Independent phases inside a function (guard validation, state computation, external invocation) are separated by exactly one blank line.
3. **Not at block boundaries.** No blank line right after `{` and none right before `}`.
4. **The review enforces it.** No formatter enforces this spacing, neither Biome nor `gofmt`; the review does (Authority 7). The stack files add their own spacing rules: around `await` in `04-typescript`, after props destructuring in `05-react`.

```ts
// Do:
const title = messages[locale].title

if (!title) return fallbackTitle

return title

// Don't:
const title = messages[locale].title
if (!title) return fallbackTitle
return title
```

```go
// Do:
runtime, err := loadRuntime(path)
if err != nil {
	return nil, fmt.Errorf("load runtime: %w", err)
}

return runtime, nil
```

## Multi-field literals

1. **One field per line.** A literal with more than one field puts every field on its own line: an object literal in JavaScript or TypeScript, a struct or map literal in Go. Break after `{`, put each field on its own line indented one level, and close with `}` on its own line. In a Go test, each case of a test table may stay on one line, so the table reads as rows; a literal outside a table follows the rule.

```ts
// Do:
const user = {
	id: 'usr_123',
	name: 'Alice',
}

// Don't:
const user = { id: 'usr_123', name: 'Alice' }
```

```go
// Do:
user := User{
	ID:   "usr_123",
	Name: "Alice",
}

// Don't:
user := User{ID: "usr_123", Name: "Alice"}
```

## Readable names

1. **Descriptive names, never one letter.** No one-letter variable, parameter, receiver or loop variable (`x`, `e`, `p`, `t`, `i`). In a Go test, the `t *testing.T` parameter keeps the name the toolchain's convention gives it.
2. **Descriptive type parameters.** Generic type parameters are meaningful words, never single letters and never prefixed with `T`: `Mapper<Input, Output>` in TypeScript, `Map[Input, Output any]` in Go.
3. **Name a result before using it.** A result that will be consumed lands in a named binding, and a result that can be missing or invalid is handled before dependent values use it.

## External input and errors

1. **External input.** Environment variables, request and response bodies, network calls, filesystem reads and third-party SDK results are external. The operation that uses the value is the defensive boundary, and the value is read at that operation.
2. **No env module.** No module, registry, service, getter or schema object whose job is to read or check environment variables. A local check beside the use is fine when the API needs a string or the failure must name the variable; it is never extracted into a shared utility.
3. **No invented values.** A missing or malformed external value fails inside the operation that needs it; no hardcoded fallback stands in for it. Tests, guardrails, build scripts and startup validators do not replace this boundary.
4. **Document keys, never values.** Each app documents every required environment variable by key, never by value, in its environment contract.
5. **Errors are handled at the external boundary.** The handler translates the failure into the operation's result, reports it through the existing observability path, applies a product-approved fallback, or returns or rethrows it with context that leaks no secret. It never swallows the error: no empty `catch`, no error discarded with `_`.
6. **Only at the boundary.** A pure internal function fails normally; an internal invariant is not an external failure.
