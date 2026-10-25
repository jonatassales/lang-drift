import type { DetailListItem } from '@langdrift/web/ui/detail-list'
import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { InfoCards } from '@langdrift/web/ui/info-cards'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { createTranslator, hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { EarlyAccessSection } from '@components/early-access/section'
import { BudgetLedgerDemo } from '@components/finops/budget-ledger-demo'
import { DeliveryReceipt } from '@components/finops/delivery-receipt'
import { PillarScenarios } from '@components/product-page/pillar-scenarios'
import { ProductFaqSection } from '@components/product-page/product-faq-section'
import { ProductSection } from '@components/product-page/product-section'
import { RecordCard } from '@components/product-page/record-card'
import { SectionOutro } from '@components/product-page/section-outro'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import {
	authorizationCardIds,
	budgetCardIds,
	chainPillarIds,
	chainScenarioIds,
	envelope,
	envelopeSegments,
	finopsFaqIds,
	workItems,
} from '@data/finops-content'
import { finopsMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
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
	const formatReport = createTranslator({
		locale,
		messages: copy.report,
	})
	const percent = new Intl.NumberFormat(locale, {
		style: 'percent',
		maximumFractionDigits: 0,
	})
	const segments = envelopeSegments(envelope)
	const observedItems = workItems.flatMap((item) =>
		item.observed === null
			? []
			: [
					{
						id: item.id,
						observed: item.observed,
						retries: item.retries,
					},
				],
	)
	const largestConsumer = observedItems.reduce((largest, item) =>
		item.observed > largest.observed ? item : largest,
	)
	const unknownCharges = workItems.length - observedItems.length
	const reportRows: readonly DetailListItem[] = [
		{
			detail: formatReport('rows.envelopeUsed.detail', {
				used: percent.format(
					(segments.spent + segments.reserved) / envelope.authorized,
				),
				available: percent.format(segments.available / envelope.authorized),
			}),
			id: 'envelope',
			term: copy.report.rows.envelopeUsed.term,
		},
		{
			detail: formatReport('rows.largestConsumer.detail', {
				item: copy.ledger.items[largestConsumer.id].title,
				count: largestConsumer.retries + 1,
			}),
			id: 'largest',
			term: copy.report.rows.largestConsumer.term,
		},
		{
			detail: formatReport('rows.unknownCharges.detail', {
				count: unknownCharges,
			}),
			id: 'unknown',
			term: copy.report.rows.unknownCharges.term,
			tone: unknownCharges > 0 ? 'warning' : 'success',
		},
		{
			detail: copy.report.rows.attention.detail,
			id: 'attention',
			term: copy.report.rows.attention.term,
			tone: 'danger',
		},
	]

	return (
		<>
			<SkipLink href="#authorization">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
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
					<div
						aria-describedby="finops-hero-caption finops-hero-summary"
						className="mt-8 grid gap-6 lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)]"
					>
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
					id="chain"
					lead={copy.chain.lead}
					title={copy.chain.title}
				>
					<PillarScenarios
						pillars={chainPillarIds.map((id) => ({
							id,
							label: copy.chain.pillars[id],
						}))}
						scenarios={chainScenarioIds.map((id) => {
							const scenario = copy.chain.scenarios[id]

							return {
								cells: scenario,
								context: scenario.context,
								id,
								tabLabel: scenario.tabLabel,
							}
						})}
						selectorLabel={copy.chain.selectorLabel}
					/>
					<SectionOutro
						fragment="authorization"
						link={copy.chain.outroLink}
						text={copy.chain.outro}
					/>
				</ProductSection>

				<ProductSection
					id="budget"
					kicker={copy.budget.kicker}
					lead={copy.budget.lead}
					title={copy.budget.title}
				>
					<InfoCards
						className="mt-11"
						items={budgetCardIds.map((id) => ({
							...copy.budget.cards[id],
							id,
						}))}
					/>
				</ProductSection>

				<ProductSection
					id="authorization"
					kicker={copy.authorization.kicker}
					lead={copy.authorization.lead}
					title={copy.authorization.title}
				>
					<InfoCards
						className="mt-11"
						items={authorizationCardIds.map((id) => ({
							...copy.authorization.cards[id],
							id,
						}))}
						ordered
					/>
				</ProductSection>

				<ProductSection
					id="report"
					kicker={copy.report.kicker}
					lead={copy.report.lead}
					title={copy.report.title}
				>
					<RecordCard
						items={reportRows}
						period={copy.report.period}
						title={copy.report.weekly}
					/>
				</ProductSection>

				<ProductSection
					id="sources"
					kicker={copy.sources.kicker}
					lead={copy.sources.lead}
					title={copy.sources.title}
				>
					<InfoCards
						className="mt-11"
						columns={4}
						items={[
							{
								...copy.sources.items.providers,
								id: 'providers',
							},
							{
								description: copy.sources.items.microruntimes.description,
								footer: (
									<Link
										className="smooth-text-link"
										href="/platform/microruntimes"
										locale={locale}
									>
										{copy.sources.items.microruntimes.link}
									</Link>
								),
								id: 'microruntimes',
								title: copy.sources.items.microruntimes.title,
							},
							{
								...copy.sources.items.evidence,
								id: 'evidence',
							},
							{
								...copy.sources.items.billing,
								id: 'billing',
							},
						]}
					/>
					<p className="mt-6 mb-0 max-w-[760px] text-[13px] leading-[1.6] text-muted-foreground">
						{copy.sources.intendedModel}
					</p>
				</ProductSection>

				<ProductFaqSection
					description={copy.faq.description}
					items={finopsFaqIds.map((id) => ({
						...copy.faq.items[id],
						id,
					}))}
					title={copy.faq.title}
				/>
				<EarlyAccessSection locale={locale} source="landing" />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
