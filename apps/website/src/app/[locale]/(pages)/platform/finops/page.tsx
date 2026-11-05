import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { finopsMessages, headerMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'
import { gateExample } from '@routes/platform/finops/content/finops-content.data'
import { gateFigures } from '@routes/platform/finops/content/gate-figures.compute'
import { BudgetLedgerDemo } from '@routes/platform/finops/sections/budget-ledger-demo'
import { BudgetTerms } from '@routes/platform/finops/sections/budget-terms'
import { DeliveryReceipt } from '@routes/platform/finops/sections/delivery-receipt'
import { FailureScenarios } from '@routes/platform/finops/sections/failure-scenarios'
import { GateFlow } from '@routes/platform/finops/sections/gate-flow'
import { SourceList } from '@routes/platform/finops/sections/source-list'
import { ProductHeading } from '@routes/platform/product-heading'
import { ProductSection } from '@routes/platform/product-page/product-section'
import { ProductTagline } from '@routes/platform/product-tagline'
import { FinopsIcon } from 'web/ui/icons/platform'

export { generateMetadata } from './page-metadata.map'

interface FinopsPageProps {
	params: Promise<{ locale: string }>
}

export default async function FinopsPage(props: FinopsPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = finopsMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="finops" />
			<SkipLink href="#authorization">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-describedby="finops-hero-caption finops-hero-summary"
					aria-labelledby="finops-hero-heading"
					className="mx-auto max-w-[1280px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<ProductHeading icon={FinopsIcon} id="finops-hero-heading">
						{headerMessages[locale].finops}
					</ProductHeading>
					<ProductTagline
						end={copy.hero.taglineEnd}
						start={copy.hero.taglineStart}
					/>
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
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
