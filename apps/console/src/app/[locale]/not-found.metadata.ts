import { getTranslations } from 'next-intl/server'

export async function generateMetadata() {
	const translate = await getTranslations('notFound')

	return { title: translate('title') }
}
