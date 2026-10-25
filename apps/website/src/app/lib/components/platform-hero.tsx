import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { GravityStars } from '@langdrift/web/vendors/smoothui/gravity-stars'
import { MagneticField } from '@langdrift/web/vendors/smoothui/magnetic-field'

import { platformProducts } from '@components/website-header/platform-products'
import { headerMessages, platformMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero's centerpiece is a black hole sun: a dark disc in both color
// schemes, with a thin rim, and a corona in the brand purple breathing behind it.
// The disc holds still; the smoke behind it leans and turns toward the
// pointer and brightens while the pointer is over the sun. Inside the disc,
// and only there, a field of stars drifts and falls toward the pointer.
const corona =
	'absolute -inset-20 rounded-full opacity-80 transition-opacity duration-700 group-hover:opacity-100 bg-[radial-gradient(circle_at_50%_42%,color-mix(in_srgb,var(--primary)_62%,transparent)_0%,color-mix(in_srgb,var(--primary)_26%,transparent)_38%,transparent_70%)] blur-3xl motion-safe:animate-pulse dark:bg-[radial-gradient(circle_at_50%_42%,color-mix(in_srgb,var(--foreground)_58%,transparent)_0%,color-mix(in_srgb,var(--primary)_46%,transparent)_30%,transparent_70%)]'
const rim =
	'absolute -inset-2 rounded-full bg-[radial-gradient(circle_closest-side,transparent_92%,color-mix(in_srgb,var(--primary)_85%,transparent)_96%,transparent_100%)] blur-[1.5px] dark:bg-[radial-gradient(circle_closest-side,transparent_92%,color-mix(in_srgb,var(--foreground)_85%,transparent)_96%,transparent_100%)]'
// The stars take the disc's text color, so they stay light on the dark disc
// whichever scheme is active.
const disc =
	'absolute inset-0 rounded-full bg-foreground text-background dark:bg-background dark:text-foreground'
// The products wait on the right, dimmed but readable; each lights up as the
// pointer reaches it, with a dot in the brand purple and a soft glow.
const productItem =
	"flex items-center gap-3 text-lg font-normal tracking-tight text-foreground/65 transition-all duration-500 ease-out before:size-1.5 before:rounded-full before:bg-primary before:opacity-0 before:transition-opacity before:duration-500 before:content-[''] hover:text-foreground hover:[text-shadow:0_0_28px_color-mix(in_srgb,var(--primary)_80%,transparent)] hover:before:opacity-100 focus-visible:text-foreground focus-visible:before:opacity-100 lg:text-xl"
const discShade =
	'pointer-events-none absolute inset-0 rounded-full shadow-[inset_0_0_80px_color-mix(in_srgb,var(--primary)_30%,transparent)]'

export function PlatformHero(props: PlatformHeroProps) {
	const { locale } = props

	const copy = platformMessages[locale].hero
	const names = headerMessages[locale]

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
						className="group relative size-44 sm:size-54 lg:size-60"
					>
						<div
							className={corona}
							data-magnetic
							style={{ animationDuration: '7s' }}
						/>
						<div className={rim} />
						<GravityStars
							className={disc}
							color="currentColor"
							connectDistance={56}
							count={320}
							glow={5}
							gravity={1.2}
							speed={0.6}
							starSize={1.1}
							tint={0.3}
						>
							<div aria-hidden="true" className={discShade} />
						</GravityStars>
					</div>
				</div>
				<ul
					aria-label={names.platformMenu}
					className="order-3 flex flex-col items-start gap-4 lg:order-3 lg:mr-[10vw] lg:justify-self-end"
				>
					{platformProducts.map((product) => {
						const Icon = product.icon
						const content = (
							<>
								<Icon className="size-4 shrink-0" />
								{names[product.label]}
							</>
						)

						return (
							<li key={product.label}>
								{'href' in product ? (
									<Link
										className={productItem}
										href={product.href}
										locale={locale}
									>
										{content}
									</Link>
								) : (
									<span className={productItem}>{content}</span>
								)}
							</li>
						)
					})}
				</ul>
			</MagneticField>
		</section>
	)
}
