'use client'
import { useFormatter, useTranslations } from 'next-intl'

import type { RuntimeWidgetSnapshot } from '@data/types'

interface RuntimeFinopsProps {
	snapshot: RuntimeWidgetSnapshot
}
export function RuntimeFinops(props: RuntimeFinopsProps) {
	const { snapshot } = props

	const t = useTranslations('ui')
	const status = useTranslations('status')
	const format = useFormatter()
	const { runtime } = snapshot

	return (
		<section className="rounded-lg border border-border p-4">
			<h2 className="text-sm font-medium">{t('finopsReport')}</h2>
			<div className="mt-4 grid gap-4 text-xs">
				<p className="leading-relaxed text-muted-foreground">
					{t('finopsBasis')}
				</p>
				<dl className="grid gap-2 sm:grid-cols-2">
					<dt>{t('configuredLimits')}</dt>
					<dd>
						{format.number(runtime.configuration.cpu)} CPU ·{' '}
						{format.number(runtime.configuration.memoryMiB)} MiB
					</dd>
					<dt>{t('llmBudget')}</dt>
					<dd>
						{runtime.configuration.llmBudgetUsd === undefined
							? t('unknown')
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
							{t('measurementWindow')}: {execution.measuredFrom ?? t('unknown')}{' '}
							→ {execution.observedAt ?? t('unknown')}
						</p>
						<p>
							{t('downtime')}:{' '}
							{execution.downtimeMilliseconds === null
								? t('unknown')
								: t('downtimeReport', {
										milliseconds: execution.downtimeMilliseconds,
										quality:
											execution.downtimeQuality === 'unknown'
												? t('unknown')
												: status(execution.downtimeQuality),
									})}
						</p>
					</section>
				))}
			</div>
		</section>
	)
}
