'use client'

import { Button } from '@langdrift/web/ui/button'
import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'
import { Suspense, useEffect, useRef, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { EvidenceSkeleton } from '@components/route-state/evidence-skeleton'
import { ConsoleError, consoleRequest } from '@data/api'
import { type ReadResult, readResult } from '@data/read-result'
import type {
	Component,
	Episode,
	Recommendation,
	RuntimeConfiguration,
} from '@data/types'
import { EpisodeEvidence } from '@views/history/episode-evidence'

import { RecommendationAccounting } from './recommendation-accounting'

interface RecommendationPanelProps {
	endpoint: string
	productId: string
	disabled?: boolean
	components?: Component[]
	agent?: {
		model?: string
		effort?: RuntimeConfiguration['effort']
	}
	onApply?: (values: Partial<RuntimeConfiguration>) => void
}

export function RecommendationPanel(props: RecommendationPanelProps) {
	const {
		endpoint,
		productId,
		disabled = false,
		onApply,
		components,
		agent,
	} = props

	const translate = useTranslations('ui')
	const [recommendation, setRecommendation] = useState<Recommendation | null>(
		null,
	)
	const [evidence, setEvidence] = useState<Promise<
		ReadResult<{ episodes: Episode[] }>
	> | null>(null)
	const [ticket, setTicket] = useState<
		| {
				name: string
				markdown: string
		  }
		| undefined
	>()
	const readingVersion = useRef(0)
	const request = useRef<AbortController | null>(null)
	const [reading, setReading] = useState(false)
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()
	useEffect(() => {
		return () => {
			request.current?.abort()
			readingVersion.current++
		}
	}, [])

	return (
		<Card className="block p-5">
			<h2 className="text-sm font-medium">{translate('recommendations')}</h2>
			<p className="mt-2 text-xs leading-relaxed text-muted-foreground">
				{translate('recommendationIntro')}
			</p>
			<label className={`${labelClass} mt-4`}>
				{translate('ticketMarkdown')}
				<input
					className={fieldClass}
					type="file"
					accept=".md,text/markdown"
					disabled={pending || disabled || reading}
					onChange={async (event) => {
						const version = ++readingVersion.current
						request.current?.abort()
						setEvidence(null)
						setTicket(undefined)
						setRecommendation(null)
						setError(null)
						setReading(true)
						try {
							const file = event.target.files?.[0]
							if (!file) return
							if (file.size > 16000)
								throw new ConsoleError('invalid_ticket', 400)
							const markdown = await file.text()

							if (version === readingVersion.current)
								setTicket({
									name: file.name,
									markdown,
								})
						} catch (failure) {
							if (version === readingVersion.current)
								setError(
									failure instanceof ConsoleError
										? failure
										: new ConsoleError('invalid_ticket', 400),
								)
						} finally {
							if (version === readingVersion.current) setReading(false)
						}
					}}
				/>
				<span className="text-[11px] text-muted-foreground">
					{translate('ticketMemoryHint')}
				</span>
			</label>
			<Button
				className="mt-4"
				disabled={pending || disabled || reading}
				onClick={() => {
					request.current?.abort()
					const controller = new AbortController()
					request.current = controller
					setEvidence(null)
					setError(null)
					startTransition(async () => {
						try {
							const result = await consoleRequest<{
								recommendation: Recommendation
							}>(endpoint, {
								method: 'POST',
								signal: controller.signal,
								body: JSON.stringify({
									...(components
										? {
												components: components.map(({ path, name, kind }) => ({
													path,
													name,
													kind,
												})),
											}
										: {}),
									...(agent ? { agent } : {}),
									...(ticket ? { ticket } : {}),
								}),
							})

							if (controller.signal.aborted) return
							setRecommendation(result.recommendation)
							const ids = result.recommendation.episodeIds
							setEvidence(
								ids.length
									? readResult(
											`products/${productId}/episodes?${ids.map((id) => `id=${encodeURIComponent(id)}`).join('&')}`,
											controller.signal,
										)
									: null,
							)
						} catch (failure) {
							if (!controller.signal.aborted) setError(failure)
						}
					})
				}}
				type="button"
				variant="outline"
			>
				{translate(pending ? 'loading' : 'getSuggestions')}
			</Button>
			{error ? (
				<div className="mt-4">
					<ErrorNotice error={error} />
				</div>
			) : null}
			{recommendation ? (
				<div className="mt-4 grid gap-3 text-xs">
					{recommendation.status !== 'suggested' ? (
						<p className="text-muted-foreground">
							{translate(`recommendationReasons.${recommendation.reasonCode}`)}
						</p>
					) : (
						<>
							<h3 className="font-medium text-primary">
								{translate('suggestion')}
							</h3>
							{recommendation.ticket ? (
								<p className="text-muted-foreground">
									{translate('ticketBasis')}: {recommendation.ticket.name}
								</p>
							) : null}
							<dl className="grid grid-cols-2 gap-2">
								{recommendation.values?.llmBudgetUsd !== undefined ? (
									<>
										<dt>{translate('estimatedBudget')}</dt>
										<dd>{recommendation.values.llmBudgetUsd} USD</dd>
									</>
								) : null}
								{recommendation.values?.cpu !== undefined ? (
									<>
										<dt className="text-muted-foreground">
											{translate('cpu')}
										</dt>
										<dd>{recommendation.values.cpu}</dd>
									</>
								) : null}
								{recommendation.values?.memoryMiB !== undefined ? (
									<>
										<dt className="text-muted-foreground">
											{translate('memory')}
										</dt>
										<dd>{recommendation.values.memoryMiB}</dd>
									</>
								) : null}
								{recommendation.values?.model ? (
									<>
										<dt className="text-muted-foreground">
											{translate('model')}
										</dt>
										<dd className="break-all">{recommendation.values.model}</dd>
									</>
								) : null}
								{recommendation.values?.effort ? (
									<>
										<dt className="text-muted-foreground">
											{translate('effort')}
										</dt>
										<dd>{translate(recommendation.values.effort)}</dd>
									</>
								) : null}
							</dl>
							<p className="whitespace-pre-wrap text-muted-foreground">
								{translate(
									`recommendationReasons.${recommendation.reasonCode}`,
								)}
							</p>
							<div>
								<h4 className="mb-2 font-medium">{translate('episodes')}</h4>
								{evidence ? (
									<Suspense fallback={<EvidenceSkeleton />}>
										<EpisodeEvidence
											initial={evidence}
											ids={recommendation.episodeIds}
										/>
									</Suspense>
								) : null}
							</div>
							{recommendation.limitationCodes?.length ? (
								<ul className="list-inside list-disc text-xs text-muted-foreground">
									{recommendation.limitationCodes.map((limitation) => (
										<li key={limitation}>
											{translate(`recommendationLimits.${limitation}`)}
										</li>
									))}
								</ul>
							) : null}
							{recommendation.evidence ? (
								<dl className="grid gap-2 text-[11px]">
									{Object.entries(recommendation.evidence).map(
										([field, evidence]) => (
											<div key={field}>
												<dt className="font-medium">
													{translate(field === 'memoryMiB' ? 'memory' : field)}
												</dt>
												<dd className="text-muted-foreground">
													{translate(
														evidence.basis === 'observed_resource'
															? 'observedResource'
															: evidence.basis === 'ticket_inference'
																? 'ticketBasis'
																: 'historicalConfiguration',
													)}
												</dd>
												<dd className="break-all font-mono text-[10px]">
													{evidence.episodeIds.join(', ')}
												</dd>
											</div>
										),
									)}
								</dl>
							) : null}
							{onApply && recommendation.values ? (
								<Button
									disabled={disabled}
									onClick={() => {
										if (recommendation.values) onApply(recommendation.values)
									}}
									type="button"
									variant="outline"
								>
									{translate('applySuggestions')}
								</Button>
							) : null}
						</>
					)}
					{recommendation.generationModel ||
					recommendation.usage ||
					recommendation.limits ? (
						<RecommendationAccounting recommendation={recommendation} />
					) : null}
				</div>
			) : null}
		</Card>
	)
}
