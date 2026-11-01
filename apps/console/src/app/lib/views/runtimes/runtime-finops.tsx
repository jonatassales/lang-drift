'use client'
import { useFormatter, useTranslations } from 'next-intl'

import type { RuntimeWidgetSnapshot } from '@data/types'

interface RuntimeFinopsProps {
	snapshot: RuntimeWidgetSnapshot
}
export function RuntimeFinops(props: RuntimeFinopsProps) {
	const { snapshot } = props

	const translate = useTranslations('ui')
	const status = useTranslations('status')
	const format = useFormatter()
	const { runtime } = snapshot

	return (
		<section className="rounded-lg border border-border p-4">
			<h2 className="text-sm font-medium">{translate('finopsReport')}</h2>
			<div className="mt-4 grid gap-4 text-xs">
				<p className="leading-relaxed text-muted-foreground">
					{translate('finopsBasis')}
				</p>
				<dl className="grid gap-2 sm:grid-cols-2">
					<dt>{translate('configuredLimits')}</dt>
					<dd>
						{format.number(runtime.configuration.cpu)} CPU ·{' '}
						{format.number(runtime.configuration.memoryMiB)} MiB
					</dd>
					<dt>{translate('llmBudget')}</dt>
					<dd>
						{runtime.configuration.llmBudgetUsd === undefined
							? translate('unknown')
							: format.number(runtime.configuration.llmBudgetUsd, {
									style: 'currency',
									currency: 'USD',
								})}
					</dd>
				</dl>
				{snapshot.executions.map((execution) => (
					<section
						key={execution.id}
						className="grid gap-2 border-t border-border pt-3"
					>
						<p className="break-all font-mono text-[10px]">
							{execution.environmentId} / {execution.id}
						</p>
						<p>
							{translate('measurementWindow')}:{' '}
							{execution.measuredFrom ?? translate('unknown')} →{' '}
							{execution.observedAt ?? translate('unknown')}
						</p>
						<p>
							{translate('downtime')}:{' '}
							{execution.downtimeMilliseconds === null
								? translate('unknown')
								: translate('downtimeReport', {
										milliseconds: execution.downtimeMilliseconds,
										quality:
											execution.downtimeQuality === 'unknown'
												? translate('unknown')
												: status(execution.downtimeQuality),
									})}
						</p>
					</section>
				))}
			</div>
		</section>
	)
}
