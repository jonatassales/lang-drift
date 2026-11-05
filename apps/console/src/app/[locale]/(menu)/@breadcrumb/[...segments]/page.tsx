import { ProductBreadcrumb } from '@components/navigation/product-breadcrumb'

interface BreadcrumbRouteParams {
	locale: string
	segments: string[]
}

interface BreadcrumbPageProps {
	params: Promise<BreadcrumbRouteParams>
}

export default async function BreadcrumbPage(props: BreadcrumbPageProps) {
	const { params } = props

	const { locale, segments } = await params

	return <ProductBreadcrumb locale={locale} segments={segments} />
}
