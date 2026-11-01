'use client'
import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import { ErrorNotice } from '@components/error-notice/error-notice'
import { useResource } from '@data/resource.hook'
import type { LifecycleEvent } from '@data/types'

interface EventHistoryProps {
	initial: Promise<{ events: LifecycleEvent[] }>
	runtimeId: string
}
export function EventHistory(props: EventHistoryProps) {
	const { runtimeId, initial } = props

	const t = useTranslations('ui')
	const causes = useTranslations('cause')
	const kinds = useTranslations('event')
	const dateTime = useDateTime()
	const resource = useResource<{ events: LifecycleEvent[] }>(
		`runtimes/${runtimeId}/history`,
		5000,
		initial,
	)
	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />

	return (
		<Card className="block min-w-0 py-0">
			{!resource.data?.events.length ? (
				<p className="p-6 text-sm text-muted-foreground">{t('empty')}</p>
			) : (
				<div className="overflow-x-auto">
					<table className="w-full text-left text-xs">
						<thead className="border-b border-border text-muted-foreground">
							<tr>
								{[
									'event',
									'cause',
									'actor',
									'environments',
									'occurredAt',
									'receivedAt',
								].map((key) => (
									<th
										className="whitespace-nowrap px-4 py-3 font-medium"
										key={key}
									>
										{t(key)}
									</th>
								))}
							</tr>
						</thead>
						<tbody className="divide-y divide-border">
							{resource.data.events.map((event) => (
								<tr key={event.id}>
									<td className="px-4 py-4">
										<span className="font-medium">
											{kinds.has(event.kind) ? kinds(event.kind) : event.kind}
										</span>
										<p className="mt-1 max-w-44 break-all font-mono text-[10px] text-muted-foreground">
											{event.id}
										</p>
										<details className="mt-2">
											<summary className="cursor-pointer text-[11px] text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring">
												{t('evidence')}
											</summary>
											<dl className="mt-2 grid gap-1 text-[10px]">
												{event.revision !== undefined ? (
													<>
														<dt>{t('revision')}</dt>
														<dd>{event.revision}</dd>
													</>
												) : null}
												{event.operationId ? (
													<>
														<dt>{t('operation')}</dt>
														<dd className="break-all font-mono">
															{event.operationId}
														</dd>
													</>
												) : null}
												{event.requestedBy ? (
													<>
														<dt>{t('actor')}</dt>
														<dd className="break-all">{event.requestedBy}</dd>
													</>
												) : null}
												{event.sourceRevision ? (
													<>
														<dt>{t('sourceRevision')}</dt>
														<dd className="break-all font-mono">
															{event.sourceRevision}
														</dd>
													</>
												) : null}
												{event.concreteHash ? (
													<>
														<dt>{t('configurationHash')}</dt>
														<dd className="break-all font-mono">
															{event.concreteHash}
														</dd>
													</>
												) : null}
												{event.markdownHash ? (
													<>
														<dt>{t('markdownHash')}</dt>
														<dd className="break-all font-mono">
															{event.markdownHash}
														</dd>
													</>
												) : null}
											</dl>
										</details>
										{event.detail ? (
											<p className="mt-2 max-w-sm whitespace-pre-wrap text-[11px] text-muted-foreground">
												{event.detail}
											</p>
										) : null}
									</td>
									<td className="px-4 py-4">
										{causes.has(event.cause)
											? causes(event.cause)
											: event.cause}
									</td>
									<td className="max-w-36 break-all px-4 py-4 text-[10px]">
										{event.actor}
									</td>
									<td className="max-w-36 break-all px-4 py-4 font-mono text-[10px]">
										{event.environmentId ?? '—'}
										{event.executionId ? (
											<p className="mt-2 text-muted-foreground">
												{event.executionId}
											</p>
										) : null}
									</td>
									<td className="whitespace-nowrap px-4 py-4 text-muted-foreground">
										{dateTime(new Date(event.occurredAt), {
											dateStyle: 'short',
											timeStyle: 'medium',
										})}
									</td>
									<td className="whitespace-nowrap px-4 py-4 text-muted-foreground">
										{dateTime(new Date(event.receivedAt), {
											dateStyle: 'short',
											timeStyle: 'medium',
										})}
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			)}
		</Card>
	)
}
