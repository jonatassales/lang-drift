'use client'

import { createTranslator } from 'next-intl'
import { lazy, Suspense, useEffect, useId, useRef, useState } from 'react'

import {
	aggregateReviewEvents,
	classificationFilters,
	type ExecutiveReviewMessages,
	getReviewsForEvents,
	initialReviewFilters,
	type QuestionId,
	questionPresets,
	type ReviewFilters,
	reviewDistribution,
	reviewSources,
	selectReviewEvents,
	type TargetId,
	targetIds,
} from '@data/executive-review-data'
import type { WebsiteLocale } from '@i18n/routing'

import { EvidenceMatrix } from './executive-review/evidence-matrix'
import { ReviewTimeline } from './executive-review/review-timeline'

const questions: QuestionId[] = [
	'overview',
	'no-linked-evidence',
	'open-review',
]

const DeferredTargetDistributionChart = lazy(() =>
	import('@langdrift/web/ui/target-distribution-chart').then((module) => ({
		default: module.TargetDistributionChart,
	})),
)

function reveal(element: HTMLElement | null) {
	if (!element) return
	element.focus({ preventScroll: true })
	element.scrollIntoView({
		block: 'start',
		behavior: window.matchMedia('(prefers-reduced-motion: reduce)').matches
			? 'instant'
			: 'smooth',
	})
}

interface ExecutiveReviewDemoProps {
	copy: ExecutiveReviewMessages
	locale: WebsiteLocale
}

