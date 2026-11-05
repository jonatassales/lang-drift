'use client'

import { ErrorNotice } from '@components/error-notice/error-notice'
import type { ProductEnvelope } from '@data/product-envelopes.types'
import { useResource } from '@data/resource.hook'
import type { RuntimeListSnapshot } from '@data/types'
import { RuntimeList } from '@views/runtimes/runtime-list'

import { ProductCreate } from './product-create'

interface ProductDetailProps {
	initial: Promise<ProductEnvelope>
	id: string
	runtimes: Promise<RuntimeListSnapshot>
}

export function ProductDetail(props: ProductDetailProps) {
	const { id, initial, runtimes } = props

	const resource = useResource<ProductEnvelope>(`products/${id}`, 0, initial)
	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (!resource.data) return null
	// Legacy imports still need explicit folder association before local operation.
	if (!resource.data.product.source) return <ProductCreate productId={id} />

	return <RuntimeList productId={id} initial={runtimes} />
}
