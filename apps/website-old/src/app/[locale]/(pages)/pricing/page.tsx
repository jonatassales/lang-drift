import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { EarlyAccessSection } from '@components/early-access/section'
import { PlansSection } from '@components/plans-section'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { conversionMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface PricingPageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: PricingPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('pricing', locale, await socialImages(parent))
}

export default async function PricingPage(props: PricingPageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)

	return (
		<>
			<a
				className="fixed top-0 left-3 z-200 max-w-[calc(100vw-1.5rem)] -translate-y-full bg-foreground px-2.5 py-2 whitespace-normal text-card focus-visible:top-3 focus-visible:translate-y-0"
				href="#plans-main"
			>
				{conversionMessages[locale].plans.skip}
			</a>
			<WebsiteHeader />
			<main className="plans-page" id="plans-main">
				<PlansSection locale={locale} />
				<EarlyAccessSection locale={locale} source="pricing" compact />
			</main>
			<WebsiteFooter locale={locale} compact />
		</>
	)
}
