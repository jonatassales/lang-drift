import { createTranslator } from 'next-intl'
import { Fragment } from 'react'

import type { WebsiteLocale } from '@i18n/routing'
import {
	dollars,
	type FinopsMessages,
	type GateFigures,
	type GateStepId,
	gateStepIds,
} from '@routes/platform/finops/content/finops-content'
import { Eyebrow } from 'web/ui/eyebrow'
import { ArrowRightIcon } from 'web/ui/icons/arrow'

interface GateAmount {
	amount: number
	label: string
}

// From the lg breakpoint the steps share four rows through a subgrid, so the
// eyebrows, the amounts, the titles and the descriptions line up across the
// row although the last step carries two amounts.
const arrow = (
	<li
		aria-hidden="true"
		className="flex items-center justify-center text-muted-foreground lg:row-span-4"
	>
		<ArrowRightIcon className="rotate-90 lg:rotate-0" size={20} />
	</li>
)

// How a paid call is admitted, as one flow: the agent recommends, the gate
// reserves, the evidence closes. Each step carries the amount one work item
// had at that point, so the amounts move across the row.
interface GateFlowProps {
	copy: FinopsMessages['authorization']
	exampleTitle: string
	figures: GateFigures
	locale: WebsiteLocale
}

export function GateFlow(props: GateFlowProps) {
	const { copy, exampleTitle, figures, locale } = props

	const money = dollars(locale)
	const message = createTranslator({
		locale,
		messages: copy,
	})
	const amounts: Record<GateStepId, readonly GateAmount[]> = {
		recommend: [
			{
				amount: figures.proposed,
				label: copy.steps.recommend.figure,
			},
		],
		gate: [
			{
				amount: figures.reserved,
				label: copy.steps.gate.figure,
			},
		],
		reconcile: [
			{
				amount: figures.observed,
				label: copy.steps.reconcile.figure,
			},
			{
				amount: figures.released,
				label: copy.released,
			},
		],
	}

	return (
		<>
			<ol className="m-0 mt-11 grid list-none gap-3 p-0 lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)_auto_minmax(0,1fr)]">
				{gateStepIds.map((id, index) => (
					<Fragment key={id}>
						{index > 0 && arrow}
						<li className="flex flex-col gap-2 rounded-2xl border border-border bg-card p-6 lg:row-span-4 lg:grid lg:grid-rows-subgrid">
							<Eyebrow>{copy.steps[id].eyebrow}</Eyebrow>
							<dl className="m-0 grid content-start gap-4 py-2">
								{amounts[id].map((figure) => (
									<div className="grid gap-1.5" key={figure.label}>
										<Eyebrow as="dt">{figure.label}</Eyebrow>
										<dd className="m-0 font-mono text-[28px] leading-none font-medium tracking-[-0.02em] tabular-nums">
											{money.format(figure.amount)}
										</dd>
									</div>
								))}
							</dl>
							<h3 className="m-0 text-lg leading-snug font-medium tracking-[-0.02em]">
								{copy.steps[id].title}
							</h3>
							<p className="m-0 text-[13px] leading-[1.55] text-muted-foreground">
								{copy.steps[id].description}
							</p>
						</li>
					</Fragment>
				))}
			</ol>
			<p className="mt-6 mb-0 max-w-[760px] text-[13px] leading-[1.6] text-muted-foreground">
				{message('source', { item: exampleTitle })}
			</p>
		</>
	)
}
