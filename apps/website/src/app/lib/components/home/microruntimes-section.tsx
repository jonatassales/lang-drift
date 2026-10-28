import { HomeProductSection } from '@components/home/home-product-section'
import { SessionCards } from '@components/microruntimes/session-cards'
import { SessionTerminal } from '@components/microruntimes/session-terminal'
import { platformProduct } from '@components/website-header/platform-products'
import { microruntimesMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface MicroruntimesSectionProps {
	locale: WebsiteLocale
}

// The Microruntimes page's planning session on the home, under a question of
// its own: the terminal plays the session into the ticket, plan and workflow
// cards beside it, as on that page.
export function MicroruntimesSection(props: MicroruntimesSectionProps) {
	const { locale } = props

	const copy = microruntimesMessages[locale]
	const sections = platformMessages[locale].sections.microruntimes

	return (
		<HomeProductSection
			description={sections.description}
			headline={sections.title}
			id="microruntimes"
			link={{
				href: '/platform/microruntimes',
				label: sections.link,
			}}
			locale={locale}
			product={platformProduct('microruntimes')}
			summary={copy.hero.summary}
		>
			<div className="grid gap-6 lg:grid-cols-[fit-content(50%)_1fr]">
				<SessionTerminal
					control={copy.session.card.control}
					copy={copy.session.terminal}
				/>
				<SessionCards copy={copy.session.card} locale={locale} />
			</div>
		</HomeProductSection>
	)
}
