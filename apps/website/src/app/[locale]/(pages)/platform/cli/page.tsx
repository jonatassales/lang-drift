import Image from 'next/image'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { CommandLine } from '@app/[locale]/(pages)/platform/command-line'
import { ProductIconTile } from '@app/[locale]/(pages)/platform/product-icon-tile'
import { ProductSection } from '@components/product-page/product-section'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { cliMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'
import { CliIcon } from 'web/ui/icons/platform'
import { InfoCards } from 'web/ui/info-cards'

import {
	cliFlowStepIds,
	cliGovernIds,
	cliLinks,
	cliPlanImage,
} from './cli-content.data'

export { generateMetadata } from './page-metadata.map'

const externalLink =
	'inline-flex min-h-11 items-center text-sm font-medium text-foreground underline decoration-border underline-offset-4 hover:decoration-primary focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary'

interface CliPageParams {
	locale: string
}

interface CliPageProps {
	params: Promise<CliPageParams>
}

export default async function CliPage(props: CliPageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	setRequestLocale(locale)

	const copy = cliMessages[locale]

	return (
		<>
			<WebsiteStructuredData locale={locale} page="cli" />
			<SkipLink href="#install">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<section
					aria-labelledby="cli-hero-heading"
					className="mx-auto max-w-[1180px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
				>
					<h1
						className="m-0 flex items-center gap-[0.35em] text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance"
						id="cli-hero-heading"
					>
						<ProductIconTile className="size-[0.74em] rounded-[0.2em]">
							<CliIcon className="size-[0.46em]" />
						</ProductIconTile>
						{copy.hero.title}
					</h1>
					<p className="mt-6 mb-0 max-w-190 text-[15px] leading-[1.65] text-muted-foreground">
						{copy.hero.lead}
					</p>
					<div className="mt-8">
						<CommandLine command={copy.hero.command} />
					</div>
					<p className="mt-4 mb-0 max-w-190 text-xs leading-[1.65] text-muted-foreground">
						{copy.hero.requirements}
					</p>
				</section>

				<ProductSection
					id="install"
					kicker={copy.install.kicker}
					lead={copy.install.lead}
					title={copy.install.title}
				>
					<div className="mt-8 grid gap-3">
						<CommandLine command={copy.install.commands.npx} />
						<CommandLine command={copy.install.commands.pnpm} />
					</div>
					<figure className="m-0 mt-10">
						<Image
							alt={copy.install.image.alt}
							className="h-auto w-full rounded-2xl border border-border"
							height={cliPlanImage.height}
							sizes="(min-width: 1180px) 1108px, 100vw"
							src={cliPlanImage.src}
							width={cliPlanImage.width}
						/>
						<figcaption className="mt-3 text-xs leading-relaxed text-muted-foreground">
							{copy.install.image.caption}
						</figcaption>
					</figure>
				</ProductSection>

				<ProductSection
					id="flow"
					kicker={copy.flow.kicker}
					lead={copy.flow.lead}
					title={copy.flow.title}
				>
					<InfoCards
						className="mt-11"
						columns={4}
						items={cliFlowStepIds.map((id) => ({
							...copy.flow.steps[id],
							id,
						}))}
						ordered
					/>
				</ProductSection>

				<ProductSection
					id="governs"
					kicker={copy.governs.kicker}
					lead={copy.governs.lead}
					title={copy.governs.title}
				>
					<InfoCards
						className="mt-11"
						columns={4}
						items={cliGovernIds.map((id) => ({
							...copy.governs.items[id],
							id,
						}))}
					/>
					<ul className="m-0 mt-10 flex list-none flex-wrap gap-x-8 gap-y-2 p-0">
						<li>
							<a className={externalLink} href={cliLinks.npm}>
								{copy.links.npm}
							</a>
						</li>
						<li>
							<a className={externalLink} href={cliLinks.github}>
								{copy.links.github}
							</a>
						</li>
						<li>
							<a className={externalLink} href={cliLinks.license}>
								{copy.links.license}
							</a>
						</li>
					</ul>
				</ProductSection>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
