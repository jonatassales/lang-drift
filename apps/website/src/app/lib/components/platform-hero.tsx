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
// headline at the top left and, under it, one frame that holds the terminal
// and the ticket, plan and workflow cards it plays into, side by side, as a
// single window. From `lg` the frame takes the height the first screen leaves
// under the header and the headline, and the cards scroll inside it; below
// `lg`, the terminal keeps its own height and the cards stack under it.
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
			<div className="relative mx-auto flex w-full max-w-7xl flex-col items-start px-4.5 pt-16 pb-20 text-left sm:px-9 sm:pt-20">
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
				<div className="mt-10 grid w-full min-w-0 gap-2 rounded-2xl border border-border bg-background/70 p-2 shadow-2xl backdrop-blur-xl lg:h-[max(24rem,calc(100dvh-20rem))] lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
					<div className="min-h-0 min-w-0">
						<SessionTerminal
							className="h-[clamp(220px,34vh,300px)] shadow-none"
							control={session.card.control}
							copy={session.terminal}
						/>
					</div>
					<div className="min-h-0 min-w-0 lg:overflow-y-auto">
						<SessionCards copy={session.card} locale={locale} />
					</div>
				</div>
			</div>
			<p className="sr-only" id="platform-hero-summary">
				{microruntimes.hero.summary}
			</p>
		</section>
	)
}
