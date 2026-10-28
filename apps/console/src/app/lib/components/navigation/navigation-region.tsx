import { RegionRetry } from '@components/route-state/region-retry'
import { readConsole } from '@data/server'
import type { ProductNavigation as NavigationData } from '@data/types'

import { Navigation } from './navigation'
import { ProductNavigation } from './product-navigation'

export async function NavigationRegion() {
	try {
		const products = await readConsole<NavigationData>('products')

		return <ProductNavigation initial={Promise.resolve(products)} />
	} catch {
		return (
			<>
				<Navigation products={[]} />
				<RegionRetry />
			</>
		)
	}
}
