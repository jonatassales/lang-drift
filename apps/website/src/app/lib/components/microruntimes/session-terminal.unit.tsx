// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { planSession } from '@data/microruntimes-content'
import { microruntimesMessages } from '@i18n/messages/index'
import {
	initialMicroruntimeSession,
	microruntimeSessionAtom,
	playWorkflowAtom,
	recordPlayedCountAtom,
} from '@state'

import { SessionTerminal } from './session-terminal'

class IntersectingObserver {
	callback: IntersectionObserverCallback

	constructor(callback: IntersectionObserverCallback) {
		this.callback = callback
	}

	observe(node: Element) {
		this.callback(
			[
				{
					isIntersecting: true,
					target: node,
				} as IntersectionObserverEntry,
			],
			this as unknown as IntersectionObserver,
		)
	}

	disconnect() {}

	unobserve() {}
}

beforeEach(() => {
	vi.useFakeTimers()
	vi.stubGlobal('IntersectionObserver', IntersectingObserver)
	vi.stubGlobal('matchMedia', (media: string) => ({
		addEventListener: () => {},
		addListener: () => {},
		dispatchEvent: () => false,
		matches: false,
		media,
		onchange: null,
		removeEventListener: () => {},
		removeListener: () => {},
	}))
})

afterEach(() => {
	cleanup()
	vi.useRealTimers()
	vi.unstubAllGlobals()
})

function advance(milliseconds: number) {
	for (let elapsed = 0; elapsed < milliseconds; elapsed += 40) {
		act(() => {
			vi.advanceTimersByTime(40)
		})
	}
}

function mount() {
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

describe('Microruntimes session terminal', () => {
	it('writes the number of completed lines to the session as it plays', () => {
		const store = mount()
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(0)

		advance(2000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(1)

		advance(600)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(2)
	})

	it('pauses and resumes the session from its title bar, and replays it once it has ended', () => {
		const store = mount()
		advance(2600)
		fireEvent.click(screen.getByRole('button', { name: 'Pause' }))
		expect(store.get(microruntimeSessionAtom).paused).toBe(true)

		advance(2000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(2)

		fireEvent.click(screen.getByRole('button', { name: 'Resume' }))
		advance(1000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBeGreaterThan(2)

		act(() => {
			store.set(recordPlayedCountAtom, planSession.length)
		})
		fireEvent.click(screen.getByRole('button', { name: 'Replay' }))
		expect(store.get(microruntimeSessionAtom)).toEqual({
			...initialMicroruntimeSession,
			run: 1,
		})
	})

	it('scrolls up to the window when the workflow plays, starts the run a second later and records its lines', () => {
		const scrollIntoView = vi.fn()
		Element.prototype.scrollIntoView = scrollIntoView
		const store = mount()
		act(() => {
			store.set(recordPlayedCountAtom, planSession.length)
			store.set(playWorkflowAtom)
		})
		expect(store.get(microruntimeSessionAtom).stage).toBe('launching')
		expect(scrollIntoView).toHaveBeenCalledWith(
			expect.objectContaining({ block: 'start' }),
		)

		advance(960)
		expect(store.get(microruntimeSessionAtom).stage).toBe('launching')

		advance(80)
		expect(store.get(microruntimeSessionAtom)).toMatchObject({
			playedCount: 0,
			run: 1,
			stage: 'running',
		})
		expect(scrollIntoView).toHaveBeenCalledOnce()

		advance(1500)
		expect(store.get(microruntimeSessionAtom).playedCount).toBeGreaterThan(0)
		expect(screen.getByRole('button', { name: 'Pause' })).toBeTruthy()
	})
})
