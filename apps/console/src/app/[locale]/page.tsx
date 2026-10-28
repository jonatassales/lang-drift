import { redirect } from '@i18n/navigation'

interface HomeProps {
	params: Promise<{ locale: string }>
}

export default async function Home(props: HomeProps) {
	const { params } = props

	const { locale } = await params

	redirect({
		href: '/home',
		locale,
	})
}
