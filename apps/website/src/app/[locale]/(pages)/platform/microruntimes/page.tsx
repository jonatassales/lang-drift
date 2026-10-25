import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { MacTerminal } from '@langdrift/web/vendors/smoothui/mac-terminal'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { createTranslator, hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { EarlyAccessSection } from '@components/early-access/section'
import { AuthoringFlow } from '@components/microruntimes/authoring-flow'
import { LifecycleSteps } from '@components/microruntimes/lifecycle-steps'
import { RuntimeDescriptorCard } from '@components/microruntimes/runtime-descriptor-card'
import { StageEvidence } from '@components/microruntimes/stage-evidence'
import { PillarScenarios } from '@components/product-page/pillar-scenarios'
import { ProductFaqSection } from '@components/product-page/product-faq-section'
import { ProductHero } from '@components/product-page/product-hero'
import { ProductSection } from '@components/product-page/product-section'
import { RecordCard } from '@components/product-page/record-card'
import { SectionOutro } from '@components/product-page/section-outro'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import {
	executionLedger,
	lifecyclePhaseIds,
	lifecycleScenarioIds,
	microruntimesFaqIds,
	ownershipRowIds,
	type StateCardId,
	stateCardIds,
	terminalSession,
} from '@data/microruntimes-content'
import { microruntimesMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface MicroruntimesPageProps {
	params: Promise<{ locale: string }>
}

const stateTones: Record<StateCardId, string> = {
	requested: 'text-foreground',
	observed: 'text-status-success',
	unknown: 'text-status-warning',
}

const headerCell = 'px-5 py-3 font-normal'

export async function generateMetadata(
	{ params }: MicroruntimesPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const images = await socialImages(parent)

	return pageMetadata('microruntimes', locale, images)
}

export default async function MicroruntimesPage(props: MicroruntimesPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = microruntimesMessages[locale]
	const formatEvidence = createTranslator({
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
			<SkipLink href="#lifecycle">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<ProductHero
					copy={copy.hero}
					headingId="microruntimes-hero-heading"
					primaryFragment="lifecycle"
					secondaryFragment="ownership"
				>
					<div className="grid gap-6 lg:grid-cols-[fit-content(50%)_1fr]">
						<MacTerminal
							label={copy.lifecycle.terminal.label}
							lines={terminalSession.map((line) => ({
								...line,
								text: copy.lifecycle.terminal.lines[line.id],
							}))}
							maxHeight={320}
							prompt="~/shop %"
							title={copy.lifecycle.terminal.title}
						/>
					</div>
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
						scenarios={lifecycleScenarioIds.map((id) => ({
							cells: copy.scenarios.items[id],
							context: copy.scenarios.items[id].context,
							id,
							tabLabel: copy.scenarios.items[id].tabLabel,
						}))}
						selectorLabel={copy.scenarios.selectorLabel}
					/>
					<SectionOutro
						fragment="state"
						link={copy.scenarios.outroLink}
						text={copy.scenarios.outro}
					/>
				</ProductSection>

				<ProductSection
					id="runs"
					kicker={copy.descriptor.kicker}
					lead={copy.descriptor.lead}
					title={copy.descriptor.title}
				>
					<RuntimeDescriptorCard copy={copy.descriptor} locale={locale} />
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
									<th className={headerCell} scope="col">
										<Eyebrow>{copy.ownership.concern}</Eyebrow>
									</th>
									<th className={headerCell} scope="col">
										<Eyebrow>{copy.ownership.consumer}</Eyebrow>
									</th>
									<th className={headerCell} scope="col">
										<Eyebrow className="text-primary">
											{copy.ownership.platform}
										</Eyebrow>
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
								<Eyebrow>{copy.state.cards[id].eyebrow}</Eyebrow>
								<strong
									className={cn(
										'text-[22px] leading-tight font-medium tracking-[-0.02em]',
										stateTones[id],
									)}
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
					<RecordCard
						items={[
							{
								detail: formatEvidence('executionValue', {
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
								detail: formatEvidence('resources', {
									count: executionLedger.ownedResources,
								}),
								id: 'owned',
								term: copy.evidence.rows.owned,
							},
							{
								detail: formatEvidence('resources', {
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
						title={copy.evidence.ledgerTitle}
					/>
					<LifecycleSteps copy={copy.lifecycle} />
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
