import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { ProductHeading } from '@app/[locale]/(pages)/platform/product-heading'
import { ProductTagline } from '@app/[locale]/(pages)/platform/product-tagline'
import { ProductBadge } from '@components/product-badge'
import { ProductSection } from '@components/product-page/product-section'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { headerMessages, safeSeatMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { DetailList } from 'web/ui/detail-list'
import { SafeSeatIcon } from 'web/ui/icons/platform'
import { InfoCards } from 'web/ui/info-cards'

import { EvidenceCards } from './evidence-cards'
import {
	categoryIds,
	flowStepIds,
	limitIds,
	principleIds,
	whyIds,
} from './safe-seat-content.data'
import { SourceCitations } from './source-citations'

export { generateMetadata } from './page-metadata.map'

interface SafeSeatPageParams {
	locale: string
}

interface SafeSeatPageProps {
	params: Promise<SafeSeatPageParams>
}

export default async function SafeSeatPage(props: SafeSeatPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	setRequestLocale(locale)

	const copy = safeSeatMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="safe-seat" />
			<SkipLink href="#why">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-labelledby="safe-seat-hero-heading"
					className="mx-auto max-w-[1180px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<div className="flex flex-wrap items-center gap-3">
						<ProductHeading icon={SafeSeatIcon} id="safe-seat-hero-heading">
							{headerMessages[locale].safeSeat}
						</ProductHeading>
						<ProductBadge>{copy.hero.comingSoon}</ProductBadge>
					</div>
					<ProductTagline
						end={copy.hero.taglineEnd}
						start={copy.hero.taglineStart}
					/>
					<p className="mt-6 mb-0 max-w-190 text-[15px] leading-[1.65] text-muted-foreground">
						{copy.hero.lead}
					</p>
					<p className="mt-4 mb-0 max-w-190 text-xs leading-[1.65] text-muted-foreground">
						{copy.hero.note}
					</p>
				</section>

				<ProductSection
					id="why"
					kicker={copy.why.kicker}
					lead={copy.why.lead}
					title={copy.why.title}
				>
					<InfoCards
						className="mt-11"
						items={whyIds.map((id) => ({
							...copy.why.items[id],
							id,
						}))}
					/>
				</ProductSection>

				<ProductSection
					id="evidence"
					kicker={copy.evidence.kicker}
					lead={copy.evidence.lead}
					title={copy.evidence.title}
				>
					<EvidenceCards copy={copy.evidence} />
				</ProductSection>

				<ProductSection
					id="reserved"
					kicker={copy.reserved.kicker}
					lead={copy.reserved.lead}
					title={copy.reserved.title}
				>
					<InfoCards
						className="mt-11"
						items={categoryIds.map((id) => ({
							...copy.reserved.categories[id],
							id,
						}))}
					/>
					<p className="mt-6 mb-0 max-w-190 text-[13px] leading-[1.6] text-muted-foreground">
						{copy.reserved.never}
					</p>
				</ProductSection>

				<ProductSection
					id="flow"
					kicker={copy.flow.kicker}
					lead={copy.flow.lead}
					title={copy.flow.title}
				>
					<InfoCards
						className="mt-11"
						items={flowStepIds.map((id) => ({
							...copy.flow.steps[id],
							id,
						}))}
						ordered
					/>
				</ProductSection>

				<ProductSection
					id="principles"
					kicker={copy.principles.kicker}
					lead={copy.principles.lead}
					title={copy.principles.title}
				>
					<DetailList
						className="mt-11 max-w-240"
						items={principleIds.map((id) => ({
							detail: copy.principles.items[id].detail,
							id,
							term: copy.principles.items[id].term,
						}))}
					/>
				</ProductSection>

				<ProductSection
					id="limits"
					kicker={copy.limits.kicker}
					lead={copy.limits.lead}
					title={copy.limits.title}
				>
					<DetailList
						className="mt-11 max-w-240"
						items={limitIds.map((id) => ({
							detail: copy.limits.items[id].detail,
							id,
							term: copy.limits.items[id].term,
						}))}
					/>
				</ProductSection>

				<ProductSection
					id="sources"
					kicker={copy.sources.kicker}
					lead={copy.sources.lead}
					title={copy.sources.title}
				>
					<SourceCitations copy={copy.sources} />
				</ProductSection>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
