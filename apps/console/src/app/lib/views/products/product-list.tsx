'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { linkClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import type { ProductsEnvelope } from '@data/product-envelopes.types'
import { useResource } from '@data/resource.hook'
import { Link } from '@i18n/navigation'
import { Card } from 'web/ui/card'
import { GitBranch, Plus } from 'web/ui/icons/client'

import { ProductCreate } from './product-create'

interface ProductListProps {
	initial: Promise<ProductsEnvelope>
}

export function ProductList(props: ProductListProps) {
	const { initial } = props

	const translate = useTranslations('ui')
	const resource = useResource<ProductsEnvelope>('products', 0, initial)

	if (!resource.error && resource.data?.products.length === 0)
		return <ProductCreate first />

	return (
		<>
			<PageHeading
				title={translate('projects')}
				description={translate('localProjectsIntro')}
			>
				<Link href="/products/new" className={linkClass}>
					<Plus size={15} />
					{translate('createProject')}
				</Link>
			</PageHeading>
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.refresh} />
			) : resource.data?.products.length ? (
				<div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-3">
					{resource.data.products.map((product) => (
						<Card key={product.id} className="block p-5">
							<Link
								href={`/products/${product.id}`}
								className="text-sm font-medium hover:text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
							>
								{product.name}
							</Link>
							<p className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
								<GitBranch size={14} />
								{product.source?.path ?? product.repository?.fullName}
							</p>
							<div className="mt-5 flex items-center justify-between border-t border-border pt-3 text-[11px] text-muted-foreground">
								<span>
									{product.components.length} {translate('components')}
								</span>
								<code className="max-w-40 truncate">{product.ref}</code>
							</div>
						</Card>
					))}
				</div>
			) : (
				<Card className="block p-10 text-center text-sm text-muted-foreground">
					{translate('projectsEmpty')}
				</Card>
			)}
		</>
	)
}
