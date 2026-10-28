import { useEffect, useRef, useState } from 'react'

// A menu that opens from a trigger and closes on Escape, on a press outside
// it or when focus leaves it. Escape returns focus to the trigger.
export function useDismissableMenu() {
	const [open, setOpen] = useState(false)
	const containerRef = useRef<HTMLDivElement>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)

	useEffect(() => {
		if (!open) return

		function dismiss(event: Event) {
			if (containerRef.current?.contains(event.target as Node)) return
			setOpen(false)
		}
		function handleEscape(event: KeyboardEvent) {
			if (event.key !== 'Escape') return
			setOpen(false)
			triggerRef.current?.focus()
		}
		document.addEventListener('pointerdown', dismiss)
		document.addEventListener('focusin', dismiss)
		document.addEventListener('keydown', handleEscape)

		return () => {
			document.removeEventListener('pointerdown', dismiss)
			document.removeEventListener('focusin', dismiss)
			document.removeEventListener('keydown', handleEscape)
		}
	}, [open])

	return {
		containerRef,
		open,
		setOpen,
		triggerRef,
	}
}
