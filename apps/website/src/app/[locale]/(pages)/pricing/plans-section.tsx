import { conversionMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import { ProductBadge } from '@layout/product-badge'
import { Button } from 'web/ui/button'
import { TextLink } from 'web/ui/text-link'

const planIds = ['starter', 'plus', 'pro'] as const

const planCtaClass =
	'flex min-h-11 justify-between gap-3 px-[15px] py-3 leading-normal'

const availabilityLinkClass =
	'inline-flex min-h-11 items-center gap-2.5 text-xs focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary'

interface PlansSectionProps {
	locale: WebsiteLocale
}

export function PlansSection(props: PlansSectionProps) {
	const { locale } = props

	const copy = conversionMessages[locale].plans

	return (
		<section
			aria-labelledby="plans-heading"
			className="mx-auto max-w-[1240px] px-[18px] pt-[38px] pb-5 sm:px-7 sm:pt-12 sm:pb-6 [&_:is(h1,h2,h3,p,strong,li,a)]:[overflow-wrap:anywhere]"
		>
			<div className="mx-auto mb-8 max-w-[740px] text-left sm:text-center">
				<div className="flex flex-wrap items-center gap-2 sm:justify-center">
					<span className="font-mono text-[10px] tracking-[0.11em] uppercase">
						{copy.eyebrow}
					</span>
					<ProductBadge>{copy.comingSoon}</ProductBadge>
				</div>
				<h1
					className="my-4 text-[clamp(34px,3.6vw,50px)] font-medium leading-[1.12] tracking-[-0.055em] [text-wrap:balance]"
					id="plans-heading"
				>
					{copy.title}
				</h1>
				<p className="m-0 text-base leading-[1.65] text-muted-foreground">
					{copy.introduction}
				</p>
			</div>
			<div className="mx-auto grid max-w-[640px] grid-cols-[minmax(0,1fr)] justify-center gap-4.5 lg:mx-0 lg:max-w-none lg:grid-cols-[repeat(auto-fit,minmax(280px,390px))]">
				{planIds.map((id) => {
					const plan = copy.cards[id]

					return (
						<article
							aria-labelledby={`plan-${id}-heading`}
							className="flex min-w-0 flex-col rounded-[15px] border border-primary/65 bg-card p-6 shadow-[0_0_0_1px_color-mix(in_srgb,var(--primary)_6%,transparent),0_8px_30px_color-mix(in_srgb,var(--primary)_4%,transparent)] sm:p-[26px]"
							data-plan={id}
							key={id}
						>
							<h2
								className="m-0 text-[21px] font-medium leading-[1.35] tracking-[-0.03em]"
								id={`plan-${id}-heading`}
							>
								{plan.name}
							</h2>
							<div className="mt-[22px] lg:mt-7">
								<strong className="block text-3xl font-medium leading-[1.25] tracking-[-0.04em]">
									{plan.value}
								</strong>
								<p className="mt-2.5 text-xs leading-[1.55] text-muted-foreground">
									{plan.support}
								</p>
							</div>
							<p className="my-[22px] text-sm leading-[1.65]">
								{plan.description}
							</p>
							<div className="mt-0 border-t border-border pt-5 pb-6 lg:mt-auto">
								<h3 className="mb-3 text-[11px] font-normal leading-normal text-muted-foreground">
									{plan.listLabel}
								</h3>
								<ul className="m-0 grid gap-[9px] pl-[17px] text-[13px] leading-[1.55] marker:text-muted-foreground">
									{plan.items.map((item) => (
										<li key={item}>{item}</li>
									))}
								</ul>
							</div>
							<Button asChild className={planCtaClass} variant="pill">
								<Link href="/book-a-demo">
									{plan.cta}
									<span aria-hidden="true" className="shrink-0">
										↗
									</span>
								</Link>
							</Button>
						</article>
					)
				})}
			</div>
			<div className="mx-auto mt-[22px] flex max-w-[640px] flex-wrap items-baseline justify-between gap-x-6 gap-y-3 lg:mx-0 lg:max-w-none">
				<p className="m-0 max-w-[730px] text-xs leading-[1.65] text-muted-foreground">
					{copy.availability}
				</p>
				<TextLink asChild className={availabilityLinkClass}>
					<Link href="/#faq">
						{copy.faqLink}
						<span aria-hidden="true">→</span>
					</Link>
				</TextLink>
			</div>
		</section>
	)
}
