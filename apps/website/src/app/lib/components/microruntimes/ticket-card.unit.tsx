// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
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
			<TicketCard copy={microruntimesMessages.en.session.card} locale="en" />
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

	it('resolves the components one by one with their figures and a running total', () => {
		const { play } = mount()
		play('bounds')
		expect(screen.queryByText('services/signup-api')).toBeNull()

		play('resolve')
		expect(screen.getAllByText('Resolving')).toHaveLength(6)

		play('go')
		expect(screen.getByText('services/signup-api')).toBeTruthy()
		expect(screen.getByText('Go')).toBeTruthy()
		expect(screen.getAllByText('api')).toHaveLength(1)
		expect(screen.getByText('250m')).toBeTruthy()
		expect(screen.getByText('256 MiB')).toBeTruthy()
		expect(screen.getByText('1 of 5 resolved · 250m · 256 MiB')).toBeTruthy()
		expect(screen.getAllByText('Resolving')).toHaveLength(4)

		play('next')
		expect(screen.getByText('apps/signup')).toBeTruthy()
		expect(screen.getByText('5 components · 1,450m · 2,304 MiB')).toBeTruthy()
		expect(screen.queryAllByText('Resolving')).toHaveLength(0)
	})

	it('fills the hardware table and the report as the plan is bounded, sized and drafted', () => {
		const { play } = mount()
		expect(screen.queryByText('2,000m')).toBeNull()

		play('bounds')
		expect(screen.getAllByText('2,000m')).toHaveLength(1)
		expect(screen.getByText('2,000m · 3,072 MiB')).toBeTruthy()

		play('size')
		expect(screen.getAllByText('1,450m')).toHaveLength(1)
		expect(screen.getByText('1,450m · 2,304 MiB')).toBeTruthy()

		play('limits')
		expect(screen.getAllByText('none')).toHaveLength(4)
		expect(screen.getAllByText('not a runtime.md field')).toHaveLength(2)

		play('price')
		expect(
			screen.getByText('price basis: none configured · units only'),
		).toBeTruthy()

		play('draft')
		expect(screen.getByText('5 one-replica Deployments')).toBeTruthy()
		expect(
			screen.getByText('1,450m · 2,304 MiB · pending team authorization'),
		).toBeTruthy()
		expect(
			screen.getByText(
				'runtime.md draft · revision 1 · proposed to the console',
			),
		).toBeTruthy()
	})

	it('pauses the session and replays it from the start', () => {
		const { play, store } = mount()
		play('go')
		fireEvent.click(screen.getByRole('button', { name: 'Pause' }))
		expect(store.get(microruntimeSessionAtom).paused).toBe(true)
		expect(screen.getByRole('button', { name: 'Resume' })).toBeTruthy()

		play('handoff')
		fireEvent.click(screen.getByRole('button', { name: 'Replay' }))
		expect(store.get(microruntimeSessionAtom)).toEqual({
			paused: false,
			playedCount: 0,
			run: 1,
		})
		expect(
			screen.getByRole('article', { hidden: true }).getAttribute('aria-hidden'),
		).toBe('true')
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
