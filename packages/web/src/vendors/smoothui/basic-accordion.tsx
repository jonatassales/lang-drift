'use client'

import { CaretDown } from '@phosphor-icons/react'
import type { ReactNode } from 'react'
import { useId, useState } from 'react'

import './basic-accordion.css'

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
export function BasicAccordion({
	items,
	allowMultiple = false,
	defaultExpandedIds = [],
	headingLevel = 5,
	layoutId,
	className = '',
}: BasicAccordionProps) {
	const generatedId = useId()
	const prefix = `${layoutId ?? 'smooth-accordion'}-${generatedId}`
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

	return (
		<div className={`smooth-basic-accordion ${className}`}>
			{items.map((item) => {
				const expanded = expandedIds.includes(item.id)
				const triggerId = `${prefix}-trigger-${item.id}`
				const regionId = `${prefix}-region-${item.id}`

				return (
					<div className="smooth-accordion-item" key={item.id}>
						<Heading className="smooth-accordion-heading">
							<button
								aria-controls={regionId}
								aria-expanded={expanded}
								className="smooth-accordion-trigger"
								id={triggerId}
								onClick={() => toggleItem(item.id)}
								type="button"
							>
								<span>{item.title}</span>
								<span
									aria-hidden="true"
									className="smooth-accordion-chevron"
									data-expanded={expanded}
								>
									<CaretDown size={16} />
								</span>
							</button>
						</Heading>
						<section
							aria-hidden={!expanded}
							aria-labelledby={triggerId}
							className="smooth-accordion-region"
							hidden={!expanded}
							id={regionId}
						>
							<div className="smooth-accordion-content">{item.content}</div>
						</section>
					</div>
				)
			})}
		</div>
	)
}
