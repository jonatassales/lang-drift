'use client'

import { useResource } from '@data/resource.hook'
import type { ProductNavigation } from '@data/types'
import { usePathname } from '@i18n/navigation'

/** URLs own product context, including legacy runtime deep links. */
export function useProductContext(initial: Promise<ProductNavigation>) {
	const pathname = usePathname()
	const segments = pathname.split('/').filter(Boolean)
	const resource = useResource<ProductNavigation>('products', 5000, initial)
	const productId =
		segments[0] === 'products' && !['new', 'import'].includes(segments[1])
			? segments[1]
			: resource.data?.runtimeProducts[segments[1]]
	const products = resource.data?.products ?? []

	return {
		products,
		product: products.find((item) => item.id === productId),
	}
}
