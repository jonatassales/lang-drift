'use client'
import { useTranslations } from 'next-intl'
import { use } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import type { EpisodesEnvelope } from '@data/read-envelopes.types'
import type { ReadResult } from '@data/read-result'
import { Link } from '@i18n/navigation'

interface EpisodeEvidenceProps {
	initial: Promise<ReadResult<EpisodesEnvelope>>
	ids: string[]
}
export function EpisodeEvidence(props: EpisodeEvidenceProps) {
	const { initial, ids } = props

	const translate = useTranslations('ui')
	const resource = use(initial)
	if (!resource.data) return <ErrorNotice error={resource.error} />

	return (
		<div className="grid gap-2">
			{ids.map((id) => {
				const episode = resource.data?.episodes.find((entry) => entry.id === id)

				return (
					<details key={id} className="rounded-md border border-border p-3">
						<summary className="cursor-pointer break-all font-mono text-[10px] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring">
							{id}
						</summary>
						{episode ? (
							<>
								<Link
									href={`/runtimes/${episode.runtimeId}/history`}
									className="mt-3 block text-xs text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
								>
									{translate('history')}
								</Link>
								<pre className="mt-2 max-h-64 overflow-auto whitespace-pre-wrap break-all text-[10px]">
									{JSON.stringify(episode, null, 2)}
								</pre>
							</>
						) : (
							<p className="mt-2 text-xs text-muted-foreground">
								{translate('unknown')}
							</p>
						)}
					</details>
				)
			})}
		</div>
	)
}
