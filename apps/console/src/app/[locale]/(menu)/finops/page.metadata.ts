import { getTranslations } from 'next-intl/server'

export async function generateMetadata() {
	const translate = await getTranslations('ui')

	return { title: translate('costOverview') }
}
