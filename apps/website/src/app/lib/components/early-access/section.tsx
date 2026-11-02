// import { EarlyAccessForm } from '@components/early-access/form'
import { earlyAccessMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface EarlyAccessSectionProps {
	locale: WebsiteLocale
	compact?: boolean
}

// The pricing page's early-access band. Sign-up is paused until feature
// flags have a proper home, so the section shows the return-soon notice. To
// open it, restore the form import, the helper id and the commented section
// below, which used `compact`.
export function EarlyAccessSection(props: EarlyAccessSectionProps) {
	const { locale } = props

	const copy = earlyAccessMessages[locale]
	const headingId = 'early-access-pricing-heading'
	// const helperId = 'early-access-pricing-helper'

	return (
		<section
			id="early-access"
			className="mx-auto max-w-[1240px] scroll-mt-25 px-4.5 pt-10 pb-14 sm:px-7 sm:pt-17 sm:pb-21.5"
			aria-labelledby={headingId}
		>
			<div className="rounded-2xl border border-border border-t-3 border-t-primary bg-card p-6.5 text-foreground sm:p-11">
				<h2
					id={headingId}
					className="m-0 text-3xl leading-tight font-medium tracking-tight text-balance"
				>
					{copy.pausedHeading}
				</h2>
				<p className="mt-3 mb-0 max-w-3xl text-base leading-relaxed text-muted-foreground">
					{copy.pausedDescription}
				</p>
			</div>
		</section>
	)

	// return (
	// 	<section
	// 		id="early-access"
	// 		className={`mx-auto max-w-[1240px] scroll-mt-25 ${compact ? 'py-6 sm:py-7.5 md:py-9.5' : 'px-4.5 pt-10 pb-14 sm:px-7 sm:pt-17 sm:pb-21.5'}`}
	// 		aria-labelledby={headingId}
	// 	>
	// 		<div className={`grid grid-cols-1 items-center rounded-2xl border border-border border-t-3 border-t-primary bg-card md:grid-cols-2 ${compact ? 'gap-8 p-7 sm:p-8' : 'gap-7 p-6.5 sm:p-11 md:gap-12'}`}>
	// 			<div className="min-w-0">
	// 				<h2 id={headingId} className={`m-0 font-medium leading-[1.3] tracking-[-0.035em] [text-wrap:balance] ${compact ? 'mb-2.5 text-[25px]' : 'mb-4.5 text-[clamp(26px,2.6vw,34px)]'}`}>{copy.pricingHeading}</h2>
	// 				<p className="m-0 text-base leading-[1.7] text-muted-foreground">{copy.description}</p>
	// 				<p id={helperId} className="mt-4 text-[13px] text-muted-foreground">
	// 					{copy.helper}
	// 				</p>
	// 			</div>
	// 			<EarlyAccessForm
	// 				copy={copy}
	// 				helperId={helperId}
	// 				locale={locale}
	// 				source="pricing"
	// 			/>
	// 		</div>
	// 	</section>
	// )
}
