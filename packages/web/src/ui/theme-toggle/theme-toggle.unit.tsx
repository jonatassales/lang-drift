import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { ThemeToggle } from './theme-toggle'

afterEach(cleanup)

describe('web/ui/theme-toggle', () => {
	it('names the scheme it switches to and reports a press', () => {
		const onToggle = vi.fn()

		render(
			<ThemeToggle label="Switch to light" onToggle={onToggle} theme="dark" />,
		)

		fireEvent.click(screen.getByRole('button', { name: 'Switch to light' }))

		expect(onToggle).toHaveBeenCalledOnce()
	})
})
