import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { projectTitle } from '@app/metadata/entity-title'
import { FormSkeleton } from '@components/route-state/form-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { RuntimeCreate } from '@views/runtimes/runtime-create'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'createRuntime')

	return { title }
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
