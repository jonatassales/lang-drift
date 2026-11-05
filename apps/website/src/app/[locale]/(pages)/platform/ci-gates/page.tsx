import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { CommandLine } from '@app/[locale]/(pages)/platform/command-line'
import { ProductHeading } from '@app/[locale]/(pages)/platform/product-heading'
import { ProductSection } from '@components/product-page/product-section'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { ciGatesMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { CiGatesIcon } from 'web/ui/icons/platform'
import { InfoCards } from 'web/ui/info-cards'

import {
	ciGatesScreenshots,
	finopsGateIds,
	postureIds,
} from './ci-gates-content.data'
import { ThemedScreenshot } from './themed-screenshot'

export { generateMetadata } from './page-metadata.map'

interface CiGatesPageParams {
	locale: string
}

interface CiGatesPageProps {
	params: Promise<CiGatesPageParams>
}

export default async function CiGatesPage(props: CiGatesPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	setRequestLocale(locale)

	const copy = ciGatesMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="ci-gates" />
			<SkipLink href="#install">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-labelledby="ci-gates-hero-heading"
					className="mx-auto max-w-[1180px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<ProductHeading icon={CiGatesIcon} id="ci-gates-hero-heading">
						{copy.hero.title}
					</ProductHeading>
					<p className="mt-6 mb-0 max-w-190 text-[15px] leading-[1.65] text-muted-foreground">
						{copy.hero.lead}
					</p>
					<p className="mt-4 mb-0 max-w-190 text-xs leading-[1.65] text-muted-foreground">
						{copy.hero.note}
					</p>
				</section>

				<ProductSection
					id="install"
					kicker={copy.install.kicker}
					lead={copy.install.lead}
					title={copy.install.title}
				>
					<div className="mt-8">
						<CommandLine command={copy.install.command} />
					</div>
				</ProductSection>

				<ProductSection
					id="report"
					kicker={copy.report.kicker}
					lead={copy.report.lead}
					title={copy.report.title}
				>
					<div className="mt-10 grid gap-10">
						<ThemedScreenshot
							copy={copy.report.summary}
							screenshot={ciGatesScreenshots.summary}
						/>
						<ThemedScreenshot
							copy={copy.report.dashboard}
							screenshot={ciGatesScreenshots.dashboard}
						/>
					</div>
				</ProductSection>

				<ProductSection
					id="posture"
					kicker={copy.posture.kicker}
					lead={copy.posture.lead}
					title={copy.posture.title}
				>
					<InfoCards
						className="mt-11"
						columns={4}
						items={postureIds.map((id) => ({
							...copy.posture.items[id],
							id,
						}))}
					/>
				</ProductSection>

				<ProductSection
					id="gates"
					kicker={copy.gates.kicker}
					lead={copy.gates.lead}
					title={copy.gates.title}
				>
					<InfoCards
						className="mt-11"
						items={finopsGateIds.map((id) => ({
							...copy.gates.items[id],
							id,
						}))}
					/>
				</ProductSection>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
