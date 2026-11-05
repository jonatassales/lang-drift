import { render } from '@testing-library/react'

import { finopsMessages } from '@i18n/messages/index'
import {
	gateExample,
	gateFigures,
} from '@routes/platform/finops/content/finops-content'
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
