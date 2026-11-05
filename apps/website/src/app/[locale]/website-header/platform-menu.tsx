'use client'

import { useTranslations } from 'next-intl'
import type { FocusEvent } from 'react'
import { useEffect, useId, useRef, useState } from 'react'

import { Link } from '@i18n/navigation'
import { ProductBadge } from '@layout/product-badge'
import { cn } from 'web/utilities'

import { panelSurface } from './panel-surface.data'
import { platformProducts } from './platform-products.data'

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

	const visible = (open || focusWithin) && !dismissed
	const hoverVisibility = dismissed ? 'hidden' : 'hidden group-hover:block'
	const panelVisibility = visible ? 'block' : hoverVisibility

	// Dismiss pinned platform menu on outside press or Escape key.
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
					'absolute top-full left-1/2 z-50 w-176 max-w-[calc(100vw-2rem)] -translate-x-1/2 pt-3',
					panelVisibility,
				)}
				id={panelId}
			>
				<div className={cn(panelSurface, 'p-2')}>
					<ul
						aria-label={translate('platformMenu')}
						className="relative grid gap-1 sm:grid-cols-2 sm:gap-x-4 sm:before:pointer-events-none sm:before:absolute sm:before:inset-y-4 sm:before:left-1/2 sm:before:w-px sm:before:-translate-x-1/2 sm:before:bg-linear-to-b sm:before:from-transparent sm:before:via-border sm:before:to-transparent"
					>
						{platformProducts.map((product) => {
							const Icon = product.icon

							// The link reads like the desktop navigation's links around it;
							// one full-width block inside it keeps this layout.
							const copy = (
								<span className="flex w-full items-start gap-3 whitespace-normal">
									<span className="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-lg border border-border bg-card text-foreground">
										<Icon className="size-4" />
									</span>
									<span className="flex min-w-0 flex-col gap-1">
										<span className="flex items-center gap-2">
											<span className="text-sm font-medium text-foreground">
												{translate(product.label)}
											</span>
											{'badge' in product && product.badge ? (
												<ProductBadge>{translate(product.badge)}</ProductBadge>
											) : null}
										</span>
										<span className="text-xs leading-relaxed font-normal text-muted-foreground">
											{translate(product.description)}
										</span>
									</span>
								</span>
							)

							return (
								<li key={product.label}>
									<Link
										className={cn(
											item,
											'inline-flex min-h-11 items-center text-[13px] font-medium whitespace-nowrap text-muted-foreground hover:bg-card hover:text-foreground focus-visible:bg-card focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
										)}
										href={product.href}
										onBlur={trackFocusLeaving}
										onClick={handleItemClick}
									>
										{copy}
									</Link>
								</li>
							)
						})}
					</ul>
				</div>
			</div>
		</div>
	)
}
