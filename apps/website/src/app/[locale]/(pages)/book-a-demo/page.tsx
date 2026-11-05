import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { bookDemoMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'

// Demo requests return with the form below, once demos open again.
// import { BookDemoForm } from './book-demo-form'
import { DemoNotice } from './demo-notice'

export { generateMetadata } from './page-metadata.map'

interface BookDemoPageProps {
	params: Promise<{ locale: string }>
}

export default async function BookDemoPage(props: BookDemoPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = bookDemoMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="book-a-demo" />
			<SkipLink href="#book-a-demo-form">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-labelledby="book-a-demo-heading"
					className="mx-auto max-w-[1280px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<div className="grid gap-12 lg:grid-cols-2 lg:gap-16">
						<div>
							<h1
								className="m-0 text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance"
								id="book-a-demo-heading"
							>
								{copy.title}
							</h1>
							<p className="mt-[23px] mb-0 max-w-[560px] text-base leading-[1.65] text-pretty text-muted-foreground">
								{copy.lead}
							</p>
						</div>
						<div
							className="scroll-mt-24 rounded-2xl border border-border bg-card p-6 sm:p-8"
							id="book-a-demo-form"
						>
							<DemoNotice copy={copy.paused} />
							<Link
								className="mt-6 inline-flex min-h-11 items-center text-sm font-medium text-foreground underline decoration-border underline-offset-4 hover:decoration-primary focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
								href="/"
							>
								{copy.paused.home}
							</Link>
							{/* Demo requests return here once demos open again. The form's zod
							schema then reaches the client bundle, and zod probes for eval:
							set it jitless (z.config) first, so the CSP reports nothing. */}
							{/* <BookDemoForm copy={copy} locale={locale} /> */}
						</div>
					</div>
				</section>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
