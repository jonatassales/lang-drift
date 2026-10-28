import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { NextIntlClientProvider } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'
import type { ReactNode } from 'react'

import '@app/template/tailwind/app.css'
import { isLocale, languageTags } from '@i18n/config'

export const metadata: Metadata = {
	title: {
		default: 'LangDrift Console',
		template: '%s · LangDrift Console',
	},
	robots: {
		index: false,
		follow: false,
	},
}

interface LocaleLayoutProps {
	children: ReactNode
	params: Promise<{ locale: string }>
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
			<body className="min-h-screen bg-background font-sans text-foreground antialiased [&_:focus-visible]:outline-2 [&_:focus-visible]:outline-offset-2 [&_:focus-visible]:outline-ring">
				<NextIntlClientProvider>{children}</NextIntlClientProvider>
			</body>
		</html>
	)
}
