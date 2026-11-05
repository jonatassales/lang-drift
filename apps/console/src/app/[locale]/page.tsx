import { redirect } from '@i18n/navigation'

interface LocaleRouteParams {
	locale: string
}

interface HomeProps {
	params: Promise<LocaleRouteParams>
}

export default async function Home(props: HomeProps) {
	const { params } = props

	const { locale } = await params

	redirect({
		href: '/home',
		locale,
	})
}
