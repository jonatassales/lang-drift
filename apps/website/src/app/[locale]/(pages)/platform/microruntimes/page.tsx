import { DetailList } from '@langdrift/web/ui/detail-list'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { createTranslator, hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { EarlyAccessSection } from '@components/early-access/section'
import { AuthoringFlow } from '@components/microruntimes/authoring-flow'
import { LifecycleDemo } from '@components/microruntimes/lifecycle-demo'
import { RuntimeDescriptorCard } from '@components/microruntimes/runtime-descriptor-card'
import { StageEvidence } from '@components/microruntimes/stage-evidence'
import { PillarScenarios } from '@components/product-page/pillar-scenarios'
import { ProductFaqSection } from '@components/product-page/product-faq-section'
import { ProductHero } from '@components/product-page/product-hero'
import { ProductSection } from '@components/product-page/product-section'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import {
	executionLedger,
	lifecyclePhaseIds,
	lifecycleScenarioIds,
	microruntimesFaqIds,
	ownershipRowIds,
	stateCardIds,
} from '@data/microruntimes-content'
import { microruntimesMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface MicroruntimesPageProps {
	params: Promise<{ locale: string }>
}

const stateTones = {
	requested: 'text-foreground',
	observed: 'text-status-success',
	unknown: 'text-status-warning',
} as const

export async function generateMetadata(
	{ params }: MicroruntimesPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('microruntimes', locale, await socialImages(parent))
}

export default async function MicroruntimesPage(props: MicroruntimesPageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = microruntimesMessages[locale]
	const evidence = createTranslator({
		locale,
		messages: copy.evidence,
	})
	const dateTime = new Intl.DateTimeFormat(locale, {
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
		month: 'short',
		timeZone: 'UTC',
		timeZoneName: 'short',
		year: 'numeric',
	})

	return (
		<>
			<a
				className="fixed top-0 left-3 z-200 max-w-[calc(100vw-1.5rem)] -translate-y-full bg-foreground px-2.5 py-2 whitespace-normal text-card focus-visible:top-3 focus-visible:translate-y-0"
				href="#lifecycle"
			>
				{copy.skip}
			</a>
			<WebsiteHeader />
			<main id="top">
				<ProductHero
					caption={copy.hero.caption}
					headingId="microruntimes-hero-heading"
					lead={copy.hero.lead}
					primaryCta={{
						href: '#lifecycle',
						label: copy.hero.primaryCta,
					}}
					secondaryCta={{
						href: '#ownership',
						label: copy.hero.secondaryCta,
					}}
					summary={copy.hero.summary}
					titleEnd={copy.hero.titleEnd}
					titleHighlight={copy.hero.titleHighlight}
					titleStart={copy.hero.titleStart}
				>
					<LifecycleDemo copy={copy.lifecycle} locale={locale} />
				</ProductHero>

				<ProductSection
					id="lifecycle"
					lead={copy.scenarios.lead}
					title={copy.scenarios.title}
				>
					<PillarScenarios
						pillars={lifecyclePhaseIds.map((id) => ({
							id,
							label: copy.scenarios.phases[id],
						}))}
						scenarios={lifecycleScenarioIds.map((id) => {
							const scenario = copy.scenarios.items[id]

							return {
								cells: {
									compose: scenario.compose,
									run: scenario.run,
									inspect: scenario.inspect,
									cleanup: scenario.cleanup,
								},
								context: scenario.context,
								id,
								tabLabel: scenario.tabLabel,
							}
						})}
						selectorLabel={copy.scenarios.selectorLabel}
					/>
					<div className="mt-8 flex flex-wrap items-center justify-between gap-x-8 gap-y-3">
						<p className="m-0 max-w-[640px] text-[13px] leading-[1.6] text-muted-foreground">
							{copy.scenarios.outro}
						</p>
						<a className="smooth-text-link" href="#state">
							{copy.scenarios.outroLink}
							<span aria-hidden="true"> ↗</span>
						</a>
					</div>
				</ProductSection>

				<ProductSection
					id="runs"
					kicker={copy.descriptor.kicker}
					lead={copy.descriptor.lead}
					title={copy.descriptor.title}
				>
					<RuntimeDescriptorCard copy={copy.descriptor} />
				</ProductSection>

				<ProductSection
					id="ownership"
					kicker={copy.ownership.kicker}
					lead={copy.ownership.lead}
					title={copy.ownership.title}
				>
					<div className="mt-11 overflow-x-auto rounded-2xl border border-border bg-card">
						<table className="w-full min-w-[640px] border-collapse text-left text-[13px] leading-[1.55]">
							<thead>
								<tr className="border-b border-border">
									<th
										className="px-5 py-3 font-mono text-[9px] font-normal tracking-[0.08em] text-muted-foreground uppercase"
										scope="col"
									>
										{copy.ownership.concern}
									</th>
									<th
										className="px-5 py-3 font-mono text-[9px] font-normal tracking-[0.08em] text-muted-foreground uppercase"
										scope="col"
									>
										{copy.ownership.consumer}
									</th>
									<th
										className="px-5 py-3 font-mono text-[9px] font-normal tracking-[0.08em] text-primary uppercase"
										scope="col"
									>
										{copy.ownership.platform}
									</th>
								</tr>
							</thead>
							<tbody>
								{ownershipRowIds.map((id) => (
									<tr
										className="border-t border-border first:border-t-0"
										key={id}
									>
										<th
											className="px-5 py-4 align-top text-[14px] font-medium"
											scope="row"
										>
											{copy.ownership.rows[id].concern}
										</th>
										<td className="px-5 py-4 align-top text-muted-foreground">
											{copy.ownership.rows[id].consumer}
										</td>
										<td className="px-5 py-4 align-top">
											{copy.ownership.rows[id].platform}
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>
				</ProductSection>

				<ProductSection
					id="state"
					kicker={copy.state.kicker}
					lead={copy.state.lead}
					title={copy.state.title}
				>
					<ul className="m-0 mt-11 grid list-none gap-3 p-0 sm:grid-cols-2 lg:grid-cols-3">
						{stateCardIds.map((id) => (
							<li
								className="flex flex-col gap-3 rounded-2xl border border-border bg-card p-6"
								key={id}
							>
								<span className="font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase">
									{copy.state.cards[id].eyebrow}
								</span>
								<strong
									className={`text-[22px] leading-tight font-medium tracking-[-0.02em] ${stateTones[id]}`}
								>
									{copy.state.cards[id].title}
								</strong>
								<p className="m-0 text-[13px] leading-[1.55] text-muted-foreground">
									{copy.state.cards[id].description}
								</p>
							</li>
						))}
					</ul>
				</ProductSection>

				<ProductSection
					id="evidence"
					kicker={copy.evidence.kicker}
					lead={copy.evidence.lead}
					title={copy.evidence.title}
				>
					<article className="mt-11 rounded-2xl border border-border bg-card p-6 sm:p-8">
						<div className="flex flex-wrap items-baseline justify-between gap-x-6 gap-y-2 border-b border-border pb-4">
							<strong className="text-lg font-medium tracking-[-0.02em]">
								{copy.evidence.ledgerTitle}
							</strong>
							<span className="font-mono text-[10px] tracking-[0.08em] text-muted-foreground uppercase">
								{copy.evidence.ledgerPeriod}
							</span>
						</div>
						<DetailList
							className="mt-1"
							items={[
								{
									detail: evidence('executionValue', {
										execution: executionLedger.executionId,
										revision: executionLedger.revision,
									}),
									id: 'execution',
									term: copy.evidence.rows.execution,
								},
								{
									detail: dateTime.format(executionLedger.startedAt),
									id: 'started',
									term: copy.evidence.rows.started,
								},
								{
									detail: dateTime.format(executionLedger.stoppedAt),
									id: 'stopped',
									term: copy.evidence.rows.stopped,
								},
								{
									detail: evidence('resources', {
										count: executionLedger.ownedResources,
									}),
									id: 'owned',
									term: copy.evidence.rows.owned,
								},
								{
									detail: evidence('resources', {
										count: executionLedger.removedResources,
									}),
									id: 'removed',
									term: copy.evidence.rows.removed,
									tone: 'success',
								},
								{
									detail: copy.evidence.noLeftovers,
									id: 'leftovers',
									term: copy.evidence.rows.leftovers,
									tone: 'success',
								},
								{
									detail: (
										<code className="font-mono text-[12px] break-all">
											{executionLedger.evidencePath}
										</code>
									),
									id: 'evidence',
									term: copy.evidence.rows.evidence,
								},
							]}
						/>
					</article>
					<StageEvidence
						copy={copy.evidence}
						stageLabels={copy.lifecycle.stages}
					/>
				</ProductSection>

				<ProductSection
					id="flow"
					kicker={copy.flow.kicker}
					lead={copy.flow.lead}
					title={copy.flow.title}
				>
					<AuthoringFlow copy={copy.flow} />
				</ProductSection>

				<ProductFaqSection
					description={copy.faq.description}
					items={microruntimesFaqIds.map((id) => ({
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
