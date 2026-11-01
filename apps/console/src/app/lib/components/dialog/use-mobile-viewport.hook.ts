import { useSyncExternalStore } from 'react'

const mobileQuery = '(max-width: 639px)'

// Whether the viewport is below Tailwind's sm breakpoint, where the console's
// navigation becomes a drawer and its dialogs turn modal. The server renders
// the wide layout.
export function useMobileViewport() {
	function subscribe(onChange: () => void) {
		const media = window.matchMedia(mobileQuery)

		media.addEventListener('change', onChange)

		return function unsubscribe() {
			media.removeEventListener('change', onChange)
		}
	}

	function readViewport() {
		return window.matchMedia(mobileQuery).matches
	}

	function readServerViewport() {
		return false
	}

	return useSyncExternalStore(subscribe, readViewport, readServerViewport)
}
