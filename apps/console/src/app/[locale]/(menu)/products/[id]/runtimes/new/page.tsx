import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { readConsole } from '@data/server'
import { RuntimeCreate } from '@views/runtimes/runtime-create'

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

	return (
		<Suspense fallback={<FormSkeleton />}>
			<RuntimeCreate productId={id} initial={readConsole(`products/${id}`)} />
		</Suspense>
	)
}
