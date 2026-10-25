import { DetailList, type DetailListItem } from '@langdrift/web/ui/detail-list'
import { InfoCards } from '@langdrift/web/ui/info-cards'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { createTranslator, hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { EarlyAccessSection } from '@components/early-access/section'
import { ProductFaqSection } from '@components/product-page/product-faq-section'
import { ProductHero } from '@components/product-page/product-hero'
import { ProductSection } from '@components/product-page/product-section'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import {
	envelope,
	envelopeSegments,
	type FinopsFaqId,
	workItems,
} from '@data/finops-content'
import { finopsMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import { routing } from '@i18n/routing'

interface FinopsPageProps {
	params: Promise<{ locale: string }>
}

const budgetCards = ['envelope', 'allocation', 'reservation'] as const
const authorizationCards = ['recommend', 'gate', 'reconcile'] as const
const faqItems: readonly FinopsFaqId[] = [
	'budget',
	'spending',
	'unknown',
	'authorization',
	'availability',
]

export async function generateMetadata(
	{ params }: FinopsPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('finops', locale, await socialImages(parent))
}

export default async function FinopsPage(props: FinopsPageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = finopsMessages[locale]
	const report = createTranslator({
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
			detail: report('rows.envelopeUsed.detail', {
				used: percent.format(
					(segments.spent + segments.reserved) / envelope.authorized,
				),
				available: percent.format(segments.available / envelope.authorized),
			}),
			id: 'envelope',
			term: copy.report.rows.envelopeUsed.term,
		},
		{
			detail: report('rows.largestConsumer.detail', {
				item: copy.ledger.items[largestConsumer.id].title,
				count: largestConsumer.retries + 1,
			}),
			id: 'largest',
			term: copy.report.rows.largestConsumer.term,
		},
		{
			detail: report('rows.unknownCharges.detail', { count: unknownCharges }),
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
			<a
				className="fixed top-0 left-3 z-200 max-w-[calc(100vw-1.5rem)] -translate-y-full bg-foreground px-2.5 py-2 whitespace-normal text-card focus-visible:top-3 focus-visible:translate-y-0"
				href="#budget"
			>
				{copy.skip}
			</a>
			<WebsiteHeader />
			<main id="top">
				<ProductHero
					caption={copy.hero.caption}
					headingId="finops-hero-heading"
					lead={copy.hero.lead}
					primaryCta={{
						href: '#budget',
						label: copy.hero.primaryCta,
					}}
					secondaryCta={{
						href: '#chain',
						label: copy.hero.secondaryCta,
					}}
					summary={copy.hero.summary}
					titleEnd={copy.hero.titleEnd}
					titleHighlight={copy.hero.titleHighlight}
					titleStart={copy.hero.titleStart}
				/>

				<ProductSection
					id="budget"
					kicker={copy.budget.kicker}
					lead={copy.budget.lead}
					title={copy.budget.title}
				>
					<InfoCards
						className="mt-11"
						items={budgetCards.map((id) => ({
							description: copy.budget.cards[id].description,
							eyebrow: copy.budget.cards[id].eyebrow,
							id,
							title: copy.budget.cards[id].title,
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
						items={authorizationCards.map((id) => ({
							description: copy.authorization.cards[id].description,
							eyebrow: copy.authorization.cards[id].eyebrow,
							id,
							title: copy.authorization.cards[id].title,
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
					<article className="mt-11 rounded-2xl border border-border bg-card p-6 sm:p-8">
						<div className="flex flex-wrap items-baseline justify-between gap-x-6 gap-y-2 border-b border-border pb-4">
							<strong className="text-lg font-medium tracking-[-0.02em]">
								{copy.report.weekly}
							</strong>
							<span className="font-mono text-[10px] tracking-[0.08em] text-muted-foreground uppercase">
								{copy.report.period}
							</span>
						</div>
						<DetailList className="mt-1" items={reportRows} />
					</article>
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
								description: copy.sources.items.providers.description,
								id: 'providers',
								title: copy.sources.items.providers.title,
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
								description: copy.sources.items.evidence.description,
								id: 'evidence',
								title: copy.sources.items.evidence.title,
							},
							{
								description: copy.sources.items.billing.description,
								id: 'billing',
								title: copy.sources.items.billing.title,
							},
						]}
					/>
					<p className="mt-6 mb-0 max-w-[760px] text-[13px] leading-[1.6] text-muted-foreground">
						{copy.sources.intendedModel}
					</p>
				</ProductSection>

				<ProductFaqSection
					description={copy.faq.description}
					items={faqItems.map((id) => ({
						answer: copy.faq.items[id].answer,
						id,
						question: copy.faq.items[id].question,
					}))}
					title={copy.faq.title}
				/>
				<EarlyAccessSection locale={locale} source="landing" />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
