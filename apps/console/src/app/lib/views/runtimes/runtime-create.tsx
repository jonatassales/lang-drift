'use client'

import { useTranslations } from 'next-intl'
import { useState } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { linkClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import type {
	ProductEnvelope,
	SavedRuntimeEnvelope,
} from '@data/read-envelopes.types'
import { useResource } from '@data/resource.hook'
import type { Runtime } from '@data/types'
import { Link } from '@i18n/navigation'

import { ConnectionInstructions } from './connection-instructions'
import { RuntimeForm } from './runtime-form'

interface RuntimeCreateProps {
	initial: Promise<ProductEnvelope>
	productId: string
}

export function RuntimeCreate(props: RuntimeCreateProps) {
	const { productId, initial } = props

	const [created, setCreated] = useState<SavedRuntimeEnvelope | null>(null)

	const product = useResource<ProductEnvelope>(
		`products/${productId}`,
		0,
		initial,
	)
	const translate = useTranslations('ui')

	function handleSaved(runtime: Runtime, secret?: string) {
		setCreated({
			runtime,
			secret,
		})
	}

	if (product.error)
		return <ErrorNotice error={product.error} retry={product.refresh} />
	if (!product.data) return null

	return (
		<>
			<PageHeading
				title={created ? created.runtime.name : translate('createRuntime')}
				description={translate('configurationIntro')}
			/>
			{created ? (
				<div className="grid max-w-2xl gap-5">
					<ConnectionInstructions
						secret={created.secret}
						local={!!product.data.product.source}
					/>
					<Link href={`/runtimes/${created.runtime.id}`} className={linkClass}>
						{translate('viewRuntime')}
					</Link>
				</div>
			) : (
				<RuntimeForm
					productId={productId}
					components={product.data.product.components}
					onSaved={handleSaved}
				/>
			)}
		</>
	)
}
