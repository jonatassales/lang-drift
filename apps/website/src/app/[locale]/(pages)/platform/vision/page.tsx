import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'
import { driftSign } from '@langdrift/web/utilities'
import { AnimatedAvatarGroup } from '@langdrift/web/vendors/smoothui/animated-avatar-group'
import { ChromaText } from '@langdrift/web/vendors/smoothui/chroma-text'
import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import {
	getFormatter,
	getTranslations,
	setRequestLocale,
} from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { BookDemoBand } from '@components/book-a-demo/band'
import { ExecutiveReviewSection } from '@components/executive-review-section'
import { FaqSection } from '@components/faq-section'
import { HeroChart } from '@components/hero-chart'
import { HeroSinapsiGraph } from '@components/hero-sinapsi-graph'
import { IntegrationsSection } from '@components/integrations-section'
import { RoiSection } from '@components/roi-section'
import { SelectedMovement } from '@components/selected-movement'
import { SkipLink } from '@components/skip-link'
import { VisionLoopEvidenceSection } from '@components/vision-loop-evidence-section'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { getHeroPoints } from '@data/hero-demo-data'
import { routing } from '@i18n/routing'

interface VisionPageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: VisionPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('vision', locale, await socialImages(parent))
}

export default async function VisionPage(props: VisionPageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const t = await getTranslations('home')
	const demo = await getTranslations('demo')
	const format = await getFormatter()
	const attribution = getHeroPoints(demo, locale)
		.flatMap((point) => (point.event ? [point.event] : []))
		.filter((event) => event.classification !== 'baseline')

	return (
		<>
			<SkipLink href="#why">{t('skip')}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section aria-labelledby="hero-heading" className="website-hero">
					<div className="website-hero-intro">
						<h1 id="hero-heading">
							{t.rich('hero.title', {
								movement: (chunks) => <ChromaText>{chunks}</ChromaText>,
							})}
						</h1>
						<p>{t('hero.description')}</p>
						<div className="website-hero-actions">
							<a className="smooth-primary-button" href="#why">
								{t('hero.seeMovement')}
								<ArrowRightIcon />
							</a>
							<a className="smooth-text-link" href="#attribution">
								{t('hero.howItWorks')}
							</a>
						</div>
					</div>
					<div
						className="website-hero-demo"
						aria-describedby="hero-demo-caption hero-demo-summary"
					>
						<div className="website-demo-grid">
							<HeroChart />
							<HeroSinapsiGraph />
						</div>
						<p className="website-demo-caption" id="hero-demo-caption">
							{t('hero.caption')}
						</p>
						<p className="website-sr-only" id="hero-demo-summary">
							{t('hero.summary')}
						</p>
					</div>
				</section>

				<VisionLoopEvidenceSection locale={locale} />

				<section className="why-section" id="movement-details">
					<div className="why-title">
						<span className="section-kicker">{t('why.kicker')}</span>
						<h2>{t('why.title')}</h2>
						<p>{t('why.description')}</p>
					</div>
					<SelectedMovement />
				</section>

				<section className="attribution-section" id="attribution">
					<div className="features-heading">
						<span className="section-kicker">{t('attribution.kicker')}</span>
						<h2>{t('attribution.title')}</h2>
						<p>{t('attribution.description')}</p>
					</div>
					<div className="attribution-grid">
						{attribution.map((event) => (
							<article key={event.id}>
								<div className="attribution-meta">
									<span>{event.date}</span>
									<span>{event.productArea}</span>
								</div>
								<h3>{event.title}</h3>
								<div className="attribution-bottom">
									<AnimatedAvatarGroup
										people={event.actors.map((actor) => ({
											initials: actor.initials,
											name: actor.name,
											role: actor.team,
											src: actor.src,
										}))}
									/>
									<strong data-drift-sign={driftSign(event.delta)}>
										{format.number(event.delta, { signDisplay: 'exceptZero' })}
									</strong>
									<span>{demo(`classifications.${event.classification}`)}</span>
								</div>
							</article>
						))}
					</div>
				</section>

				<section className="intent-section" id="intent">
					<div>
						<span className="section-kicker">{t('intent.kicker')}</span>
						<h2>{t('intent.title')}</h2>
					</div>
					<div className="intent-grid">
						<article>
							<span>{demo('classifications.intentional')}</span>
							<strong>{t('intent.points', { count: 14 })}</strong>
							<p>{t('intent.intentional')}</p>
						</article>
						<article>
							<span>{demo('classifications.unexplained')}</span>
							<strong>{t('intent.points', { count: 4 })}</strong>
							<p>{t('intent.unexplained')}</p>
						</article>
						<article>
							<span>{demo('classifications.review')}</span>
							<strong>{t('intent.events', { count: 1 })}</strong>
							<p>{t('intent.review')}</p>
						</article>
					</div>
				</section>

				<section className="report-section" id="report">
					<div className="report-copy">
						<span className="section-kicker">{t('report.kicker')}</span>
						<h2>{t('report.title')}</h2>
					</div>
					<article className="report-preview">
						<div className="report-preview-head">
							<span>{t('report.weekly')}</span>
							<strong>
								{t('report.vision', {
									from: format.number(81),
									to: format.number(76),
								})}
							</strong>
						</div>
						<dl>
							<div>
								<dt>{t('report.largest')}</dt>
								<dd data-drift-sign="negative">
									{t('report.movement', { delta: format.number(-3) })}
								</dd>
							</div>
							<div>
								<dt>{t('report.why')}</dt>
								<dd>{t('report.rationale')}</dd>
							</div>
							<div>
								<dt>{t('report.who')}</dt>
								<dd>{t('report.actor')}</dd>
							</div>
							<div>
								<dt>{t('report.status')}</dt>
								<dd>{demo('classifications.unexplained')}</dd>
							</div>
							<div>
								<dt>{t('report.attention')}</dt>
								<dd>{t('report.action')}</dd>
							</div>
						</dl>
					</article>
				</section>

				<IntegrationsSection locale={locale} />
				<ExecutiveReviewSection locale={locale} />
				<RoiSection locale={locale} />
				<FaqSection locale={locale} />
				<BookDemoBand locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
