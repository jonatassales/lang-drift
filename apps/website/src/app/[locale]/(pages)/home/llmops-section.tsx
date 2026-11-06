import { llmopsMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import { platformProduct } from '@layout/website-header/platform-product.compute'
import { BudgetLedgerDemo } from '@routes/platform/llmops/sections/budget-ledger-demo'
import { DeliveryReceipt } from '@routes/platform/llmops/sections/delivery-receipt'

import { HomeProductSection } from './home-product-section'

// The LLMOps page's opening on the home: the envelope statement, each
// sentence on its own line, over the budget ledger with the delivery receipt
// docked beside it from `lg`, the receipt following the selected work item.
interface LlmopsSectionProps {
	locale: WebsiteLocale
}

export function LlmopsSection(props: LlmopsSectionProps) {
	const { locale } = props

	const copy = llmopsMessages[locale]
	const sections = platformMessages[locale].sections.llmops

	return (
		<HomeProductSection
			description={sections.description}
			headline={
				<>
					<span className="block">{copy.hero.taglineStart}</span>{' '}
					<span className="block">{copy.hero.taglineEnd}</span>
				</>
			}
			id="llmops"
			link={sections.link}
			locale={locale}
			product={platformProduct('llmops')}
			summary={copy.hero.summary}
		>
			<div className="grid gap-6 lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)]">
				<BudgetLedgerDemo copy={copy.ledger} locale={locale} />
				<DeliveryReceipt
					copy={copy.receipt}
					headingLevel={3}
					items={copy.ledger.items}
					locale={locale}
				/>
			</div>
		</HomeProductSection>
	)
}
