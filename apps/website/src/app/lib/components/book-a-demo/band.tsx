import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'

import { bookDemoMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface BookDemoBandProps {
	/** A product page's own sentence; the band's own otherwise. */
	description?: string
	locale: WebsiteLocale
}

// The way to the Book a demo page, closing the home and the product pages.
export function BookDemoBand(props: BookDemoBandProps) {
	const { description, locale } = props

	const copy = bookDemoMessages[locale].band

	return (
		<section
			aria-labelledby="book-a-demo-heading"
			className="mx-auto max-w-[1240px] scroll-mt-25 px-4.5 pt-10 pb-14 sm:px-7 sm:pt-17 sm:pb-21.5"
			id="book-a-demo"
		>
			<div className="rounded-2xl border border-border border-t-3 border-t-primary bg-card p-6.5 text-foreground sm:p-11">
				<h2
					className="m-0 text-3xl leading-tight font-medium tracking-tight text-balance sm:text-4xl"
					id="book-a-demo-heading"
				>
					{copy.title}
				</h2>
				<p className="mt-5 mb-0 max-w-3xl text-base leading-relaxed text-muted-foreground">
					{description ?? copy.description}
				</p>
				<p className="mt-8 mb-0">
					<Link
						className="inline-flex min-h-11 items-center gap-2 rounded-full bg-foreground px-5 text-sm font-semibold whitespace-nowrap text-background transition-opacity hover:opacity-85"
						href="/book-a-demo"
						locale={locale}
					>
						{copy.cta}
						<ArrowRightIcon />
					</Link>
				</p>
			</div>
		</section>
	)
}
