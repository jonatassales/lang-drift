// @vitest-environment jsdom

import { act, cleanup, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, describe, expect, it } from 'vitest'

import { type PlanLineId, planLineIndex } from '@data/microruntimes-content'
import { microruntimesMessages } from '@i18n/messages/index'
import { initialMicroruntimeSession, microruntimeSessionAtom } from '@state'

import { TicketCard } from './ticket-card'

afterEach(cleanup)

function mount() {
	const store = createStore()
	render(
		<Provider store={store}>
			<TicketCard copy={microruntimesMessages.en.session.card} />
		</Provider>,
	)

	const play = (lineId: PlanLineId) =>
		act(() => {
			store.set(microruntimeSessionAtom, {
				...initialMicroruntimeSession,
				playedCount: planLineIndex(lineId) + 1,
			})
		})

	return {
		play,
		store,
	}
}

describe('Microruntimes ticket card', () => {
	it('stays hidden with the artifacts and dashes in place until the command is entered', () => {
		const { play } = mount()
		const card = screen.getByRole('article', { hidden: true })
		expect(card.getAttribute('aria-hidden')).toBe('true')
		expect(screen.getByText('0031')).toBeTruthy()
		expect(screen.getByText('0002-PRD-early-access.prd.md')).toBeTruthy()
		expect(screen.getAllByText('Not set yet').length).toBeGreaterThanOrEqual(3)

		play('plan')
		expect(card.getAttribute('aria-hidden')).toBeNull()
	})

	it('turns the suggestion into a skeleton while the console searches, then shows it', () => {
		const { play } = mount()
		play('memory')
		expect(screen.getAllByText('Resolving')).toHaveLength(2)

		play('suggest')
		expect(screen.getByText('operator/allowed-model')).toBeTruthy()
		expect(screen.getByText('medium')).toBeTruthy()
		expect(
			screen.getByText('recommendation · the team authorizes'),
		).toBeTruthy()
	})

	it('marks the artifacts read and fills the guardrail as the agent reads them', () => {
		const { play } = mount()
		expect(screen.queryAllByText('Read')).toHaveLength(0)

		play('reads')
		expect(screen.getAllByText('Read').length).toBeGreaterThanOrEqual(4)

		play('guardrail')
		expect(screen.getByText('0001-finops.gr.md · FIN-02 · FIN-04')).toBeTruthy()
		expect(
			screen.getByText('review rule · no automatic enforcement'),
		).toBeTruthy()
	})
})
