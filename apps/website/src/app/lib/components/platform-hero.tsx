import { ChromaText } from '@langdrift/web/ui/chroma-text'
import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { RotatingText } from '@langdrift/web/ui/rotating-text'
import { cn } from '@langdrift/web/utilities'

import { SessionCards } from '@components/microruntimes/session-cards'
import { SessionTerminal } from '@components/microruntimes/session-terminal'
import { microruntimesMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformHeroProps {
	locale: WebsiteLocale
}

// The hero is the moving headline alone, centered on the first screen over a
// black-and-white shader, with the rainbow flowing through the word
// "Platform" as the page's one touch of color. The shader's brand roles are
// swapped for neutral ones, so it draws in grays: the light scheme takes the
// secondary text gray at under half opacity, so its veils show over the
// pale field while it stays closer to white than to silver, and the dark
// one takes the surface and border grays at under half opacity, so its
// veils barely lift from the black. Under the headline, one
// frame holds the Microruntimes terminal and the ticket, plan and workflow
// cards it plays into, side by side; its top 21rem show at the bottom of the first screen.
// The hero starts under the transparent header, so the two read as one
// background. From `lg` the frame is at least 40rem tall and grows with the
// cards, so they never scroll; the terminal stretches to the same height and
// scrolls its transcript inside; below `lg`, the cards stack under the
// terminal.
// The frame sits above the page: a lit top edge and a deep, soft shadow
// under it, in the foreground's gray over the light field and in the
// background's black over the dark one, where a lit edge does the lifting.
const elevated =
	'shadow-[inset_0_1px_0_color-mix(in_srgb,var(--foreground)_14%,transparent),0_32px_80px_-20px_color-mix(in_srgb,var(--foreground)_38%,transparent)] dark:shadow-[inset_0_1px_0_color-mix(in_srgb,var(--foreground)_16%,transparent),0_32px_80px_-16px_color-mix(in_srgb,var(--background)_95%,transparent)]'

const monochrome =
	'opacity-45 [--primary-hover:var(--muted-foreground)] [--primary:var(--muted-foreground)] dark:opacity-45 dark:[--primary-hover:var(--border)] dark:[--primary:var(--muted)]'

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
			<HeroShader className={monochrome} />
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
				<div
					className={cn(
						'grid w-full min-w-0 gap-2 rounded-2xl border border-border bg-background/70 p-2 backdrop-blur-xl lg:min-h-160 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]',
						elevated,
					)}
				>
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
				</div>
			</div>
			<p className="sr-only" id="platform-hero-summary">
				{microruntimes.hero.summary}
			</p>
		</section>
	)
}
