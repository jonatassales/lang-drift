import { act, cleanup, render } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { ScrollHint } from './scroll-hint'

let report: IntersectionObserverCallback = () => {}

class RecordingObserver {
	constructor(callback: IntersectionObserverCallback) {
		report = callback
	}

	observe() {}

	disconnect() {}

	unobserve() {}
}

afterEach(() => {
	cleanup()
	vi.unstubAllGlobals()
})

function intersect(isIntersecting: boolean) {
	act(() => {
		report(
			[{ isIntersecting } as IntersectionObserverEntry],
			{} as IntersectionObserver,
		)
	})
}

describe('ScrollHint', () => {
	it('shows the band until the end of the content is on screen', () => {
		vi.stubGlobal('IntersectionObserver', RecordingObserver)
		const { container } = render(<ScrollHint />)
		const band = container.querySelector('[data-at-end]')
		expect(band?.getAttribute('data-at-end')).toBe('false')

		intersect(true)
		expect(band?.getAttribute('data-at-end')).toBe('true')

		intersect(false)
		expect(band?.getAttribute('data-at-end')).toBe('false')
	})
})
