import { render } from '@testing-library/react'
import { createStore, Provider } from 'jotai'

import { finopsMessages } from '@i18n/messages/index'
import { ComputeLedgerDemo } from '@routes/platform/finops/sections/compute-ledger-demo'
import { EnvironmentReceipt } from '@routes/platform/finops/sections/environment-receipt'

export function mountComputeLedgerDemo() {
	const store = createStore()
	const copy = finopsMessages.en

	render(
		<Provider store={store}>
			<ComputeLedgerDemo copy={copy.ledger} locale="en" />
			<EnvironmentReceipt
				copy={copy.receipt}
				environmentsCopy={copy.ledger.environments}
				locale="en"
			/>
		</Provider>,
	)

	return store
}
