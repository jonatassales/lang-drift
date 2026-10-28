import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'
import type { ReactNode } from 'react'

import type { PlatformProduct } from '@components/website-header/platform-products'
import { headerMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface HomeProductSectionProps {
	/** The illustrative note under the demo. */
	caption?: string
	/** The product's demo, with its own grid. */
	children: ReactNode
	description: string
	/** A string, or the Vision question with its chroma span. */
	headline: ReactNode
	id: 'vision' | 'microruntimes' | 'finops' | 'ai-voice'
	/** The product page; absent while the product has none. */
	link?: { href: string; label: string }
	locale: WebsiteLocale
	product: PlatformProduct
	/** What the demo shows, for screen readers. */
	summary?: string
}

// One product on the home page: the product's icon and name on a row of
// their own at the right, the question and its answer at the left, the demo,
// and the way to the product page last. The heading and lead carry the
// Vision hero's type, so the four sections and that page read as one family.
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

	const Icon = product.icon
	const headingId = `${id}-heading`
	const captionId = `${id}-caption`
	const summaryId = `${id}-summary`
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
			<p className="m-0 flex items-center justify-end gap-2.5 text-sm font-medium tracking-tight text-foreground">
				<span
					aria-hidden="true"
					className="inline-flex size-7 shrink-0 items-center justify-center rounded-md border border-border bg-card"
				>
					<Icon className="size-4" />
				</span>
				{headerMessages[locale][product.label]}
			</p>
			<h2
				className="m-0 mt-6 text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance text-foreground"
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
			{link ? (
				<p className="mt-8 mb-0">
					<Link
						className="smooth-text-link inline-flex min-h-11 items-center gap-2"
						href={link.href}
						locale={locale}
					>
						{link.label}
						<ArrowRightIcon />
					</Link>
				</p>
			) : null}
		</section>
	)
}
