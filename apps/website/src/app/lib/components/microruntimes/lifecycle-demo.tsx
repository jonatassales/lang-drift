'use client'

import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useAtom } from 'jotai'
import { createTranslator } from 'next-intl'

import { ChoiceButton } from '@components/product-page/choice-button'
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

// The Microruntimes hero demo: an illustrative runtime's last execution, its
// six lifecycle steps, and what the selected step does and leaves behind.
// Selecting a step drives the evidence section further down the page.
export function LifecycleDemo(props: LifecycleDemoProps) {
	const { copy, locale } = props

	const [selectedId, setSelectedId] = useAtom(selectedLifecycleStageAtom)
	const lifecycle = createTranslator({
		locale,
		messages: copy,
	})
	const number = new Intl.NumberFormat(locale)
	const stage = copy.stages[selectedId]
	const summary = [
		{
			id: 'components',
			term: copy.summary.components,
			value: lifecycle('componentsValue', {
				count: runtimeDescriptor.components.length,
			}),
		},
		{
			id: 'cpu',
			term: copy.summary.cpu,
			value: lifecycle('cpuValue', {
				value: number.format(runtimeDescriptor.cpuMillis),
			}),
		},
		{
			id: 'memory',
			term: copy.summary.memory,
			value: lifecycle('memoryValue', {
				value: number.format(runtimeDescriptor.memoryMiB),
			}),
		},
		{
			id: 'desired',
			term: copy.summary.desired,
			value: copy.states.stopped,
		},
	]

	return (
		<article className="@container rounded-2xl border border-border bg-card p-5 text-left sm:p-7">
			<span className="font-mono text-[12px] text-foreground">
				{copy.label}
			</span>

			<dl className="m-0 mt-5 grid grid-cols-2 gap-x-6 gap-y-4 border-y border-border py-4 @lg:grid-cols-3 @4xl:grid-cols-5">
				{summary.map((entry) => (
					<div className="grid gap-1" key={entry.id}>
						<Eyebrow as="dt">{entry.term}</Eyebrow>
						<dd className="m-0 font-mono text-[14px] tabular-nums">
							{entry.value}
						</dd>
					</div>
				))}
				<div className="col-span-2 grid gap-1 @lg:col-span-1">
					<Eyebrow as="dt">{copy.summary.observed}</Eyebrow>
					<dd className="m-0 grid text-[14px] text-status-success">
						{copy.states.cleanedUp}
						<span className="text-[11px]">{copy.observedDetail}</span>
					</dd>
				</div>
			</dl>

			<ol
				aria-label={copy.stepperLabel}
				className="m-0 mt-5 grid list-none grid-cols-2 gap-2 p-0 @md:grid-cols-3 @4xl:grid-cols-6"
			>
				{lifecycleStageIds.map((id, index) => {
					const selected = id === selectedId

					return (
						<li key={id}>
							<ChoiceButton
								className="flex items-center gap-2.5 border-border px-3 py-2.5 text-[13px]"
								kind="step"
								onSelect={() => setSelectedId(id)}
								selected={selected}
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
							</ChoiceButton>
						</li>
					)
				})}
			</ol>

			<div
				aria-atomic="true"
				aria-live="polite"
				className="mt-5 grid gap-4 rounded-xl bg-foreground/[0.03] p-4 @lg:grid-cols-2 sm:p-5"
			>
				<div className="grid content-start gap-1.5">
					<Eyebrow>{copy.whatItDoes}</Eyebrow>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.does}</p>
				</div>
				<div className="grid content-start gap-1.5">
					<Eyebrow>{copy.evidenceItLeaves}</Eyebrow>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.leaves}</p>
				</div>
			</div>
		</article>
	)
}
