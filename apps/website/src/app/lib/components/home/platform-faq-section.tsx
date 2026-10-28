import { FaqTabs } from '@components/home/faq-tabs'
import { FaqSectionShell } from '@components/product-page/faq-section-shell'
import { finopsFaqIds } from '@data/finops-content'
import { microruntimesFaqIds } from '@data/microruntimes-content'
import {
	finopsMessages,
	headerMessages,
	microruntimesMessages,
	platformMessages,
} from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface PlatformFaqSectionProps {
	locale: WebsiteLocale
}

// The platform's FAQ: the questions each product page answers, one product
// per tab, named as the header names them.
export function PlatformFaqSection(props: PlatformFaqSectionProps) {
	const { locale } = props

	const copy = platformMessages[locale].faq
	const names = headerMessages[locale]
	const microruntimes = microruntimesMessages[locale].faq.items
	const finops = finopsMessages[locale].faq.items
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
			id: 'finops',
			items: finopsFaqIds.map((id) => ({
				id,
				...finops[id],
			})),
			label: names.finops,
		},
	]

	return (
		<FaqSectionShell description={copy.description} title={copy.title}>
			<FaqTabs groups={groups} tabsLabel={copy.tabsLabel} />
		</FaqSectionShell>
	)
}
