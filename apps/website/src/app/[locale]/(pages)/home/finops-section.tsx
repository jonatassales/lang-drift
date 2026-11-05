import { finopsMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'
import { platformProduct } from '@layout/website-header/platform-products'
import { BudgetLedgerDemo } from '@routes/platform/finops/sections/budget-ledger-demo'
import { DeliveryReceipt } from '@routes/platform/finops/sections/delivery-receipt'

import { HomeProductSection } from './home-product-section'

// The FinOps page's opening on the home: the envelope statement, each
// sentence on its own line, over the budget ledger with the delivery receipt
// docked beside it from `lg`, the receipt following the selected work item.
interface FinopsSectionProps {
	locale: WebsiteLocale
}

export function FinopsSection(props: FinopsSectionProps) {
	const { locale } = props

	const copy = finopsMessages[locale]
	const sections = platformMessages[locale].sections.finops

	return (
		<HomeProductSection
			caption={copy.hero.caption}
			description={sections.description}
			headline={
				<>
					<span className="block">{copy.hero.taglineStart}</span>{' '}
					<span className="block">{copy.hero.taglineEnd}</span>
				</>
			}
			id="finops"
			link={sections.link}
			locale={locale}
			product={platformProduct('finops')}
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
