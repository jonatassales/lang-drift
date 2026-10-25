'use client'

import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useAtomValue } from 'jotai'

import {
	hasPlayed,
	type MicroruntimesMessages,
	type PlanRevealId,
	planReveals,
	ticketPlan,
} from '@data/microruntimes-content'
import { microruntimeSessionAtom } from '@state'

import { CardValue, type CardValueState } from './card-value'
import { ReadMark } from './read-mark'

interface TicketCardProps {
	copy: MicroruntimesMessages['session']['card']
}

const artifactIds = ['prd', 'designDoc', 'spec'] as const

// The ticket the simulated session plans, printed like the ticket itself:
// absent until the command is entered, then filled in line by line from the
// terminal's progress. It owns no clock; every field derives from how many
// lines the terminal has completed.
export function TicketCard(props: TicketCardProps) {
	const { copy } = props

	const { playedCount } = useAtomValue(microruntimeSessionAtom)
	const played = (id: PlanRevealId) => hasPlayed(playedCount, planReveals[id])
	const visible = played('card')
	const artifactsRead = played('artifactsRead')
	const suggestion: CardValueState = played('suggestion')
		? 'set'
		: played('suggestionSkeleton')
			? 'resolving'
			: 'unset'
	const guardrail = ticketPlan.guardrail

	return (
		<article
			aria-hidden={visible ? undefined : 'true'}
			className={cn(
				'rounded-2xl border border-border bg-card p-5 text-left motion-safe:transition-opacity motion-safe:duration-500 sm:p-6',
				visible ? 'visible opacity-100' : 'invisible opacity-0',
			)}
		>
			<header className="grid gap-2">
				<div className="flex items-center justify-between gap-4">
					<Eyebrow>{copy.ticket}</Eyebrow>
					<span className="rounded-full border border-border px-2 py-[3px] font-mono text-[11px]">
						{copy.status}
					</span>
				</div>
				<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
					<strong className="font-mono text-[44px] leading-none font-medium tracking-[-0.03em] max-sm:text-[36px]">
						{ticketPlan.ticket.number}
					</strong>
					<span className="font-mono text-[13px] text-muted-foreground">
						{ticketPlan.ticket.slug}
					</span>
					{artifactsRead && <ReadMark label={copy.readMark} />}
				</p>
				<p className="m-0 text-[15px] leading-snug font-medium">{copy.title}</p>
			</header>

			<dl className="m-0 mt-5 grid gap-2.5 border-t border-border pt-4 text-[13px]">
				{artifactIds.map((id) => (
					<div
						className="grid gap-1 sm:grid-cols-[minmax(0,0.3fr)_minmax(0,1fr)] sm:gap-4"
						key={id}
					>
						<Eyebrow as="dt">{copy.artifacts[id]}</Eyebrow>
						<dd className="m-0 flex flex-wrap items-center gap-x-2 font-mono text-[12px]">
							{ticketPlan.artifacts[id]}
							{artifactsRead && <ReadMark label={copy.readMark} />}
						</dd>
					</div>
				))}
				<div className="grid gap-1 sm:grid-cols-[minmax(0,0.3fr)_minmax(0,1fr)] sm:gap-4">
					<Eyebrow as="dt">{copy.artifacts.guardrail}</Eyebrow>
					<dd className="m-0 grid gap-0.5 font-mono text-[12px]">
						<CardValue
							labels={copy}
							state={played('guardrail') ? 'set' : 'unset'}
							width="w-40"
						>
							{[guardrail.file, ...guardrail.criteria].join(' · ')}
						</CardValue>
						{played('guardrail') && (
							<span className="font-sans text-[11px] text-muted-foreground">
								{copy.guardrailNote}
							</span>
						)}
					</dd>
				</div>
			</dl>

			<div className="mt-5 grid grid-cols-2 gap-x-6 gap-y-3 border-t border-border pt-4">
				<div className="grid gap-1.5">
					<Eyebrow>{copy.suggestedModel}</Eyebrow>
					<CardValue
						className="font-mono text-[14px]"
						labels={copy}
						state={suggestion}
						width="w-32"
					>
						{ticketPlan.model}
					</CardValue>
				</div>
				<div className="grid gap-1.5">
					<Eyebrow>{copy.suggestedEffort}</Eyebrow>
					<CardValue
						className="font-mono text-[14px]"
						labels={copy}
						state={suggestion}
						width="w-16"
					>
						{copy.efforts[ticketPlan.effort]}
					</CardValue>
				</div>
				{suggestion === 'set' && (
					<p className="col-span-2 m-0 text-[11px] text-muted-foreground">
						{copy.recommendationNote}
					</p>
				)}
			</div>
		</article>
	)
}
