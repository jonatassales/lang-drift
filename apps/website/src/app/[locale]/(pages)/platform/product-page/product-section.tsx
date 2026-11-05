import type { ReactNode } from 'react'

import { SectionHeading } from 'web/ui/section-heading'

// One section of a product page: the shell every section shares (width,
// top border, padding, scroll margin for the header) around its heading.
interface ProductSectionProps {
	children: ReactNode
	id: string
	kicker?: string
	lead?: string
	title: string
}

export function ProductSection(props: ProductSectionProps) {
	const { children, id, kicker, lead, title } = props

	const headingId = `${id}-heading`

	return (
		<section
			aria-labelledby={headingId}
			className="mx-auto max-w-[1180px] scroll-mt-24 border-t border-border px-[clamp(18px,3vw,36px)] py-16 sm:py-20"
			id={id}
		>
			<SectionHeading
				headingId={headingId}
				kicker={kicker}
				lead={lead}
				title={title}
			/>
			{children}
		</section>
	)
}
