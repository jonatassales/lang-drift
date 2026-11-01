import { runtimeTitle } from '@data/entity-title'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await runtimeTitle(id, 'configuration')

	return { title }
}
