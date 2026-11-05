import { render } from '@testing-library/react'

import { finopsMessages } from '@i18n/messages/index'
import { gateExample } from '@routes/platform/finops/content/finops-content.data'
import { gateFigures } from '@routes/platform/finops/content/gate-figures.compute'
import { GateFlow } from '@routes/platform/finops/sections/gate-flow'

export function mountGateFlow() {
	const copy = finopsMessages.en

	return render(
		<GateFlow
			copy={copy.authorization}
			exampleTitle={copy.ledger.items[gateExample.id].title}
			figures={gateFigures(gateExample)}
			locale="en"
		/>,
	)
}
