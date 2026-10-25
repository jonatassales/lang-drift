import type { Metadata, Viewport } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale, NextIntlClientProvider } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'
import { type ReactNode, Suspense } from 'react'

import site from '@app/metadata/site.json'
import { LocaleAvailabilityNotice } from '@components/locale-availability-notice'
import { getWebsiteMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import '@pack/react/styles.css'
import '@app/template/tailwind/app.css'
import '@app/globals.css'
import '@app/responsive.css'
import '@app/hero.css'
import '@components/hero-demo.css'

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
			className="dark"
			style={{ colorScheme: 'dark' }}
			suppressHydrationWarning
		>
			<body suppressHydrationWarning>
				<NextIntlClientProvider
					locale={locale}
					messages={getWebsiteMessages(locale)}
					timeZone="UTC"
				>
					<Suspense>
						<LocaleAvailabilityNotice />
					</Suspense>
					{children}
				</NextIntlClientProvider>
			</body>
		</html>
	)
}
