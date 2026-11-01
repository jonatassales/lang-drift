import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { RuntimeDetailSkeleton } from '@components/route-state/runtime-detail-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type { RuntimeDetailSnapshot } from '@data/types'
import { RuntimeDetail } from '@views/runtimes/runtime-detail'

export { generateMetadata } from './page.metadata'

interface PageProps {
	params: Promise<{ id: string }>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()
	const detail = readConsole<RuntimeDetailSnapshot>(`runtimes/${id}`)

	return (
		<Suspense fallback={<RuntimeDetailSkeleton />}>
			<RuntimeDetail
				id={id}
				initial={detail}
				productInitial={detail.then((value) =>
					readConsole(`products/${value.runtime.productId}`),
				)}
			/>
		</Suspense>
	)
}
