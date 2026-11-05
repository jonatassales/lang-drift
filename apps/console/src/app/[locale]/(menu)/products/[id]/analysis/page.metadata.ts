import { projectTitle } from '@data/entity-title'

interface ProductRouteParams {
	id: string
}

interface PageProps {
	params: Promise<ProductRouteParams>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'tickets')

	return { title }
}
