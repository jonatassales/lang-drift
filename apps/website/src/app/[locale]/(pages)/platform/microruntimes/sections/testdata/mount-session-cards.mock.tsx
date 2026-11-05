import { act, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'

import { microruntimesMessages } from '@i18n/messages/index'
import type { PlanLineId } from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { planLineIndex } from '@routes/platform/microruntimes/content/plan-line-index.compute'
import { SessionCards } from '@routes/platform/microruntimes/sections/session-cards'
import { recordPlayedCountAtom } from '@state'

type CardLabel = 'Ticket' | 'Analysis' | 'Plan' | 'Workflow'

export function mountSessionCards() {
	const store = createStore()

	render(
		<Provider store={store}>
			<SessionCards copy={microruntimesMessages.en.session.card} locale="en" />
		</Provider>,
	)

	// Records the count the way the terminal does, so an act change takes
	// the cards back from the visitor.
	const play = (lineId: PlanLineId) =>
		act(() => {
			store.set(recordPlayedCountAtom, planLineIndex(lineId) + 1)
		})

	const header = (label: CardLabel) =>
		screen.getByRole('button', {
			hidden: true,
			name: new RegExp(`^${label}`),
		})

	const expanded = (label: CardLabel) =>
		header(label).getAttribute('aria-expanded') === 'true'

	return {
		expanded,
		header,
		play,
		store,
	}
}
