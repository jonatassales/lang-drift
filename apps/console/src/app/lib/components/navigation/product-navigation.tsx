'use client'

import type { ProductNavigation as NavigationData } from '@data/types'
import { useProductContext } from '@views/products/product-context.hook'

import { Navigation } from './navigation'

interface ProductNavigationProps {
	initial: Promise<NavigationData>
}

export function ProductNavigation(props: ProductNavigationProps) {
	const { initial } = props

	const { product, products } = useProductContext(initial)

	return <Navigation product={product} products={products} />
}
