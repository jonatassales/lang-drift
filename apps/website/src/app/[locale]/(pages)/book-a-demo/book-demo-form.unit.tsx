// @vitest-environment jsdom
import { cleanup, fireEvent, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { bookDemoCopy } from '@routes/book-a-demo/testdata/book-demo-copy.data'
import { fillBookDemoForm } from '@routes/book-a-demo/testdata/fill-book-demo-form.mock'
import { mountBookDemoForm } from '@routes/book-a-demo/testdata/mount-book-demo-form.mock'
import { submitBookDemoForm } from '@routes/book-a-demo/testdata/submit-book-demo-form.mock'

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

describe('website/book-a-demo/book-demo-form', () => {
	it('asks for the three fields and sends nothing when they are empty', async () => {
		mountBookDemoForm()
		submitBookDemoForm()

		expect(
			await screen.findByText(bookDemoCopy.errors.emailRequired),
		).toBeTruthy()

		expect(screen.getByText(bookDemoCopy.errors.companyRequired)).toBeTruthy()
		expect(screen.getByText(bookDemoCopy.errors.positionRequired)).toBeTruthy()
		expect(
			screen
				.getByLabelText(bookDemoCopy.form.email)
				.getAttribute('aria-invalid'),
		).toBe('true')
		expect(transport).not.toHaveBeenCalled()
	})

	it('refuses a free-mail address as such', async () => {
		mountBookDemoForm()
		fillBookDemoForm({
			company: 'Acme',
			email: 'jane@gmail.com',
			position: 'cto',
		})
		submitBookDemoForm()

		expect(
			await screen.findByText(bookDemoCopy.errors.emailFreeMail),
		).toBeTruthy()

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
		mountBookDemoForm()
		fillBookDemoForm({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submitBookDemoForm()

		expect(
			await screen.findByRole('heading', {
				name: bookDemoCopy.form.sent.title,
			}),
		).toBeTruthy()

		expect(screen.getByText(bookDemoCopy.form.sent.description)).toBeTruthy()
		expect(screen.queryByLabelText(bookDemoCopy.form.email)).toBeNull()
	})

	it('posts a valid request with the locale and shows the return-soon notice when requests are paused', async () => {
		mountBookDemoForm()
		fillBookDemoForm({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submitBookDemoForm()

		expect(
			await screen.findByRole('heading', {
				name: bookDemoCopy.form.unavailable.title,
			}),
		).toBeTruthy()

		expect(
			screen.getByText(bookDemoCopy.form.unavailable.description),
		).toBeTruthy()
		expect(screen.queryByLabelText(bookDemoCopy.form.email)).toBeNull()
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
		mountBookDemoForm()
		fillBookDemoForm({
			company: 'Acme',
			email: 'jane@acme.io',
			position: 'cto',
		})
		submitBookDemoForm()

		const toast = await screen.findByRole('alert')

		expect(toast.textContent).toContain(bookDemoCopy.form.failed)

		await waitFor(() => {
			expect(
				(screen.getByLabelText(bookDemoCopy.form.email) as HTMLInputElement)
					.value,
			).toBe('jane@acme.io')
		})
		expect(
			screen.queryByRole('heading', {
				name: bookDemoCopy.form.unavailable.title,
			}),
		).toBeNull()

		const closeButton = screen.getByRole('button', {
			name: 'Close notification',
		})
		fireEvent.click(closeButton)
		expect(screen.queryByRole('alert')).toBeNull()
	})
})
