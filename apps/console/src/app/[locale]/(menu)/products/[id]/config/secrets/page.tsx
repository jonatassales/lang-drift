import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { FormSkeleton } from '@components/route-state/form-skeleton'
import { readConsole } from '@data/server'
import { ProductSecrets } from '@views/products/product-secrets'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('secrets'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	return (
		<Suspense fallback={<FormSkeleton />}>
			<ProductSecrets
				productId={id}
				initial={readConsole(`products/${id}/secrets`)}
			/>
		</Suspense>
	)
}
