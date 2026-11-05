// @vitest-environment jsdom
import {
	cleanup,
	fireEvent,
	render,
	screen,
	waitFor,
} from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { bookDemoMessages } from '@i18n/messages/index'

import { BookDemoForm } from './form'

const copy = bookDemoMessages.en
let transport: ReturnType<typeof vi.fn>

beforeEach(() => {
	transport = vi.fn(async () =>
		Response.json(
			{
				code: 'UNAVAILABLE',
				ok: false,
			},
			{
				headers: { 'Retry-After': '30' },
				status: 503,
			},
		),
	)
	vi.stubGlobal('fetch', transport)
})
afterEach(() => {
	cleanup()
	vi.unstubAllGlobals()
})

function mount() {
	return render(<BookDemoForm copy={copy} locale="en" />)
}

function fill(fields: { company?: string; email?: string; position?: string }) {
	if (fields.email !== undefined) {
		fireEvent.change(screen.getByLabelText(copy.form.email), {
			target: { value: fields.email },
		})
	}
	if (fields.company !== undefined) {
		fireEvent.change(screen.getByLabelText(copy.form.company), {
			target: { value: fields.company },
		})
	}
	if (fields.position !== undefined) {
		fireEvent.change(screen.getByLabelText(copy.form.position), {
			target: { value: fields.position },
		})
	}
}

function submit() {
	fireEvent.click(screen.getByRole('button', { name: copy.form.submit }))
}

describe('website/lib/components/book-a-demo/form', () => {
	it('asks for the three fields and sends nothing when they are empty', async () => {
		mount()
		submit()

		expect(await screen.findByText(copy.errors.emailRequired)).toBeTruthy()
		expect(screen.getByText(copy.errors.companyRequired)).toBeTruthy()
		expect(screen.getByText(copy.errors.positionRequired)).toBeTruthy()
		expect(
			screen.getByLabelText(copy.form.email).getAttribute('aria-invalid'),
		).toBe('true')
		expect(transport).not.toHaveBeenCalled()
	})

	it('refuses a free-mail address as such', async () => {
		mount()
		fill({
			company: 'Acme',
			email: 'jane@gmail.com',
			position: 'cto',
		})
		submit()

		expect(await screen.findByText(copy.errors.emailFreeMail)).toBeTruthy()
		expect(transport).not.toHaveBeenCalled()
	})

	it('shows the sent notice in place of the form on successful submission', async () => {
		transport.mockResolvedValue(
			Response.json(
				{
					code: 'OK',
					ok: true,
				},
				{ status: 200 },
			),
		)
		mount()
		fill({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submit()

		expect(
			await screen.findByRole('heading', {
				name: copy.form.sent.title,
			}),
		).toBeTruthy()

		expect(screen.getByText(copy.form.sent.description)).toBeTruthy()
		expect(screen.queryByLabelText(copy.form.email)).toBeNull()
	})

	it('posts a valid request with the locale and shows the return-soon notice when requests are paused', async () => {
		mount()
		fill({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submit()

		expect(
			await screen.findByRole('heading', {
				name: copy.form.unavailable.title,
			}),
		).toBeTruthy()

		expect(screen.getByText(copy.form.unavailable.description)).toBeTruthy()
		expect(screen.queryByLabelText(copy.form.email)).toBeNull()
		expect(transport).toHaveBeenCalledTimes(1)
		const [url, options] = transport.mock.calls[0] as [string, RequestInit]
		expect(url).toBe('/api/book-a-demo/')
		expect(options.method).toBe('POST')
		expect(JSON.parse(String(options.body))).toEqual({
			company: 'Acme',
			email: 'jane@acme.io',
			locale: 'en',
			position: 'cto',
		})
	})

	it('shows an error toast and keeps the form when request fails', async () => {
		transport.mockResolvedValue(
			Response.json(
				{
					code: 'BOOM',
					ok: false,
				},
				{ status: 500 },
			),
		)
		mount()
		fill({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submit()

		const toast = await screen.findByRole('alert')
		expect(toast.textContent).toContain(copy.form.failed)
		await waitFor(() => {
			expect(
				(screen.getByLabelText(copy.form.email) as HTMLInputElement).value,
			).toBe('jane@acme.io')
		})
		expect(
			screen.queryByRole('heading', { name: copy.form.unavailable.title }),
		).toBeNull()

		const closeButton = screen.getByRole('button', {
			name: 'Close notification',
		})
		fireEvent.click(closeButton)
		expect(screen.queryByRole('alert')).toBeNull()
	})
})
