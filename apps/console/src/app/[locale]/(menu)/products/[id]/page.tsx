import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { projectTitle } from '@app/metadata/entity-title'
import { RuntimeListSkeleton } from '@components/route-state/runtime-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { ProductDetail } from '@views/products/product-detail'

interface ProductPageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: ProductPageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'runtimes')

	return { title }
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
				runtimes={readConsole(`runtimes?productId=${id}`)}
			/>
		</Suspense>
	)
}
