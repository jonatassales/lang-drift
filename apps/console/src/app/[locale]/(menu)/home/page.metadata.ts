import { getTranslations } from 'next-intl/server'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('home') }
}
