import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import { AuthForm } from '@views/auth/auth-form'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('signup') }
}
export default async function Page() {
	const meta = await readConsole<{ source: string }>('meta')

	return <AuthForm source={meta.source} signup />
}
