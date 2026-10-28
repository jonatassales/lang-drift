import { ProductBreadcrumb } from '@components/navigation/product-breadcrumb'

interface BreadcrumbPageProps {
	params: Promise<{
		locale: string
		segments: string[]
	}>
}

export default async function BreadcrumbPage(props: BreadcrumbPageProps) {
	const { params } = props

	const { locale, segments } = await params

	return <ProductBreadcrumb locale={locale} segments={segments} />
}
