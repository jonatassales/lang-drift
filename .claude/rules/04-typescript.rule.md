# TypeScript rules

TypeScript and JavaScript in apps and packages. Loaded in every session. The general rules apply first; this file adds what is specific to TypeScript.

## Formatting

1. **Biome formatter.** Tabs (`indentStyle: "tab"`) with a visual width of 2 (`indentWidth: 2`), never spaces; a line width of 80 (`lineWidth: 80`), breaking statements, parameter lists, JSX tags and ternaries as they approach it; single quotes in JavaScript, TypeScript and CSS (`quoteStyle: 'single'`) and double quotes in JSX attributes (`jsxQuoteStyle: "double"`); semicolons as needed (`semicolons: "asNeeded"`), only where a line starting with `(` or `[` requires one.
2. **Multi-field literals.** One property per line (General, Multi-field literals).

## Statement spacing

**Around `await`.** Every statement that contains `await` outside a nested function has a blank line before it and a blank line after it. Two consecutive `await` statements are also separated by a blank line. The blank line is not required only where the statement is the first or the last statement of its block. This adds to the general spacing rules (blank line before `return`, between phases, none at block boundaries).

```ts
// Do:
const { locale } = await params

if (!hasLocale(routing.locales, locale)) notFound()

const title = messages[locale].title

const images = await socialImages(parent)

return pageMetadata(title, images)

// Don't:
const { locale } = await params
if (!hasLocale(routing.locales, locale)) notFound()
const title = messages[locale].title
const images = await socialImages(parent)
return pageMetadata(title, images)
```

## Imports

1. **Aliases across boundaries.** An import that does not reach a sibling file uses a configured absolute alias (`@app/*`, `@components/*`, `@state`, `@langdrift/web/*`); see the filesystem rules for `../` and `./`.
2. **Aliases are developer-owned.** Make navigation and ownership meaningful rather than imposing a universal vocabulary. Meaningful absolute subpaths are allowed.
3. **Barrels re-export only.** A barrel (`index.ts`) holds named re-exports and nothing else (`export { MyComponent }`), never `export * from './module'`. Barrels are intentional and optional.
4. **Import order.** Three groups separated by exactly one blank line (Biome's `:BLANK_LINE:`): third-party packages; workspace packages and absolute aliases; relative sibling imports.
5. **Aliases live in `tsconfig.paths.json` only**, extended by `tsconfig.json`: no `baseUrl` and no duplicated `paths`.
6. **Biome's import matcher follows the real aliases.** The import-group matcher list in `biome.json` stays aligned with the applications' actual aliases and first-party package names. It is configuration, not an approved alias vocabulary. Biome preserves side-effect import order and detached-comment boundaries; it does not prove that every import belongs to one globally ordered block. `biome lint` does not run import-organizing assists or JSON formatting; check JSON configuration with `biome format` without `--write`.

```ts
// Do:
import type { ReactNode } from 'react'
import { z } from 'zod'

import { Button } from '@langdrift/web/ui'
import { toCurrency } from '@app/template/formatters/to-currency.fmt'

import { mapFormValues } from './map-form-values.map'

// Don't:
import { loginFormState } from './form-state.data'
import { cn } from '@langdrift/web/utilities'
import Link from 'next/link'
import { DemoButton } from './website-header/demo-button'
```

## File suffixes

A filename's suffix names the file's responsibility. With one function per file (General, Single responsibility), every function suffix is singular.

| Suffix | Responsibility |
| --- | --- |
| `.action` | One server action or command-style function. |
| `.adapter` | One adapter that implements a port or integrates an external system. |
| `.atom` | The Jotai atoms of one concern. |
| `.compute` | One pure computation or derived calculation. |
| `.controller` | One imperative controller for an external client-side integration, such as a challenge widget. |
| `.data` | Static data, defaults or initial state. |
| `.enum` | Enum definitions. |
| `.factory` | One factory that builds an instance (`src/factories/`). |
| `.fmt` | One formatting function (string, number, date, currency). |
| `.guard` | One route guard. |
| `.handler` | One event or interaction handler. |
| `.hook` | One custom React hook. |
| `.loader` | One loader of an external client library or script. |
| `.map` | One mapping function or state mutation. |
| `.mock` | One mock or fixture. |
| `.port` | The interface of one real boundary. |
| `.route` | The implementation a route handler (`app/api/**/route.ts`) delegates to. |
| `.service` | One service class. |
| `.state` | One state container or state contract. |
| `.type` | One type definition. |
| `.types` | Several type definitions of one concern. |
| `.validation` | One Zod schema. |

**`.service` files are classes.** A `.service` file exports one class as its implementation. A function-only module uses another suffix, such as `.compute`, `.data`, `.fmt` or `.map`.

## External input

1. **Read `process.env.NAME` at the operation**, as a static property access, so Next can inline `NEXT_PUBLIC_*` values. No `getEnv`, `requireEnv`, `readEnv`, Zod env object, hydration step or proxy.
2. **`try/catch` is the boundary handler.** It wraps the external operation and follows the general error rules; it never wraps a pure internal function.

```ts
// Do:
const secret = process.env.CRYPTO_SECRET

if (!secret) throw new Error('Missing CRYPTO_SECRET')

// Don't:
const apiKey = process.env.API_KEY ?? 'development-key'
```

## Bindings

1. **`await` and optional chaining never share a statement.** An awaited value that is consumed lands in a named binding first. Continue, use the specified fallback or throw according to the application contract; this is not a ban on optional chaining elsewhere, nor a requirement to throw for every absent value.
2. **The translator is `translate`.** The next-intl translator is named after what it does, never `t`.

```ts
// Do:
const translate = useTranslations('status')

// Don't:
const t = useTranslations('status')
```

## Types and contracts

1. **`interface` over `type` for object shapes.** Reserve `type` for unions, intersections, primitives, tuples, mapped types and other shapes an `interface` cannot express.
2. **No inline object types in signatures.** When a function, class or generic needs a non-trivial object type, define a named `interface` directly above it in the same module.
3. **Domain types stay with their owner** and are imported with `import type`. Component props and types derived only from props and used only in that file stay local. Package-owned types keep their package ownership; a view barrel is not their new owner.

```ts
// Do:
export interface TokenPayload {
	sub: string
	email: string
	sid: string
}

// Don't:
export type TokenPayload = {
	sub: string
	email: string
	sid: string
}
```
