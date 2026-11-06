import { render } from '@testing-library/react'
import { createStore, Provider } from 'jotai'

import { llmopsMessages } from '@i18n/messages/index'
import { BudgetLedgerDemo } from '@routes/platform/llmops/sections/budget-ledger-demo'
import { DeliveryReceipt } from '@routes/platform/llmops/sections/delivery-receipt'

export function mountBudgetLedgerDemo() {
	const store = createStore()
	const copy = llmopsMessages.en
	render(
		<Provider store={store}>
			<BudgetLedgerDemo copy={copy.ledger} locale="en" />
			<DeliveryReceipt
				copy={copy.receipt}
				items={copy.ledger.items}
				locale="en"
			/>
		</Provider>,
	)

	return store
}
