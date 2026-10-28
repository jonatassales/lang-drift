// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, describe, expect, it } from 'vitest'

import { type PlanLineId, planLineIndex } from '@data/microruntimes-content'
import { microruntimesMessages } from '@i18n/messages/index'
import {
	initialMicroruntimeSession,
	microruntimeSessionAtom,
	recordPlayedCountAtom,
	replaySessionAtom,
} from '@state'

import { SessionCards } from './session-cards'

afterEach(cleanup)

function mount() {
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

	const header = (label: 'Ticket' | 'Analysis' | 'Plan') =>
		screen.getByRole('button', {
			hidden: true,
			name: new RegExp(`^${label}`),
		})

	const expanded = (label: 'Ticket' | 'Analysis' | 'Plan') =>
		header(label).getAttribute('aria-expanded') === 'true'

	return {
		expanded,
		header,
		play,
		store,
	}
}

describe('Microruntimes session cards', () => {
	it('stays hidden with the ticket open until the command is entered', () => {
		const { expanded, play } = mount()
		const stack = screen.getByText('0031').closest('[aria-hidden]')
		expect(stack?.getAttribute('aria-hidden')).toBe('true')
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
		expect(screen.getAllByText('Not set yet').length).toBeGreaterThanOrEqual(4)

		play('plan')
		expect(stack?.getAttribute('aria-hidden')).toBeNull()

		play('links')
		expect(screen.getAllByText('0002-PRD-early-access.prd.md')).toHaveLength(2)
	})

	it('opens the card of each act as the session reaches it', () => {
		const { expanded, play } = mount()
		play('links')
		expect(expanded('Ticket')).toBe(true)

		play('analysis')
		expect(expanded('Ticket')).toBe(false)
		expect(expanded('Analysis')).toBe(true)

		play('planning')
		expect(expanded('Analysis')).toBe(false)
		expect(expanded('Plan')).toBe(true)
	})

	it('lets the visitor open a card by hand until the next act takes it back', () => {
		const { expanded, header, play, store } = mount()
		play('readPrd')
		fireEvent.click(header('Ticket'))
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
		expect(store.get(microruntimeSessionAtom).openedCard).toBe('ticket')

		play('readSpec')
		expect(expanded('Ticket')).toBe(true)

		play('planning')
		expect(expanded('Plan')).toBe(true)
		expect(store.get(microruntimeSessionAtom).openedCard).toBeNull()
	})

	it('reads the documents and judges the episodes line by line', () => {
		const { play } = mount()
		expect(screen.queryAllByText('Read')).toHaveLength(0)

		play('readSpec')
		expect(screen.getAllByText('Read')).toHaveLength(3)
		expect(
			screen.getByText('5 components · 3 acceptance criteria'),
		).toBeTruthy()

		play('judgeGuardrail')
		expect(screen.getAllByText('Read')).toHaveLength(4)
		expect(
			screen.getByText(
				"FIN-02: a paid action needs a gate · FIN-04: the budget is the team's",
			),
		).toBeTruthy()

		play('vectorSearch')
		expect(screen.getAllByText('Resolving')).toHaveLength(8)

		play('candidates')
		expect(screen.getByText('checkout-preview-sizing')).toBeTruthy()
		expect(screen.getByText('0.81')).toBeTruthy()
		expect(screen.queryByText('kept')).toBeNull()

		play('drop')
		expect(screen.getAllByText('kept')).toHaveLength(3)
		expect(screen.getByText('dropped · docs only')).toBeTruthy()

		play('suggest')
		expect(screen.getByText('operator/allowed-model')).toBeTruthy()
		expect(
			screen.getByText('recommendation · the team authorizes'),
		).toBeTruthy()
	})

	it('carries each act’s latest figure in its collapsed header', () => {
		const { header, play } = mount()
		expect(header('Ticket').textContent).toContain(
			'0031 · early-access-runtime',
		)

		play('reads')
		expect(header('Analysis').textContent).toContain('4 documents read')

		play('drop')
		expect(header('Analysis').textContent).toContain('3 of 6 episodes kept')

		play('suggest')
		expect(header('Analysis').textContent).toContain(
			'operator/allowed-model · effort medium',
		)

		play('go')
		expect(header('Plan').textContent).toContain('250m · 256 MiB')

		play('draft')
		expect(header('Plan').textContent).toContain(
			'1,450m · 2,304 MiB · pending team authorization',
		)
	})

	it('steps the plan’s figures up with the components and closes the report', () => {
		const { play } = mount()
		play('bounds')
		expect(screen.getByText('2,000m · 3,072 MiB')).toBeTruthy()

		play('resolve')
		expect(screen.getAllByText('Resolving').length).toBeGreaterThanOrEqual(8)

		play('go')
		expect(screen.getAllByText('250m · 256 MiB').length).toBeGreaterThanOrEqual(
			1,
		)
		expect(screen.getByText('250m of 2,000m')).toBeTruthy()
		expect(screen.getByText('1 of 5 resolved · 250m · 256 MiB')).toBeTruthy()

		play('crossCheck')
		expect(screen.getByText('5 components · 1,450m · 2,304 MiB')).toBeTruthy()
		expect(
			screen.getByText(
				'ports 8080 · 8081 · 4222 · 6379 · 3000 · distinct · no collisions',
			),
		).toBeTruthy()

		play('size')
		expect(screen.getByText('550m · 768 MiB')).toBeTruthy()
		expect(screen.getByText('5 one-replica Deployments')).toBeTruthy()

		play('price')
		expect(
			screen.getByText('no price source configured · units only'),
		).toBeTruthy()

		play('validate')
		expect(screen.getByText('schema · 0 errors')).toBeTruthy()

		play('draft')
		expect(screen.getAllByText('pending team authorization')).toHaveLength(1)
		expect(
			screen.getByText(
				'runtime.md draft · revision 1 · proposed to the console',
			),
		).toBeTruthy()
	})

	it('opens the ticket again when the session replays', () => {
		const { expanded, header, play, store } = mount()
		play('handoff')
		fireEvent.click(header('Analysis'))
		expect(expanded('Analysis')).toBe(true)

		act(() => {
			store.set(replaySessionAtom)
		})
		expect(store.get(microruntimeSessionAtom)).toEqual({
			...initialMicroruntimeSession,
			run: 1,
		})
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
	})
})
