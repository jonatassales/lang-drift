import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { platformMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'
import { FinopsSection } from '@routes/home/finops-section'
import { PlatformFaqSection } from '@routes/home/platform-faq-section'
import { PlatformHero } from '@routes/home/platform-hero'
import { ScrollHint } from 'web/ui/scroll-hint'

export { generateMetadata } from './page.metadata'

interface WebsitePageProps {
	params: Promise<{ locale: string }>
}

export default async function WebsitePage(props: WebsitePageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = platformMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="home" />
			<SkipLink href="#top">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<PlatformHero locale={locale} />
				<FinopsSection locale={locale} />
				<PlatformFaqSection locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
			<ScrollHint />
		</>
	)
}
