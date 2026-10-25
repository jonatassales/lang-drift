# Code rules

Rules for TypeScript and JavaScript in apps and packages. Loaded in every session.

1. **External input is read where it is used.** Environment variables, response bodies, network calls, filesystem reads and SDK results are external, and the operation that uses the value is the defensive boundary. Read `process.env.NAME` at that operation, as a static property access so Next can inline `NEXT_PUBLIC_*` values. No env module, registry, `getEnv`, `requireEnv`, Zod env object or proxy exists to read or check environment variables. A local check beside the use is fine when the API needs a string or the failure must name the variable; it is never extracted into a shared utility.

2. **No invented values.** A missing or malformed external value fails inside the operation that needs it. `process.env.API_KEY ?? 'development-key'` does not happen.

3. **`try/catch` marks an external boundary.** The catch translates the failure into the operation's result, reports it through the existing observability path, applies a product-approved fallback or rethrows with context that leaks no secret. It never swallows the error. A pure internal function fails normally; an internal invariant is not an external failure.

4. **Guard clauses and early returns.** Leave a function, loop or branch as soon as the result is known. After `return`, `throw`, `continue` or `break`, the rest of the function is the next statement: no `else`.

5. **Name compound predicates.** When a condition combines several business or runtime rules, give it a name: `if (!shouldPersistLocale(request, preference, locale)) return response`.

6. **At most two conditional levels.** Beyond that, use a guard, a named function, a `switch` or an explicit state, never a flattened boolean expression that reads worse than the branches.

7. **Readable bindings.** Descriptive names, never one letter. An awaited value that is consumed lands in a named binding, and a result that can be missing is handled before dependent values use it; `await` and `?.` never share a statement. A blank line surrounds an `await` statement and precedes a `return` that follows another statement. Biome does not enforce this spacing; reviewers do.
