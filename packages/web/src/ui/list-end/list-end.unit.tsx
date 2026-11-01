import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ListEnd } from './list-end'

afterEach(cleanup)
afterEach(vi.unstubAllGlobals)

beforeEach(() => {
	vi.stubGlobal(
		'IntersectionObserver',
		class {
			disconnect() {}
			observe() {}
		},
	)
})

describe('web/ui/list-end', () => {
	it('renders nothing once the list has no more pages', () => {
		const { container } = render(
			<ListEnd
				failed={false}
				hasMore={false}
				label="Load more"
				loading={false}
				onEnd={vi.fn()}
			/>,
		)

		expect(container.innerHTML).toBe('')
	})

	it('asks for the next page when its button is pressed', () => {
		const onEnd = vi.fn()

		render(
			<ListEnd
				failed={false}
				hasMore
				label="Load more"
				loading={false}
				onEnd={onEnd}
			/>,
		)

		fireEvent.click(screen.getByRole('button', { name: 'Load more' }))

		expect(onEnd).toHaveBeenCalledOnce()
	})
})
