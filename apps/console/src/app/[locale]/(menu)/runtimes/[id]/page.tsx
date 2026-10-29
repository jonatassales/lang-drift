import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { runtimeTitle } from '@app/metadata/entity-title'
import { RuntimeDetailSkeleton } from '@components/route-state/runtime-detail-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type { RuntimeDetailSnapshot } from '@data/types'
import { RuntimeDetail } from '@views/runtimes/runtime-detail'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await runtimeTitle(id)

	return { title }
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
