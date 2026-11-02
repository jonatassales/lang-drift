'use client'

import { AnimatedTabs } from '@langdrift/web/vendors/smoothui/animated-tabs'
import { BasicAccordion } from '@langdrift/web/vendors/smoothui/basic-accordion'
import { useId, useState } from 'react'

import type { ProductFaqItem } from '@components/product-page/product-faq-section'

interface FaqGroup {
	id: string
	items: readonly ProductFaqItem[]
	label: string
}

interface FaqTabsProps {
	groups: readonly FaqGroup[]
	tabsLabel: string
}

// One accordion switched between the products' questions. The accordion is
// mounted only inside the active panel, so a switch shows it collapsed.
export function FaqTabs(props: FaqTabsProps) {
	const { groups, tabsLabel } = props

	const layoutId = useId()
	const [activeId, setActiveId] = useState(groups[0]?.id ?? '')

	function handleTabChange(tabId: string) {
		setActiveId(tabId)
	}

	return (
		<div className="mt-10">
			<AnimatedTabs
				activeTab={activeId}
				ariaLabel={tabsLabel}
				layoutId={layoutId}
				onChange={handleTabChange}
				tabs={groups.map((group) => ({
					id: group.id,
					label: group.label,
				}))}
				variant="underline"
			/>
			{groups.map((group) => {
				const active = group.id === activeId

				return (
					<div
						aria-labelledby={`${layoutId}-tab-${group.id}`}
						hidden={!active}
						id={`${layoutId}-panel-${group.id}`}
						key={group.id}
						role="tabpanel"
						tabIndex={active ? 0 : -1}
					>
						{active ? (
							<BasicAccordion
								className="mt-6"
								headingLevel={3}
								items={group.items.map((item) => ({
									content: (
										<p className="m-0 text-[13px] leading-[1.6] text-muted-foreground">
											{item.answer}
										</p>
									),
									id: item.id,
									title: item.question,
								}))}
								layoutId={`platform-faq-${group.id}`}
							/>
						) : null}
					</div>
				)
			})}
		</div>
	)
}
