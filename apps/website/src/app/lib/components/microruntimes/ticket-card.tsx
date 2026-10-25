'use client'

import { DetailList } from '@langdrift/web/ui/detail-list'
import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useAtomValue } from 'jotai'

import { revealed, ticketPlan } from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'
import { microruntimeSessionAtom } from '@state'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import { ComponentsManifest } from './components-manifest'
import { HardwareTable } from './hardware-table'
import { PlanReport } from './plan-report'
import { ReadMark } from './read-mark'
import { SessionControl } from './session-control'
import {
	createTicketCardFormat,
	type TicketCardCopy,
} from './ticket-card-format'

interface TicketCardProps {
	copy: TicketCardCopy
	locale: WebsiteLocale
}

const artifactIds = ['prd', 'designDoc', 'spec'] as const
const fileName = 'font-mono text-[12px]'

// The ticket the simulated session plans, printed like the ticket itself:
// absent until the command is entered, then filled in line by line from the
// terminal's progress. It owns no clock; every field derives from how many
// lines the terminal has completed.
export function TicketCard(props: TicketCardProps) {
	const { copy, locale } = props

	const { playedCount } = useAtomValue(microruntimeSessionAtom)
	const format = createTicketCardFormat(copy, locale)
	const visible = revealed(playedCount, 'card')
	const artifactsRead = revealed(playedCount, 'artifactsRead')
	const guardrailState = revealState(revealed(playedCount, 'guardrail'))
	const suggestionState = revealState(
		revealed(playedCount, 'suggestion'),
		revealed(playedCount, 'suggestionSkeleton'),
	)
	const guardrail = ticketPlan.guardrail
	const artifacts = [
		...artifactIds.map((id) => ({
			detail: (
				<span className="flex flex-wrap items-center gap-x-2">
					<span className={fileName}>{ticketPlan.artifacts[id]}</span>
					{artifactsRead && <ReadMark label={copy.readMark} />}
				</span>
			),
			id,
			term: copy.artifacts[id],
		})),
		{
			detail: (
				<span className="grid gap-0.5">
					<CardValue
						className={fileName}
						labels={copy}
						state={guardrailState}
						width="w-40"
					>
						{[guardrail.file, ...guardrail.criteria].join(' · ')}
					</CardValue>
					{guardrailState === 'set' && (
						<span className="text-[11px] text-muted-foreground">
							{copy.guardrailNote}
						</span>
					)}
				</span>
			),
			id: 'guardrail',
			term: copy.artifacts.guardrail,
		},
	]

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

			<CardSection>
				<DetailList items={artifacts} />
			</CardSection>

			<CardSection className="grid grid-cols-2 gap-x-6 gap-y-3">
				<div className="grid gap-1.5">
					<Eyebrow>{copy.suggestedModel}</Eyebrow>
					<CardValue
						className="font-mono text-[14px]"
						labels={copy}
						state={suggestionState}
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
						state={suggestionState}
						width="w-16"
					>
						{copy.efforts[ticketPlan.effort]}
					</CardValue>
				</div>
				{suggestionState === 'set' && (
					<p className="col-span-2 m-0 text-[11px] text-muted-foreground">
						{copy.recommendationNote}
					</p>
				)}
			</CardSection>

			<ComponentsManifest
				copy={copy}
				format={format}
				playedCount={playedCount}
			/>
			<HardwareTable copy={copy} format={format} playedCount={playedCount} />
			<PlanReport copy={copy} format={format} playedCount={playedCount} />

			<div className="mt-3 flex justify-end">
				<SessionControl copy={copy.control} />
			</div>
		</article>
	)
}
