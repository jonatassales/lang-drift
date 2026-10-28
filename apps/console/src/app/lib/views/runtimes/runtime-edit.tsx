'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { useResource } from '@data/resource.hook'
import type { Execution, Product, Runtime } from '@data/types'
import { useRouter } from '@i18n/navigation'

import { RuntimeForm } from './runtime-form'

interface RuntimeEditProps {
	initial: Promise<{
		runtime: Runtime
		executions: Execution[]
	}>
	id: string
	productInitial: Promise<{ product: Product }>
}

export function RuntimeEdit(props: RuntimeEditProps) {
	const { id, initial, productInitial } = props

	const t = useTranslations('ui')
	const router = useRouter()
	const resource = useResource<{
		runtime: Runtime
		executions: Execution[]
	}>(`runtimes/${id}`, 0, initial)
	const product = useResource<{ product: Product }>(
		resource.data ? `products/${resource.data.runtime.productId}` : null,
		0,
		productInitial,
	)
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
				title={t('configuration')}
				description={resource.data.runtime.name}
			/>
			<RuntimeForm
				productId={resource.data.runtime.productId}
				components={product.data.product.components}
				runtime={resource.data.runtime}
				locked={locked}
				onSaved={() => router.push(`/runtimes/${id}`)}
			/>
		</>
	)
}
