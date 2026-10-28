'use client'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { useResource } from '@data/resource.hook'
import type { Product, RuntimeListSnapshot } from '@data/types'
import { RuntimeList } from '@views/runtimes/runtime-list'

import { ProductCreate } from './product-create'

interface ProductDetailProps {
	initial: Promise<{ product: Product }>
	id: string
	runtimes: Promise<RuntimeListSnapshot>
}

export function ProductDetail(props: ProductDetailProps) {
	const { id, initial, runtimes } = props

	const resource = useResource<{ product: Product }>(
		`products/${id}`,
		0,
		initial,
	)
	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (!resource.data) return null
	// Legacy imports still need explicit folder association before local operation.
	if (!resource.data.product.source) return <ProductCreate productId={id} />

	return <RuntimeList productId={id} initial={runtimes} />
}
