import { Analytics } from '@vercel/analytics/next'
import type { Viewport } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale, NextIntlClientProvider } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'
import { type ReactNode, Suspense } from 'react'

import 'web/styles.css'
import '@app/template/tailwind/app.css'
import { getWebsiteMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import { LocaleAvailabilityNotice } from './locale-availability-notice'
import { ScrollToTop } from './scroll-to-top'

export { generateStaticParams } from './generate-static-params.compute'
export { generateMetadata } from './layout-metadata.map'

export const viewport: Viewport = {
	colorScheme: 'dark light',
	themeColor: '#0b0b0c',
}

interface LocaleRouteParams {
	locale: string
}

interface LocaleLayoutProps {
	children: ReactNode
	params: Promise<LocaleRouteParams>
}

export default async function LocaleLayout(props: LocaleLayoutProps) {
	const { children, params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)

	return (
		<html
			lang={locale}
			className="dark scroll-pt-24 scroll-smooth bg-background font-sans text-foreground scheme-light motion-reduce:scroll-auto dark:scheme-dark"
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
