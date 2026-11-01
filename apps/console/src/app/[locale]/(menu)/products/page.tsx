import { Suspense } from 'react'

import { ProductListSkeleton } from '@components/route-state/product-list-skeleton'
import { readConsole } from '@data/server'
import { ProductList } from '@views/products/product-list'

export { generateMetadata } from './page.metadata'

export default async function Page() {
	return (
		<Suspense fallback={<ProductListSkeleton />}>
			<ProductList initial={readConsole('products')} />
		</Suspense>
	)
}
