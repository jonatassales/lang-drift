'use client'

import { hasLocale, useLocale, useTranslations } from 'next-intl'

import { routing } from '@i18n/routing'
import { RouteState } from '@layout/route-state'

interface ErrorPageProps {
	retry: () => void
}

export default function ErrorPage(props: ErrorPageProps) {
	const { retry } = props

	const locale = useLocale()
	const translate = useTranslations('routeState')

	function handleRetry() {
		retry()
	}

	if (!hasLocale(routing.locales, locale)) {
		throw new Error(`Unsupported request locale: ${locale}`)
	}

	return (
		<RouteState
			description={translate('errorDescription')}
			homeLabel={translate('home')}
			imageSrc="/route-state/error.svg"
			imageAlt={translate('errorImage')}
			locale={locale}
			title={translate('error')}
		>
			<button
				className="inline-flex min-h-10 items-center rounded-full border border-border px-5 whitespace-nowrap text-foreground transition hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
				onClick={handleRetry}
				type="button"
			>
				{translate('retry')}
			</button>
		</RouteState>
	)
}
