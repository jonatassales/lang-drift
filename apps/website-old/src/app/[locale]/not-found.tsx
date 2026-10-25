import { getTranslations } from 'next-intl/server'

import { Link } from '@i18n/navigation'

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
