import { act, cleanup, render } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import type { TerminalLine } from './mac-terminal'
import { MacTerminal } from './mac-terminal'

const lines: TerminalLine[] = [
	{
		id: 'doctor',
		text: 'ex doctor',
		type: 'command',
	},
	{
		id: 'tools',
		text: 'docker 27.3',
		type: 'output',
	},
]

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

// motion reads the preference once and then follows "change" events, so the
// stub keeps one listener set alive across tests and notifies it.
const motionPreference = {
	listeners: new Set<() => void>(),
	reduced: false,
}

function stubBrowser(prefersReducedMotion: boolean) {
	motionPreference.reduced = prefersReducedMotion
	for (const listener of motionPreference.listeners) listener()

	vi.stubGlobal('IntersectionObserver', IntersectingObserver)
	vi.stubGlobal('matchMedia', (media: string) => ({
		addEventListener: (_type: string, listener: () => void) =>
			motionPreference.listeners.add(listener),
		addListener: () => {},
		dispatchEvent: () => false,
		get matches() {
			return motionPreference.reduced && media.includes('reduced-motion')
		},
		media,
		onchange: null,
		removeEventListener: () => {},
		removeListener: () => {},
	}))
}

function advance(milliseconds: number) {
	for (let elapsed = 0; elapsed < milliseconds; elapsed += 40) {
		act(() => {
			vi.advanceTimersByTime(40)
		})
	}
}

beforeEach(() => {
	vi.useFakeTimers()
})

afterEach(() => {
	cleanup()
	vi.useRealTimers()
	vi.unstubAllGlobals()
})

describe('MacTerminal progress', () => {
	it('reports each completed line in order and never reports zero', () => {
		stubBrowser(false)
		const onProgress = vi.fn()
		render(
			<MacTerminal
				label="Terminal"
				lines={lines}
				onProgress={onProgress}
				title="ex"
			/>,
		)

		advance(5000)

		expect(onProgress.mock.calls).toEqual([[1], [2]])
	})

	it('stops while paused and resumes where it was', () => {
		stubBrowser(false)
		const onProgress = vi.fn()
		const view = render(
			<MacTerminal
				label="Terminal"
				lines={lines}
				onProgress={onProgress}
				title="ex"
			/>,
		)

		advance(600)
		expect(onProgress).toHaveBeenLastCalledWith(1)

		view.rerender(
			<MacTerminal
				label="Terminal"
				lines={lines}
				onProgress={onProgress}
				paused
				title="ex"
			/>,
		)
		advance(5000)
		expect(onProgress).toHaveBeenCalledTimes(1)

		view.rerender(
			<MacTerminal
				label="Terminal"
				lines={lines}
				onProgress={onProgress}
				paused={false}
				title="ex"
			/>,
		)
		advance(1000)
		expect(onProgress).toHaveBeenLastCalledWith(2)
	})

	it('reports every line at once under reduced motion', () => {
		stubBrowser(true)
		const onProgress = vi.fn()
		render(
			<MacTerminal
				label="Terminal"
				lines={lines}
				onProgress={onProgress}
				title="ex"
			/>,
		)

		expect(onProgress.mock.calls).toEqual([[2]])
	})
})
