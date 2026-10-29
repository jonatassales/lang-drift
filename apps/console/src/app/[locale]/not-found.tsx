import { getTranslations } from 'next-intl/server'

import { Link } from '@i18n/navigation'

export async function generateMetadata() {
	const t = await getTranslations('notFound')

	return { title: t('title') }
}

export default async function NotFound() {
	const [t, ui] = await Promise.all([
		getTranslations('notFound'),
		getTranslations('ui'),
	])

	return (
		<main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-5 px-6">
			<h1 className="text-3xl font-semibold">{t('title')}</h1>
			<p className="text-muted-foreground">{t('description')}</p>
			<Link
				className="w-fit text-primary underline underline-offset-4"
				href="/home"
			>
				{ui('home')}
			</Link>
		</main>
	)
}
