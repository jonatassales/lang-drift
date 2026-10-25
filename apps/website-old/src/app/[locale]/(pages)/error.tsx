'use client'

import { useTranslations } from 'next-intl'

import { Link } from '@i18n/navigation'

export default function ErrorPage({ reset }: { reset: () => void }) {
	const t = useTranslations('routeState')
	return (
		<main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-5 px-6 text-foreground">
			<h1 className="text-3xl font-semibold">{t('error')}</h1>
			<p className="text-muted-foreground">{t('errorDescription')}</p>
			<button
				className="min-h-11 w-fit rounded-md border border-current px-4"
				onClick={reset}
				type="button"
			>
				{t('retry')}
			</button>
			<Link className="w-fit underline underline-offset-4" href="/">
				{t('home')}
			</Link>
		</main>
	)
}
