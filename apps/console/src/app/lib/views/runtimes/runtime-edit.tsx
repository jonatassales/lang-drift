'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import type {
	ProductEnvelope,
	RuntimeExecutionsEnvelope,
} from '@data/read-envelopes.types'
import { useResource } from '@data/resource.hook'
import { useRouter } from '@i18n/navigation'

import { RuntimeForm } from './runtime-form'

interface RuntimeEditProps {
	initial: Promise<RuntimeExecutionsEnvelope>
	id: string
	productInitial: Promise<ProductEnvelope>
}

export function RuntimeEdit(props: RuntimeEditProps) {
	const { id, initial, productInitial } = props

	const router = useRouter()
	const resource = useResource<RuntimeExecutionsEnvelope>(
		`runtimes/${id}`,
		0,
		initial,
	)
	const product = useResource<ProductEnvelope>(
		resource.data ? `products/${resource.data.runtime.productId}` : null,
		0,
		productInitial,
	)

	const translate = useTranslations('ui')

	function handleSaved() {
		router.push(`/runtimes/${id}`)
	}

	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (product.error)
		return <ErrorNotice error={product.error} retry={product.refresh} />
	if (!resource.data || !product.data) return null
	const locked =
		resource.data.runtime.desiredState !== 'stopped' ||
		resource.data.executions.some((execution) => execution.state !== 'stopped')

	return (
		<>
			<PageHeading
				title={translate('configuration')}
				description={resource.data.runtime.name}
			/>
			<RuntimeForm
				productId={resource.data.runtime.productId}
				components={product.data.product.components}
				runtime={resource.data.runtime}
				locked={locked}
				onSaved={handleSaved}
			/>
		</>
	)
}
