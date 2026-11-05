'use client'

import { useAtomValue, useSetAtom } from 'jotai'
import type { MouseEvent, ReactNode } from 'react'

import type { WebsiteLocale } from '@i18n/routing'
import {
	type PlanPhaseId,
	planPhaseIds,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { revealed } from '@routes/platform/microruntimes/content/revealed.compute'
import {
	microruntimeSessionAtom,
	openCardAtom,
	openCardByHandAtom,
	planPlayedCountAtom,
	runScriptAtom,
	workflowChoicesAtom,
} from '@state'
import { cn } from 'web/utilities'

import { AnalysisCard, analysisSummary } from './analysis-card'
import { PlanCard, planSummary } from './plan-card'
import { SessionCard } from './session-card'
import { TicketCard, ticketSummary } from './ticket-card'
import {
	type TicketCardCopy,
	ticketCardFormat,
} from './ticket-card-format.factory'
import {
	runOutcome,
	WorkflowCard,
	workflowStatus,
	workflowSummary,
} from './workflow-card'

// The four cards the session fills beside the terminal, one per act: the
// ticket, the analysis of its documents and the project's episodes, the
// plan, and the workflow that plays the run. The act opens its card and
// collapses the others to their summary line; a click opens a card by hand
// until the next act. The stack is absent until the command is entered,
// and its height is the row's, so the terminal ends where the open card
// ends. The plan cards keep their end while a run plays; one status line
// announces only how the run ended. Every value derives from how many lines
// the terminal has completed.
interface SessionCardsProps {
	copy: TicketCardCopy
	locale: WebsiteLocale
}

export function SessionCards(props: SessionCardsProps) {
	const { copy, locale } = props

	const session = useAtomValue(microruntimeSessionAtom)
	const playedCount = useAtomValue(planPlayedCountAtom)
	const choices = useAtomValue(workflowChoicesAtom)
	const script = useAtomValue(runScriptAtom)
	const openCard = useAtomValue(openCardAtom)
	const openCardByHand = useSetAtom(openCardByHandAtom)

	function handleOpenCard(event: MouseEvent<HTMLButtonElement>) {
		const phase = event.currentTarget.dataset.phase as PlanPhaseId | undefined

		if (phase) {
			openCardByHand(phase)
		}
	}

	const format = ticketCardFormat(copy, locale)
	const visible = revealed(playedCount, 'cards')
	const outcome = runOutcome(session, script)
	const cards: Record<PlanPhaseId, { body: ReactNode; summary: ReactNode }> = {
		ticket: {
			body: <TicketCard copy={copy} playedCount={playedCount} />,
			summary: ticketSummary(),
		},
		analysis: {
			body: (
				<AnalysisCard copy={copy} format={format} playedCount={playedCount} />
			),
			summary: analysisSummary(copy, format, playedCount),
		},
		plan: {
			body: <PlanCard copy={copy} format={format} playedCount={playedCount} />,
			summary: planSummary(copy, format, playedCount),
		},
		workflow: {
			body: <WorkflowCard copy={copy} format={format} />,
			summary: workflowSummary(copy, format, choices, outcome),
		},
	}

	return (
		<div
			aria-hidden={visible ? undefined : 'true'}
			className={cn(
				'grid min-w-0 content-start gap-3 motion-safe:transition-opacity motion-safe:duration-500',
				visible ? 'visible opacity-100' : 'invisible opacity-0',
			)}
		>
			{planPhaseIds.map((phase) => (
				<SessionCard
					body={cards[phase].body}
					expanded={openCard === phase}
					key={phase}
					label={copy.labels[phase]}
					onOpen={handleOpenCard}
					phase={phase}
					summary={cards[phase].summary}
				/>
			))}
			<p className="sr-only" role="status">
				{workflowStatus(copy, format, outcome)}
			</p>
		</div>
	)
}
