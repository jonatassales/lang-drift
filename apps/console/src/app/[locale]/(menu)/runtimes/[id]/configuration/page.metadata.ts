import { runtimeTitle } from '@data/entity-title'

interface RuntimeRouteParams {
	id: string
}

interface PageProps {
	params: Promise<RuntimeRouteParams>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await runtimeTitle(id, 'configuration')

	return { title }
}
