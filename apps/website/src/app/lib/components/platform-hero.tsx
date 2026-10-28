import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'
import { RotatingText } from '@langdrift/web/ui/rotating-text'
import { GravityStars } from '@langdrift/web/vendors/smoothui/gravity-stars'
import { MagneticField } from '@langdrift/web/vendors/smoothui/magnetic-field'
import { ShortSlideDown } from '@langdrift/web/vendors/smoothui/short-slide-down'

import {
	platformProduct,
	platformProducts,
} from '@components/website-header/platform-products'
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
// The products drop in on the right one after another, each name with one
// line about it beneath. They stay solid foreground, white in the dark scheme
// and black in the light one, and do not react to the pointer.
const productItem =
	'grid grid-cols-[auto_1fr] items-center gap-x-3 gap-y-1 text-foreground'
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
						className="text-[min(3rem,12vw)] leading-none font-normal tracking-[-0.055em] text-foreground sm:text-6xl lg:text-[clamp(3rem,calc((100vw-25rem)/13.2),4rem)]"
						id="platform-hero-heading"
					>
						<span className="block">{copy.titleStart}</span>
						<RotatingText
							className="block whitespace-nowrap"
							suffix={copy.titleEnd}
							words={copy.words}
						/>
					</h1>
					<div className="mt-8 flex flex-wrap items-center gap-4 sm:gap-6">
						<Link
							className="inline-flex min-h-11 items-center gap-2 rounded-full bg-foreground px-5 text-sm font-semibold whitespace-nowrap text-background transition-opacity hover:opacity-85"
							href={platformProduct('finops').href}
							locale={locale}
						>
							{copy.primaryCta}
							<ArrowRightIcon />
						</Link>
						<Link
							className="smooth-text-link"
							href="/book-a-demo"
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
				<ShortSlideDown
					className="order-3 flex flex-col items-start gap-4 lg:order-3 lg:justify-self-end xl:mr-[4vw]"
					items={platformProducts.map((product) => {
						const Icon = product.icon
						const content = (
							<>
								<Icon className="size-4 shrink-0" />
								<span className="text-lg font-normal tracking-tight lg:text-xl">
									{names[product.label]}
								</span>
								<span className="col-start-2 text-sm whitespace-nowrap">
									{copy.products[product.label]}
								</span>
							</>
						)

						return {
							content: (
								<Link
									className={productItem}
									href={product.href}
									locale={locale}
								>
									{content}
								</Link>
							),
							id: product.label,
						}
					})}
					label={names.productList}
				/>
			</MagneticField>
		</section>
	)
}
