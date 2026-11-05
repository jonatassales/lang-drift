'use client'
import { useTranslations } from 'next-intl'

import type { RuntimeWidgetSnapshot } from '@data/types'
import { useDateTime } from '@i18n/use-date-time.hook'

import { RuntimeMeter } from './runtime-meter'

interface RuntimeExecutionMetricsProps {
	snapshot: RuntimeWidgetSnapshot
}

export function RuntimeExecutionMetrics(props: RuntimeExecutionMetricsProps) {
	const { snapshot } = props

	const translate = useTranslations('ui')
	const dateTime = useDateTime()
	const { runtime, executions, environments, startedAt } = snapshot

	return (
		<section className="grid gap-3">
			<h2 className="text-sm font-medium">{translate('executionMetrics')}</h2>
			{!executions.length ? (
				<p className="text-xs text-muted-foreground">
					{translate('noExecutionObservations')}
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
								label={translate('cpu')}
								value={execution.cpuUsed}
								limit={runtime.configuration.cpu}
							/>
							<RuntimeMeter
								label={translate('memory')}
								value={execution.memoryUsedMiB}
								limit={runtime.configuration.memoryMiB}
							/>
							<RuntimeMeter
								label={translate('llmCost')}
								value={execution.llmCostUsd ?? null}
								limit={runtime.configuration.llmBudgetUsd ?? null}
							/>
						</div>
						<p className="mt-3 text-[11px] text-muted-foreground">
							{translate('runningTime')}:{' '}
							{minutes === null
								? translate('unknown')
								: translate('durationMinutes', { minutes })}
						</p>
						<p className="mt-1 text-[11px] text-muted-foreground">
							{translate('lastSeen')}:{' '}
							{execution.observedAt
								? dateTime(new Date(execution.observedAt), {
										timeStyle: 'medium',
									})
								: translate('unknown')}
						</p>
					</div>
				)
			})}
		</section>
	)
}
