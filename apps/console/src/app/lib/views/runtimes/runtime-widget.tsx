'use client'
import { useTranslations } from 'next-intl'

import { StatusLabel } from '@components/status-label/status-label'
import type { Execution, RuntimeWidgetSnapshot } from '@data/types'
import { Link } from '@i18n/navigation'

import { RuntimeFinops } from './runtime-finops'
import { RuntimeLoadTest } from './runtime-load-test'
import { RuntimeMeter } from './runtime-meter'
import { RuntimePlayback } from './runtime-playback'
import { RuntimeStackIcons } from './runtime-stack-icons'

interface RuntimeWidgetProps {
	snapshot: RuntimeWidgetSnapshot
	onRefresh: () => void
}

/** Highest available observation per execution, never a sum across machines. */
function peak(
	executions: Execution[],
	key: 'cpuUsed' | 'memoryUsedMiB' | 'llmCostUsd',
) {
	const values = executions
		.filter((execution) => execution.state === 'running')
		.map((execution) => execution[key])
		.filter((value): value is number => value != null)

	return values.length ? Math.max(...values) : null
}

export function RuntimeWidget(props: RuntimeWidgetProps) {
	const { snapshot, onRefresh } = props

	const t = useTranslations('ui')
	const { runtime, executions } = snapshot
	const state = executions.some((execution) => execution.state === 'running')
		? 'running'
		: (executions[0]?.state ??
			(runtime.desiredState === 'running' ? 'unknown' : runtime.desiredState))
	const stacks = Array.from(
		new Set(runtime.components.flatMap((component) => component.stacks ?? [])),
	)

	return (
		<article
			data-runtime-widget={runtime.id}
			className="min-w-0 rounded-xl border border-border bg-card px-4 py-3"
		>
			<header className="flex flex-wrap items-center justify-between gap-2">
				<Link
					href={`/runtimes/${runtime.id}`}
					className="min-w-0 truncate text-sm font-medium hover:underline"
				>
					{runtime.name}
				</Link>
				<div className="flex items-center gap-2">
					<StatusLabel value={state} />
					<RuntimePlayback snapshot={snapshot} onRefresh={onRefresh} />
				</div>
			</header>
			<div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px] text-muted-foreground">
				<RuntimeStackIcons stacks={stacks} size={14} />
				<span>{t('componentCount', { count: runtime.components.length })}</span>
				<span>
					{runtime.configuration.model} · {t(runtime.configuration.effort)}
				</span>
				{snapshot.partial ? <span>{t('partialTelemetry')}</span> : null}
			</div>
			<div className="mt-3 grid gap-2 sm:grid-cols-3 sm:gap-5">
				<RuntimeMeter
					label={t('cpu')}
					value={peak(executions, 'cpuUsed')}
					limit={runtime.configuration.cpu}
				/>
				<RuntimeMeter
					label={t('memory')}
					value={peak(executions, 'memoryUsedMiB')}
					limit={runtime.configuration.memoryMiB}
				/>
				<RuntimeMeter
					label={t('llmCost')}
					value={peak(executions, 'llmCostUsd')}
					limit={runtime.configuration.llmBudgetUsd ?? null}
				/>
			</div>
			<p className="mt-2 text-[10px] text-muted-foreground">
				{t('peakPerExecution')}
			</p>
			<details className="mt-2 border-t border-border pt-2">
				<summary className="cursor-pointer text-[11px] text-muted-foreground">
					{t('quickReport')}
				</summary>
				<div className="mt-3 grid gap-3 lg:grid-cols-2">
					<RuntimeFinops snapshot={snapshot} />
					<RuntimeLoadTest key={runtime.id} snapshot={snapshot} />
				</div>
			</details>
		</article>
	)
}
