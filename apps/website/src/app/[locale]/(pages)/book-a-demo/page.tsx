import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { BookDemoForm } from '@components/book-a-demo/form'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { bookDemoMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface BookDemoPageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: BookDemoPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const images = await socialImages(parent)

	return pageMetadata('book-a-demo', locale, images)
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
							<BookDemoForm copy={copy} locale={locale} />
						</div>
					</div>
				</section>
			</main>
			<WebsiteFooter compact locale={locale} />
		</>
	)
}
