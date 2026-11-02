// @vitest-environment jsdom
import { render } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const pushMock = vi.fn()
vi.mock('@i18n/navigation', () => ({
	useRouter: () => ({ push: pushMock }),
}))

import { AutoRedirect } from './auto-redirect'

describe('website/early-access-sent/auto-redirect', () => {
	beforeEach(() => {
		vi.useFakeTimers()
		pushMock.mockReset()
	})

	afterEach(() => {
		vi.useRealTimers()
	})

	it('redirects to root destination after 5000ms by default', () => {
		render(<AutoRedirect />)

		expect(pushMock).not.toHaveBeenCalled()

		vi.advanceTimersByTime(4999)
		expect(pushMock).not.toHaveBeenCalled()

		vi.advanceTimersByTime(1)
		expect(pushMock).toHaveBeenCalledWith('/')
	})

	it('cleans up timeout on unmount', () => {
		const { unmount } = render(<AutoRedirect delayMs={5000} />)

		unmount()

		vi.advanceTimersByTime(6000)
		expect(pushMock).not.toHaveBeenCalled()
	})
})
