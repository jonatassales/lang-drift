// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'

vi.mock('@i18n/navigation', () => ({
	Link: (props: { children: ReactNode; href: string }) => (
		<a href={props.href}>{props.children}</a>
	),
}))

import { bookDemoMessages } from '@i18n/messages/index'

import { BookDemoBand } from './band'

afterEach(() => {
	cleanup()
})

describe('the Book a demo band', () => {
	it('invites to the page with its own sentence', () => {
		render(<BookDemoBand locale="pt-BR" />)

		const copy = bookDemoMessages['pt-BR'].band

		expect(screen.getByRole('heading', { name: copy.title })).toBeTruthy()
		expect(screen.getByText(copy.description)).toBeTruthy()
		expect(
			screen.getByRole('link', { name: copy.cta }).getAttribute('href'),
		).toBe('/book-a-demo')
	})

	it('shows the sentence a product page passes instead of its own', () => {
		render(
			<BookDemoBand
				description="See a team envelope for your own work."
				locale="en"
			/>,
		)

		expect(
			screen.getByText('See a team envelope for your own work.'),
		).toBeTruthy()
		expect(screen.queryByText(bookDemoMessages.en.band.description)).toBeNull()
	})
})
