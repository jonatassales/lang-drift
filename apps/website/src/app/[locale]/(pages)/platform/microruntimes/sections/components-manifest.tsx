import {
	hasPlayed,
	revealed,
	suggestedTotals,
	ticketPlan,
} from '@routes/platform/microruntimes/content/microruntimes-content'
import { Eyebrow } from 'web/ui/eyebrow'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import { ComponentRow } from './component-row'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

// The components the plan resolves from the spec: a dash until the console
// starts, five skeleton rows while it resolves, then one row per terminal
// line, under a total that counts and sums the rows as they arrive.
interface ComponentsManifestProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

export function ComponentsManifest(props: ComponentsManifestProps) {
	const { copy, format, playedCount } = props

	const resolved = ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
	const resolvedTotals = suggestedTotals(resolved)
	const total = ticketPlan.components.length
	const allResolved = resolved.length === total
	const state = revealState(
		resolved.length > 0,
		revealed(playedCount, 'componentsSkeleton'),
	)

	return (
		<CardSection>
			<div className="grid grid-cols-[minmax(0,1fr)_auto_auto] gap-x-3">
				<Eyebrow>{copy.plan.components}</Eyebrow>
				<Eyebrow>{copy.plan.columns.cpu}</Eyebrow>
				<Eyebrow>{copy.plan.columns.memory}</Eyebrow>
			</div>
			{state === 'unset' ? (
				<CardValue
					className="mt-2 block font-mono text-[12px]"
					labels={copy}
					state="unset"
				/>
			) : (
				<>
					<ul className="m-0 mt-1 list-none p-0">
						{ticketPlan.components.map((component) => (
							<ComponentRow
								component={component}
								cpu={format.cpu(component.cpuMillis)}
								key={component.id}
								kind={copy.plan.kinds[component.kind]}
								labels={copy}
								memory={format.memory(component.memoryMiB)}
								share={component.cpuMillis / ticketPlan.bounded.cpuMillis}
								state={
									hasPlayed(playedCount, component.id) ? 'set' : 'resolving'
								}
							/>
						))}
					</ul>
					<CardValue
						className="mt-2 block border-t border-border pt-2 text-[12px] text-muted-foreground"
						labels={copy}
						state={state}
						width="w-44"
					>
						{allResolved
							? format.message('plan.componentsTotal', {
									count: resolved.length,
									cpu: format.formatNumber.format(resolvedTotals.cpuMillis),
									memory: format.formatNumber.format(resolvedTotals.memoryMiB),
								})
							: format.message('plan.resolvedProgress', {
									resolved: resolved.length,
									total,
									cpu: format.formatNumber.format(resolvedTotals.cpuMillis),
									memory: format.formatNumber.format(resolvedTotals.memoryMiB),
								})}
					</CardValue>
				</>
			)}
		</CardSection>
	)
}
