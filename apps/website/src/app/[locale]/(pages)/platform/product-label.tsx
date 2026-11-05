import type { ReactNode } from 'react'

import { ProductIconTile } from './product-icon-tile'

interface ProductLabelProps {
	icon: ReactNode
	name: string
}

// The product a platform page belongs to, named above its title with the
// product's icon.
export function ProductLabel(props: ProductLabelProps) {
	const { icon, name } = props

	return (
		<span className="inline-flex items-center gap-2.5 text-sm font-medium text-foreground">
			<ProductIconTile className="size-7 rounded-lg">{icon}</ProductIconTile>
			{name}
		</span>
	)
}
