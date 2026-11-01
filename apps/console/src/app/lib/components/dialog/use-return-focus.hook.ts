import { useRef } from 'react'

export interface ReturnFocus {
	remember: () => void
	restore: (event: Event) => void
}

// Radix gives focus back to a Dialog.Trigger, but the console's dialogs open
// from controls outside their Dialog. Passed to DialogContent's
// onOpenAutoFocus and onCloseAutoFocus, these remember the element that had
// focus as the dialog opened and give focus back to it as the dialog closes.
// A route change that closes a dialog focuses the page instead, and a dialog
// whose content Radix re-creates, when it turns modal or back, keeps focus.
export function useReturnFocus(): ReturnFocus {
	const trigger = useRef<HTMLElement | null>(null)
	const openedAt = useRef('')

	function remember() {
		if (trigger.current) return

		const focused = document.activeElement

		if (!(focused instanceof HTMLElement)) return

		trigger.current = focused
		openedAt.current = window.location.pathname
	}

	function restore(event: Event) {
		event.preventDefault()

		const stillOpen = document.activeElement?.closest('[role="dialog"]')

		if (stillOpen) return

		const element = trigger.current
		const navigated = window.location.pathname !== openedAt.current

		trigger.current = null

		if (navigated) return
		if (!element?.isConnected) return

		element.focus()
	}

	return {
		remember,
		restore,
	}
}
