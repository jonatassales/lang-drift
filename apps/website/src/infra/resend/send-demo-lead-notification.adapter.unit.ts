import type { Resend } from 'resend'
import { describe, expect, it, vi } from 'vitest'

import { sendDemoLeadNotification } from './send-demo-lead-notification.adapter'

describe('website/infra/resend/send-demo-lead-notification.adapter', () => {
	it('throws when RESEND_API_KEY is missing and no client is provided', async () => {
		const originalKey = process.env.RESEND_API_KEY
		delete process.env.RESEND_API_KEY

		try {
			await expect(
				sendDemoLeadNotification({
					company: 'Acme Corp',
					email: 'lead@acme.com',
					locale: 'en',
					position: 'cto',
				}),
			).rejects.toThrow('Missing RESEND_API_KEY')
		} finally {
			process.env.RESEND_API_KEY = originalKey
		}
	})

	it('sends internal notification email with lead details to jonatas.langdrift@gmail.com', async () => {
		const sendMock = vi.fn().mockResolvedValue({
			data: { id: 'email_lead_999' },
			error: null,
		})
		const clientMock = {
			emails: {
				send: sendMock,
			},
		} as unknown as Resend

		const result = await sendDemoLeadNotification(
			{
				company: 'Acme Corp',
				email: 'lead@acme.com',
				locale: 'en',
				position: 'cto',
				submittedAt: '2026-11-02T10:00:00.000Z',
			},
			clientMock,
		)

		expect(result).toEqual({ id: 'email_lead_999' })
		expect(sendMock).toHaveBeenCalledWith(
			expect.objectContaining({
				from: 'LangDrift <no-reply@langdrift.md>',
				subject: '[LangDrift] New Demo Request: Acme Corp',
				to: 'jonatas.langdrift@gmail.com',
			}),
		)
		const sentBody = sendMock.mock.calls[0][0].text
		expect(sentBody).toContain('Company: Acme Corp')
		expect(sentBody).toContain('Contact Email: lead@acme.com')
		expect(sentBody).toContain('Position: cto')
		expect(sentBody).toContain('Locale: en')
		expect(sentBody).toContain('Submitted At: 2026-11-02T10:00:00.000Z')
	})

	it('throws error when Resend API returns an error', async () => {
		const sendMock = vi.fn().mockResolvedValue({
			data: null,
			error: { message: 'Quota exceeded' },
		})
		const clientMock = {
			emails: {
				send: sendMock,
			},
		} as unknown as Resend

		await expect(
			sendDemoLeadNotification(
				{
					company: 'Acme Corp',
					email: 'lead@acme.com',
					locale: 'en',
					position: 'cto',
				},
				clientMock,
			),
		).rejects.toThrow('Quota exceeded')
	})
})
