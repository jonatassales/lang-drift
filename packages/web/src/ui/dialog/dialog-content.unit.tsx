import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { Dialog } from './dialog'
import { DialogContent } from './dialog-content'
import { DialogTitle } from './dialog-title'

afterEach(cleanup)

describe('web/ui/dialog/dialog-content', () => {
	it('names the dialog with its title', () => {
		render(
			<Dialog open>
				<DialogContent>
					<DialogTitle>Choose a folder</DialogTitle>
				</DialogContent>
			</Dialog>,
		)

		expect(screen.getByRole('dialog', { name: 'Choose a folder' })).toBeTruthy()
	})

	it('shows the corner close button only with a label', () => {
		const { rerender } = render(
			<Dialog open>
				<DialogContent aria-label="Secret" />
			</Dialog>,
		)

		expect(screen.queryByRole('button')).toBeNull()

		rerender(
			<Dialog open>
				<DialogContent aria-label="Secret" closeLabel="Close" />
			</Dialog>,
		)

		expect(screen.getByRole('button', { name: 'Close' })).toBeTruthy()
	})
})
