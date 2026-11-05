'use client'

import type { FocusEvent, KeyboardEvent, MouseEvent, ReactNode } from 'react'
import { useId, useRef, useState } from 'react'

import { cn } from 'web/utilities'

export type AnimatedTabsVariant = 'underline' | 'pill' | 'segment'

export interface AnimatedTab {
	id: string
	label: string
	icon?: ReactNode
}

export interface AnimatedTabsProps {
	tabs: readonly AnimatedTab[]
	ariaLabel: string
	activeTab?: string
	defaultTab?: string
	onChange?: (tabId: string) => void
	variant?: AnimatedTabsVariant
	layoutId?: string
	className?: string
}

// The list: an underline sits on a bottom border, a pill or a segment on a
// rounded border filled with the background role.
const listVariants: Record<AnimatedTabsVariant, string> = {
	underline: 'border-b',
	pill: 'rounded-full border bg-background',
	segment: 'gap-0 rounded-lg border bg-background',
}

// Each tab's corners follow the list it sits in.
const tabVariants: Record<AnimatedTabsVariant, string> = {
	underline: 'rounded-t-md',
	pill: 'rounded-full',
	segment: 'rounded-md',
}

// The selected tab's indicator: a 2px primary bar under an underline tab, a
// card-filled outline behind a pill or a segment tab.
const indicatorVariants: Record<AnimatedTabsVariant, string> = {
	underline: 'inset-x-0 bottom-0 h-0.5 bg-primary',
	pill: 'inset-0 rounded-[inherit] border border-border bg-card',
	segment: 'inset-0 rounded-[inherit] border border-border bg-card',
}

/** SmoothUI Animated Tabs, with localized naming and explicit panel relations. */
export function AnimatedTabs(props: AnimatedTabsProps) {
	const {
		activeTab: controlledActiveTab,
		ariaLabel,
		className = '',
		defaultTab,
		layoutId: customLayoutId,
		onChange,
		tabs,
		variant = 'underline',
	} = props

	const generatedId = useId()
	const layoutId = customLayoutId ?? `animated-tabs-${generatedId}`
	const [internalActiveTab, setInternalActiveTab] = useState(
		defaultTab ?? tabs[0]?.id ?? '',
	)
	const controlled = controlledActiveTab !== undefined
	const requestedTab = controlled ? controlledActiveTab : internalActiveTab
	const activeTab = tabs.some((tab) => tab.id === requestedTab)
		? requestedTab
		: tabs[0]?.id
	const containerRef = useRef<HTMLDivElement>(null)
	const tabRefs = useRef(new Map<string, HTMLButtonElement>())

	function selectTab(id: string) {
		if (id === activeTab) return

		if (!controlled) setInternalActiveTab(id)
		onChange?.(id)
	}

	function revealTab(button: HTMLButtonElement) {
		const container = containerRef.current
		if (!container) return

		const tabBounds = button.getBoundingClientRect()
		const bounds = container.getBoundingClientRect()
		if (tabBounds.left < bounds.left + 4) {
			container.scrollLeft += tabBounds.left - bounds.left - 4
		} else if (tabBounds.right > bounds.right - 4) {
			container.scrollLeft += tabBounds.right - bounds.right + 4
		}
	}

	function handleTabClick(event: MouseEvent<HTMLButtonElement>) {
		const tabId = event.currentTarget.dataset.tabId
		if (tabId) {
			selectTab(tabId)
		}
	}

	function handleTabFocus(event: FocusEvent<HTMLButtonElement>) {
		revealTab(event.currentTarget)
	}

	function handleTabKeyDown(event: KeyboardEvent<HTMLButtonElement>) {
		const indexValue = event.currentTarget.dataset.tabIndex
		if (indexValue === undefined) return
		const index = Number(indexValue)

		let nextIndex: number
		switch (event.key) {
			case 'ArrowRight':
				nextIndex = (index + 1) % tabs.length
				break
			case 'ArrowLeft':
				nextIndex = (index - 1 + tabs.length) % tabs.length
				break
			case 'Home':
				nextIndex = 0
				break
			case 'End':
				nextIndex = tabs.length - 1
				break
			default:
				return
		}

		event.preventDefault()
		const tab = tabs[nextIndex]
		if (!tab) return

		selectTab(tab.id)
		tabRefs.current.get(tab.id)?.focus({ preventScroll: true })
	}

	return (
		<div
			aria-label={ariaLabel}
			className={cn(
				'relative box-border inline-flex max-w-full min-w-0 gap-1 overflow-x-auto overscroll-x-contain border-border p-1 align-top [scrollbar-width:thin]',
				listVariants[variant],
				className,
			)}
			data-variant={variant}
			ref={containerRef}
			role="tablist"
		>
			{tabs.map((tab, index) => {
				const selected = activeTab === tab.id

				return (
					<button
						aria-controls={`${layoutId}-panel-${tab.id}`}
						aria-selected={selected}
						className={cn(
							'relative inline-flex min-h-11 flex-none cursor-pointer items-center justify-center gap-2 border-0 bg-transparent px-4 py-2.5 font-medium whitespace-nowrap text-muted-foreground transition-colors duration-150 hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-foreground aria-selected:text-foreground motion-reduce:transition-none',
							tabVariants[variant],
						)}
						data-tab-id={tab.id}
						data-tab-index={index}
						id={`${layoutId}-tab-${tab.id}`}
						key={tab.id}
						onClick={handleTabClick}
						onFocus={handleTabFocus}
						onKeyDown={handleTabKeyDown}
						ref={(button) => {
							if (button) tabRefs.current.set(tab.id, button)
							else tabRefs.current.delete(tab.id)
						}}
						role="tab"
						tabIndex={selected ? 0 : -1}
						type="button"
					>
						{selected ? (
							<span
								aria-hidden="true"
								className={cn(
									'pointer-events-none absolute transition-opacity duration-150 motion-reduce:transition-none',
									indicatorVariants[variant],
								)}
							/>
						) : null}
						{tab.icon ? (
							<span
								aria-hidden="true"
								className="relative inline-flex items-center"
							>
								{tab.icon}
							</span>
						) : null}
						<span className="relative">{tab.label}</span>
					</button>
				)
			})}
		</div>
	)
}
