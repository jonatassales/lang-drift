'use client'

import { cn } from '@langdrift/web/utilities'
import { useAtom } from 'jotai'

import {
	type EnvelopeSegmentId,
	envelope,
	envelopeSegmentIds,
	envelopeSegments,
	type FinopsMessages,
	workItems,
} from '@data/finops-content'
import type { WebsiteLocale } from '@i18n/routing'
import { selectedWorkItemAtom } from '@state'

interface BudgetLedgerDemoProps {
	copy: FinopsMessages['ledger']
	locale: WebsiteLocale
}

const segmentColors: Record<EnvelopeSegmentId, string> = {
	spent: 'bg-primary',
	reserved: 'bg-primary/45',
	protected: 'bg-muted-foreground/45',
	available: 'bg-foreground/10',
}

const cell = 'font-mono text-[13px] tabular-nums sm:text-right'

// The FinOps hero demo: an illustrative team envelope as a stacked bar, then
// the work items drawing on it. Selecting an item drives the receipt below.
export function BudgetLedgerDemo(props: BudgetLedgerDemoProps) {
	const { copy, locale } = props

	const [selectedId, setSelectedId] = useAtom(selectedWorkItemAtom)
	const money = new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 0,
		style: 'currency',
	})
	const percent = new Intl.NumberFormat(locale, {
		maximumFractionDigits: 0,
		style: 'percent',
	})
	const segments = envelopeSegments(envelope)
	const barLabel = envelopeSegmentIds
		.map((id) => `${copy.segments[id]} ${money.format(segments[id])}`)
		.join(', ')

	return (
		<article className="mx-auto max-w-[1040px] rounded-2xl border border-border bg-card p-5 text-left sm:p-7">
			<header className="flex flex-wrap items-end justify-between gap-4">
				<div className="grid gap-1.5">
					<span className="font-mono text-[10px] tracking-[0.08em] text-muted-foreground uppercase">
						{copy.label}
					</span>
					<span className="text-[13px] text-muted-foreground">
						{copy.period}
					</span>
				</div>
				<div className="grid justify-items-end gap-1.5">
					<span className="font-mono text-[10px] tracking-[0.08em] text-muted-foreground uppercase">
						{copy.authorized}
					</span>
					<strong className="font-mono text-[28px] leading-none font-medium tracking-[-0.02em]">
						{money.format(envelope.authorized)}
					</strong>
				</div>
			</header>

			<div
				aria-label={barLabel}
				className="mt-6 flex h-3 overflow-hidden rounded-full bg-foreground/5"
				role="img"
			>
				{envelopeSegmentIds.map((id) => (
					<span
						className={cn('h-full', segmentColors[id])}
						key={id}
						style={{
							flexBasis: `${(segments[id] / envelope.authorized) * 100}%`,
						}}
					/>
				))}
			</div>
			<ul className="m-0 mt-4 grid list-none grid-cols-2 gap-x-6 gap-y-3 p-0 sm:grid-cols-4">
				{envelopeSegmentIds.map((id) => (
					<li className="grid gap-1" key={id}>
						<span className="flex items-center gap-2 text-[12px] text-muted-foreground">
							<span
								aria-hidden="true"
								className={cn('size-2 rounded-full', segmentColors[id])}
							/>
							{copy.segments[id]}
						</span>
						<span className="font-mono text-[14px] tabular-nums">
							{money.format(segments[id])}
							<span className="ml-2 text-[11px] text-muted-foreground">
								{percent.format(segments[id] / envelope.authorized)}
							</span>
						</span>
					</li>
				))}
			</ul>

			<div className="mt-7 border-t border-border pt-4">
				<div
					aria-hidden="true"
					className="hidden grid-cols-[minmax(0,1.6fr)_repeat(3,minmax(0,0.8fr))] gap-4 px-3 pb-2 font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase sm:grid"
				>
					<span>{copy.columns.item}</span>
					<span className="text-right">{copy.columns.estimate}</span>
					<span className="text-right">{copy.columns.reserved}</span>
					<span className="text-right">{copy.columns.observed}</span>
				</div>
				<ul className="m-0 grid list-none gap-1 p-0">
					{workItems.map((item) => {
						const selected = item.id === selectedId
						const names = copy.items[item.id]

						return (
							<li key={item.id}>
								<button
									aria-pressed={selected}
									className={cn(
										'grid w-full cursor-pointer grid-cols-2 gap-x-4 gap-y-1.5 rounded-xl border border-transparent px-3 py-3 text-left text-foreground motion-safe:transition-colors sm:grid-cols-[minmax(0,1.6fr)_repeat(3,minmax(0,0.8fr))] sm:items-center',
										'hover:bg-foreground/5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring',
										selected &&
											'border-primary/60 bg-primary/10 hover:bg-primary/10',
									)}
									onClick={() => setSelectedId(item.id)}
									type="button"
								>
									<span className="col-span-2 grid gap-0.5 sm:col-span-1">
										<span className="text-[14px] font-medium">
											{names.title}
										</span>
										<span className="text-[12px] text-muted-foreground">
											{names.owner}
										</span>
									</span>
									<span className={cell}>
										<span className="text-muted-foreground sm:sr-only">
											{copy.columns.estimate}{' '}
										</span>
										{money.format(item.estimate)}
									</span>
									<span className={cell}>
										<span className="text-muted-foreground sm:sr-only">
											{copy.columns.reserved}{' '}
										</span>
										{money.format(item.reserved)}
									</span>
									<span className={cn(cell, 'col-span-2 sm:col-span-1')}>
										<span className="text-muted-foreground sm:sr-only">
											{copy.columns.observed}{' '}
										</span>
										{item.observed === null ? (
											<span className="inline-grid font-sans text-status-warning sm:justify-items-end">
												<span>{copy.unknown}</span>
												<span className="text-[11px]">
													{copy.awaitingReconciliation}
												</span>
											</span>
										) : (
											money.format(item.observed)
										)}
									</span>
								</button>
							</li>
						)
					})}
				</ul>
				<p className="mt-3 mb-0 px-3 text-[12px] text-muted-foreground">
					{copy.selectHint}
				</p>
			</div>
		</article>
	)
}
