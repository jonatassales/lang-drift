import { Eyebrow } from '@langdrift/web/ui/eyebrow'

import {
	revealed,
	suggestedTotals,
	ticketPlan,
} from '@data/microruntimes-content'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import { ShareBar } from './share-bar'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

interface HardwareTableProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

const headCell = 'pb-2 text-left font-normal'
const rowHead = 'py-2 pr-3 text-left font-normal'
const note = 'block font-sans text-[11px] text-muted-foreground'

// The hardware the plan is bounded by and the hardware it suggests: CPU and
// memory with the headroom left in the envelope, and GPU and disk, which
// runtime.md does not bound, without a number.
export function HardwareTable(props: HardwareTableProps) {
	const { copy, format, playedCount } = props

	const bounded = ticketPlan.bounded
	const suggested = suggestedTotals(ticketPlan.components)
	const boundedState = revealState(revealed(playedCount, 'bounded'))
	const suggestedState = revealState(revealed(playedCount, 'suggested'))
	const limitsKnown = revealed(playedCount, 'limits')
	const limitsState = revealState(limitsKnown)
	const resourceRows = [
		{
			id: 'cpu',
			label: copy.hardware.cpu,
			bounded: format.cpu(bounded.cpuMillis),
			suggested: format.cpu(suggested.cpuMillis),
			share: suggested.cpuMillis / bounded.cpuMillis,
		},
		{
			id: 'memory',
			label: copy.hardware.memory,
			bounded: format.memory(bounded.memoryMiB),
			suggested: format.memory(suggested.memoryMiB),
			share: suggested.memoryMiB / bounded.memoryMiB,
		},
	]
	const limitRows = [
		{
			id: 'gpu',
			label: copy.hardware.gpu,
			value: copy.none,
		},
		{
			id: 'disk',
			label: copy.hardware.disk,
			value: (
				<>
					<span aria-hidden="true">—</span>
					<span className="sr-only">{copy.none}</span>
				</>
			),
		},
	]

	return (
		<CardSection>
			<table className="w-full border-collapse text-[12px]">
				<thead>
					<tr>
						<th className={headCell} scope="col">
							<Eyebrow>{copy.hardware.resource}</Eyebrow>
						</th>
						<th className={headCell} scope="col">
							<Eyebrow>{copy.hardware.bounded}</Eyebrow>
						</th>
						<th className={headCell} scope="col">
							<Eyebrow>{copy.hardware.suggested}</Eyebrow>
						</th>
					</tr>
				</thead>
				<tbody>
					{resourceRows.map((row) => (
						<tr className="border-t border-border" key={row.id}>
							<th className={rowHead} scope="row">
								{row.label}
							</th>
							<td className="py-2 pr-3 font-mono tabular-nums">
								<CardValue labels={copy} state={boundedState}>
									{row.bounded}
								</CardValue>
							</td>
							<td className="py-2 font-mono tabular-nums">
								<CardValue labels={copy} state={suggestedState}>
									{row.suggested}
								</CardValue>
								<ShareBar
									className="mt-1"
									share={suggestedState === 'set' ? row.share : 0}
								/>
							</td>
						</tr>
					))}
					{limitRows.map((row) => (
						<tr className="border-t border-border" key={row.id}>
							<th className={rowHead} scope="row">
								{row.label}
							</th>
							<td className="py-2 pr-3 font-mono">
								<CardValue labels={copy} state={limitsState}>
									{row.value}
								</CardValue>
							</td>
							<td className="py-2 font-mono">
								<CardValue labels={copy} state={limitsState}>
									{row.value}
								</CardValue>
								{limitsKnown && (
									<span className={note}>{copy.notRuntimeField}</span>
								)}
							</td>
						</tr>
					))}
				</tbody>
			</table>
			{revealed(playedCount, 'priceBasis') && (
				<p className="mt-2 mb-0 text-[11px] text-muted-foreground">
					{copy.priceBasis}
				</p>
			)}
		</CardSection>
	)
}
