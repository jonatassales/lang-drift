import { getTranslations } from 'next-intl/server'

import { RouteStateNotice } from '@components/route-state/route-state-notice'
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
		<RouteStateNotice title={t('title')} description={t('description')}>
			<Link
				className="w-fit text-primary underline underline-offset-4"
				href="/home"
			>
				{ui('home')}
			</Link>
		</RouteStateNotice>
	)
}
