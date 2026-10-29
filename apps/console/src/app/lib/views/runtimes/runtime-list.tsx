'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { linkClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import { useResource } from '@data/resource.hook'
import type { RuntimeListSnapshot } from '@data/types'
import { Link } from '@i18n/navigation'
import { Card } from '@ui/card'

import { RuntimeWidget } from './runtime-widget'
import { useStubStream } from './stub-stream.hook'

interface RuntimeListProps {
	initial: Promise<RuntimeListSnapshot>
	productId?: string
	embedded?: boolean
}

export function RuntimeList(props: RuntimeListProps) {
	const { productId, embedded = false, initial } = props

	const t = useTranslations('ui')
	const resource = useResource<RuntimeListSnapshot>(
		`runtimes${productId ? `?productId=${encodeURIComponent(productId)}` : ''}`,
		5000,
		initial,
	)

	const stream = useStubStream(resource.data?.source === 'stub', productId)

	return (
		<>
			{embedded ? (
				<h2 className="mb-4 text-base font-medium">{t('runtimes')}</h2>
			) : (
				<PageHeading title={t('runtimes')} description={t('runtimesIntro')}>
					<Link
						href={
							productId ? `/products/${productId}/runtimes/new` : '/products'
						}
						className={linkClass}
					>
						{t('createRuntime')}
					</Link>
				</PageHeading>
			)}
			{resource.data?.source === 'stub' ? (
				<p role="status" className="mb-4 text-xs text-primary">
					{t(stream.connected ? 'webrtcConnected' : 'webrtcFallback')}
				</p>
			) : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.refresh} />
			) : resource.data?.runtimes.length ? (
				<div className="grid gap-3">
					{(stream.widgets ?? resource.data.widgets).map((snapshot) => (
						<RuntimeWidget
							key={snapshot.runtime.id}
							snapshot={snapshot}
							onRefresh={resource.refresh}
						/>
					))}
				</div>
			) : (
				<Card className="p-10 text-center text-sm text-muted-foreground">
					<p>{t('runtimeEmpty')}</p>
					{productId ? (
						<Link
							className={`${linkClass} mt-4`}
							href={`/products/${productId}/runtimes/new`}
						>
							{t('firstRuntime')}
						</Link>
					) : null}
				</Card>
			)}
		</>
	)
}
