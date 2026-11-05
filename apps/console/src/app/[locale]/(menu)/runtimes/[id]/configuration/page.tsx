import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type { RuntimeDetailSnapshot } from '@data/types'
import { RuntimeEdit } from '@views/runtimes/runtime-edit'

export { generateMetadata } from './page.metadata'

interface RuntimeRouteParams {
	id: string
}

interface PageProps {
	params: Promise<RuntimeRouteParams>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()
	const detail = readConsole<RuntimeDetailSnapshot>(`runtimes/${id}`)

	return (
		<Suspense fallback={<FormSkeleton />}>
			<RuntimeEdit
				id={id}
				initial={detail}
				productInitial={detail.then((value) =>
					readConsole(`products/${value.runtime.productId}`),
				)}
			/>
		</Suspense>
	)
}
