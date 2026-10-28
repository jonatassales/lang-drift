import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { RotatingText } from '@langdrift/web/ui/rotating-text'

import { SessionCards } from '@components/microruntimes/session-cards'
import { SessionTerminal } from '@components/microruntimes/session-terminal'
import { microruntimesMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero is the Microruntimes planning session over the shader: the moving
// headline at the top left with the terminal under it, and the ticket, plan
// and workflow cards the terminal plays into on the right. From `lg` the
// terminal's height is what the first screen leaves under the header, the
// headline and the paddings, so the whole window shows without scrolling;
// below `lg`, headline, terminal and cards stack in that order.
export function PlatformHero(props: PlatformHeroProps) {
	const { locale } = props

	const copy = platformMessages[locale].hero
	const microruntimes = microruntimesMessages[locale]
	const session = microruntimes.session

	return (
		<section
			aria-describedby="platform-hero-summary"
			aria-labelledby="platform-hero-heading"
			className="relative isolate flex min-h-[calc(100dvh-5rem)] flex-col justify-start overflow-hidden bg-background"
		>
			<HeroShader />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="relative mx-auto grid w-full max-w-7xl gap-6 px-4.5 pt-16 pb-20 sm:px-9 sm:pt-20 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-x-10">
				<div className="flex min-w-0 flex-col items-start text-left">
					<h1
						className="text-[min(2.5rem,10vw)] leading-none font-normal tracking-[-0.055em] text-foreground sm:text-5xl lg:text-[clamp(2.5rem,calc((100vw-25rem)/16),3.25rem)]"
						id="platform-hero-heading"
					>
						<span className="block">{copy.titleStart}</span>
						<RotatingText
							className="block whitespace-nowrap"
							suffix={copy.titleEnd}
							words={copy.words}
						/>
					</h1>
					<div className="mt-10 w-full min-w-0">
						<SessionTerminal
							className="h-[clamp(220px,34vh,300px)] lg:h-[max(16rem,calc(100dvh-21rem))]"
							control={session.card.control}
							copy={session.terminal}
						/>
					</div>
				</div>
				<div className="min-w-0 lg:self-start">
					<SessionCards copy={session.card} locale={locale} />
				</div>
			</div>
			<p className="sr-only" id="platform-hero-summary">
				{microruntimes.hero.summary}
			</p>
		</section>
	)
}
