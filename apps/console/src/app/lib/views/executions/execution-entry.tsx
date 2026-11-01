'use client'

import { useFormatter, useTranslations } from 'next-intl'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import { StatusLabel } from '@components/status-label/status-label'
import type { ProductExecution } from '@data/types'
import { Link } from '@i18n/navigation'

interface ExecutionEntryProps {
	item: ProductExecution
}

export function ExecutionEntry(props: ExecutionEntryProps) {
	const { item } = props

	const translate = useTranslations('ui')
	const format = useFormatter()
	const dateTime = useDateTime()
	const { execution } = item

	return (
		<article className="grid grid-cols-1 gap-3 rounded-xl border border-border bg-card p-4 text-xs">
			<header className="flex flex-wrap items-center justify-between gap-3">
				<Link
					prefetch={false}
					className="max-w-full break-words font-medium hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
					href={`/runtimes/${execution.runtimeId}`}
				>
					{item.runtimeName}
				</Link>
				<StatusLabel value={execution.state} />
			</header>
			<p className="break-all text-muted-foreground">
				{item.environmentName ?? execution.environmentId} · {execution.id}
			</p>
			<dl className="grid grid-cols-2 gap-2 sm:grid-cols-4">
				<div>
					<dt className="text-muted-foreground">{translate('cpuUsage')}</dt>
					<dd>
						{execution.cpuUsed === null
							? translate('unknown')
							: format.number(execution.cpuUsed, {
									maximumFractionDigits: 3,
								})}
					</dd>
				</div>
				<div>
					<dt className="text-muted-foreground">{translate('memoryUsage')}</dt>
					<dd>
						{execution.memoryUsedMiB === null
							? translate('unknown')
							: format.number(execution.memoryUsedMiB, {
									maximumFractionDigits: 1,
								})}
					</dd>
				</div>
				<div>
					<dt className="text-muted-foreground">{translate('revision')}</dt>
					<dd>{execution.revision}</dd>
				</div>
				<div>
					<dt className="text-muted-foreground">{translate('observedAt')}</dt>
					<dd>
						{execution.observedAt
							? dateTime(new Date(execution.observedAt), {
									dateStyle: 'short',
									timeStyle: 'medium',
								})
							: translate('unknown')}
					</dd>
				</div>
			</dl>
		</article>
	)
}
