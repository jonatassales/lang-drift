import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { conversionMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'

import { PlansSection } from './plans-section'

export { generateMetadata } from './page-metadata.map'

interface LocaleRouteParams {
	locale: string
}

interface PricingPageProps {
	params: Promise<LocaleRouteParams>
}

export default async function PricingPage(props: PricingPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)

	return (
		<>
			<WebsiteStructuredData locale={locale} page="pricing" />
			<SkipLink href="#plans-main">
				{conversionMessages[locale].plans.skip}
			</SkipLink>
			<WebsiteHeader />
			<main className="min-w-0" id="plans-main">
				<PlansSection locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
