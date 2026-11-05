import type { ReactNode } from 'react'

import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import type { PlatformProduct } from '@layout/website-header/platform-products.data'
import { ArrowRightIcon } from 'web/ui/icons/arrow'

// One product on the home page: the question and its answer, the demo, and
// the way to the product page last. Every section shares the heading and
// lead type.
interface HomeProductSectionProps {
	/** The note under the demo. */
	caption?: string
	/** The product's demo, with its own grid. */
	children: ReactNode
	description: string
	/** A string, or a heading with inline markup. */
	headline: ReactNode
	id: PlatformProduct['label']
	/** The label of the way to the product page. */
	link?: string
	locale: WebsiteLocale
	product: PlatformProduct
	/** What the demo shows, for screen readers. */
	summary?: string
}

export function HomeProductSection(props: HomeProductSectionProps) {
	const {
		caption,
		children,
		description,
		headline,
		id,
		link,
		locale,
		product,
		summary,
	} = props

	const headingId = `${id}-heading`
	const captionId = `${id}-caption`
	const summaryId = `${id}-summary`
	const href = 'href' in product ? product.href : undefined
	const describedBy = [caption && captionId, summary && summaryId]
		.filter(Boolean)
		.join(' ')

	return (
		<section
			aria-describedby={describedBy || undefined}
			aria-labelledby={headingId}
			className="mx-auto max-w-[1280px] scroll-mt-24 border-t border-border px-[clamp(18px,3vw,36px)] py-16 sm:py-20"
			id={id}
		>
			<h2
				className="m-0 text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance text-foreground"
				id={headingId}
			>
				{headline}
			</h2>
			<p className="mt-[23px] mb-0 max-w-[750px] text-base leading-[1.65] text-pretty text-muted-foreground">
				{description}
			</p>
			<div className="mt-12">{children}</div>
			{caption ? (
				<p
					className="mt-5 mb-0 max-w-[880px] text-xs leading-[1.65] text-muted-foreground"
					id={captionId}
				>
					{caption}
				</p>
			) : null}
			{summary ? (
				<p className="sr-only" id={summaryId}>
					{summary}
				</p>
			) : null}
			{link && href ? (
				<p className="mt-8 mb-0">
					<Link
						className="smooth-text-link inline-flex min-h-11 items-center gap-2 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
						href={href}
						locale={locale}
					>
						{link}
						<ArrowRightIcon />
					</Link>
				</p>
			) : null}
		</section>
	)
}
