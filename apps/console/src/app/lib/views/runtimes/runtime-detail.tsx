'use client'

import { useFormatter, useTranslations } from 'next-intl'
import { Suspense, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { linkClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import { FieldSkeleton } from '@components/route-state/field-skeleton'
import { StatusLabel } from '@components/status-label/status-label'
import { mutation } from '@data/api'
import type { ProductEnvelope } from '@data/read-envelopes.types'
import { useResource } from '@data/resource.hook'
import type { RuntimeDetailSnapshot } from '@data/types'
import { Link } from '@i18n/navigation'
import { useDateTime } from '@i18n/use-date-time.hook'
import { ExecutionTable } from '@views/executions/execution-table'
import { Button } from 'web/ui/button'
import { Card } from 'web/ui/card'

import { RuntimeConnection } from './runtime-connection'
import { RuntimeExecutionMetrics } from './runtime-execution-metrics'
import { RuntimeFinops } from './runtime-finops'
import { RuntimeGraph } from './runtime-graph'
import { RuntimeLoadTest } from './runtime-load-test'
import { RuntimePlayback } from './runtime-playback'
import { RuntimeStackIcons } from './runtime-stack-icons'
import { useStubStream } from './stub-stream.hook'

interface RuntimeDetailProps {
	initial: Promise<RuntimeDetailSnapshot>
	id: string
	productInitial: Promise<ProductEnvelope>
}

export function RuntimeDetail(props: RuntimeDetailProps) {
	const { id, initial, productInitial } = props

	const [error, setError] = useState<unknown>(null)
	const [requested, setRequested] = useState(false)
	const [pending, startTransition] = useTransition()

	const dateTime = useDateTime()
	const resource = useResource<RuntimeDetailSnapshot>(
		`runtimes/${id}`,
		5000,
		initial,
	)
	const stream = useStubStream(
		resource.data?.source === 'stub',
		resource.data?.runtime.productId,
	)

	const translate = useTranslations('ui')
	const format = useFormatter()

	function operate(action: 'revoke' | 'delete') {
		if (action === 'delete' && !window.confirm(translate('deleteConfirm')))
			return
		if (action === 'revoke' && !window.confirm(translate('revokeConfirm')))
			return
		setError(null)
		setRequested(false)
		startTransition(async () => {
			try {
				await mutation(
					`runtimes/${id}${action === 'delete' ? '' : `/${action}`}`,
					action === 'delete' ? 'DELETE' : 'POST',
				)

				setRequested(true)
				resource.refresh()
			} catch (failure) {
				setError(failure)
			}
		})
	}

	function handleDelete() {
		operate('delete')
	}

	function handleRevoke() {
		operate('revoke')
	}

	const snapshot =
		stream.widgets?.find((item) => item.runtime.id === id) ??
		resource.data?.widget
	if (resource.error)
		return <ErrorNotice error={resource.error} retry={resource.refresh} />
	if (!snapshot) return null
	const { runtime, environments, executions } = snapshot
	const deleted = runtime.desiredState === 'deleted'
	const locked =
		runtime.desiredState !== 'stopped' ||
		executions.some((execution) => execution.state !== 'stopped')

	return (
		<>
			<PageHeading title={runtime.name} description={runtime.id}>
				<div className="flex flex-wrap gap-2">
					<RuntimePlayback snapshot={snapshot} onRefresh={resource.refresh} />
					{!deleted ? (
						<Link href={`/runtimes/${id}/configuration`} className={linkClass}>
							{translate('configure')}
						</Link>
					) : null}
					<Link href={`/runtimes/${id}/history`} className={linkClass}>
						{translate('history')}
					</Link>
				</div>
			</PageHeading>
			<div className="grid gap-6">
				{error ? <ErrorNotice error={error} /> : null}
				{requested ? (
					<p role="status" className="text-sm text-muted-foreground">
						{translate('operationRequested')}
					</p>
				) : null}
				<section data-runtime-graph>
					<div className="mb-3 flex flex-wrap items-center justify-between gap-2">
						<h2 className="text-sm font-medium">
							{translate('componentGraph')}
						</h2>
						{snapshot.source === 'stub' ? (
							<p role="status" className="text-xs text-muted-foreground">
								{translate(
									stream.connected ? 'webrtcConnected' : 'webrtcFallback',
								)}
							</p>
						) : null}
					</div>
					<RuntimeGraph
						components={runtime.components}
						active={executions.some(
							(execution) => execution.state === 'running',
						)}
						tone={snapshot.tone}
						links={
							snapshot.source === 'stub' ? snapshot.componentLinks : undefined
						}
					/>
					<p className="mt-2 text-xs text-muted-foreground">
						{translate(
							snapshot.source === 'stub' && snapshot.componentLinks?.length
								? 'topologySynthetic'
								: 'topologyUnavailable',
						)}
					</p>
					<p className="mt-1 text-[11px] text-muted-foreground">
						{translate('graphColorLegend')}
					</p>
				</section>
				<RuntimeExecutionMetrics snapshot={snapshot} />
				<Card className="block p-5">
					<dl className="grid gap-5 sm:grid-cols-2 xl:grid-cols-5">
						<div>
							<dt className="mb-2 text-[11px] text-muted-foreground">
								{translate('desiredState')}
							</dt>
							<dd>
								<StatusLabel value={runtime.desiredState} />
							</dd>
						</div>
						<div>
							<dt className="mb-2 text-[11px] text-muted-foreground">
								{translate('cpu')}
							</dt>
							<dd className="text-xl tabular-nums">
								{format.number(runtime.configuration.cpu)}
							</dd>
						</div>
						<div>
							<dt className="mb-2 text-[11px] text-muted-foreground">
								{translate('memory')}
							</dt>
							<dd className="text-xl tabular-nums">
								{format.number(runtime.configuration.memoryMiB)}
							</dd>
						</div>
						<div>
							<dt className="mb-2 text-[11px] text-muted-foreground">
								{translate('model')}
							</dt>
							<dd className="break-all text-xs">
								{runtime.configuration.model}
								<p className="mt-1 text-muted-foreground">
									{translate(runtime.configuration.effort)}
								</p>
							</dd>
						</div>
						<div>
							<dt className="mb-2 text-[11px] text-muted-foreground">
								{translate('revision')}
							</dt>
							<dd className="text-xl">{runtime.revision}</dd>
						</div>
					</dl>
					{locked && !deleted ? (
						<p className="mt-5 border-t border-border pt-4 text-xs text-muted-foreground">
							{translate('configurationLocked')}
						</p>
					) : null}
					{runtime.credentialRevoked ? (
						<p className="mt-4 text-xs text-status-danger">
							{translate('revoked')}
						</p>
					) : null}
				</Card>
				<section>
					<h2 className="mb-3 text-sm font-medium">
						{translate('environments')}
					</h2>
					{environments.length ? (
						<div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
							{environments.map((environment) => (
								<Card key={environment.id} className="block p-4">
									<div className="flex flex-wrap items-center justify-between gap-2">
										<strong className="text-sm font-medium">
											{environment.name}
										</strong>
										<StatusLabel value={environment.connection} />
									</div>
									<p className="mt-3 break-all font-mono text-[10px] text-muted-foreground">
										{environment.id}
									</p>
									<p className="mt-2 text-[11px] text-muted-foreground">
										{translate('lastSeen')}:{' '}
										{dateTime(new Date(environment.lastSeenAt), {
											dateStyle: 'short',
											timeStyle: 'medium',
										})}
									</p>
									{environment.cleanupPending ? (
										<p className="mt-2 text-xs text-status-danger">
											{translate('cleanupPending')}
										</p>
									) : null}
								</Card>
							))}
						</div>
					) : (
						<Card className="block p-5 text-sm text-muted-foreground">
							{translate('environmentsEmpty')}
						</Card>
					)}
				</section>
				<section>
					<h2 className="mb-3 text-sm font-medium">
						{translate('executions')}
					</h2>
					<ExecutionTable
						executions={
							snapshot.source === 'stub'
								? executions
								: (resource.data?.executions ?? executions)
						}
						environments={environments}
					/>
				</section>
				<section>
					<h2 className="mb-3 text-sm font-medium">
						{translate('components')} · {runtime.components.length}
					</h2>
					<ul className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
						{runtime.components.map((component) => (
							<li
								key={component.path}
								className="min-w-0 rounded-lg border border-border bg-card p-4 text-xs"
							>
								<div className="flex items-center gap-2 font-medium">
									<RuntimeStackIcons stacks={component.stacks} size={18} />
									{component.name}
								</div>
								<p className="mt-2 break-all text-muted-foreground">
									{component.kind} · {component.path}
								</p>
								<p className="mt-1 break-all text-muted-foreground">
									{component.image ?? translate('unknown')} ·{' '}
									{component.port ?? translate('unknown')}
								</p>
							</li>
						))}
					</ul>
				</section>
				<RuntimeFinops snapshot={snapshot} />
				<RuntimeLoadTest key={runtime.id} snapshot={snapshot} />
				<Card className="block p-5">
					<h2 className="text-sm font-medium">{translate('markdown')}</h2>
					<pre className="mt-4 max-h-96 overflow-auto rounded-md bg-background p-4 text-xs leading-relaxed">
						<code>{runtime.markdown}</code>
					</pre>
				</Card>
				{!deleted && !runtime.credentialRevoked ? (
					<Suspense fallback={<FieldSkeleton />}>
						<RuntimeConnection initial={productInitial} />
					</Suspense>
				) : null}
				<div className="flex flex-wrap gap-3 border-t border-border pt-5">
					{!deleted ? (
						<Button
							disabled={pending}
							onClick={handleDelete}
							type="button"
							variant="destructive"
						>
							{translate('delete')}
						</Button>
					) : null}
					{!deleted && !runtime.credentialRevoked ? (
						<Button
							disabled={pending}
							onClick={handleRevoke}
							type="button"
							variant="outline"
						>
							{translate('revoke')}
						</Button>
					) : null}
				</div>
			</div>
		</>
	)
}
