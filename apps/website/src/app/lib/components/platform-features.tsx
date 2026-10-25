import { FeaturesGrid } from '@langdrift/web/vendors/smoothui/features-grid'

import { platformProducts } from '@components/website-header/platform-products'
import { headerMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformFeaturesProps {
	locale: WebsiteLocale
}

// The second section of the home page: the four platform products as cards,
// named and described with the header's copy.
export function PlatformFeatures(props: PlatformFeaturesProps) {
	const { locale } = props

	const copy = platformMessages[locale].features
	const names = headerMessages[locale]
	const items = platformProducts.map((product) => {
		const Icon = product.icon

		return {
			description: names[product.description],
			icon: <Icon className="size-6" />,
			id: product.label,
			title: names[product.label],
		}
	})

	return (
		<FeaturesGrid
			headingId="platform-features-heading"
			items={items}
			lead={copy.lead}
			title={copy.title}
		/>
	)
}
