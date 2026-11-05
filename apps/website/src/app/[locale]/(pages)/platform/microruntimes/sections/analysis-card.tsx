import { hasPlayed } from '@routes/platform/microruntimes/content/has-played.compute'
import {
	documentIds,
	documentReads,
	episodes,
	similarityThreshold,
	ticketPlan,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { revealed } from '@routes/platform/microruntimes/content/revealed.compute'
import { workflowModel } from '@routes/platform/microruntimes/content/workflow-model.compute'
import { Eyebrow } from 'web/ui/eyebrow'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import { DocumentRow } from './document-row'
import { EpisodeRow } from './episode-row'
import { extractedFact } from './extracted-fact.compute'
import type {
	TicketCardCopy,
	TicketCardFormat,
} from './ticket-card-format.factory'

export { analysisSummary } from './analysis-summary.compute'

const aside = 'font-mono text-[11px] text-muted-foreground'
const suggestedModel = workflowModel(ticketPlan.model)

// What the agent learned before it planned: each linked document with what
// its read extracted, the project's episodes ranked by similarity and
// judged against the threshold, the resources none of them measured, and
// the model and effort it suggests from all of that.
interface AnalysisCardProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

export function AnalysisCard(props: AnalysisCardProps) {
	const { copy, format, playedCount } = props

	const { analysis } = copy
	const documentsRead = revealed(playedCount, 'documentsRead')
	const guardrailJudged = revealed(playedCount, 'guardrailJudged')
	const episodesState = revealState(
		revealed(playedCount, 'candidates'),
		revealed(playedCount, 'episodesSkeleton'),
	)
	const judged = revealed(playedCount, 'verdicts')
	const citedState = revealState(revealed(playedCount, 'cited'))
	const measured = revealed(playedCount, 'measured')
	const suggestionState = revealState(
		revealed(playedCount, 'suggestion'),
		revealed(playedCount, 'suggestionSkeleton'),
	)

	return (
		<div className="grid">
			<div>
				<div className="flex items-baseline justify-between gap-3">
					<Eyebrow>{analysis.documents}</Eyebrow>
					{documentsRead && (
						<span className={aside}>
							{format.message('analysis.documentsSummary', {
								count: documentIds.length,
							})}
						</span>
					)}
				</div>
				<ul className="m-0 mt-1 list-none p-0">
					{documentIds.map((id) => (
						<DocumentRow
							extracted={extractedFact(id, copy, format)}
							file={ticketPlan.artifacts[id]}
							judgment={
								id === 'guardrail' && guardrailJudged
									? analysis.guardrailJudgment
									: undefined
							}
							key={id}
							labels={copy}
							read={hasPlayed(playedCount, documentReads[id])}
						/>
					))}
				</ul>
			</div>

			<CardSection className="mt-4 pt-3">
				<div className="flex items-baseline justify-between gap-3">
					<Eyebrow>{analysis.episodes}</Eyebrow>
					<span className={aside}>
						{format.message('analysis.threshold', {
							value: format.formatNumber.format(similarityThreshold),
						})}
					</span>
				</div>
				{episodesState === 'unset' ? (
					<CardValue
						className="mt-2 block font-mono text-[12px]"
						labels={copy}
						state="unset"
					/>
				) : (
					<ul className="m-0 mt-1 list-none p-0">
						{episodes.map((episode) => (
							<EpisodeRow
								episode={episode}
								key={episode.id}
								labels={copy}
								similarity={format.formatNumber.format(episode.similarity)}
								state={episodesState}
								verdict={
									judged
										? {
												kind: episode.verdict,
												label: analysis.verdicts[episode.verdict],
												reason: episode.reason
													? analysis.reasons[episode.reason]
													: undefined,
											}
										: undefined
								}
							/>
						))}
					</ul>
				)}
				<div className="mt-1 flex flex-wrap items-baseline justify-between gap-x-3 border-t border-border pt-2">
					<Eyebrow>{analysis.measured}</Eyebrow>
					<span className="flex flex-wrap items-baseline gap-x-2">
						<CardValue
							className="font-mono text-[12px]"
							labels={copy}
							state={citedState}
							width="w-24"
						>
							{analysis.measuredNone}
						</CardValue>
						{measured && (
							<span className="text-[11px] text-muted-foreground motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0">
								{analysis.measuredReasoning}
							</span>
						)}
					</span>
				</div>
			</CardSection>

			<CardSection className="mt-4 grid grid-cols-2 gap-x-6 gap-y-2 pt-3">
				<div className="grid gap-1.5">
					<Eyebrow>{analysis.suggestedModel}</Eyebrow>
					<CardValue
						className="grid gap-0.5"
						labels={copy}
						state={suggestionState}
						width="w-32"
					>
						<span className="text-[14px]">{suggestedModel.name}</span>
						<span className="font-mono text-[11px] text-muted-foreground">
							{suggestedModel.id}
						</span>
					</CardValue>
				</div>
				<div className="grid gap-1.5">
					<Eyebrow>{analysis.suggestedEffort}</Eyebrow>
					<CardValue
						className="font-mono text-[14px]"
						labels={copy}
						state={suggestionState}
						width="w-16"
					>
						{analysis.efforts[ticketPlan.effort]}
					</CardValue>
				</div>
				{suggestionState === 'set' && (
					<p className="col-span-2 m-0 text-[11px] text-muted-foreground">
						{analysis.recommendationNote}
					</p>
				)}
			</CardSection>
		</div>
	)
}
