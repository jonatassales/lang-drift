'use client'

import { useAtomValue } from 'jotai'
import { createTranslator } from 'next-intl'

import { dollars } from '@app/template/formatters/dollars.fmt'
import type { WebsiteLocale } from '@i18n/routing'
import { environments } from '@routes/platform/finops/content/finops-content.data'
import type { FinopsMessages } from '@routes/platform/finops/content/finops-content.types'
import { selectedEnvironmentAtom } from '@state'
import { DetailList } from 'web/ui/detail-list'
import { Eyebrow } from 'web/ui/eyebrow'

// The receipt of the environment selected in the ledger: what it asked for,
// what it consumed, how long it sat idle, and whether it was released. A
// usage its meter did not report reads unknown, never zero.
interface EnvironmentReceiptProps {
	copy: FinopsMessages['receipt']
	environmentsCopy: FinopsMessages['ledger']['environments']
	locale: WebsiteLocale
}

export function EnvironmentReceipt(props: EnvironmentReceiptProps) {
	const { copy, environmentsCopy, locale } = props

	const selectedId = useAtomValue(selectedEnvironmentAtom)

	function measured(value: number | null, render: (known: number) => string) {
		if (value === null) return copy.unknownUsage

		return render(value)
	}

	const environment = environments.find((entry) => entry.id === selectedId)

	if (!environment) return null

	const translate = createTranslator({
		locale,
		messages: copy,
	})
	const money = dollars(locale)
	const environmentCopy = environmentsCopy[environment.id]
	const cpu = measured(environment.cpuHours, (hours) =>
		translate('cpuHours', { hours }),
	)
	const memory = measured(environment.memoryGbHours, (hours) =>
		translate('memoryHours', { hours }),
	)
	const gpu = measured(environment.gpuHours, (hours) =>
		translate('gpuHours', { hours }),
	)
	const peak = measured(environment.peakCores, (cores) =>
		translate('peakCores', { cores }),
	)
	const idle = measured(environment.idleHours, (hours) =>
		translate('idleValue', {
			cost: money.format(environment.idle),
			hours,
		}),
	)

	return (
		<article
			aria-atomic="true"
			aria-live="polite"
			className="rounded-2xl border border-border bg-card p-6 max-md:p-[18px]"
		>
			<div className="grid gap-5">
				<div className="grid gap-[7px]">
					<Eyebrow>{copy.selected}</Eyebrow>
					<h2 className="m-0 text-[28px] font-medium tracking-[-0.025em]">
						{environmentCopy.title}
					</h2>
					<p className="m-0 text-[12px] text-muted-foreground">
						{environmentCopy.owner}
					</p>
				</div>
				<div className="grid gap-[7px]">
					<Eyebrow>{copy.labels.observed}</Eyebrow>
					{environment.observed === null ? (
						<strong className="text-[18px] font-medium text-status-warning">
							{copy.unknownObserved}
						</strong>
					) : (
						<strong className="font-mono text-[28px] font-medium tabular-nums">
							{money.format(environment.observed)}
						</strong>
					)}
				</div>
			</div>
			<DetailList
				className="mt-[22px] border-t border-border pt-1"
				items={[
					{
						detail: money.format(environment.ceiling),
						id: 'ceiling',
						term: copy.labels.ceiling,
					},
					{
						detail: cpu,
						id: 'cpu',
						term: copy.labels.cpu,
					},
					{
						detail: memory,
						id: 'memory',
						term: copy.labels.memory,
					},
					{
						detail: gpu,
						id: 'gpu',
						term: copy.labels.gpu,
					},
					{
						detail: peak,
						id: 'peak',
						term: copy.labels.peak,
					},
					{
						detail: idle,
						id: 'idle',
						term: copy.labels.idle,
					},
					{
						detail: copy.cleanup[environment.cleanup],
						id: 'cleanup',
						term: copy.labels.cleanup,
						tone: environment.cleanup === 'done' ? 'default' : 'warning',
					},
					{
						detail: (
							<code className="font-mono text-[12px] break-all">
								{environment.source}
							</code>
						),
						id: 'source',
						term: copy.labels.source,
					},
				]}
			/>
		</article>
	)
}
