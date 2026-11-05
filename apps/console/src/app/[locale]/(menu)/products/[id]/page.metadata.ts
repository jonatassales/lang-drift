import { projectTitle } from '@data/entity-title'

interface ProductRouteParams {
	id: string
}

interface ProductPageProps {
	params: Promise<ProductRouteParams>
}

export async function generateMetadata(props: ProductPageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'runtimes')

	return { title }
}
