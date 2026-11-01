import { ScrollHint } from '@langdrift/web/ui/scroll-hint'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { FinopsSection } from '@components/home/finops-section'
import { PlatformFaqSection } from '@components/home/platform-faq-section'
import { PlatformHero } from '@components/platform-hero'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { platformMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

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
