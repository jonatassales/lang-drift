import { redirect } from '@i18n/navigation'

interface ConfigPageProps {
	params: Promise<{
		locale: string
		id: string
	}>
}

export default async function ConfigPage(props: ConfigPageProps) {
	const { params } = props

	const { locale, id } = await params

	return redirect({
		href: `/products/${id}/config/secrets`,
		locale,
	})
}
