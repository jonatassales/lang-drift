'use client'

import type { FocusEvent, KeyboardEvent, MouseEvent, ReactNode } from 'react'
import { useId, useRef, useState } from 'react'

import './animated-tabs.css'

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
	variant?: 'underline' | 'pill' | 'segment'
	layoutId?: string
	className?: string
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
	const layoutId = customLayoutId ?? `smooth-tabs-${generatedId}`
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
			className={`smooth-animated-tabs ${className}`}
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
						className="smooth-animated-tab"
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
							<span aria-hidden="true" className="smooth-tab-indicator" />
						) : null}
						{tab.icon ? (
							<span aria-hidden="true" className="smooth-tab-icon">
								{tab.icon}
							</span>
						) : null}
						<span className="smooth-tab-label">{tab.label}</span>
					</button>
				)
			})}
		</div>
	)
}
