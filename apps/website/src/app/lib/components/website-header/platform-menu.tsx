'use client'

import { cn } from '@langdrift/web/utilities'
import { useTranslations } from 'next-intl'
import type { FocusEvent } from 'react'
import { useEffect, useId, useRef, useState } from 'react'

import { Link } from '@i18n/navigation'

import { platformProducts } from './platform-products'
import { panelSurface } from './surface'

const item = 'block w-full rounded-xl px-3 py-2.5 text-left'

// The mouse opens the panel by hovering, through CSS alone. A click pins it,
// keyboard focus shows it, Escape dismisses it until the focus or the pointer
// leaves the trigger, and a press outside unpins it.
export function PlatformMenu() {
	const [open, setOpen] = useState(false)
	const [focusWithin, setFocusWithin] = useState(false)
	const [dismissed, setDismissed] = useState(false)
	const panelId = useId()
	const rootRef = useRef<HTMLDivElement>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)

	const translate = useTranslations('header')

	const visible = (open || focusWithin) && !dismissed

	useEffect(() => {
		if (!visible) return

		function unpinOnOutsidePress(event: PointerEvent) {
			const inside =
				event.target instanceof Node && rootRef.current?.contains(event.target)

			if (inside) return

			setOpen(false)
		}

		function dismissOnEscape(event: KeyboardEvent) {
			if (event.key !== 'Escape') return

			setOpen(false)
			setDismissed(true)
			triggerRef.current?.focus()
		}

		document.addEventListener('pointerdown', unpinOnOutsidePress)
		document.addEventListener('keydown', dismissOnEscape)

		return () => {
			document.removeEventListener('pointerdown', unpinOnOutsidePress)
			document.removeEventListener('keydown', dismissOnEscape)
		}
	}, [visible])

	function handleToggle() {
		setDismissed(false)
		setOpen((value) => !value)
	}

	function handleFocus() {
		setFocusWithin(true)
	}

	function handlePointerEnter() {
		setDismissed(false)
	}

	function handleItemClick() {
		setOpen(false)
	}

	function trackFocusLeaving(event: FocusEvent<HTMLElement>) {
		if (rootRef.current?.contains(event.relatedTarget)) return

		setFocusWithin(false)
		setDismissed(false)
	}

	function panelVisibility() {
		if (visible) return 'block'
		if (dismissed) return 'hidden'

		return 'hidden group-hover:block'
	}

	return (
		<div className="group relative" ref={rootRef}>
			<button
				aria-controls={panelId}
				aria-expanded={visible}
				className="inline-flex min-h-11 items-center gap-1 whitespace-nowrap text-muted-foreground group-hover:text-foreground aria-expanded:text-foreground focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
				onBlur={trackFocusLeaving}
				onClick={handleToggle}
				onFocus={handleFocus}
				onPointerEnter={handlePointerEnter}
				ref={triggerRef}
				type="button"
			>
				{translate('platform')}
				<svg
					aria-hidden="true"
					className="size-3.5 transition-transform group-hover:rotate-180"
					fill="none"
					viewBox="0 0 24 24"
				>
					<path
						d="m6 9 6 6 6-6"
						stroke="currentColor"
						strokeLinecap="round"
						strokeLinejoin="round"
						strokeWidth="1.6"
					/>
				</svg>
			</button>
			<div
				className={cn(
					'absolute top-full left-1/2 z-50 w-104 -translate-x-1/2 pt-3',
					panelVisibility(),
				)}
				id={panelId}
			>
				<div className={cn(panelSurface, 'p-2')}>
					<ul aria-label={translate('platformMenu')} className="grid gap-1">
						{platformProducts.map((product) => {
							const Icon = product.icon

							// The header's legacy anchor rule reaches the link inside the
							// panel; one full-width block inside it keeps this layout.
							const copy = (
								<span className="flex w-full items-start gap-3 whitespace-normal">
									<span className="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-lg border border-border bg-card text-foreground">
										<Icon className="size-4" />
									</span>
									<span className="flex min-w-0 flex-col gap-1">
										<span className="text-sm font-medium text-foreground">
											{translate(product.label)}
										</span>
										<span className="text-xs leading-relaxed font-normal text-muted-foreground">
											{translate(product.description)}
										</span>
									</span>
								</span>
							)

							return (
								<li key={product.label}>
									{'href' in product ? (
										<Link
											className={cn(
												item,
												'hover:bg-card focus-visible:bg-card focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
											)}
											href={product.href}
											onBlur={trackFocusLeaving}
											onClick={handleItemClick}
										>
											{copy}
										</Link>
									) : (
										<div className={item}>{copy}</div>
									)}
								</li>
							)
						})}
					</ul>
				</div>
			</div>
		</div>
	)
}
