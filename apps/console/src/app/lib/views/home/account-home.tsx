'use client'

import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { useResource } from '@data/resource.hook'
import type { HomeOverview } from '@data/types'
import { Link } from '@i18n/navigation'
import { ProductCreate } from '@views/products/product-create'

interface AccountHomeProps {
	initial: Promise<HomeOverview>
}

export function AccountHome(props: AccountHomeProps) {
	const { initial } = props

	const resource = useResource<HomeOverview>('home', 5000, initial)
	const translate = useTranslations('ui')

	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (!resource.data) return null
	if (!resource.data.products.length) return <ProductCreate first />

	return (
		<>
			<PageHeading
				title={translate('home')}
				description={translate('runningRuntimes')}
			/>
			{resource.data.unknownCount ? (
				<p role="status" className="mb-5 text-sm text-muted-foreground">
					{translate('unknownExecutions', {
						count: resource.data.unknownCount,
					})}
				</p>
			) : null}
			<div className="grid gap-4 lg:grid-cols-2">
				{resource.data.running.map(({ runtime, productName, executions }) => (
					<Card key={runtime.id} className="block p-5">
						<Link
							className="text-xs text-muted-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
							href={`/products/${runtime.productId}`}
						>
							{productName}
						</Link>
						<Link
							className="mt-2 block text-base font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
							href={`/runtimes/${runtime.id}`}
						>
							{runtime.name}
						</Link>
						<p className="mt-4 text-xs text-muted-foreground">
							{translate('runningExecutions', {
								count: executions.length,
							})}
						</p>
					</Card>
				))}
			</div>
			{!resource.data.running.length ? (
				<p className="text-sm text-muted-foreground">
					{translate('noRunningRuntimes')}
				</p>
			) : null}
		</>
	)
}
