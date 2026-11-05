import { microruntimesMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import { SessionCards } from '@routes/platform/microruntimes/sections/session-cards'
import { SessionTerminal } from '@routes/platform/microruntimes/sections/session-terminal'
import { ChromaText } from 'web/ui/chroma-text'
import { ElevatedFrame } from 'web/ui/elevated-frame'
import { HeroShader } from 'web/ui/hero-shader'
import { RotatingText } from 'web/ui/rotating-text'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero is the moving headline alone, centered on the first screen over
// the shader in its muted tone, with the rainbow flowing through the word
// "Platform" as the page's one touch of color. Under the headline, one
// elevated frame holds the Microruntimes terminal and the ticket, plan and
// workflow cards it plays into, side by side; its top 21rem show at the bottom
// of the first screen. The hero starts under the transparent header, so the
// two read as one background. From `lg` the frame is at least 40rem tall and
// grows with the cards, so they never scroll; the terminal stretches to the
// same height and scrolls its transcript inside; below `lg`, the cards stack
// under the terminal.
export function PlatformHero(props: PlatformHeroProps) {
	const { locale } = props

	const copy = platformMessages[locale].hero
	const microruntimes = microruntimesMessages[locale]
	const session = microruntimes.session

	return (
		<section
			aria-describedby="platform-hero-summary"
			aria-labelledby="platform-hero-heading"
			className="relative isolate -mt-[70px] overflow-hidden bg-background pt-[70px] min-[761px]:-mt-20 min-[761px]:pt-20"
		>
			<HeroShader tone="muted" />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="relative mx-auto w-full max-w-7xl px-4.5 pb-20 sm:px-9">
				<div className="flex min-h-[calc(100dvh-26rem)] items-center justify-center text-center">
					<h1
						className="text-[min(3rem,11vw)] leading-none font-normal tracking-[-0.055em] text-foreground sm:text-6xl lg:text-7xl"
						id="platform-hero-heading"
					>
						<span className="block">{copy.titleStart}</span>
						<RotatingText
							className="block whitespace-nowrap"
							suffix={<ChromaText>{copy.titleEnd}</ChromaText>}
							words={copy.words}
						/>
					</h1>
				</div>
				<ElevatedFrame className="grid w-full min-w-0 gap-2 p-2 lg:min-h-160 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
					<div className="min-h-0 min-w-0">
						<SessionTerminal
							className="h-[clamp(220px,34vh,300px)] shadow-none"
							control={session.card.control}
							copy={session.terminal}
						/>
					</div>
					<div className="min-h-0 min-w-0">
						<SessionCards copy={session.card} locale={locale} />
					</div>
				</ElevatedFrame>
			</div>
			<p className="sr-only" id="platform-hero-summary">
				{microruntimes.hero.summary}
			</p>
		</section>
	)
}
