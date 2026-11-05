'use client'

import { useAtom } from 'jotai'
import type { MouseEvent } from 'react'

import type { WebsiteLocale } from '@i18n/routing'
import { dollars } from '@routes/platform/finops/content/dollars.fmt'
import { envelopeSegments } from '@routes/platform/finops/content/envelope-segments.compute'
import {
	type EnvelopeSegmentId,
	envelope,
	envelopeSegmentIds,
	type FinopsMessages,
	type WorkItemId,
	workItems,
} from '@routes/platform/finops/content/finops-content.data'
import { ChoiceButton } from '@routes/platform/product-page/choice-button'
import { selectedWorkItemAtom } from '@state'
import { Eyebrow } from 'web/ui/eyebrow'
import { cn } from 'web/utilities'

const segmentColors: Record<EnvelopeSegmentId, string> = {
	spent: 'bg-primary',
	reserved: 'bg-primary/45',
	protected: 'bg-muted-foreground/45',
	available: 'bg-foreground/10',
}

const amountColumn = 'font-mono text-[13px] tabular-nums sm:text-right'

// The FinOps hero demo: a simulated team envelope as a stacked bar, then
// the work items drawing on it. Selecting an item drives the receipt docked
// beside it.
interface BudgetLedgerDemoProps {
	copy: FinopsMessages['ledger']
	locale: WebsiteLocale
}

export function BudgetLedgerDemo(props: BudgetLedgerDemoProps) {
	const { copy, locale } = props

	const [selectedId, setSelectedId] = useAtom(selectedWorkItemAtom)

	function selectWorkItem(event: MouseEvent<HTMLButtonElement>) {
		const targetId = event.currentTarget.dataset.id as WorkItemId | undefined

		if (targetId) {
			setSelectedId(targetId)
		}
	}

	const money = dollars(locale)
	const percent = new Intl.NumberFormat(locale, {
		maximumFractionDigits: 0,
		style: 'percent',
	})
	const segments = envelopeSegments(envelope)
	const barLabel = envelopeSegmentIds
		.map((id) => `${copy.segments[id]} ${money.format(segments[id])}`)
		.join(', ')

	return (
		<article className="rounded-2xl border border-border bg-card p-5 sm:p-7">
			<header className="flex flex-wrap items-end justify-between gap-4">
				<div className="grid gap-1.5">
					<Eyebrow className="text-[10px]">{copy.label}</Eyebrow>
					<span className="text-[13px] text-muted-foreground">
						{copy.period}
					</span>
				</div>
				<div className="grid justify-items-end gap-1.5">
					<Eyebrow className="text-[10px]">{copy.authorized}</Eyebrow>
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
					className="hidden grid-cols-[minmax(0,1.6fr)_repeat(3,minmax(0,0.8fr))] gap-4 px-3 pb-2 sm:grid"
				>
					<Eyebrow>{copy.columns.item}</Eyebrow>
					<Eyebrow className="text-right">{copy.columns.estimate}</Eyebrow>
					<Eyebrow className="text-right">{copy.columns.reserved}</Eyebrow>
					<Eyebrow className="text-right">{copy.columns.observed}</Eyebrow>
				</div>
				<ul className="m-0 grid list-none gap-1 p-0">
					{workItems.map((item) => {
						const itemCopy = copy.items[item.id]

						return (
							<li key={item.id}>
								<ChoiceButton
									className="grid grid-cols-2 gap-x-4 gap-y-1.5 border-transparent px-3 py-3 sm:grid-cols-[minmax(0,1.6fr)_repeat(3,minmax(0,0.8fr))] sm:items-center"
									data-id={item.id}
									kind="toggle"
									onSelect={selectWorkItem}
									selected={item.id === selectedId}
								>
									<span className="col-span-2 grid gap-0.5 sm:col-span-1">
										<span className="text-[14px] font-medium">
											{itemCopy.title}
										</span>
										<span className="text-[12px] text-muted-foreground">
											{itemCopy.owner}
										</span>
									</span>
									<span className={amountColumn}>
										<span className="block text-muted-foreground sm:sr-only">
											{copy.columns.estimate}{' '}
										</span>
										{money.format(item.estimate)}
									</span>
									<span className={amountColumn}>
										<span className="block text-muted-foreground sm:sr-only">
											{copy.columns.reserved}{' '}
										</span>
										{money.format(item.reserved)}
									</span>
									<span
										className={cn(amountColumn, 'col-span-2 sm:col-span-1')}
									>
										<span className="block text-muted-foreground sm:sr-only">
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
								</ChoiceButton>
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
