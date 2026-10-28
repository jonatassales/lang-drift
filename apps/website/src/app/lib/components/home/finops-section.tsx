import { BudgetLedgerDemo } from '@components/finops/budget-ledger-demo'
import { DeliveryReceipt } from '@components/finops/delivery-receipt'
import { HomeProductSection } from '@components/home/home-product-section'
import { platformProduct } from '@components/website-header/platform-products'
import { finopsMessages, platformMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface FinopsSectionProps {
	locale: WebsiteLocale
}

// The FinOps page's opening on the home: the envelope statement, each
// sentence on its own line, over the budget ledger with the delivery receipt
// docked beside it from `lg`, the receipt following the selected work item.
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
					<span className="block">{copy.hero.titleStart}</span>{' '}
					<span className="block">{copy.hero.titleEnd}</span>
				</>
			}
			id="finops"
			link={{
				href: '/platform/finops',
				label: sections.link,
			}}
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
