import { getTranslations } from 'next-intl/server'

import { ProductCreate } from '@views/products/product-create'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('createProject'),
	}
}

export default function Page() {
	return <ProductCreate />
}
