# React rules

React components in apps and in `packages/web`. Loaded in every session.

1. **Server Components by default.** Add `'use client'` only for browser APIs, interactivity, client state or a client-only library, at the narrowest leaf boundary that needs it.
2. **One component per file, one hook per `.hook.ts`.** As the general rules require of every file. Related files are colocated by concern.
3. **Function declarations.** Components and the functions inside component scope use function declarations, never arrow functions.
4. **Props are a named interface.** A component with props declares `interface <Name>Props` directly above the component, separated by exactly one blank line. The component takes a single `props` argument and destructures it on the first line of the body, followed by a blank line. The interface stays with its component unless it has a real shared owner. An anonymous object annotation, a props type alias or a destructured parameter list does not express this contract.
5. **Function scope order.** Props destructuring; React hooks (`useState`, `useReducer`, `useContext`, `useMemo`, `useCallback`, `useTransition`, `useId`); custom hooks; third-party hooks (`useAtomValue`, `useTranslations`); named event handlers; derived `const` values and memoized values; the JSX return.
6. **Named event handlers.** Handlers passed to JSX props are named functions defined in the component scope, never inline.
7. **Effects are a last resort.** Derive state during render, handle work in event handlers or encapsulate a subscription in a custom hook. Use `useEffect` only when strictly necessary, such as an external subscription or a DOM measurement, and document why in the code or the pull request.
8. **Transitions, not loading flags.** Use `useTransition` or `startTransition` for a non-urgent update that must keep the interface responsive. A `useState` loading flag (`const [isLoading, setIsLoading] = useState(false)`) is for a real application state, not for a transition React already models with `isPending`.
9. **Scoped `<Suspense>` and route boundaries.** `<Suspense>` wraps an asynchronous section that can render on its own: not the whole page for one section, and not every component. A route that can suspend during navigation gets a `loading.tsx`; a segment that needs a runtime error boundary gets an `error.tsx`. `error.tsx` handles render failures and `try/catch` handles an imperative external call; they are different boundaries. `template.tsx` exists only when a route must remount on navigation; an effect or a key hack never imitates it.
10. **Routes compose pages.** Independently meaningful UI is composed at its owning page, layout or feature: the filesystem owns page identity. A route composes its page instead of passing a page discriminator into a universal view, and no central component selects whole pages through a mode or section prop. Navigation and title models may represent useful data without owning route dispatch. Configurable components, render props and dependency injection stay allowed.
11. **Shared primitives live in `packages/web`.** Generic UI, the icon family and their dependencies have one owner there; use its SmoothUI and shadcn foundations instead of rebuilding app-local base primitives. An app never adds its own icon dependency, so every app uses one icon identity. Apps consume the shared language control through `packages/web` and keep their business composition, messages and behavior.
12. **No Context for shared state.** No React Context provider holds shared application state; shared client state uses Jotai atoms (state rules). Framework providers, such as the next-intl provider, are not application state.

```tsx
// Do:
interface ButtonProps extends React.ComponentProps<'button'> {
	variant?: 'primary' | 'secondary'
}

export function Button(props: ButtonProps) {
	const { variant = 'primary', className, children, onClick, ...rest } = props

	function handleClick(event: React.MouseEvent<HTMLButtonElement>) {
		onClick?.(event)
	}

	return (
		<button
			className={cn(
				'inline-flex items-center gap-2',
				variant === 'primary' && 'bg-primary',
				className,
			)}
			onClick={handleClick}
			{...rest}
		>
			{children}
		</button>
	)
}

// Don't:
type ButtonProps = { variant?: string }
export const Button = ({ variant, ...rest }: ButtonProps) => (
	<button onClick={(event) => event.preventDefault()} {...rest} />
)
```

```tsx
// Do:
export function Status() {
	const [isOnline] = useNetworkStatus()

	return <span>{isOnline ? 'Online' : 'Offline'}</span>
}

// Don't:
export function Status() {
	const [isOnline, setIsOnline] = useState(false)

	useEffect(() => {
		window.addEventListener('online', () => setIsOnline(true))
	}, [])

	return <span>{isOnline ? 'Online' : 'Offline'}</span>
}
```
