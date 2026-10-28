'use client'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { mutation } from '@data/api'
import type { RuntimeWidgetSnapshot } from '@data/types'
import { Button } from '@ui/button'
import { Pause, Play } from '@ui/icons'

interface RuntimePlaybackProps {
	snapshot: RuntimeWidgetSnapshot
	onRefresh?: () => void
}

export function RuntimePlayback(props: RuntimePlaybackProps) {
	const { snapshot, onRefresh } = props

	const t = useTranslations('ui')
	const [pending, startTransition] = useTransition()
	const [error, setError] = useState<unknown>(null)
	const { runtime, executions } = snapshot
	const stopping =
		runtime.desiredState === 'running' ||
		executions.some(
			(execution) => !['stopped', 'failed'].includes(execution.state),
		)
	const unavailable =
		runtime.desiredState === 'deleted' ||
		(!stopping && runtime.credentialRevoked)

	return (
		<div className="flex flex-wrap items-center gap-2">
			<Button
				disabled={pending || unavailable}
				aria-label={t(stopping ? 'pauseNamed' : 'playNamed', {
					name: runtime.name,
				})}
				onClick={() =>
					startTransition(async () => {
						setError(null)
						try {
							await mutation(
								`runtimes/${runtime.id}/${stopping ? 'stop' : 'start'}`,
								'POST',
							)

							onRefresh?.()
						} catch (failure) {
							setError(failure)
						}
					})
				}
			>
				{stopping ? <Pause size={14} /> : <Play size={14} />}
				{t(stopping ? 'pause' : 'start')}
			</Button>
			{error ? <ErrorNotice error={error} /> : null}
		</div>
	)
}
