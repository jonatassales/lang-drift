import { executiveReviewMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

import { ExecutiveReviewDemo } from './executive-review-demo'
import './executive-review.css'

interface ExecutiveReviewSectionProps {
	locale: WebsiteLocale
}

export function ExecutiveReviewSection(props: ExecutiveReviewSectionProps) {
	const { locale } = props

	const copy = executiveReviewMessages[locale]

	return (
		<section
			id="executive-review"
			className="executive-review-section"
			aria-labelledby="executive-review-heading"
		>
			<div className="executive-review-intro">
				<span className="section-kicker">{copy.eyebrow}</span>
				<h2 id="executive-review-heading">{copy.title}</h2>
				<p>{copy.description}</p>
			</div>
			<ExecutiveReviewDemo copy={copy} locale={locale} />
		</section>
	)
}
