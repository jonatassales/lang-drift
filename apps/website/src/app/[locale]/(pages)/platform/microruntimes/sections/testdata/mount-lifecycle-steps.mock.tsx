import { render } from '@testing-library/react'
import { createStore, Provider } from 'jotai'

import { microruntimesMessages } from '@i18n/messages/index'
import { LifecycleSteps } from '@routes/platform/microruntimes/sections/lifecycle-steps'
import { StageEvidence } from '@routes/platform/microruntimes/sections/stage-evidence'

export function mountLifecycleSteps() {
	const store = createStore()
	const copy = microruntimesMessages.en
	render(
		<Provider store={store}>
			<LifecycleSteps copy={copy.lifecycle} />
			<StageEvidence copy={copy.evidence} stageLabels={copy.lifecycle.stages} />
		</Provider>,
	)

	return store
}
