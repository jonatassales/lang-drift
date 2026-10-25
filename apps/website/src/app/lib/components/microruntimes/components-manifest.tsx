import { Eyebrow } from '@langdrift/web/ui/eyebrow'

import {
	hasPlayed,
	revealed,
	suggestedTotals,
	ticketPlan,
} from '@data/microruntimes-content'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import { ComponentRow } from './component-row'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

interface ComponentsManifestProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

// The components the plan resolves from the spec: a dash until the console
// starts, five skeleton rows while it resolves, then one row per terminal
// line, under a total that counts and sums the rows as they arrive.
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
				<Eyebrow>{copy.components}</Eyebrow>
				<Eyebrow>{copy.columns.cpu}</Eyebrow>
				<Eyebrow>{copy.columns.memory}</Eyebrow>
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
								kind={copy.kinds[component.kind]}
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
							? format.message('componentsTotal', {
									count: resolved.length,
									cpu: format.formatNumber.format(resolvedTotals.cpuMillis),
									memory: format.formatNumber.format(resolvedTotals.memoryMiB),
								})
							: format.message('resolvedProgress', {
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
