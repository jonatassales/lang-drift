import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { readConsole } from '@data/server'
import type { RuntimeDetailSnapshot } from '@data/types'
import { RuntimeEdit } from '@views/runtimes/runtime-edit'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('configuration'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

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
