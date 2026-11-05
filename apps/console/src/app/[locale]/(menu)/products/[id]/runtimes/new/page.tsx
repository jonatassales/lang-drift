import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { RuntimeCreate } from '@views/runtimes/runtime-create'

export { generateMetadata } from './page.metadata'

interface ProductRouteParams {
	id: string
}

interface PageProps {
	params: Promise<ProductRouteParams>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()

	return (
		<Suspense fallback={<FormSkeleton />}>
			<RuntimeCreate productId={id} initial={readConsole(`products/${id}`)} />
		</Suspense>
	)
}
