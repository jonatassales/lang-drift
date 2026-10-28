'use client'
import { useTranslations } from 'next-intl'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import type { RuntimeWidgetSnapshot } from '@data/types'

import { RuntimeMeter } from './runtime-meter'

interface RuntimeExecutionMetricsProps {
	snapshot: RuntimeWidgetSnapshot
}

export function RuntimeExecutionMetrics(props: RuntimeExecutionMetricsProps) {
	const { snapshot } = props

	const t = useTranslations('ui')
	const dateTime = useDateTime()
	const { runtime, executions, environments, startedAt } = snapshot

	return (
		<section className="grid gap-3">
			<h2 className="text-sm font-medium">{t('executionMetrics')}</h2>
			{!executions.length ? (
				<p className="text-xs text-muted-foreground">
					{t('noExecutionObservations')}
				</p>
			) : null}
			{executions.map((execution) => {
				const start = startedAt[execution.id]
				const minutes =
					start && execution.observedAt
						? Math.max(
								0,
								Math.floor(
									(Date.parse(execution.observedAt) - Date.parse(start)) /
										60000,
								),
							)
						: null

				return (
					<div
						key={execution.id}
						className="rounded-lg border border-border bg-card p-4"
					>
						<p className="mb-3 text-xs font-medium">
							{environments.find(
								(environment) => environment.id === execution.environmentId,
							)?.name ?? execution.environmentId}
						</p>
						<div className="grid gap-4 sm:grid-cols-3">
							<RuntimeMeter
								label={t('cpu')}
								value={execution.cpuUsed}
								limit={runtime.configuration.cpu}
							/>
							<RuntimeMeter
								label={t('memory')}
								value={execution.memoryUsedMiB}
								limit={runtime.configuration.memoryMiB}
							/>
							<RuntimeMeter
								label={t('llmCost')}
								value={execution.llmCostUsd ?? null}
								limit={runtime.configuration.llmBudgetUsd ?? null}
							/>
						</div>
						<p className="mt-3 text-[11px] text-muted-foreground">
							{t('runningTime')}:{' '}
							{minutes === null
								? t('unknown')
								: t('durationMinutes', { minutes })}
						</p>
						<p className="mt-1 text-[11px] text-muted-foreground">
							{t('lastSeen')}:{' '}
							{execution.observedAt
								? dateTime(new Date(execution.observedAt), {
										timeStyle: 'medium',
									})
								: t('unknown')}
						</p>
					</div>
				)
			})}
		</section>
	)
}
