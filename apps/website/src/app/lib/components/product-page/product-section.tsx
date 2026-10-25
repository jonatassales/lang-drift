import { SectionHeading } from '@langdrift/web/ui/section-heading'
import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface ProductSectionProps {
	children: ReactNode
	className?: string
	id: string
	kicker: string
	lead?: string
	title: string
}

// One section of a product page: the shell every section shares (width,
// top border, padding, scroll margin for the header) around its heading.
export function ProductSection(props: ProductSectionProps) {
	const { children, className, id, kicker, lead, title } = props

	const headingId = `${id}-heading`

	return (
		<section
			aria-labelledby={headingId}
			className={cn(
				'mx-auto max-w-[1180px] scroll-mt-24 border-t border-border px-[clamp(18px,3vw,36px)] py-16 sm:py-20',
				className,
			)}
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
