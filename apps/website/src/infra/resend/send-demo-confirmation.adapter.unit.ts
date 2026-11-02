import type { Resend } from 'resend'
import { describe, expect, it, vi } from 'vitest'

import { sendDemoConfirmation } from './send-demo-confirmation.adapter'

describe('website/infra/resend/send-demo-confirmation.adapter', () => {
	it('throws when RESEND_API_KEY is missing and no client is provided', async () => {
		const originalKey = process.env.RESEND_API_KEY
		delete process.env.RESEND_API_KEY

		try {
			await expect(
				sendDemoConfirmation({
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

	it('sends confirmation email with English content for en locale', async () => {
		const sendMock = vi.fn().mockResolvedValue({
			data: { id: 'email_123' },
			error: null,
		})
		const clientMock = {
			emails: {
				send: sendMock,
			},
		} as unknown as Resend

		const result = await sendDemoConfirmation(
			{
				company: 'Acme Corp',
				email: 'lead@acme.com',
				locale: 'en',
				position: 'cto',
			},
			clientMock,
		)

		expect(result).toEqual({ id: 'email_123' })
		expect(sendMock).toHaveBeenCalledWith(
			expect.objectContaining({
				from: 'LangDrift <no-reply@langdrift.md>',
				subject: 'We received your demo request — LangDrift',
				to: 'lead@acme.com',
			}),
		)
	})

	it('sends confirmation email with Portuguese content for pt-BR locale', async () => {
		const sendMock = vi.fn().mockResolvedValue({
			data: { id: 'email_456' },
			error: null,
		})
		const clientMock = {
			emails: {
				send: sendMock,
			},
		} as unknown as Resend

		const result = await sendDemoConfirmation(
			{
				company: 'Empresa Teste',
				email: 'contato@empresa.com.br',
				locale: 'pt-BR',
				position: 'founder',
			},
			clientMock,
		)

		expect(result).toEqual({ id: 'email_456' })
		expect(sendMock).toHaveBeenCalledWith(
			expect.objectContaining({
				from: 'LangDrift <no-reply@langdrift.md>',
				subject: 'Recebemos sua solicitação de demo — LangDrift',
				to: 'contato@empresa.com.br',
			}),
		)
	})

	it('throws error when Resend API returns an error', async () => {
		const sendMock = vi.fn().mockResolvedValue({
			data: null,
			error: { message: 'Invalid API key' },
		})
		const clientMock = {
			emails: {
				send: sendMock,
			},
		} as unknown as Resend

		await expect(
			sendDemoConfirmation(
				{
					company: 'Acme Corp',
					email: 'lead@acme.com',
					locale: 'en',
					position: 'cto',
				},
				clientMock,
			),
		).rejects.toThrow('Invalid API key')
	})
})
