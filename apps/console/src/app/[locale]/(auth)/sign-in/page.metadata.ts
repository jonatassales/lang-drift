import { getTranslations } from 'next-intl/server'

export async function generateMetadata() {
	const translate = await getTranslations('ui')

	const title = translate('login')

	return { title }
}
