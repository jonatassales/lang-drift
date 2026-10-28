import { ChromaText } from '@langdrift/web/vendors/smoothui/chroma-text'
import { createTranslator } from 'next-intl'

import { HeroChart } from '@components/hero-chart'
import { HeroSinapsiGraph } from '@components/hero-sinapsi-graph'
import { HomeProductSection } from '@components/home/home-product-section'
import { platformProduct } from '@components/website-header/platform-products'
import { homeMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface VisionSectionProps {
	locale: WebsiteLocale
}

// The Vision page's hero on the home: its question, description, curve and
// graph, caption and summary as they are. The demo grid is the Tailwind twin
// of that page's `website-demo-grid`, two columns from `lg`.
export function VisionSection(props: VisionSectionProps) {
	const { locale } = props

	const copy = homeMessages[locale].hero
	const sections = platformMessages[locale].sections
	const t = createTranslator({
		locale,
		messages: { title: copy.title },
	})

	return (
		<HomeProductSection
			caption={copy.caption}
			description={copy.description}
			headline={t.rich('title', {
				movement: (chunks) => <ChromaText>{chunks}</ChromaText>,
			})}
			id="vision"
			link={sections.vision.link}
			locale={locale}
			product={platformProduct('vision')}
			summary={copy.summary}
		>
			<div className="grid gap-6 *:mx-auto *:w-full *:min-w-0 *:max-w-[780px] lg:grid-cols-[minmax(460px,0.42fr)_minmax(0,0.58fr)] lg:*:mx-0 lg:*:max-w-none">
				<HeroChart />
				<HeroSinapsiGraph headingLevel={3} />
			</div>
		</HomeProductSection>
	)
}
