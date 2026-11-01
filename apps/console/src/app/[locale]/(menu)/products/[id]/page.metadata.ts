import { projectTitle } from '@data/entity-title'

interface ProductPageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: ProductPageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'runtimes')

	return { title }
}
