import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { RuntimeListSkeleton } from '@components/route-state/runtime-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { ProductDetail } from '@views/products/product-detail'

export { generateMetadata } from './page.metadata'

interface ProductRouteParams {
	id: string
}

interface ProductPageProps {
	params: Promise<ProductRouteParams>
}

export default async function Page(props: ProductPageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()

	return (
		<Suspense fallback={<RuntimeListSkeleton />}>
			<ProductDetail
				id={id}
				initial={readConsole(`products/${id}`)}
				runtimes={readConsole(`runtimes?productId=${encodeURIComponent(id)}`)}
			/>
		</Suspense>
	)
}
