import { render } from '@testing-library/react'
import { createStore, Provider } from 'jotai'

import { microruntimesMessages } from '@i18n/messages/index'
import { SessionTerminal } from '@routes/platform/microruntimes/sections/session-terminal'

export function mountSessionTerminal() {
	const store = createStore()
	render(
		<Provider store={store}>
			<SessionTerminal
				control={microruntimesMessages.en.session.card.control}
				copy={microruntimesMessages.en.session.terminal}
			/>
		</Provider>,
	)

	return store
}
