import { Provider } from 'jotai'
import { notFound } from 'next/navigation'
import { createTranslator, hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { microruntimesMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { SkipLink } from '@layout/skip-link'
import { WebsiteFooter } from '@layout/website-footer'
import { WebsiteHeader } from '@layout/website-header/website-header'
import { WebsiteStructuredData } from '@layout/website-structured-data'
import {
	executionLedger,
	lifecyclePhaseIds,
	lifecycleScenarioIds,
	ownershipRowIds,
	stateCardIds,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import type { StateCardId } from '@routes/platform/microruntimes/content/microruntimes-content.types'
import { AuthoringFlow } from '@routes/platform/microruntimes/sections/authoring-flow'
import { LifecycleSteps } from '@routes/platform/microruntimes/sections/lifecycle-steps'
import { RuntimeDescriptorCard } from '@routes/platform/microruntimes/sections/runtime-descriptor-card'
import { SessionCards } from '@routes/platform/microruntimes/sections/session-cards'
import { SessionTerminal } from '@routes/platform/microruntimes/sections/session-terminal'
import { StageEvidence } from '@routes/platform/microruntimes/sections/stage-evidence'
import { ProductHeading } from '@routes/platform/product-heading'
import { PillarScenarios } from '@routes/platform/product-page/pillar-scenarios'
import { ProductSection } from '@routes/platform/product-page/product-section'
import { RecordCard } from '@routes/platform/product-page/record-card'
import { SectionOutro } from '@routes/platform/product-page/section-outro'
import { Eyebrow } from 'web/ui/eyebrow'
import { MicroruntimesIcon } from 'web/ui/icons/platform'
import { cn } from 'web/utilities'

export { generateMetadata } from './page-metadata.map'

const stateTones: Record<StateCardId, string> = {
	requested: 'text-foreground',
	observed: 'text-status-success',
	unknown: 'text-status-warning',
}

const headerCell = 'px-5 py-3 font-normal'

interface LocaleRouteParams {
	locale: string
}

interface MicroruntimesPageProps {
	params: Promise<LocaleRouteParams>
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
			<WebsiteStructuredData locale={locale} page="microruntimes" />
			<SkipLink href="#lifecycle">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-labelledby="microruntimes-hero-heading"
					className="mx-auto max-w-[1280px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<ProductHeading
						icon={MicroruntimesIcon}
						id="microruntimes-hero-heading"
					>
						{copy.hero.title}
					</ProductHeading>
					<Provider>
						<div
							aria-describedby="microruntimes-hero-summary"
							className="mt-8 grid gap-6 lg:grid-cols-[fit-content(50%)_1fr]"
						>
							<SessionTerminal
								className="lg:h-full"
								control={copy.session.card.control}
								copy={copy.session.terminal}
							/>
							<SessionCards copy={copy.session.card} locale={locale} />
						</div>
					</Provider>
					<p className="sr-only" id="microruntimes-hero-summary">
						{copy.hero.summary}
					</p>
				</section>

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
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
