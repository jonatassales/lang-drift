import { Analytics } from '@vercel/analytics/next'
import type { Metadata, Viewport } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale, NextIntlClientProvider } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'
import { type ReactNode, Suspense } from 'react'

import site from '@app/metadata/site.json'
import { LocaleAvailabilityNotice } from '@components/locale-availability-notice'
import { ScrollToTop } from '@components/scroll-to-top'
import { getWebsiteMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import '@langdrift/web/styles.css'
import '@app/template/tailwind/app.css'

interface LayoutProps {
	children: ReactNode
	params: Promise<{ locale: string }>
}

export function generateStaticParams() {
	return routing.locales.map((locale) => ({ locale }))
}

export async function generateMetadata({
	params,
}: Omit<LayoutProps, 'children'>): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return {
		metadataBase: new URL(site.url),
		applicationName: site.name,
	}
}

export const viewport: Viewport = {
	colorScheme: 'dark light',
	themeColor: '#0b0b0c',
}

export default async function LocaleLayout(props: LayoutProps) {
	const { children, params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)

	return (
		<html
			lang={locale}
			className="dark scroll-pt-24 scroll-smooth bg-background font-sans text-foreground motion-reduce:scroll-auto"
			style={{ colorScheme: 'dark' }}
			suppressHydrationWarning
		>
			<body
				className="min-h-screen bg-background text-foreground transition-[background,color] duration-280 ease-[ease] motion-reduce:transition-none"
				suppressHydrationWarning
			>
				<NextIntlClientProvider
					locale={locale}
					messages={getWebsiteMessages(locale)}
					timeZone="UTC"
				>
					<Suspense>
						<LocaleAvailabilityNotice />
					</Suspense>
					<ScrollToTop />
					{children}
				</NextIntlClientProvider>
				<Analytics />
			</body>
		</html>
	)
}
