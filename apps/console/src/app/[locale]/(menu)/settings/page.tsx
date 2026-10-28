import { getTranslations } from 'next-intl/server'

import { ConsoleSettings } from '@views/settings/console-settings'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('settings') }
}
export default function Page() {
	return <ConsoleSettings />
}
