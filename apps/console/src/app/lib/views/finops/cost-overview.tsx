import { getFormatter, getTranslations } from 'next-intl/server'

import { PageHeading } from '@components/page-heading/page-heading'
import { readConsole } from '@data/server'
import type { ProductNavigation, RuntimeListSnapshot } from '@data/types'
import { Link } from '@i18n/navigation'
import { Card } from '@ui/card'
import { RuntimeFinops } from '@views/runtimes/runtime-finops'

import { costTotals } from './cost-totals'

// FinOps' opening page: what every runtime reported against its limits and
// budget, read from the stub. A reading no execution reported stays unknown.
export async function CostOverview() {
	const [list, navigation] = await Promise.all([
		readConsole<RuntimeListSnapshot>('runtimes'),
		readConsole<ProductNavigation>('products'),
	])

	const t = await getTranslations('ui')
	const format = await getFormatter()
	const totals = costTotals(list.widgets)
	const projectNames = new Map(
		navigation.products.map((product) => [product.id, product.name]),
	)
	const usd = (value: number) =>
		format.number(value, {
			style: 'currency',
			currency: 'USD',
		})
	const measured = (value: number | null, show: (value: number) => string) =>
		value === null ? t('unknown') : show(value)

	return (
		<>
			<PageHeading
				title={t('costOverview')}
				description={t('costOverviewDescription')}
			/>
			<section className="mb-8 grid gap-3 sm:grid-cols-2">
				<Card className="p-4">
					<p className="text-xs text-muted-foreground">{t('knownLlmCost')}</p>
					<p className="mt-2 text-2xl font-medium">
						{usd(totals.knownLlmCostUsd)}
					</p>
					<p className="mt-1 text-xs text-muted-foreground">
						{t('unreportedLlmCosts', {
							count: totals.unknownLlmCostCount,
							total: totals.executionCount,
						})}
					</p>
				</Card>
				<Card className="p-4">
					<p className="text-xs text-muted-foreground">{t('llmBudgets')}</p>
					<p className="mt-2 text-2xl font-medium">
						{usd(totals.budgetedLlmUsd)}
					</p>
					<p className="mt-1 text-xs text-muted-foreground">
						{t('unbudgetedRuntimes', { count: totals.unbudgetedRuntimeCount })}
					</p>
				</Card>
			</section>
			{totals.runtimes.length === 0 ? (
				<p className="text-sm text-muted-foreground">{t('noRuntimes')}</p>
			) : (
				<div className="grid gap-4">
					{totals.runtimes.map((cost) => {
						const { runtime } = cost.snapshot
						const { configuration } = runtime

						return (
							<Card className="p-4" key={runtime.id}>
								<div className="flex flex-wrap items-baseline justify-between gap-2">
									<h2 className="text-sm font-medium">
										<Link href={`/runtimes/${runtime.id}`}>{runtime.name}</Link>
									</h2>
									<span className="text-xs text-muted-foreground">
										{projectNames.get(runtime.productId) ?? t('unknown')}
									</span>
								</div>
								<dl className="mt-3 grid gap-3 text-xs sm:grid-cols-3">
									<div className="grid gap-1">
										<dt className="text-muted-foreground">
											{t('llmCostBudget')}
										</dt>
										<dd>
											{measured(cost.knownLlmCostUsd, usd)} /{' '}
											{configuration.llmBudgetUsd === undefined
												? t('unknown')
												: usd(configuration.llmBudgetUsd)}
										</dd>
									</div>
									<div className="grid gap-1">
										<dt className="text-muted-foreground">
											{t('cpuAgainstLimit')}
										</dt>
										<dd>
											{measured(cost.cpuUsed, (value) => format.number(value))}{' '}
											/ {format.number(configuration.cpu)}
										</dd>
									</div>
									<div className="grid gap-1">
										<dt className="text-muted-foreground">
											{t('memoryAgainstLimit')}
										</dt>
										<dd>
											{measured(cost.memoryUsedMiB, (value) =>
												format.number(value),
											)}{' '}
											/ {format.number(configuration.memoryMiB)}
										</dd>
									</div>
								</dl>
								<details className="mt-4">
									<summary className="cursor-pointer text-xs">
										{t('finopsReport')}
									</summary>
									<div className="mt-3">
										<RuntimeFinops snapshot={cost.snapshot} />
									</div>
								</details>
								<Link
									className="mt-3 inline-block text-xs"
									href={`/runtimes/${runtime.id}/history`}
								>
									{t('openRecommendations')}
								</Link>
							</Card>
						)
					})}
				</div>
			)}
		</>
	)
}