export function ExecutiveReviewDemo(props: ExecutiveReviewDemoProps) {
	const { copy, locale } = props

	const prefix = useId()
	const inspectorId = `${prefix}-source-context`
	const t = createTranslator({
		locale,
		messages: copy,
	})
	const number = new Intl.NumberFormat(locale)
	const [selection, setSelection] = useState<{
		filters: ReviewFilters
		sourceId: string | null
	}>({
		filters: initialReviewFilters,
		sourceId: null,
	})
	const [hasInteracted, setHasInteracted] = useState(false)
	const [chartReady, setChartReady] = useState(false)
	const rootRef = useRef<HTMLDivElement>(null)
	const inspectorRef = useRef<HTMLDivElement>(null)
	const inspectorHeadingRef = useRef<HTMLHeadingElement>(null)
	const sourceTriggerRef = useRef<HTMLButtonElement | null>(null)
	const revealSourceRef = useRef(false)
	const { filters, sourceId } = selection
	const visibleEvents = selectReviewEvents(filters)
	const totals = aggregateReviewEvents(visibleEvents)
	const reviews = getReviewsForEvents(visibleEvents)
	const selectedSource = reviewSources.find(
		(source) =>
			source.id === sourceId &&
			visibleEvents.some((event) => event.id === source.eventId),
	)
	const activeQuestion = filters.focus === 'all' ? 'overview' : filters.focus
	const chartData = reviewDistribution.map((row) => ({
		...row,
		label: copy.targets[row.targetId].label,
	}))

	useEffect(() => {
		if (sourceId && revealSourceRef.current) {
			revealSourceRef.current = false
			reveal(inspectorHeadingRef.current)
		}
	}, [sourceId])

	useEffect(() => {
		const root = rootRef.current
		if (!root || typeof IntersectionObserver === 'undefined') {
			setChartReady(true)
			return
		}

		const observer = new IntersectionObserver(
			(entries) => {
				if (!entries.some((entry) => entry.isIntersecting)) return
				setChartReady(true)
				observer.disconnect()
			},
			{ rootMargin: '300px' },
		)
		observer.observe(root)
		return () => observer.disconnect()
	}, [])

	function applyFilters(next: ReviewFilters, reset = false) {
		setHasInteracted(true)
		revealSourceRef.current = false
		setSelection((current) => ({
			filters: next,
			sourceId:
				!reset &&
				selectReviewEvents(next).some((event) =>
					event.sourceIds.includes(current.sourceId ?? ''),
				)
					? current.sourceId
					: null,
		}))
	}

	function openSource(id: string, trigger: HTMLButtonElement) {
		if (!visibleEvents.some((event) => event.sourceIds.includes(id))) return
		sourceTriggerRef.current = trigger
		if (sourceId === id) reveal(inspectorHeadingRef.current)
		else {
			revealSourceRef.current = true
			setSelection((current) => ({
				...current,
				sourceId: id,
			}))
		}
	}

	function closeSource() {
		if (inspectorRef.current?.contains(document.activeElement)) {
			const lastTrigger = sourceTriggerRef.current
			const visibleTrigger =
				lastTrigger?.isConnected && lastTrigger.getClientRects().length
					? lastTrigger
					: [
							...(rootRef.current?.querySelectorAll<HTMLButtonElement>(
								'[data-review-source]',
							) ?? []),
						].find(
							(button) =>
								button.dataset.reviewSource === sourceId &&
								button.getClientRects().length,
						)
			visibleTrigger?.focus({ preventScroll: true })
			visibleTrigger?.scrollIntoView({
				block: 'nearest',
				behavior: 'instant',
			})
		}
		revealSourceRef.current = false
		setSelection((current) => ({
			...current,
			sourceId: null,
		}))
	}

	function renderSource(id: string) {
		const source = reviewSources.find((item) => item.id === id)
		if (!source || !visibleEvents.some((event) => event.id === source.eventId))
			return null

		return (
			<button
				key={id}
				type="button"
				className="review-source-link"
				data-review-source={id}
				aria-controls={inspectorId}
				aria-expanded={sourceId === id}
				aria-label={t('source.referenceLabel', {
					id,
					eventTitle: copy.events[source.eventId].title,
					sourceKind: copy.kinds[source.kind],
				})}
				onClick={(event) => openSource(id, event.currentTarget)}
			>
				{id}
				<span aria-hidden="true">↗</span>
			</button>
		)
	}

	function revealReference(id: string) {
		reveal(
			rootRef.current?.querySelector<HTMLElement>(`[id="${prefix}-${id}"]`) ??
				null,
		)
	}

	const metrics = [
		{
			id: 'events',
			label: copy.metrics.events,
			value: number.format(totals.events),
		},
		{
			id: 'targets',
			label: copy.metrics.targets,
			value: number.format(totals.targets),
		},
		{
			id: 'reviews',
			label: copy.metrics.openReviews,
			value: number.format(totals.openReviews),
		},
		{
			id: 'evidence',
			label: copy.metrics.eventsWithEvidence,
			value: totals.events
				? t('metrics.evidenceRatio', {
						linked: totals.eventsWithEvidence,
						total: totals.events,
					})
				: '—',
		},
	]

	return (
		<div className="executive-review-demo" ref={rootRef}>
			<div className="review-question-heading">
				<span>{copy.questionLabel}</span>
				<span className="review-period">{copy.periodLabel}</span>
			</div>
			<fieldset className="review-questions" aria-label={copy.questionLabel}>
				{questions.map((question, index) => (
					<button
						className="min-w-0 whitespace-normal!"
						key={question}
						type="button"
						data-review-question={question}
						aria-pressed={activeQuestion === question}
						onClick={() => applyFilters(questionPresets[question], true)}
					>
						<span className="review-question-number" aria-hidden="true">
							0{index + 1}
						</span>
						<span className="min-w-0 flex-1 [overflow-wrap:anywhere]">
							{copy.questions[question]}
						</span>
						<span className="review-question-arrow" aria-hidden="true">
							↗
						</span>
					</button>
				))}
			</fieldset>

			<div className="review-workspace">
				<div className="review-selection-heading">
					<div>
						<span>{copy.labels.currentSelection}</span>
						<p>
							{t('labels.selectionSummary', {
								target: filters.targetId
									? copy.targets[filters.targetId].label
									: copy.labels.allTargets,
								classification: filters.classification
									? copy.classifications[filters.classification]
									: copy.labels.allClassifications,
								focus: copy.questions[activeQuestion],
							})}
						</p>
					</div>
					<button
						className="review-reset"
						type="button"
						onClick={() => applyFilters(initialReviewFilters, true)}
					>
						{copy.labels.reset}
						<span aria-hidden="true">↺</span>
					</button>
				</div>
				<output className="website-sr-only">
					{hasInteracted
						? t('labels.resultStatus', { count: totals.events })
						: ''}
				</output>
				<div
					className="review-answer"
					data-review-answer={totals.events ? activeQuestion : 'empty'}
				>
					<h3>
						{totals.events === 0
							? copy.answers.emptyTitle
							: activeQuestion === 'overview'
								? copy.answers.overviewTitle
								: activeQuestion === 'no-linked-evidence'
									? copy.answers.noEvidenceTitle
									: copy.answers.openReviewTitle}
					</h3>
					<p>
						{totals.events === 0
							? copy.answers.emptyDescription
							: activeQuestion === 'overview'
								? t('answers.overview', { ...totals })
								: activeQuestion === 'no-linked-evidence'
									? t('answers.noEvidence', {
											...totals,
											count: totals.events,
										})
									: t('answers.openReview', {
											...totals,
											count: totals.openReviews,
										})}
					</p>
					{activeQuestion === 'no-linked-evidence' && totals.events > 0 ? (
						<ul className="review-answer-events">
							{visibleEvents.map((event) => (
								<li key={event.id}>
									<button
										type="button"
										onClick={() => revealReference(`event-${event.id}`)}
										aria-label={t('labels.eventReference', {
											eventTitle: copy.events[event.id].title,
										})}
									>
										{copy.events[event.id].title}
										<span aria-hidden="true">↓</span>
									</button>
									<span>{copy.labels.noLinkedSources}</span>
									{reviews
										.filter((review) => review.eventId === event.id)
										.map((review) => (
											<p key={review.id}>
												{copy.reviews[review.id].nextDecision}
											</p>
										))}
								</li>
							))}
						</ul>
					) : null}
					{activeQuestion === 'open-review' && reviews.length > 0 ? (
						<ul className="review-answer-reviews">
							{reviews.map((review) => (
								<li key={review.id}>
									<button
										type="button"
										onClick={() => revealReference(`review-${review.id}`)}
										aria-label={t('labels.reviewReference', {
											reviewTitle: copy.reviews[review.id].title,
										})}
									>
										{copy.reviews[review.id].title}
										<span aria-hidden="true">↓</span>
									</button>
									<p>
										{review.id === 'R01'
											? copy.answers.manualSetup
											: copy.reviews[review.id].summary}
									</p>
									<div className="review-references">
										{review.checkpoints
											.flatMap((checkpoint) => checkpoint.sourceIds)
											.map(renderSource)}
									</div>
								</li>
							))}
						</ul>
					) : null}
				</div>
				<dl className="review-metrics">
					{metrics.map((metric) => (
						<div key={metric.id}>
							<dt>{metric.label}</dt>
							<dd data-review-metric={metric.id}>
								<span key={metric.value} className="review-metric-value">
									{metric.value}
								</span>
							</dd>
							{!totals.events && metric.id === 'evidence' ? (
								<p>{copy.metrics.noEvents}</p>
							) : null}
						</div>
					))}
				</dl>

				<div className="review-analysis">
					<div className="review-distribution">
						<div className="review-panel-heading">
							<h3>{copy.chart.title}</h3>
							<span className="review-chart-scope">{copy.chart.scope}</span>
							<p>{copy.chart.description}</p>
						</div>
						<fieldset className="review-filters">
							<legend>{copy.labels.target}</legend>
							<div>
								<button
									type="button"
									aria-pressed={filters.targetId === null}
									onClick={() =>
										applyFilters({
											...filters,
											targetId: null,
											classification: null,
										})
									}
								>
									{copy.labels.allTargets}
								</button>
								{targetIds.map((targetId) => (
									<button
										key={targetId}
										type="button"
										data-review-target={targetId}
										aria-pressed={filters.targetId === targetId}
										onClick={() =>
											applyFilters({
												...filters,
												targetId,
												classification: null,
											})
										}
									>
										{copy.targets[targetId].label}
									</button>
								))}
							</div>
						</fieldset>
						<fieldset className="review-filters">
							<legend>{copy.labels.classification}</legend>
							<div>
								<button
									type="button"
									aria-pressed={filters.classification === null}
									onClick={() =>
										applyFilters({
											...filters,
											classification: null,
										})
									}
								>
									{copy.labels.allClassifications}
								</button>
								{classificationFilters.map((classification) => (
									<button
										key={classification}
										type="button"
										data-review-class={classification}
										aria-pressed={filters.classification === classification}
										onClick={() =>
											applyFilters({
												...filters,
												classification,
											})
										}
									>
										<span
											className="review-class-dot"
											data-class={classification}
											aria-hidden="true"
										/>
										{copy.classifications[classification]}
									</button>
								))}
							</div>
						</fieldset>
						{chartReady ? (
							<Suspense
								fallback={
									<div
										aria-hidden="true"
										className="review-distribution-placeholder"
									/>
								}
							>
								<DeferredTargetDistributionChart
									data={chartData}
									locale={locale}
									labels={copy.classifications}
									selectedTarget={filters.targetId}
									selectedClass={filters.classification}
									onSelect={(id, classification) => {
										if (!targetIds.includes(id as TargetId)) return
										if (inspectorRef.current?.contains(document.activeElement))
											rootRef.current
												?.querySelector<HTMLButtonElement>(
													`[data-review-target="${id}"]`,
												)
												?.focus({ preventScroll: true })
										applyFilters({
											...filters,
											targetId: id as TargetId,
											classification,
										})
									}}
								/>
							</Suspense>
						) : (
							<div
								aria-hidden="true"
								className="review-distribution-placeholder"
							/>
						)}
						<p className="review-classification-note">{copy.explanation}</p>
						{filters.targetId ? (
							<p className="review-target-description">
								{copy.targets[filters.targetId].description}
							</p>
						) : null}
					</div>
					{totals.events ? (
						<EvidenceMatrix
							events={visibleEvents}
							copy={copy}
							prefix={prefix}
							selectedSourceId={sourceId}
							inspectorId={inspectorId}
							onOpen={openSource}
							t={t}
						/>
					) : (
						<div className="review-empty-matrix">
							<h3>{copy.matrix.title}</h3>
							<p>{copy.metrics.noEvents}</p>
						</div>
					)}
				</div>

				<div
					className="review-source-context"
					id={inspectorId}
					ref={inspectorRef}
					hidden={!selectedSource}
				>
					{selectedSource ? (
						<>
							<div className="review-source-heading">
								<div>
									<span>
										{copy.source.title} · {selectedSource.id}
									</span>
									<h3 tabIndex={-1} ref={inspectorHeadingRef}>
										{t('source.titleTemplate', {
											eventTitle: copy.events[selectedSource.eventId].title,
											sourceKind: copy.kinds[selectedSource.kind],
										})}
									</h3>
								</div>
								<button
									type="button"
									onClick={closeSource}
									aria-label={copy.source.close}
								>
									<span aria-hidden="true">×</span>
								</button>
							</div>
							<blockquote>{copy.sources[selectedSource.id]}</blockquote>
						</>
					) : null}
				</div>
				{totals.events ? (
					<ReviewTimeline
						reviews={reviews}
						copy={copy}
						prefix={prefix}
						renderSource={renderSource}
						t={t}
					/>
				) : null}
			</div>
			<a className="smooth-text-link review-early-access" href="#book-a-demo">
				{copy.labels.earlyAccess}
				<span aria-hidden="true">↗</span>
			</a>
		</div>
	)
}
