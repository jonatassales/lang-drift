import { redirect } from '@i18n/navigation'

interface PageProps {
	params: Promise<{ locale: string }>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { locale } = await params

	redirect({
		href: '/products',
		locale,
	})
}
