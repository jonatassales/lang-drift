import {
	getEventSources,
	reviewEvents,
	reviewItems,
} from '@data/executive-review-data'
import {
	conversionMessages,
	executiveReviewMessages,
} from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import './conversion.css'

const benefitIds = ['context', 'tradeoff', 'opportunity'] as const

interface RoiSectionProps {
	locale: WebsiteLocale
}

export function RoiSection(props: RoiSectionProps) {
	const { locale } = props

	const copy = conversionMessages[locale].roi
	const reviewCopy = executiveReviewMessages[locale]
	const event = reviewEvents.find((item) => item.id === 'E02')
	const review = reviewItems.find((item) => item.id === 'R01')
	if (!event || !review || review.eventId !== event.id) return null
	const sources = getEventSources(event).filter(
		(source) => source.kind === 'decision' || source.kind === 'verification',
	)

	return (
		<section
			className="conversion-section roi-section"
			id="roi"
			aria-labelledby="roi-heading"
		>
			<div className="roi-explanation">
				<span className="section-kicker">{copy.eyebrow}</span>
				<h2 id="roi-heading">{copy.title}</h2>
				<p className="conversion-introduction">{copy.description}</p>
				<ul className="roi-benefits">
					{benefitIds.map((id) => (
						<li key={id}>
							<h3>{copy.benefits[id].title}</h3>
							<p>{copy.benefits[id].description}</p>
						</li>
					))}
				</ul>
			</div>
			<article className="roi-brief" aria-labelledby="roi-brief-heading">
				<div className="roi-brief-label">
					<span>{copy.brief.label}</span>
					<span>
						{event.id} · {review.id}
					</span>
				</div>
				<h3 id="roi-brief-heading">
					{reviewCopy.reviews.R01.checkpoints.pending.label}
				</h3>
				<p>{copy.brief.context}</p>
				<blockquote>{copy.brief.question}</blockquote>
				<ul className="roi-evidence">
					{copy.brief.evidence.map((label) => (
						<li key={label}>{label}</li>
					))}
				</ul>
				<ul className="roi-records" aria-label={reviewCopy.matrix.title}>
					{sources.map((source) => (
						<li key={source.id}>
							<code>{source.id}</code>
							<span>{reviewCopy.kinds[source.kind]}</span>
						</li>
					))}
				</ul>
				<p className="roi-brief-limit">{reviewCopy.reviews.R01.limit}</p>
			</article>
			<div className="conversion-actions roi-actions">
				<a className="smooth-primary-button" href="#early-access">
					{copy.earlyAccess}
					<span aria-hidden="true">↗</span>
				</a>
				{/* <Link className="smooth-text-link" href="/pricing" id="pricing">
          {copy.explorePlans}
          <span aria-hidden="true">→</span>
        </Link> */}
			</div>
		</section>
	)
}
