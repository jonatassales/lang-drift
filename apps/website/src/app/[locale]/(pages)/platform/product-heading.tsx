import type { ComponentType, ReactNode } from 'react'

import type { PlatformIconProps } from 'web/ui/icons/platform'

import { ProductIconTile } from './product-icon-tile'

// The heading every platform product page opens with: the product's icon on
// its tile, then the product's name. Tile and icon are sized in em, so they
// scale with the title on every page alike.
interface ProductHeadingProps {
	children: ReactNode
	icon: ComponentType<PlatformIconProps>
	id: string
}

export function ProductHeading(props: ProductHeadingProps) {
	const { children, icon: Icon, id } = props

	return (
		<h1
			className="m-0 flex items-center gap-[0.35em] text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance"
			id={id}
		>
			<ProductIconTile className="size-[0.74em] rounded-[0.2em]">
				<Icon className="size-[0.46em]" />
			</ProductIconTile>
			{children}
		</h1>
	)
}
