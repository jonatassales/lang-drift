'use client'

import { useTranslations } from 'next-intl'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import type { ProductEvent } from '@data/types'
import { Link } from '@i18n/navigation'

interface HistoryEntryProps {
	item: ProductEvent
}

export function HistoryEntry(props: HistoryEntryProps) {
	const { item } = props

	const t = useTranslations('ui')
	const kinds = useTranslations('event')
	const causes = useTranslations('cause')
	const dateTime = useDateTime()
	const { event } = item

	return (
		<article className="grid grid-cols-1 gap-3 rounded-xl border border-border bg-card p-4 text-xs">
			<header className="flex flex-wrap items-center justify-between gap-3">
				<Link
					prefetch={false}
					className="max-w-full break-words font-medium hover:underline"
					href={`/runtimes/${event.runtimeId}/history`}
				>
					{item.runtimeName}
				</Link>
				<span>{kinds.has(event.kind) ? kinds(event.kind) : event.kind}</span>
			</header>
			<p className="break-words text-muted-foreground">
				{causes.has(event.cause) ? causes(event.cause) : event.cause} ·{' '}
				{event.actor}
			</p>
			<dl className="grid grid-cols-1 gap-2 sm:grid-cols-3">
				<div>
					<dt className="text-muted-foreground">{t('environments')}</dt>
					<dd className="break-all">
						{item.environmentName ?? event.environmentId ?? '—'}
					</dd>
				</div>
				<div>
					<dt className="text-muted-foreground">{t('occurredAt')}</dt>
					<dd>
						{dateTime(new Date(event.occurredAt), {
							dateStyle: 'short',
							timeStyle: 'medium',
						})}
					</dd>
				</div>
				<div>
					<dt className="text-muted-foreground">{t('receivedAt')}</dt>
					<dd>
						{dateTime(new Date(event.receivedAt), {
							dateStyle: 'short',
							timeStyle: 'medium',
						})}
					</dd>
				</div>
			</dl>
			<p className="break-all font-mono text-[10px] text-muted-foreground">
				{event.id}
				{event.executionId ? ` · ${event.executionId}` : ''}
			</p>
			{event.detail ? (
				<p className="whitespace-pre-wrap break-words text-muted-foreground">
					{event.detail}
				</p>
			) : null}
		</article>
	)
}
