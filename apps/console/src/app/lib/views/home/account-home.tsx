'use client'

import { useTranslations } from 'next-intl'
import { useEffect, useRef } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { useResource } from '@data/resource.hook'
import type { HomeOverview } from '@data/types'
import { Link, useRouter } from '@i18n/navigation'
import { Card } from '@ui/card'
import { ProductCreate } from '@views/products/product-create'

interface AccountHomeProps {
	initial: Promise<HomeOverview>
}

export function AccountHome(props: AccountHomeProps) {
	const { initial } = props

	const t = useTranslations('ui')
	const router = useRouter()
	const resource = useResource<HomeOverview>('home', 5000, initial)
	const decided = useRef(false)
	useEffect(() => {
		if (decided.current || !resource.data || resource.error) return
		decided.current = true
		if (resource.data.products.length && !resource.data.running.length)
			router.replace(`/products/${resource.data.products[0].id}`)
	}, [resource.data, resource.error, router])
	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (!resource.data) return null
	if (!resource.data.products.length) return <ProductCreate first />

	return (
		<>
			<PageHeading title={t('home')} description={t('runningRuntimes')} />
			{resource.data.unknownCount ? (
				<p role="status" className="mb-5 text-sm text-muted-foreground">
					{t('unknownExecutions', {
						count: resource.data.unknownCount,
					})}
				</p>
			) : null}
			<div className="grid gap-4 lg:grid-cols-2">
				{resource.data.running.map(({ runtime, productName, executions }) => (
					<Card key={runtime.id} className="p-5">
						<Link
							className="text-xs text-muted-foreground"
							href={`/products/${runtime.productId}`}
						>
							{productName}
						</Link>
						<Link
							className="mt-2 block text-base font-medium"
							href={`/runtimes/${runtime.id}`}
						>
							{runtime.name}
						</Link>
						<p className="mt-4 text-xs text-muted-foreground">
							{t('runningExecutions', {
								count: executions.length,
							})}
						</p>
					</Card>
				))}
			</div>
			{!resource.data.running.length ? (
				<p className="text-sm text-muted-foreground">
					{t('noRunningRuntimes')}
				</p>
			) : null}
		</>
	)
}
