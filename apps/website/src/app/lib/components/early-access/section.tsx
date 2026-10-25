// import { EarlyAccessForm } from '@components/early-access/form'
import { earlyAccessMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import './early-access.css'

interface EarlyAccessSectionProps {
	locale: WebsiteLocale
	source: 'landing' | 'pricing'
	compact?: boolean
	/** A product page's own sentence for the landing band. */
	description?: string
}

// Sign-up is paused until feature flags have a proper home, so the section
// shows the return-soon notice. To open it, restore the form import, the
// helper id and the commented section below, which used `compact`.
export function EarlyAccessSection(props: EarlyAccessSectionProps) {
	const { description, locale, source } = props

	const copy = earlyAccessMessages[locale]
	const headingId = `early-access-${source}-heading`
	// const helperId = `early-access-${source}-helper`

	return (
		<section
			id="early-access"
			className="mx-auto max-w-[1240px] scroll-mt-25 px-4.5 pt-10 pb-14 sm:px-7 sm:pt-17 sm:pb-21.5"
			aria-labelledby={headingId}
		>
			<div className="rounded-2xl border border-border border-t-3 border-t-primary bg-card p-6.5 text-foreground sm:p-11">
				{source === 'landing' ? (
					<>
						<h2
							id={headingId}
							className="m-0 text-3xl leading-tight font-medium tracking-tight text-balance sm:text-4xl"
						>
							{copy.demoHeading}
						</h2>
						<p className="mt-5 mb-0 max-w-3xl text-base leading-relaxed text-muted-foreground">
							{description ?? copy.demoDescription}
						</p>
						<p className="mt-8 mb-0 text-xl font-medium">
							{copy.pausedHeading}
						</p>
					</>
				) : (
					<h2
						id={headingId}
						className="m-0 text-3xl leading-tight font-medium tracking-tight text-balance"
					>
						{copy.pausedHeading}
					</h2>
				)}
				<p className="mt-3 mb-0 max-w-3xl text-base leading-relaxed text-muted-foreground">
					{copy.pausedDescription}
				</p>
			</div>
		</section>
	)

	// return (
	// 	<section
	// 		id="early-access"
	// 		className={`early-access-section${compact ? ' early-access-compact' : ''}`}
	// 		aria-labelledby={headingId}
	// 	>
	// 		<div className="early-access-surface">
	// 			<div className="early-access-intro">
	// 				<h2 id={headingId}>
	// 					{source === 'pricing' ? copy.pricingHeading : copy.homeHeading}
	// 				</h2>
	// 				<p>{copy.description}</p>
	// 				<p id={helperId} className="early-access-trial-note">
	// 					{copy.helper}
	// 				</p>
	// 			</div>
	// 			<EarlyAccessForm
	// 				copy={copy}
	// 				helperId={helperId}
	// 				locale={locale}
	// 				source={source}
	// 			/>
	// 		</div>
	// 	</section>
	// )
}
