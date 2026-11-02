// @vitest-environment jsdom
import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { Toast } from './toast'

afterEach(cleanup)

describe('web/ui/toast', () => {
	it('renders message with alert role and assertive live region', () => {
		render(<Toast>Something went wrong</Toast>)

		const alert = screen.getByRole('alert')
		expect(alert.textContent).toBe('Something went wrong')
		expect(alert.getAttribute('aria-live')).toBe('assertive')
	})

	it('calls onClose when close button is clicked', () => {
		const onCloseMock = vi.fn()
		render(<Toast onClose={onCloseMock}>Error message</Toast>)

		const closeButton = screen.getByRole('button', {
			name: 'Close notification',
		})
		fireEvent.click(closeButton)

		expect(onCloseMock).toHaveBeenCalledOnce()
	})

	it('renders without close button when onClose is omitted', () => {
		render(<Toast>Notification</Toast>)

		expect(
			screen.queryByRole('button', { name: 'Close notification' }),
		).toBeNull()
	})
})
