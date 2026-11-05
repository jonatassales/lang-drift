import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { privacyMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'

export { generateMetadata } from './page.metadata'

// Provider URLs verified against Cloudflare's published pages on 2026-09-30.

interface PrivacyPageProps {
	params: Promise<{ locale: string }>
}

const link =
	'focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary'

export default async function PrivacyPage(props: PrivacyPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = privacyMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="privacy" />
			<SkipLink href="#privacy-main">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main
				id="privacy-main"
				className="mx-auto max-w-3xl px-6 pt-36 pb-24 sm:px-8 sm:pt-44"
			>
				<article className="text-base leading-relaxed text-muted-foreground [&_h2]:mb-4 [&_h2]:text-2xl [&_h2]:font-medium [&_h2]:tracking-tight [&_h2]:text-foreground [&_p]:mb-4 [&_a]:text-foreground [&_a]:underline [&_a]:decoration-border [&_a]:underline-offset-4 hover:[&_a]:decoration-primary">
					<header className="mb-12 border-b border-border pb-8">
						<h1 className="mb-5 text-4xl font-medium tracking-tight text-foreground sm:text-5xl">
							{copy.title}
						</h1>
						<p className="text-sm">
							<time dateTime="2026-11-04">{copy.updated}</time>
						</p>
						<p>{copy.introduction}</p>
					</header>
					<section aria-labelledby="privacy-information" className="mb-10">
						<h2 id="privacy-information">{copy.overviewTitle}</h2>
						<p>{copy.overview}</p>
						<p>{copy.paused}</p>
					</section>
					<section aria-labelledby="privacy-preferences" className="mb-10">
						<h2 id="privacy-preferences">{copy.preferencesTitle}</h2>
						<p>{copy.preferences}</p>
					</section>
					<section aria-labelledby="privacy-contact" className="mb-10">
						<h2 id="privacy-contact">{copy.contactTitle}</h2>
						<p>{copy.contact}</p>
						<a className={link} href="mailto:founders@langdrift.md">
							founders@langdrift.md
						</a>
					</section>
					<section aria-labelledby="privacy-changes">
						<h2 id="privacy-changes">{copy.changesTitle}</h2>
						<p>{copy.changes}</p>
					</section>
				</article>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
