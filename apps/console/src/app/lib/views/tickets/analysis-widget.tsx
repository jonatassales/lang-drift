import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import type { ReactNode } from 'react'

import type {
	Ticket,
	TicketAnalysis,
	TicketDocumentId,
	TicketEpisode,
} from '@data/types'

import { DocumentRow } from './document-row'
import { EpisodeRow } from './episode-row'
import { PhaseSummary } from './phase-summary'
import { type TicketFormat, useTicketFormat } from './ticket-format.hook'
import { ticketDocumentIds } from './ticket-phases'
import { TicketSection } from './ticket-section'
import { UnknownValue } from './unknown-value'
import { WidgetSection } from './widget-section'

interface AnalysisWidgetProps {
	analysis: TicketAnalysis
	ticket: Ticket
}

const aside = 'font-mono text-[11px] text-muted-foreground'

// What the read of each document extracted, or null while it is unread or
// the fact is not known yet.
function extractedFact(
	id: TicketDocumentId,
	analysis: TicketAnalysis,
	format: TicketFormat,
) {
	if (!analysis.documentsRead.includes(id)) return null

	const { number, t } = format
	switch (id) {
		case 'prd':
			return analysis.useCase
				? t('analysis.extracted.prd', { useCase: analysis.useCase })
				: null
		case 'designDoc':
			if (!analysis.bounded || analysis.budgetCents === null) return null
			return t('analysis.extracted.designDoc', {
				budget: format.money(analysis.budgetCents),
				cpu: number.format(analysis.bounded.cpuMillis),
				memory: number.format(analysis.bounded.memoryMiB),
			})
		case 'spec':
			if (
				analysis.specComponents === null ||
				analysis.acceptanceCriteria === null
			)
				return null
			return t('analysis.extracted.spec', {
				components: analysis.specComponents,
				criteria: analysis.acceptanceCriteria,
			})
		case 'guardrail':
			if (!analysis.criteria) return null
			return t('analysis.extracted.guardrail', {
				criteria: analysis.criteria.join(' · '),
			})
	}
}

function verdictLabel(episode: TicketEpisode, format: TicketFormat) {
	const { t } = format
	const label = t(`analysis.verdicts.${episode.verdict}`)
	if (!episode.reason) return label

	return `${label} · ${t(`analysis.reasons.${episode.reason}`)}`
}

// The analysis as the website's Analysis panel draws it: each linked
// document with what its read extracted, the episodes ranked by similarity
// and judged against the threshold, the resources they measured, and the
// model and effort the agent suggests, which the team authorizes.
export function AnalysisWidget(props: AnalysisWidgetProps) {
	const { analysis, ticket } = props

	const format = useTicketFormat()
	const { number, t } = format
	const guardrailRead = analysis.documentsRead.includes('guardrail')

	return (
		<TicketSection
			id="analysis"
			summary={<PhaseSummary phase="analysis" ticket={ticket} />}
			title={t('labels.analysis')}
		>
			<div className="flex items-baseline justify-between gap-3">
				<Eyebrow>{t('analysis.documents')}</Eyebrow>
				<span className={aside}>
					{t('analysis.documentsSummary', {
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
								? t('analysis.guardrailJudgment')
								: undefined
						}
						key={id}
					/>
				))}
			</ul>

			<WidgetSection>
				<div className="flex items-baseline justify-between gap-3">
					<Eyebrow>{t('analysis.episodes')}</Eyebrow>
					<span className={aside}>
						{t('analysis.threshold', {
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
					<Eyebrow>{t('analysis.measured')}</Eyebrow>
					<span className="flex flex-wrap items-baseline gap-x-2">
						{measuredLine(analysis, format)}
					</span>
				</div>
			</WidgetSection>

			<WidgetSection className="grid grid-cols-2 gap-x-6 gap-y-2">
				<div className="grid gap-1.5">
					<Eyebrow>{t('analysis.suggestedModel')}</Eyebrow>
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
					<Eyebrow>{t('analysis.suggestedEffort')}</Eyebrow>
					{analysis.effort ? (
						<span className="font-mono text-[14px]">
							{t(`analysis.efforts.${analysis.effort}`)}
						</span>
					) : (
						<UnknownValue className="font-mono text-[14px]" />
					)}
				</div>
				{analysis.model && (
					<p className="col-span-2 m-0 text-[11px] text-muted-foreground">
						{t('analysis.recommendationNote')}
					</p>
				)}
			</WidgetSection>
		</TicketSection>
	)
}

// What the episodes measured: the resources, none cited once the episodes
// are judged, or a dash before that.
function measuredLine(
	analysis: TicketAnalysis,
	format: TicketFormat,
): ReactNode {
	const { t } = format
	if (analysis.measuredResources) {
		return (
			<span className="font-mono text-[12px]">
				{format.figures(analysis.measuredResources)}
			</span>
		)
	}
	if (!analysis.episodes) {
		return <UnknownValue className="font-mono text-[12px]" />
	}

	return (
		<>
			<span className="font-mono text-[12px]">
				{t('analysis.measuredNone')}
			</span>
			<span className="text-[11px] text-muted-foreground">
				{t('analysis.measuredReasoning')}
			</span>
		</>
	)
}
