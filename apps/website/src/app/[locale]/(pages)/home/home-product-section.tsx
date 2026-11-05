import type { ReactNode } from 'react'

import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import type { PlatformProduct } from '@layout/website-header/platform-products.data'
import { ArrowRightIcon } from 'web/ui/icons/arrow'
import { TextLink } from 'web/ui/text-link'

// One product on the home page: the question and its answer, the demo, and
// the way to the product page last. Every section shares the heading and
// lead type.
interface HomeProductSectionProps {
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
	const summaryId = `${id}-summary`
	const href = 'href' in product ? product.href : undefined

	return (
		<section
			aria-describedby={summary ? summaryId : undefined}
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
			{summary ? (
				<p className="sr-only" id={summaryId}>
					{summary}
				</p>
			) : null}
			{link && href ? (
				<p className="mt-8 mb-0">
					<TextLink asChild className="inline-flex min-h-11 items-center gap-2">
						<Link href={href} locale={locale}>
							{link}
							<ArrowRightIcon />
						</Link>
					</TextLink>
				</p>
			) : null}
		</section>
	)
}
