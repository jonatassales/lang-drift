'use client'

import { CaretDown } from '@phosphor-icons/react'
import type { MouseEvent, ReactNode } from 'react'
import { useId, useState } from 'react'

import { cn } from 'web/utilities'

export interface AccordionItem {
	id: string | number
	title: ReactNode
	content: ReactNode
}

export interface BasicAccordionProps {
	items: readonly AccordionItem[]
	allowMultiple?: boolean
	defaultExpandedIds?: readonly (string | number)[]
	headingLevel?: 2 | 3 | 4 | 5 | 6
	layoutId?: string
	className?: string
}

/** SmoothUI Basic Accordion, adapted for heading order and unique ARIA IDs. */
export function BasicAccordion(props: BasicAccordionProps) {
	const {
		allowMultiple = false,
		className = '',
		defaultExpandedIds = [],
		headingLevel = 5,
		items,
		layoutId,
	} = props

	const generatedId = useId()
	const prefix = `${layoutId ?? 'accordion'}-${generatedId}`
	const [expandedIds, setExpandedIds] = useState<Array<string | number>>(() =>
		allowMultiple
			? [...defaultExpandedIds]
			: [...defaultExpandedIds.slice(0, 1)],
	)
	const Heading = `h${headingLevel}` as 'h2' | 'h3' | 'h4' | 'h5' | 'h6'

	function toggleItem(id: string | number) {
		setExpandedIds((current) =>
			current.includes(id)
				? current.filter((expandedId) => expandedId !== id)
				: allowMultiple
					? [...current, id]
					: [id],
		)
	}

	function handleTriggerClick(event: MouseEvent<HTMLButtonElement>) {
		const itemIndex = event.currentTarget.dataset.itemIndex
		if (itemIndex === undefined) return

		const item = items[Number(itemIndex)]
		if (item) {
			toggleItem(item.id)
		}
	}

	return (
		<div className={cn('min-w-0', className)}>
			{items.map((item, index) => {
				const expanded = expandedIds.includes(item.id)
				const triggerId = `${prefix}-trigger-${item.id}`
				const regionId = `${prefix}-region-${item.id}`

				return (
					<div className="border-border not-first:border-t" key={item.id}>
						<Heading>
							<button
								aria-controls={regionId}
								aria-expanded={expanded}
								className="flex min-h-11 w-full cursor-pointer items-center justify-between gap-3 rounded-sm border-0 bg-transparent py-2.5 text-start font-medium text-foreground hover:underline hover:underline-offset-4 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-foreground"
								data-item-index={index}
								id={triggerId}
								onClick={handleTriggerClick}
								type="button"
							>
								<span className="min-w-0 wrap-anywhere">{item.title}</span>
								<span
									aria-hidden="true"
									className="inline-flex flex-none items-center transition-transform duration-200 data-[expanded=true]:rotate-180 motion-reduce:transition-none"
									data-expanded={expanded}
								>
									<CaretDown size={16} />
								</span>
							</button>
						</Heading>
						<section
							aria-hidden={!expanded}
							aria-labelledby={triggerId}
							className="overflow-hidden"
							hidden={!expanded}
							id={regionId}
						>
							<div className="pt-1 pb-4 wrap-anywhere">{item.content}</div>
						</section>
					</div>
				)
			})}
		</div>
	)
}
