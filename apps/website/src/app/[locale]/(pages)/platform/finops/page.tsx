import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { BudgetLedgerDemo } from '@components/finops/budget-ledger-demo'
import { BudgetTerms } from '@components/finops/budget-terms'
import { DeliveryReceipt } from '@components/finops/delivery-receipt'
import { FailureScenarios } from '@components/finops/failure-scenarios'
import { GateFlow } from '@components/finops/gate-flow'
import { SourceList } from '@components/finops/source-list'
import { ProductFaqSection } from '@components/product-page/product-faq-section'
import { ProductSection } from '@components/product-page/product-section'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { finopsFaqIds, gateExample, gateFigures } from '@data/finops-content'
import { finopsMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface FinopsPageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: FinopsPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const images = await socialImages(parent)

	return pageMetadata('finops', locale, images)
}

export default async function FinopsPage(props: FinopsPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = finopsMessages[locale]

	return (
		<>
			<SkipLink href="#authorization">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-describedby="finops-hero-caption finops-hero-summary"
					aria-labelledby="finops-hero-heading"
					className="mx-auto max-w-[1280px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<Eyebrow className="block text-[10px] tracking-[0.11em]">
						{copy.hero.eyebrow}
					</Eyebrow>
					<h1
						className="m-0 mt-4 text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance"
						id="finops-hero-heading"
					>
						<span className="block">{copy.hero.titleStart}</span>{' '}
						<span className="block">{copy.hero.titleEnd}</span>
					</h1>
					<div className="mt-8 grid gap-6 lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)]">
						<BudgetLedgerDemo copy={copy.ledger} locale={locale} />
						<DeliveryReceipt
							copy={copy.receipt}
							items={copy.ledger.items}
							locale={locale}
						/>
					</div>
					<p
						className="mt-5 mb-0 max-w-[880px] text-xs leading-[1.65] text-muted-foreground"
						id="finops-hero-caption"
					>
						{copy.hero.caption}
					</p>
					<p className="sr-only" id="finops-hero-summary">
						{copy.hero.summary}
					</p>
				</section>

				<ProductSection
					id="authorization"
					kicker={copy.authorization.kicker}
					lead={copy.authorization.lead}
					title={copy.authorization.title}
				>
					<GateFlow
						copy={copy.authorization}
						exampleTitle={copy.ledger.items[gateExample.id].title}
						figures={gateFigures(gateExample)}
						locale={locale}
					/>
				</ProductSection>

				<ProductSection
					id="failures"
					kicker={copy.failures.kicker}
					lead={copy.failures.lead}
					title={copy.failures.title}
				>
					<FailureScenarios copy={copy.failures} />
				</ProductSection>

				<ProductSection
					id="budget"
					kicker={copy.budget.kicker}
					lead={copy.budget.lead}
					title={copy.budget.title}
				>
					<BudgetTerms copy={copy.budget.terms} />
				</ProductSection>

				<ProductSection
					id="sources"
					kicker={copy.sources.kicker}
					lead={copy.sources.lead}
					title={copy.sources.title}
				>
					<SourceList copy={copy.sources} locale={locale} />
				</ProductSection>

				<ProductFaqSection
					description={copy.faq.description}
					items={finopsFaqIds.map((id) => ({
						...copy.faq.items[id],
						id,
					}))}
					title={copy.faq.title}
				/>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
