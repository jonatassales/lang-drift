import { HeroShader } from '@langdrift/web/ui/hero-shader'

import { platformMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

export function PlatformHero(props: PlatformHeroProps) {
	const { locale } = props

	const copy = platformMessages[locale].hero

	return (
		<section
			aria-labelledby="platform-hero-heading"
			className="relative isolate flex min-h-[calc(100dvh-5rem)] flex-col justify-center overflow-hidden bg-background"
		>
			<HeroShader />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="relative mx-auto flex w-full max-w-3xl flex-col items-center px-4.5 py-24 text-center sm:px-7">
				<h1
					className="text-5xl leading-[1.05] font-medium tracking-tight text-balance text-foreground sm:text-7xl"
					id="platform-hero-heading"
				>
					{copy.title}
				</h1>
				<p className="mt-6 max-w-xl text-base leading-relaxed text-pretty text-muted-foreground sm:text-lg">
					{copy.lead}
				</p>
				<div className="mt-8 flex flex-wrap items-center justify-center gap-4 sm:gap-6">
					<Link
						className="smooth-primary-button"
						href="/platform/vision"
						locale={locale}
					>
						{copy.primaryCta}
						<svg
							aria-hidden="true"
							fill="none"
							height="16"
							viewBox="0 0 24 24"
							width="16"
						>
							<path
								d="M4 12h16m-6-6 6 6-6 6"
								stroke="currentColor"
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="1.5"
							/>
						</svg>
					</Link>
					<Link
						className="smooth-text-link"
						href="/platform/vision#early-access"
						locale={locale}
					>
						{copy.secondaryCta}
					</Link>
				</div>
				<ul className="mt-14 grid w-full gap-3 text-sm text-muted-foreground sm:grid-cols-3">
					<li>{copy.claims.agents}</li>
					<li>{copy.claims.evidence}</li>
					<li>{copy.claims.budget}</li>
				</ul>
			</div>
		</section>
	)
}
