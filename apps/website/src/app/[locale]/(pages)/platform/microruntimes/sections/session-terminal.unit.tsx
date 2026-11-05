// @vitest-environment jsdom

import { act, cleanup, fireEvent, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { planSession } from '@routes/platform/microruntimes/content/microruntimes-content'
import { advanceTerminal } from '@routes/platform/microruntimes/sections/testdata/advance-terminal.mock'
import { mountSessionTerminal } from '@routes/platform/microruntimes/sections/testdata/mount-session-terminal.mock'
import {
	initialMicroruntimeSession,
	microruntimeSessionAtom,
	playWorkflowAtom,
	recordPlayedCountAtom,
} from '@state'

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

describe('Microruntimes session terminal', () => {
	it('writes the number of completed lines to the session as it plays', () => {
		const store = mountSessionTerminal()
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(0)

		advanceTerminal(2000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(1)

		advanceTerminal(600)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(2)
	})

	it('pauses and resumes the session from its title bar, and replays it once it has ended', () => {
		const store = mountSessionTerminal()
		advanceTerminal(2600)
		fireEvent.click(screen.getByRole('button', { name: 'Pause' }))
		expect(store.get(microruntimeSessionAtom).paused).toBe(true)

		advanceTerminal(2000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(2)

		fireEvent.click(screen.getByRole('button', { name: 'Resume' }))
		advanceTerminal(1000)
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
		const store = mountSessionTerminal()
		act(() => {
			store.set(recordPlayedCountAtom, planSession.length)
			store.set(playWorkflowAtom)
		})
		expect(store.get(microruntimeSessionAtom).stage).toBe('launching')
		expect(scrollIntoView).toHaveBeenCalledWith(
			expect.objectContaining({ block: 'start' }),
		)

		advanceTerminal(960)
		expect(store.get(microruntimeSessionAtom).stage).toBe('launching')

		advanceTerminal(80)
		expect(store.get(microruntimeSessionAtom)).toMatchObject({
			playedCount: 0,
			run: 1,
			stage: 'running',
		})
		expect(scrollIntoView).toHaveBeenCalledOnce()

		advanceTerminal(1500)
		expect(store.get(microruntimeSessionAtom).playedCount).toBeGreaterThan(0)
		expect(screen.getByRole('button', { name: 'Pause' })).toBeTruthy()
	})
})
