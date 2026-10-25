import type { ReactNode } from 'react'

import type {
	ExecutiveReviewMessages,
	ReviewItem,
} from '@data/executive-review-data'

type SourceRenderer = (id: string) => ReactNode
type PanelTranslate = (
	key: 'labels.week' | 'matrix.linkedCellLabel',
	values: Record<string, string | number>,
) => string

interface ReviewTimelineProps {
	reviews: readonly ReviewItem[]
	copy: ExecutiveReviewMessages
	prefix: string
	renderSource: SourceRenderer
	t: PanelTranslate
}

export function ReviewTimeline(props: ReviewTimelineProps) {
	const { reviews, copy, prefix, renderSource, t } = props

	return (
		<section
			className="review-timeline"
			aria-labelledby={`${prefix}-timeline-title`}
		>
			<div className="review-panel-heading">
				<h3 id={`${prefix}-timeline-title`}>{copy.timeline.title}</h3>
			</div>
			{reviews.length === 0 ? (
				<div className="review-no-reviews">
					<p>{copy.timeline.noReviews}</p>
					<p>{copy.timeline.noReviewsLimit}</p>
				</div>
			) : (
				<ol className="review-items">
					{reviews.map((review) => {
						const text = copy.reviews[review.id]

						return (
							<li key={review.id}>
								<article
									id={`${prefix}-review-${review.id}`}
									tabIndex={-1}
									data-review-item={review.id}
								>
									<div className="review-item-intro">
										<span className="review-event-id">{review.id}</span>
										<h4>{text.title}</h4>
										<p className="review-related-event">
											{copy.labels.relatedEvent}:{' '}
											{copy.events[review.eventId].title}
										</p>
										<p>{text.summary}</p>
									</div>
									<ol className="review-checkpoints">
										{review.checkpoints.map((checkpoint) => (
											<li
												key={checkpoint.id}
												data-pending={checkpoint.week === null || undefined}
											>
												<span className="review-checkpoint-time">
													{checkpoint.week === null
														? copy.labels.pendingCondition
														: t('labels.week', { week: checkpoint.week })}
												</span>
												<h5>{text.checkpoints[checkpoint.id].label}</h5>
												<p>{text.checkpoints[checkpoint.id].description}</p>
												{checkpoint.sourceIds.length ? (
													<div className="review-references">
														{checkpoint.sourceIds.map(renderSource)}
													</div>
												) : null}
											</li>
										))}
									</ol>
									<div className="review-next-decision">
										<span>{copy.labels.nextDecision}</span>
										<p>{text.nextDecision}</p>
									</div>
									<p className="review-limit">{text.limit}</p>
								</article>
							</li>
						)
					})}
				</ol>
			)}
		</section>
	)
}
