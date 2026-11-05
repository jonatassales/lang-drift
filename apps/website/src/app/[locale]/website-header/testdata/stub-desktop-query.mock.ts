import { vi } from 'vitest'

export interface DesktopQuery {
	grow: () => void
	listeners: Set<() => void>
}

// Replaces matchMedia with one query that starts below the desktop
// breakpoint and crosses it when the test grows the viewport.
export function stubDesktopQuery(): DesktopQuery {
	const listeners = new Set<() => void>()
	const query = {
		addEventListener(event: string, listener: () => void) {
			if (event === 'change') listeners.add(listener)
		},
		matches: false,
		removeEventListener(event: string, listener: () => void) {
			if (event === 'change') listeners.delete(listener)
		},
	}

	vi.stubGlobal('matchMedia', () => query)

	return {
		grow() {
			query.matches = true

			for (const listener of listeners) listener()
		},
		listeners,
	}
}
