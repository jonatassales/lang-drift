import { conversionMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import './conversion.css'

const planIds = ['trial', 'plus', 'pro'] as const

interface PlansSectionProps {
	locale: WebsiteLocale
}

export function PlansSection(props: PlansSectionProps) {
	const { locale } = props

	const copy = conversionMessages[locale].plans

	return (
		<section className="plans-section" aria-labelledby="plans-heading">
			<div className="plans-introduction">
				<span className="section-kicker">{copy.eyebrow}</span>
				<h1 id="plans-heading">{copy.title}</h1>
				<p>{copy.introduction}</p>
			</div>
			<div className="plans-grid">
				{planIds.map((id) => {
					const plan = copy.cards[id]

					return (
						<article
							className="website-plan"
							data-plan={id}
							key={id}
							aria-labelledby={`plan-${id}-heading`}
						>
							<div className="website-plan-header">
								<h2 id={`plan-${id}-heading`}>{plan.name}</h2>
								<span className="website-plan-badge">{plan.badge}</span>
							</div>
							<div className="website-plan-value">
								<strong>{plan.value}</strong>
								<p>{plan.support}</p>
							</div>
							<p className="website-plan-description">{plan.description}</p>
							<div className="website-plan-features">
								<h3>{plan.listLabel}</h3>
								<ul>
									{plan.items.map((item) => (
										<li key={item}>{item}</li>
									))}
								</ul>
							</div>
							<a
								className={
									id === 'trial'
										? 'smooth-primary-button website-plan-cta'
										: 'website-plan-cta website-plan-secondary'
								}
								href="#early-access"
							>
								{plan.cta}
								<span aria-hidden="true">↗</span>
							</a>
						</article>
					)
				})}
			</div>
			<div className="plans-availability">
				<p>{copy.availability}</p>
				<Link className="smooth-text-link" href="/platform/vision#faq">
					{copy.faqLink}
					<span aria-hidden="true">→</span>
				</Link>
			</div>
		</section>
	)
}
