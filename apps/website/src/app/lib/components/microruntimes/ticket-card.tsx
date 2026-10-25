'use client'

import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useAtomValue } from 'jotai'
import { createTranslator } from 'next-intl'

import {
	hasPlayed,
	type MicroruntimesMessages,
	type PlanRevealId,
	planReveals,
	suggestedTotals,
	ticketPlan,
} from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'
import { microruntimeSessionAtom } from '@state'

import { CardValue, type CardValueState } from './card-value'
import { ComponentRow } from './component-row'
import { ReadMark } from './read-mark'
import { SessionControl } from './session-control'

interface TicketCardProps {
	copy: MicroruntimesMessages['session']['card']
	locale: WebsiteLocale
}

const artifactIds = ['prd', 'designDoc', 'spec'] as const

// The ticket the simulated session plans, printed like the ticket itself:
// absent until the command is entered, then filled in line by line from the
// terminal's progress. It owns no clock; every field derives from how many
// lines the terminal has completed.
export function TicketCard(props: TicketCardProps) {
	const { copy, locale } = props

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
	const card = createTranslator({
		locale,
		messages: copy,
	})
	const number = new Intl.NumberFormat(locale)
	const cpu = (value: number) =>
		card('cpuValue', { value: number.format(value) })
	const memory = (value: number) =>
		card('memoryValue', { value: number.format(value) })
	const resolved = ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
	const resolvedTotals = suggestedTotals(resolved)
	const components: CardValueState = played('componentsSkeleton')
		? resolved.length > 0
			? 'set'
			: 'resolving'
		: 'unset'
	const bounded = ticketPlan.bounded
	const suggested = suggestedTotals(ticketPlan.components)
	const boundedState: CardValueState = played('bounded') ? 'set' : 'unset'
	const sizedState: CardValueState = played('suggested') ? 'set' : 'unset'
	const limitsKnown = played('limits')
	const reportDrafted = played('report')
	const hardwareRows = [
		{
			id: 'cpu',
			label: copy.hardware.cpu,
			bounded: cpu(bounded.cpuMillis),
			suggested: cpu(suggested.cpuMillis),
			share: suggested.cpuMillis / bounded.cpuMillis,
		},
		{
			id: 'memory',
			label: copy.hardware.memory,
			bounded: memory(bounded.memoryMiB),
			suggested: memory(suggested.memoryMiB),
			share: suggested.memoryMiB / bounded.memoryMiB,
		},
	]
	const reportRows = [
		{
			id: 'bounded',
			term: copy.report.bounded,
			state: boundedState,
			detail: `${cpu(bounded.cpuMillis)} · ${memory(bounded.memoryMiB)}`,
		},
		{
			id: 'suggested',
			term: copy.report.suggested,
			state: (resolved.length > 0 ? 'set' : 'unset') as CardValueState,
			detail: `${cpu(resolvedTotals.cpuMillis)} · ${memory(resolvedTotals.memoryMiB)}`,
		},
		{
			id: 'resource',
			term: copy.report.suggestedResource,
			state: (reportDrafted ? 'set' : 'unset') as CardValueState,
			detail: card('deployments', { count: ticketPlan.components.length }),
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

			<div className="mt-5 border-t border-border pt-4">
				<div className="grid grid-cols-[minmax(0,1fr)_auto_auto] gap-x-3">
					<Eyebrow>{copy.components}</Eyebrow>
					<Eyebrow>{copy.columns.cpu}</Eyebrow>
					<Eyebrow>{copy.columns.memory}</Eyebrow>
				</div>
				{components === 'unset' ? (
					<CardValue
						className="mt-2 block font-mono text-[12px]"
						labels={copy}
						state="unset"
					/>
				) : (
					<ul className="m-0 mt-1 list-none p-0">
						{ticketPlan.components.map((component) => (
							<ComponentRow
								component={component}
								cpu={cpu(component.cpuMillis)}
								key={component.id}
								labels={copy}
								memory={memory(component.memoryMiB)}
								share={component.cpuMillis / ticketPlan.bounded.cpuMillis}
								state={
									hasPlayed(playedCount, component.id) ? 'set' : 'resolving'
								}
							/>
						))}
					</ul>
				)}
				{components !== 'unset' && (
					<CardValue
						className="mt-2 block border-t border-border pt-2 text-[12px] text-muted-foreground"
						labels={copy}
						state={components}
						width="w-44"
					>
						{resolved.length === ticketPlan.components.length
							? card('componentsTotal', {
									count: resolved.length,
									cpu: number.format(resolvedTotals.cpuMillis),
									memory: number.format(resolvedTotals.memoryMiB),
								})
							: card('resolvedProgress', {
									resolved: resolved.length,
									total: ticketPlan.components.length,
									cpu: number.format(resolvedTotals.cpuMillis),
									memory: number.format(resolvedTotals.memoryMiB),
								})}
					</CardValue>
				)}
			</div>

			<div className="mt-5 border-t border-border pt-4">
				<table className="w-full border-collapse text-[12px]">
					<thead>
						<tr>
							<th className="pb-2 text-left font-normal" scope="col">
								<Eyebrow>{copy.hardware.resource}</Eyebrow>
							</th>
							<th className="pb-2 text-left font-normal" scope="col">
								<Eyebrow>{copy.hardware.bounded}</Eyebrow>
							</th>
							<th className="pb-2 text-left font-normal" scope="col">
								<Eyebrow>{copy.hardware.suggested}</Eyebrow>
							</th>
						</tr>
					</thead>
					<tbody>
						{hardwareRows.map((row) => (
							<tr className="border-t border-border" key={row.id}>
								<th className="py-2 pr-3 text-left font-normal" scope="row">
									{row.label}
								</th>
								<td className="py-2 pr-3 font-mono tabular-nums">
									<CardValue labels={copy} state={boundedState}>
										{row.bounded}
									</CardValue>
								</td>
								<td className="py-2 font-mono tabular-nums">
									<CardValue labels={copy} state={sizedState}>
										{row.suggested}
									</CardValue>
									{sizedState === 'set' && (
										<span
											aria-hidden="true"
											className="mt-1 block h-0.5 overflow-hidden rounded bg-foreground/5"
										>
											<span
												className="block h-full rounded bg-primary motion-safe:transition-[width] motion-safe:duration-700 motion-safe:ease-out"
												style={{ width: `${row.share * 100}%` }}
											/>
										</span>
									)}
								</td>
							</tr>
						))}
						<tr className="border-t border-border">
							<th className="py-2 pr-3 text-left font-normal" scope="row">
								{copy.hardware.gpu}
							</th>
							<td className="py-2 pr-3 font-mono">
								<CardValue labels={copy} state={limitsKnown ? 'set' : 'unset'}>
									{copy.none}
								</CardValue>
							</td>
							<td className="py-2 font-mono">
								<CardValue labels={copy} state={limitsKnown ? 'set' : 'unset'}>
									{copy.none}
								</CardValue>
								{limitsKnown && (
									<span className="block font-sans text-[11px] text-muted-foreground">
										{copy.notRuntimeField}
									</span>
								)}
							</td>
						</tr>
						<tr className="border-t border-border">
							<th className="py-2 pr-3 text-left font-normal" scope="row">
								{copy.hardware.disk}
							</th>
							<td className="py-2 pr-3 font-mono">
								{limitsKnown ? (
									<span aria-hidden="true">—</span>
								) : (
									<CardValue labels={copy} state="unset" />
								)}
							</td>
							<td className="py-2 font-mono">
								{limitsKnown ? (
									<span aria-hidden="true">—</span>
								) : (
									<CardValue labels={copy} state="unset" />
								)}
								{limitsKnown && (
									<span className="block font-sans text-[11px] text-muted-foreground">
										{copy.notRuntimeField}
									</span>
								)}
							</td>
						</tr>
					</tbody>
				</table>
				{played('priceBasis') && (
					<p className="mt-2 mb-0 text-[11px] text-muted-foreground">
						{copy.priceBasis}
					</p>
				)}
			</div>

			<div className="mt-5 border-t border-border pt-4">
				<dl className="m-0 grid gap-2.5 text-[13px]">
					{reportRows.map((row) => (
						<div
							className="grid gap-1 sm:grid-cols-[minmax(0,0.4fr)_minmax(0,1fr)] sm:gap-4"
							key={row.id}
						>
							<Eyebrow as="dt">{row.term}</Eyebrow>
							<dd className="m-0 font-mono text-[13px] tabular-nums">
								<CardValue labels={copy} state={row.state} width="w-40">
									{row.detail}
								</CardValue>
							</dd>
						</div>
					))}
				</dl>
				<div aria-atomic="true" aria-live="polite" className="grid gap-2.5">
					{reportDrafted && (
						<>
							<dl className="m-0 mt-2.5 grid gap-1 sm:grid-cols-[minmax(0,0.4fr)_minmax(0,1fr)] sm:gap-4">
								<Eyebrow as="dt">{copy.report.final}</Eyebrow>
								<dd className="m-0 font-mono text-[14px] font-medium tabular-nums">
									{`${cpu(suggested.cpuMillis)} · ${memory(suggested.memoryMiB)} · ${copy.pendingAuthorization}`}
								</dd>
							</dl>
							<p className="m-0 font-mono text-[11px] text-muted-foreground">
								{card('planPath', { revision: ticketPlan.draftRevision })}
							</p>
						</>
					)}
				</div>
				<div className="mt-3 flex justify-end">
					<SessionControl copy={copy.control} />
				</div>
			</div>
		</article>
	)
}
