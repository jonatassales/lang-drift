import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { ProductSecrets } from '@views/products/product-secrets'

export { generateMetadata } from './page.metadata'

interface PageProps {
	params: Promise<{ id: string }>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()

	return (
		<Suspense fallback={<FormSkeleton />}>
			<ProductSecrets
				productId={id}
				initial={readConsole(`products/${id}/secrets`)}
			/>
		</Suspense>
	)
}
