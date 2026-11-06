'use client'

import { useAtom } from 'jotai'
import { createTranslator } from 'next-intl'
import type { MouseEvent } from 'react'

import { dollars } from '@app/template/formatters/dollars.fmt'
import type { WebsiteLocale } from '@i18n/routing'
import { computeEnvelope } from '@routes/platform/finops/content/compute-envelope.compute'
import {
	computeSegmentIds,
	environments,
	teamComputeEnvelope,
} from '@routes/platform/finops/content/finops-content.data'
import type {
	ComputeSegmentId,
	Environment,
	EnvironmentId,
	FinopsMessages,
} from '@routes/platform/finops/content/finops-content.types'
import { ChoiceButton } from '@routes/platform/product-page/choice-button'
import { EnvelopeBar } from '@routes/platform/product-page/envelope-bar'
import { selectedEnvironmentAtom } from '@state'
import { Eyebrow } from 'web/ui/eyebrow'

const segmentColors: Record<ComputeSegmentId, string> = {
	used: 'bg-primary',
	idle: 'bg-primary/45',
	held: 'bg-muted-foreground/45',
	available: 'bg-foreground/10',
}

const amountColumn = 'font-mono text-[13px] tabular-nums sm:text-right'

// The FinOps hero demo: the team's compute envelope as a stacked bar, then the
// environments drawing on it with their CPU, memory and GPU hours. Selecting
// one drives the receipt docked beside it.
interface ComputeLedgerDemoProps {
	copy: FinopsMessages['ledger']
	locale: WebsiteLocale
}

export function ComputeLedgerDemo(props: ComputeLedgerDemoProps) {
	const { copy, locale } = props

	const [selectedId, setSelectedId] = useAtom(selectedEnvironmentAtom)

	function selectEnvironment(event: MouseEvent<HTMLButtonElement>) {
		const targetId = event.currentTarget.dataset.id as EnvironmentId | undefined

		if (targetId) {
			setSelectedId(targetId)
		}
	}

	function usage(environment: Environment) {
		const { cpuHours, gpuHours, memoryGbHours } = environment

		if (cpuHours === null || gpuHours === null || memoryGbHours === null) {
			return copy.unknown
		}

		return translate('usage', {
			cpu: cpuHours,
			gpu: gpuHours,
			memory: memoryGbHours,
		})
	}

	const translate = createTranslator({
		locale,
		messages: copy,
	})
	const money = dollars(locale)
	const { authorized } = teamComputeEnvelope
	const amounts = computeEnvelope(teamComputeEnvelope)
	const segments = computeSegmentIds.map((id) => ({
		amount: amounts[id],
		color: segmentColors[id],
		id,
		label: copy.segments[id],
	}))

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
						{money.format(authorized)}
					</strong>
				</div>
			</header>

			<EnvelopeBar
				authorized={authorized}
				locale={locale}
				segments={segments}
			/>

			<div className="mt-7 border-t border-border pt-4">
				<div
					aria-hidden="true"
					className="hidden grid-cols-[minmax(0,1.8fr)_repeat(2,minmax(0,0.9fr))] gap-4 px-3 pb-2 sm:grid"
				>
					<Eyebrow>{copy.columns.environment}</Eyebrow>
					<Eyebrow className="text-right">{copy.columns.ceiling}</Eyebrow>
					<Eyebrow className="text-right">{copy.columns.observed}</Eyebrow>
				</div>
				<ul className="m-0 grid list-none gap-1 p-0">
					{environments.map((environment) => {
						const environmentCopy = copy.environments[environment.id]

						return (
							<li key={environment.id}>
								<ChoiceButton
									className="grid grid-cols-2 gap-x-4 gap-y-1.5 border-transparent px-3 py-3 sm:grid-cols-[minmax(0,1.8fr)_repeat(2,minmax(0,0.9fr))] sm:items-center"
									data-id={environment.id}
									kind="toggle"
									onSelect={selectEnvironment}
									selected={environment.id === selectedId}
								>
									<span className="col-span-2 grid gap-0.5 sm:col-span-1">
										<span className="text-[14px] font-medium">
											{environmentCopy.title}
										</span>
										<span className="text-[12px] text-muted-foreground">
											{environmentCopy.owner}
										</span>
										<span className="font-mono text-[11px] text-muted-foreground tabular-nums">
											{usage(environment)}
										</span>
									</span>
									<span className={amountColumn}>
										<span className="block text-muted-foreground sm:sr-only">
											{copy.columns.ceiling}{' '}
										</span>
										{money.format(environment.ceiling)}
									</span>
									<span className={amountColumn}>
										<span className="block text-muted-foreground sm:sr-only">
											{copy.columns.observed}{' '}
										</span>
										{environment.observed === null ? (
											<span className="inline-grid font-sans text-status-warning sm:justify-items-end">
												<span>{copy.unknown}</span>
												<span className="text-[11px]">{copy.meterSilent}</span>
											</span>
										) : (
											money.format(environment.observed)
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
