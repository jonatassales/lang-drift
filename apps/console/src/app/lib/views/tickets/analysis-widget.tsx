import type { Ticket, TicketAnalysis } from '@data/types'
import { Eyebrow } from 'web/ui/eyebrow'

import { DocumentRow } from './document-row'
import { EpisodeRow } from './episode-row'
import { extractedFact } from './extracted-fact.compute'
import { measuredLine } from './measured-line.compute'
import { PhaseSummary } from './phase-summary'
import { useTicketFormat } from './ticket-format.hook'
import { ticketDocumentIds } from './ticket-phases'
import { TicketSection } from './ticket-section'
import { UnknownValue } from './unknown-value'
import { verdictLabel } from './verdict-label.compute'
import { WidgetSection } from './widget-section'

interface AnalysisWidgetProps {
	analysis: TicketAnalysis
	ticket: Ticket
}

const aside = 'font-mono text-[11px] text-muted-foreground'

// The analysis as the website's Analysis panel draws it: each linked
// document with what its read extracted, the episodes ranked by similarity
// and judged against the threshold, the resources they measured, and the
// model and effort the agent suggests, which the team authorizes.
export function AnalysisWidget(props: AnalysisWidgetProps) {
	const { analysis, ticket } = props

	const format = useTicketFormat()
	const { number, translate } = format
	const guardrailRead = analysis.documentsRead.includes('guardrail')

	return (
		<TicketSection
			id="analysis"
			summary={<PhaseSummary phase="analysis" ticket={ticket} />}
			title={translate('labels.analysis')}
		>
			<div className="flex items-baseline justify-between gap-3">
				<Eyebrow>{translate('analysis.documents')}</Eyebrow>
				<span className={aside}>
					{translate('analysis.documentsSummary', {
						count: analysis.documentsRead.length,
					})}
				</span>
			</div>
			<ul className="m-0 mt-1 list-none p-0">
				{ticketDocumentIds.map((id) => (
					<DocumentRow
						extracted={extractedFact(id, analysis, format)}
						file={ticket.artifacts[id]}
						judgment={
							id === 'guardrail' && guardrailRead
								? translate('analysis.guardrailJudgment')
								: undefined
						}
						key={id}
					/>
				))}
			</ul>

			<WidgetSection>
				<div className="flex items-baseline justify-between gap-3">
					<Eyebrow>{translate('analysis.episodes')}</Eyebrow>
					<span className={aside}>
						{translate('analysis.threshold', {
							value: number.format(analysis.threshold),
						})}
					</span>
				</div>
				{!analysis.episodes ? (
					<UnknownValue className="mt-2 block font-mono text-[12px]" />
				) : (
					<ul className="m-0 mt-1 list-none p-0">
						{analysis.episodes.map((episode) => (
							<EpisodeRow
								episode={episode}
								key={episode.id}
								similarity={number.format(episode.similarity)}
								verdict={verdictLabel(episode, format)}
							/>
						))}
					</ul>
				)}
				<div className="mt-1 flex flex-wrap items-baseline justify-between gap-x-3 border-t border-border pt-2">
					<Eyebrow>{translate('analysis.measured')}</Eyebrow>
					<span className="flex flex-wrap items-baseline gap-x-2">
						{measuredLine(analysis, format)}
					</span>
				</div>
			</WidgetSection>

			<WidgetSection className="grid grid-cols-2 gap-x-6 gap-y-2">
				<div className="grid gap-1.5">
					<Eyebrow>{translate('analysis.suggestedModel')}</Eyebrow>
					{analysis.model ? (
						<span className="grid gap-0.5">
							<span className="text-[14px]">{analysis.model.name}</span>
							<span className="font-mono text-[11px] text-muted-foreground">
								{analysis.model.id}
							</span>
						</span>
					) : (
						<UnknownValue className="font-mono text-[14px]" />
					)}
				</div>
				<div className="grid gap-1.5">
					<Eyebrow>{translate('analysis.suggestedEffort')}</Eyebrow>
					{analysis.effort ? (
						<span className="font-mono text-[14px]">
							{translate(`analysis.efforts.${analysis.effort}`)}
						</span>
					) : (
						<UnknownValue className="font-mono text-[14px]" />
					)}
				</div>
				{analysis.model && (
					<p className="col-span-2 m-0 text-[11px] text-muted-foreground">
						{translate('analysis.recommendationNote')}
					</p>
				)}
			</WidgetSection>
		</TicketSection>
	)
}
