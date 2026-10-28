import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'
import { RotatingText } from '@langdrift/web/ui/rotating-text'

import { platformProduct } from '@components/website-header/platform-products'
import { platformMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero is the moving headline and its two calls to action, at the top
// left, over the shader.
export function PlatformHero(props: PlatformHeroProps) {
	const { locale } = props

	const copy = platformMessages[locale].hero

	return (
		<section
			aria-labelledby="platform-hero-heading"
			className="relative isolate flex min-h-[calc(100dvh-5rem)] flex-col justify-start overflow-hidden bg-background"
		>
			<HeroShader />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="relative mx-auto flex w-full max-w-7xl flex-col items-start px-4.5 pt-16 pb-20 text-left sm:px-9 sm:pt-20">
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
		</section>
	)
}
