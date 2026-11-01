import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'

import { isLocale } from '@i18n/config'

interface LocaleMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: LocaleMetadataProps,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!isLocale(locale)) notFound()

	const translate = await getTranslations({
		locale,
		namespace: 'metadata',
	})

	const siteName = translate('siteName')

	return {
		title: {
			default: siteName,
			template: translate('titleTemplate', { siteName }),
		},
		applicationName: siteName,
		robots: {
			index: false,
			follow: false,
		},
	}
}
