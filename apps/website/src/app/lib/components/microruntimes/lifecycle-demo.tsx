'use client'

import { cn } from '@langdrift/web/utilities'
import { useAtom } from 'jotai'
import { createTranslator } from 'next-intl'

import {
	lifecycleStageIds,
	type MicroruntimesMessages,
	runtimeDescriptor,
} from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'
import { selectedLifecycleStageAtom } from '@state'

interface LifecycleDemoProps {
	copy: MicroruntimesMessages['lifecycle']
	locale: WebsiteLocale
}

const term =
	'font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase'

// The Microruntimes hero demo: an illustrative runtime's summary, its six
// lifecycle steps, and what the selected step does and leaves behind.
// Selecting a step drives the evidence section further down the page.
export function LifecycleDemo(props: LifecycleDemoProps) {
	const { copy, locale } = props

	const [selectedId, setSelectedId] = useAtom(selectedLifecycleStageAtom)
	const t = createTranslator({
		locale,
		messages: copy,
	})
	const number = new Intl.NumberFormat(locale)
	const stage = copy.stages[selectedId]

	return (
		<article className="mx-auto max-w-[1040px] rounded-2xl border border-border bg-card p-5 text-left sm:p-7">
			<header className="flex flex-wrap items-center justify-between gap-3">
				<span className="font-mono text-[12px] text-foreground">
					{copy.label}
				</span>
				<span className="inline-flex items-center gap-2 rounded-full border border-border px-2.5 py-1 text-[11px] text-status-warning">
					<span
						aria-hidden="true"
						className="size-1.5 rounded-full bg-status-warning"
					/>
					{copy.states.unknown}
				</span>
			</header>

			<dl className="m-0 mt-5 grid grid-cols-2 gap-x-6 gap-y-4 border-y border-border py-4 sm:grid-cols-5">
				<div className="grid gap-1">
					<dt className={term}>{copy.summary.components}</dt>
					<dd className="m-0 text-[14px]">
						{t('componentsValue', {
							count: runtimeDescriptor.components.length,
						})}
					</dd>
				</div>
				<div className="grid gap-1">
					<dt className={term}>{copy.summary.cpu}</dt>
					<dd className="m-0 font-mono text-[14px] tabular-nums">
						{number.format(runtimeDescriptor.cpuMillis)}m
					</dd>
				</div>
				<div className="grid gap-1">
					<dt className={term}>{copy.summary.memory}</dt>
					<dd className="m-0 font-mono text-[14px] tabular-nums">
						{number.format(runtimeDescriptor.memoryMiB)} MiB
					</dd>
				</div>
				<div className="grid gap-1">
					<dt className={term}>{copy.summary.desired}</dt>
					<dd className="m-0 text-[14px]">{copy.states.running}</dd>
				</div>
				<div className="col-span-2 grid gap-1 sm:col-span-1">
					<dt className={term}>{copy.summary.observed}</dt>
					<dd className="m-0 grid text-[14px] text-status-warning">
						{copy.states.unknown}
						<span className="text-[11px]">{copy.unknownDetail}</span>
					</dd>
				</div>
			</dl>

			<ol
				aria-label={copy.stepperLabel}
				className="m-0 mt-5 grid list-none grid-cols-2 gap-2 p-0 sm:grid-cols-3 lg:grid-cols-6"
			>
				{lifecycleStageIds.map((id, index) => {
					const selected = id === selectedId

					return (
						<li key={id}>
							<button
								aria-current={selected ? 'step' : undefined}
								className={cn(
									'flex w-full cursor-pointer items-center gap-2.5 rounded-xl border border-border px-3 py-2.5 text-left text-[13px] text-foreground motion-safe:transition-colors',
									'hover:bg-foreground/5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring',
									selected &&
										'border-primary/60 bg-primary/10 hover:bg-primary/10',
								)}
								onClick={() => setSelectedId(id)}
								type="button"
							>
								<span
									className={cn(
										'inline-flex size-6 shrink-0 items-center justify-center rounded-full border border-border font-mono text-[10px] text-muted-foreground',
										selected &&
											'border-primary bg-primary text-primary-foreground',
									)}
								>
									{index + 1}
								</span>
								{copy.stages[id].label}
							</button>
						</li>
					)
				})}
			</ol>

			<div
				aria-atomic="true"
				aria-live="polite"
				className="mt-5 grid gap-4 rounded-xl bg-foreground/[0.03] p-4 sm:grid-cols-2 sm:p-5"
			>
				<div className="grid content-start gap-1.5">
					<span className={term}>{copy.whatItDoes}</span>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.does}</p>
				</div>
				<div className="grid content-start gap-1.5">
					<span className={term}>{copy.evidenceItLeaves}</span>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.leaves}</p>
				</div>
			</div>
		</article>
	)
}
