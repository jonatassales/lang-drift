import { BasicAccordion } from '@langdrift/web/vendors/smoothui/basic-accordion'

import { conversionMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import './conversion.css'

const faqIds = [
	'meaning',
	'expected',
	'connection',
	'tools',
	// 'trial',
	// 'free',
	// 'paid',
	'availability',
] as const

interface FaqSectionProps {
	locale: WebsiteLocale
}

export function FaqSection(props: FaqSectionProps) {
	const { locale } = props

	const copy = conversionMessages[locale].faq

	return (
		<section
			className="conversion-section faq-section"
			id="faq"
			aria-labelledby="faq-heading"
		>
			<div className="faq-introduction">
				<h2 id="faq-heading">{copy.title}</h2>
				<p className="conversion-introduction">{copy.description}</p>
				<p className="faq-offer">{copy.offer}</p>
			</div>
			<BasicAccordion
				className="website-faq [&_button]:whitespace-normal!"
				headingLevel={3}
				layoutId="website-faq"
				items={faqIds.map((id) => ({
					id,
					title: copy.items[id].question,
					content: <p>{copy.items[id].answer}</p>,
				}))}
			/>
		</section>
	)
}
