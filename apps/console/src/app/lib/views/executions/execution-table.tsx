'use client'
import { Card } from '@langdrift/web/ui/card'
import { useFormatter, useTranslations } from 'next-intl'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import { StatusLabel } from '@components/status-label/status-label'
import type { Environment, Execution } from '@data/types'

interface ExecutionTableProps {
	executions: Execution[]
	environments: Environment[]
}
export function ExecutionTable(props: ExecutionTableProps) {
	const { executions, environments } = props

	const translate = useTranslations('ui')
	const format = useFormatter()
	const dateTime = useDateTime()

	return (
		<Card className="block py-0">
			{!executions.length ? (
				<p className="p-6 text-sm text-muted-foreground">
					{translate('empty')}
				</p>
			) : (
				<div className="overflow-x-auto">
					<table className="w-full text-left text-xs">
						<thead className="border-b border-border text-muted-foreground">
							<tr>
								{[
									'environments',
									'state',
									'cpuUsage',
									'memoryUsage',
									'downtime',
									'quality',
									'observedAt',
								].map((key) => (
									<th
										className="whitespace-nowrap px-4 py-3 font-medium"
										key={key}
									>
										{translate(key)}
									</th>
								))}
							</tr>
						</thead>
						<tbody className="divide-y divide-border">
							{executions.map((execution) => (
								<tr key={execution.id}>
									<td className="px-4 py-4">
										<p>
											{environments.find(
												(environment) =>
													environment.id === execution.environmentId,
											)?.name ?? execution.environmentId}
										</p>
										<p className="mt-1 max-w-44 break-all font-mono text-[10px] text-muted-foreground">
											{execution.id}
										</p>
										<p className="mt-1 text-[10px] text-muted-foreground">
											{translate('revision')} {execution.revision}
										</p>
									</td>
									<td className="px-4 py-4">
										<StatusLabel value={execution.state} />
									</td>
									<td className="px-4 py-4 tabular-nums">
										{execution.cpuUsed === null
											? translate('unknown')
											: format.number(execution.cpuUsed, {
													maximumFractionDigits: 3,
												})}
									</td>
									<td className="px-4 py-4 tabular-nums">
										{execution.memoryUsedMiB === null
											? translate('unknown')
											: format.number(execution.memoryUsedMiB, {
													maximumFractionDigits: 1,
												})}
									</td>
									<td className="px-4 py-4 tabular-nums">
										{execution.downtimeMilliseconds === null
											? translate('unknown')
											: format.number(execution.downtimeMilliseconds / 1000, {
													maximumFractionDigits: 1,
												})}
										{execution.downtimeStartedAt ? (
											<p className="mt-1 text-[10px] text-muted-foreground">
												{dateTime(new Date(execution.downtimeStartedAt), {
													dateStyle: 'short',
													timeStyle: 'medium',
												})}
											</p>
										) : null}
									</td>
									<td className="px-4 py-4">
										{execution.downtimeQuality === 'unknown' ? (
											translate('unknown')
										) : (
											<StatusLabel value={execution.downtimeQuality} />
										)}
									</td>
									<td className="whitespace-nowrap px-4 py-4 text-muted-foreground">
										{execution.observedAt
											? dateTime(new Date(execution.observedAt), {
													dateStyle: 'short',
													timeStyle: 'medium',
												})
											: translate('unknown')}
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
