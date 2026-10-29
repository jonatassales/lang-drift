import type { Metadata } from 'next'
import { hasLocale } from 'next-intl'
import { getTranslations } from 'next-intl/server'

import { getWebsiteMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import { routing } from '@i18n/routing'

interface NotFoundMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: NotFoundMetadataProps,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) return {}

	return { title: getWebsiteMessages(locale).routeState.notFound }
}

export default async function NotFound() {
	const t = await getTranslations('routeState')
	return (
		<main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-5 px-6 text-foreground">
			<h1 className="text-3xl font-semibold">{t('notFound')}</h1>
			<p className="text-muted-foreground">{t('notFoundDescription')}</p>
			<Link className="w-fit underline underline-offset-4" href="/">
				{t('home')}
			</Link>
		</main>
	)
}
