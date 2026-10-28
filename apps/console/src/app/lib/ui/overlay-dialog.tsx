'use client'

import type { ReactNode } from 'react'
import { useEffectEvent, useLayoutEffect, useRef } from 'react'

interface OverlayDialogProps {
	children: ReactNode
	className: string
	id?: string
	label?: string
	labelledBy?: string
	mobileOnly?: boolean
	modal?: boolean
	onClose: () => void
}

export function OverlayDialog(props: OverlayDialogProps) {
	const {
		children,
		className,
		id,
		label,
		labelledBy,
		mobileOnly = false,
		modal = false,
		onClose,
	} = props

	const dialogRef = useRef<HTMLDialogElement>(null)
	const modalRef = useRef(false)
	const close = useEffectEvent(onClose)

	useLayoutEffect(() => {
		const dialog = dialogRef.current
		if (!dialog) return
		const trigger =
			document.activeElement instanceof HTMLElement
				? document.activeElement
				: null
		const media = window.matchMedia('(max-width: 639px)')
		const previousOverflow = document.body.style.overflow

		function present() {
			if (!dialog) return
			if (mobileOnly && !media.matches) {
				close()

				return
			}
			dialog.close()
			modalRef.current = modal || media.matches
			if (modalRef.current) {
				document.body.style.overflow = 'hidden'
				dialog.showModal()
			} else {
				document.body.style.overflow = previousOverflow
				dialog.show()
			}
			dialog.querySelector<HTMLButtonElement>('button')?.focus()
		}

		function handleEscape(event: KeyboardEvent) {
			if (event.key === 'Escape' && !modalRef.current) {
				event.preventDefault()
				close()
			}
		}

		present()
		media.addEventListener('change', present)
		window.addEventListener('keydown', handleEscape)

		return () => {
			media.removeEventListener('change', present)
			window.removeEventListener('keydown', handleEscape)
			dialog.close()
			document.body.style.overflow = previousOverflow
			if (trigger?.isConnected) trigger.focus()
		}
	}, [mobileOnly, modal])

	return (
		<dialog
			aria-label={label}
			aria-labelledby={labelledBy}
			className={className}
			id={id}
			onCancel={(event) => {
				event.preventDefault()
				onClose()
			}}
			onClick={(event) => {
				if (event.target !== event.currentTarget || !modalRef.current) return
				const bounds = event.currentTarget.getBoundingClientRect()
				if (
					event.clientX < bounds.left ||
					event.clientX > bounds.right ||
					event.clientY < bounds.top ||
					event.clientY > bounds.bottom
				)
					onClose()
			}}
			onKeyDown={(event) => {
				if (event.key !== 'Tab' || !modalRef.current) return
				const controls = [
					...event.currentTarget.querySelectorAll<HTMLElement>(
						'button:not(:disabled), a[href], input:not(:disabled), [tabindex="0"]',
					),
				].filter((element) => element.getClientRects().length > 0)
				const first = controls[0]
				const last = controls.at(-1)
				if (event.shiftKey && document.activeElement === first) {
					event.preventDefault()
					last?.focus()
				} else if (!event.shiftKey && document.activeElement === last) {
					event.preventDefault()
					first?.focus()
				}
			}}
			ref={dialogRef}
		>
			{children}
		</dialog>
	)
}
