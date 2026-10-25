// @vitest-environment jsdom

import { act, cleanup, render } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { microruntimesMessages } from '@i18n/messages/index'
import { microruntimeSessionAtom } from '@state'

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

describe('Microruntimes session terminal', () => {
	it('writes the number of completed lines to the session as it plays', () => {
		const store = createStore()
		render(
			<Provider store={store}>
				<SessionTerminal copy={microruntimesMessages.en.session.terminal} />
			</Provider>,
		)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(0)

		advance(2000)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(1)

		advance(600)
		expect(store.get(microruntimeSessionAtom).playedCount).toBe(2)
	})
})
