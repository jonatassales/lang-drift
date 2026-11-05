// @vitest-environment jsdom
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { CopyCommandButton } from './copy-command-button'

const labels = {
	copy: 'Copy command',
	copied: 'Copied',
	failed: 'Copy failed',
}

const command = 'npx @langdrift/cli init'

function stubClipboard(writeText: (text: string) => Promise<void>) {
	Object.defineProperty(navigator, 'clipboard', {
		configurable: true,
		value: {
			writeText,
		},
	})
}

afterEach(() => {
	cleanup()
	vi.useRealTimers()
	Reflect.deleteProperty(navigator, 'clipboard')
})

describe('website/platform/copy-command-button', () => {
	it('copies the command, says so, then offers to copy again', async () => {
		vi.useFakeTimers({
			shouldAdvanceTime: true,
		})

		const writeText = vi.fn(() => Promise.resolve())

		stubClipboard(writeText)
		render(<CopyCommandButton command={command} labels={labels} />)
		fireEvent.click(screen.getByRole('button', { name: 'Copy command' }))

		const copied = await screen.findByRole('button', { name: 'Copied' })

		expect(copied).toBeDefined()
		expect(writeText).toHaveBeenCalledWith(command)

		act(() => {
			vi.advanceTimersByTime(2000)
		})

		expect(screen.getByRole('button', { name: 'Copy command' })).toBeDefined()
	})

	it('says the copy failed when the browser refuses it', async () => {
		stubClipboard(() => Promise.reject(new Error('denied')))
		render(<CopyCommandButton command={command} labels={labels} />)
		fireEvent.click(screen.getByRole('button', { name: 'Copy command' }))

		const failed = await screen.findByRole('button', { name: 'Copy failed' })

		expect(failed).toBeDefined()
	})
})
