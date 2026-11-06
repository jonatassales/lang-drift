import { render } from '@testing-library/react'

import { llmopsMessages } from '@i18n/messages/index'
import { gateFigures } from '@routes/platform/llmops/content/gate-figures.compute'
import { gateExample } from '@routes/platform/llmops/content/llmops-content.data'
import { GateFlow } from '@routes/platform/llmops/sections/gate-flow'

export function mountGateFlow() {
	const copy = llmopsMessages.en

	return render(
		<GateFlow
			copy={copy.authorization}
			exampleTitle={copy.ledger.items[gateExample.id].title}
			figures={gateFigures(gateExample)}
			locale="en"
		/>,
	)
}
