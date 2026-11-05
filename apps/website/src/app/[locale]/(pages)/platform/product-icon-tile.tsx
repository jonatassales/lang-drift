import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

interface ProductIconTileProps {
	children: ReactNode
	className?: string
}

// A product's icon on a bordered card tile, the mark a platform page opens
// with. The caller sizes the tile and its icon: in a heading they scale with
// the title, on a label line they keep a fixed size.
export function ProductIconTile(props: ProductIconTileProps) {
	const { children, className } = props

	return (
		<span
			aria-hidden="true"
			className={cn(
				'inline-flex shrink-0 items-center justify-center border border-border bg-card text-foreground',
				className,
			)}
		>
			{children}
		</span>
	)
}
