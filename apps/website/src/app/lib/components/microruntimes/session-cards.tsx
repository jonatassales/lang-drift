'use client'

import { cn } from '@langdrift/web/utilities'
import { useAtomValue, useSetAtom } from 'jotai'
import type { ReactNode } from 'react'

import {
	type PlanPhaseId,
	planPhaseIds,
	revealed,
} from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'
import {
	microruntimeSessionAtom,
	openCardAtom,
	openCardByHandAtom,
} from '@state'

import { AnalysisCard, analysisSummary } from './analysis-card'
import { PlanCard, planSummary } from './plan-card'
import { SessionCard } from './session-card'
import { TicketCard, ticketSummary } from './ticket-card'
import {
	createTicketCardFormat,
	type TicketCardCopy,
} from './ticket-card-format'

interface SessionCardsProps {
	copy: TicketCardCopy
	locale: WebsiteLocale
}

// The three cards the session fills beside the terminal, one per act: the
// ticket, the analysis of its documents and the project's episodes, and
// the plan. The act opens its card and collapses the others to their
// summary line; a click opens a card by hand until the next act. The stack
// is absent until the command is entered, and its height is the row's, so
// the terminal ends where the open card ends. It owns no clock; every
// value derives from how many lines the terminal has completed.
export function SessionCards(props: SessionCardsProps) {
	const { copy, locale } = props

	const { playedCount } = useAtomValue(microruntimeSessionAtom)
	const openCard = useAtomValue(openCardAtom)
	const openCardByHand = useSetAtom(openCardByHandAtom)
	const format = createTicketCardFormat(copy, locale)
	const visible = revealed(playedCount, 'cards')
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
	}

	return (
		<div
			aria-hidden={visible ? undefined : 'true'}
			className={cn(
				'grid content-start gap-3 motion-safe:transition-opacity motion-safe:duration-500',
				visible ? 'visible opacity-100' : 'invisible opacity-0',
			)}
		>
			{planPhaseIds.map((phase) => (
				<SessionCard
					body={cards[phase].body}
					expanded={openCard === phase}
					key={phase}
					label={copy.labels[phase]}
					onOpen={() => openCardByHand(phase)}
					summary={cards[phase].summary}
				/>
			))}
		</div>
	)
}
