import {
	finopsMessages,
	headerMessages,
	llmopsMessages,
	microruntimesMessages,
	platformMessages,
	safeSeatMessages,
} from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import { finopsFaqIds } from '@routes/platform/finops/content/finops-content.data'
import { llmopsFaqIds } from '@routes/platform/llmops/content/llmops-content.data'
import { microruntimesFaqIds } from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { FaqSectionShell } from '@routes/platform/product-page/faq-section-shell'
import { safeSeatFaqIds } from '@routes/platform/safe-seat/safe-seat-content.data'

import { FaqTabs } from './faq-tabs'

// The platform's one FAQ: each product's questions in a tab of their own,
// named as the header names the product.
interface PlatformFaqSectionProps {
	locale: WebsiteLocale
}

export function PlatformFaqSection(props: PlatformFaqSectionProps) {
	const { locale } = props

	const copy = platformMessages[locale].faq
	const names = headerMessages[locale]
	const microruntimes = microruntimesMessages[locale].faq.items
	const llmops = llmopsMessages[locale].faq.items
	const finops = finopsMessages[locale].faq.items
	const safeSeat = safeSeatMessages[locale].faq.items
	const groups = [
		{
			id: 'microruntimes',
			items: microruntimesFaqIds.map((id) => ({
				id,
				...microruntimes[id],
			})),
			label: names.microruntimes,
		},
		{
			id: 'llmops',
			items: llmopsFaqIds.map((id) => ({
				id,
				...llmops[id],
			})),
			label: names.llmops,
		},
		{
			id: 'finops',
			items: finopsFaqIds.map((id) => ({
				id,
				...finops[id],
			})),
			label: names.finops,
		},
		{
			id: 'safe-seat',
			items: safeSeatFaqIds.map((id) => ({
				id,
				...safeSeat[id],
			})),
			label: names.safeSeat,
		},
	]

	return (
		<FaqSectionShell
			className="max-w-[1280px]"
			description={copy.description}
			title={copy.title}
		>
			<FaqTabs groups={groups} tabsLabel={copy.tabsLabel} />
		</FaqSectionShell>
	)
}
