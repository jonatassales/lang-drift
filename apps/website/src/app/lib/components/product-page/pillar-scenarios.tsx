'use client'

import { AnimatedTabs } from '@langdrift/web/vendors/smoothui/animated-tabs'
import { useId, useState } from 'react'

export interface PillarScenario {
	cells: Record<string, string>
	context: string
	id: string
	tabLabel: string
}

interface PillarScenariosProps {
	pillars: readonly { id: string; label: string }[]
	scenarios: readonly PillarScenario[]
	selectorLabel: string
}

// Scenario tabs over a row of pillars: each scenario tells one story across
// the same columns (intent to outcome, compose to clean up).
export function PillarScenarios(props: PillarScenariosProps) {
	const { pillars, scenarios, selectorLabel } = props

	const layoutId = useId()
	const [activeId, setActiveId] = useState(scenarios[0]?.id ?? '')

	return (
		<div className="mt-11">
			<span className="font-mono text-[10px] tracking-[0.08em] text-muted-foreground uppercase">
				{selectorLabel}
			</span>
			<AnimatedTabs
				activeTab={activeId}
				ariaLabel={selectorLabel}
				className="mt-3"
				layoutId={layoutId}
				onChange={setActiveId}
				tabs={scenarios.map((scenario) => ({
					id: scenario.id,
					label: scenario.tabLabel,
				}))}
				variant="underline"
			/>
			{scenarios.map((scenario) => {
				const active = scenario.id === activeId

				return (
					<div
						aria-labelledby={`${layoutId}-tab-${scenario.id}`}
						hidden={!active}
						id={`${layoutId}-panel-${scenario.id}`}
						key={scenario.id}
						role="tabpanel"
						tabIndex={active ? 0 : -1}
					>
						{active ? (
							<>
								<p className="mt-6 mb-0 max-w-[760px] text-[15px] leading-[1.6] text-foreground">
									{scenario.context}
								</p>
								<ol className="m-0 mt-6 grid list-none gap-3 p-0 lg:grid-cols-4">
									{pillars.map((pillar, index) => (
										<li
											className="relative flex flex-col gap-3 rounded-2xl border border-border bg-card p-5"
											key={pillar.id}
										>
											<div className="flex items-center justify-between gap-3">
												<h3 className="m-0 flex items-center gap-2.5 text-base font-medium tracking-[-0.01em]">
													<span className="inline-flex size-6 items-center justify-center rounded-full border border-border font-mono text-[10px] text-muted-foreground">
														{index + 1}
													</span>
													{pillar.label}
												</h3>
												{index < pillars.length - 1 ? (
													<span
														aria-hidden="true"
														className="hidden text-muted-foreground lg:inline"
													>
														→
													</span>
												) : null}
											</div>
											<p className="m-0 text-[13px] leading-[1.6] text-muted-foreground">
												{scenario.cells[pillar.id]}
											</p>
										</li>
									))}
								</ol>
							</>
						) : null}
					</div>
				)
			})}
		</div>
	)
}
