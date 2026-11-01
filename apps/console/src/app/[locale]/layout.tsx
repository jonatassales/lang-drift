import { Analytics } from '@vercel/analytics/next'
import type { Metadata, Viewport } from 'next'
import { notFound } from 'next/navigation'
import { NextIntlClientProvider } from 'next-intl'
import { getTranslations, setRequestLocale } from 'next-intl/server'
import type { ReactNode } from 'react'

import '@langdrift/web/styles.css'
import '@app/template/tailwind/app.css'
import { isLocale, languageTags } from '@i18n/config'

interface LocaleLayoutProps {
	children: ReactNode
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: Omit<LocaleLayoutProps, 'children'>,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!isLocale(locale)) notFound()

	const t = await getTranslations({
		locale,
		namespace: 'metadata',
	})

	const siteName = t('siteName')

	return {
		title: {
			default: siteName,
			template: t('titleTemplate', { siteName }),
		},
		applicationName: siteName,
		robots: {
			index: false,
			follow: false,
		},
	}
}

// Dark is the console's default scheme; the theme color is its dark
// background token.
export const viewport: Viewport = {
	colorScheme: 'dark light',
	themeColor: '#0b0b0c',
}

export default async function LocaleLayout(props: LocaleLayoutProps) {
	const { children, params } = props

	const { locale } = await params

	if (!isLocale(locale)) notFound()
	setRequestLocale(locale)

	return (
		<html
			className="dark scheme-light dark:scheme-dark"
			lang={languageTags[locale]}
			suppressHydrationWarning
		>
			<body className="min-h-screen bg-background font-sans text-foreground antialiased">
				<NextIntlClientProvider>{children}</NextIntlClientProvider>
				<Analytics />
			</body>
		</html>
	)
}
