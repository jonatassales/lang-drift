import { Card } from '@langdrift/web/ui/card'
import { getFormatter, getTranslations } from 'next-intl/server'

import { PageHeading } from '@components/page-heading/page-heading'
import { readConsole } from '@data/server'
import type { ProductNavigation, RuntimeListSnapshot } from '@data/types'
import { Link } from '@i18n/navigation'
import { RecommendationPanel } from '@views/runtimes/recommendation-panel'
import { RuntimeFinops } from '@views/runtimes/runtime-finops'

import { costTotals } from './cost-totals'
import { Reading } from './reading'

// FinOps' opening page: what every runtime reported against its limits and
// budget, read from the stub. A value nothing reported stays unknown.
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
	const plain = (value: number) => format.number(value)
	const reported = (
		value: number | null | undefined,
		show: (value: number) => string = plain,
	) => (value == null ? t('unknown') : show(value))

	return (
		<>
			<PageHeading
				title={t('costOverview')}
				description={t('costOverviewDescription')}
			/>
			<Card className="block mb-8 p-4 sm:max-w-sm">
				<p className="text-xs text-muted-foreground">{t('knownLlmCost')}</p>
				<p className="mt-2 text-2xl font-medium">
					{reported(totals.knownLlmCostUsd, usd)}
				</p>
				<p className="mt-1 text-xs text-muted-foreground">
					{totals.executionCount === 0
						? t('noExecutionsYet')
						: t('unreportedLlmCosts', {
								count: totals.unknownLlmCostCount,
								total: totals.executionCount,
							})}
				</p>
			</Card>
			{totals.runtimes.length === 0 ? (
				<p className="text-sm text-muted-foreground">{t('noRuntimes')}</p>
			) : (
				<div className="grid gap-4">
					{totals.runtimes.map((cost) => {
						const { runtime } = cost.snapshot
						const { configuration } = runtime

						return (
							<Card className="block p-4" key={runtime.id}>
								<div className="flex flex-wrap items-baseline justify-between gap-2">
									<h2 className="text-sm font-medium">
										<Link
											className="focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
											href={`/runtimes/${runtime.id}`}
										>
											{runtime.name}
										</Link>
									</h2>
									<span className="text-xs text-muted-foreground">
										{projectNames.get(runtime.productId) ?? t('unknown')}
									</span>
								</div>
								<dl className="mt-3 grid gap-3 text-xs sm:grid-cols-3">
									<Reading label={t('llmCostBudget')}>
										{reported(cost.knownLlmCostUsd, usd)} /{' '}
										{reported(configuration.llmBudgetUsd, usd)}
									</Reading>
									<Reading label={t('cpuAgainstLimit')}>
										{reported(cost.cpuUsed)} /{' '}
										{format.number(configuration.cpu)}
									</Reading>
									<Reading label={t('memoryAgainstLimit')}>
										{reported(cost.memoryUsedMiB)} /{' '}
										{format.number(configuration.memoryMiB)}
									</Reading>
								</dl>
								{cost.unknownLlmCostCount > 0 ? (
									<p className="mt-2 text-xs text-muted-foreground">
										{t('unreportedLlmCosts', {
											count: cost.unknownLlmCostCount,
											total: cost.executionCount,
										})}
									</p>
								) : null}
								<details className="mt-4">
									<summary className="cursor-pointer text-xs focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring">
										{t('finopsReport')}
									</summary>
									<div className="mt-3">
										<RuntimeFinops snapshot={cost.snapshot} />
									</div>
								</details>
								<details className="mt-2">
									<summary className="cursor-pointer text-xs focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring">
										{t('recommendations')}
									</summary>
									<div className="mt-3">
										<RecommendationPanel
											endpoint={`runtimes/${runtime.id}/recommendations`}
											productId={runtime.productId}
										/>
									</div>
								</details>
							</Card>
						)
					})}
				</div>
			)}
		</>
	)
}
