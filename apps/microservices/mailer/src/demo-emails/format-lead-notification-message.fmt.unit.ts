import { describe, expect, it } from 'vitest'

import { demoEmailsMessages } from '@app/i18n/messages/index.js'

import { formatLeadNotificationMessage } from './format-lead-notification-message.fmt.js'

const submission = {
	email: 'ana@acme.com',
	company: 'Acme',
	position: 'cto',
	locale: 'pt-BR',
	submittedAt: '2026-11-03T14:05:00.000Z',
} as const

describe('mailer/demo-emails/format-lead-notification-message.fmt', () => {
	it('tells the team everything the visitor submitted, in the visitor locale', () => {
		const message = formatLeadNotificationMessage(
			demoEmailsMessages,
			'pt-BR',
			submission,
		)

		expect(message.subject).toBe('[LangDrift] Nova solicitação de demo: Acme')
		expect(message.text).toBe(
			[
				'Nova solicitação de demo enviada em langdrift.md:',
				'',
				'Empresa: Acme',
				'E-mail de contato: ana@acme.com',
				'Cargo: cto',
				'Idioma: pt-BR',
				'Enviada em: 2026-11-03T14:05:00.000Z',
			].join('\n'),
		)
	})

	it('writes English for the en locale', () => {
		const message = formatLeadNotificationMessage(demoEmailsMessages, 'en', {
			...submission,
			locale: 'en',
		})

		expect(message.subject).toBe('[LangDrift] New Demo Request: Acme')
		expect(message.text).toContain('Contact Email: ana@acme.com')
	})
})
