import { redirect } from '@i18n/navigation'

interface LocaleRouteParams {
	locale: string
}

interface PageProps {
	params: Promise<LocaleRouteParams>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { locale } = await params

	redirect({
		href: '/login',
		locale,
	})
}
