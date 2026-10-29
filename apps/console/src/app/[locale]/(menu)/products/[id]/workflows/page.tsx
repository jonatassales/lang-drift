import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { projectTitle } from '@app/metadata/entity-title'
import { PageHeading } from '@components/page-heading/page-heading'
import { TicketListSkeleton } from '@components/route-state/ticket-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { TicketListing } from '@views/tickets/ticket-listing'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'workflows')

	return { title }
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading title={t('workflows')} description={t('workflowsIntro')} />
			<Suspense fallback={<TicketListSkeleton />}>
				<TicketListing
					initial={readConsole(`products/${id}/tickets?reached=workflow`)}
					phase="workflow"
					productId={id}
				/>
			</Suspense>
		</>
	)
}
