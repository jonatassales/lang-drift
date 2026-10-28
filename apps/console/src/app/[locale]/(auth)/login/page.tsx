import { getTranslations } from 'next-intl/server'

import { AuthForm } from '@views/auth/auth-form'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('login') }
}
export default function Page() {
	return <AuthForm />
}
