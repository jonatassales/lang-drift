import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { MagneticField } from '@langdrift/web/vendors/smoothui/magnetic-field'

import { platformMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero's centerpiece is a black hole sun: a disc in the page background
// with a thin bright rim, and a corona in the brand purple breathing behind it.
// The disc holds still; the smoke behind it leans and turns toward the
// pointer and brightens while the pointer is over the sun.
const corona =
	'absolute -inset-20 rounded-full opacity-80 transition-opacity duration-700 group-hover:opacity-100 bg-[radial-gradient(circle_at_50%_42%,color-mix(in_srgb,var(--primary)_62%,transparent)_0%,color-mix(in_srgb,var(--primary)_26%,transparent)_38%,transparent_70%)] blur-3xl motion-safe:animate-pulse dark:bg-[radial-gradient(circle_at_50%_42%,color-mix(in_srgb,var(--foreground)_58%,transparent)_0%,color-mix(in_srgb,var(--primary)_46%,transparent)_30%,transparent_70%)]'
const rim =
	'absolute -inset-2 rounded-full bg-[radial-gradient(circle_closest-side,transparent_92%,color-mix(in_srgb,var(--foreground)_85%,transparent)_96%,transparent_100%)] blur-[1.5px]'
const disc =
	'absolute inset-0 rounded-full bg-background shadow-[inset_0_0_80px_color-mix(in_srgb,var(--primary)_30%,transparent)]'

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
			<MagneticField
				className="mx-auto grid w-full max-w-7xl items-center gap-12 px-4.5 py-20 sm:px-9 lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] lg:gap-10"
				radius={420}
				rotate
				strength={0.45}
			>
				<div className="order-2 flex flex-col items-start text-left lg:order-1">
					<h1
						className="max-w-sm text-5xl leading-none font-normal tracking-[-0.055em] text-balance text-foreground sm:text-6xl xl:text-[64px]"
						id="platform-hero-heading"
					>
						{copy.title}
					</h1>
					<div className="mt-8 flex flex-wrap items-center gap-4 sm:gap-6">
						<Link
							className="inline-flex min-h-11 items-center gap-2 rounded-full bg-foreground px-5 text-sm font-semibold whitespace-nowrap text-background transition-opacity hover:opacity-85"
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
				</div>
				<div className="order-1 flex justify-center lg:order-2">
					<div
						aria-hidden="true"
						className="group relative size-52 sm:size-64 lg:size-72"
					>
						<div
							className={corona}
							data-magnetic
							style={{ animationDuration: '7s' }}
						/>
						<div className={rim} />
						<div className={disc} />
					</div>
				</div>
				<div aria-hidden="true" className="hidden lg:order-3 lg:block" />
			</MagneticField>
		</section>
	)
}
