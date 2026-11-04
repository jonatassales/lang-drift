import { describe, expect, it } from 'vitest'

import { demoEmailsMessages } from '@app/i18n/messages/index.js'

import { formatConfirmationMessage } from './format-confirmation-message.fmt.js'

describe('mailer/demo-emails/format-confirmation-message.fmt', () => {
	it('writes the confirmation in the request locale with the company named', () => {
		const message = formatConfirmationMessage(demoEmailsMessages, 'pt-BR', {
			company: 'Acme',
		})

		expect(message.subject).toBe(
			'Recebemos sua solicitação de demo — LangDrift',
		)
		expect(message.text).toBe(
			[
				'Olá,',
				'',
				'Obrigado por solicitar uma demonstração do LangDrift para a Acme.',
				'',
				'Nossa equipe entrará em contato em breve para agendar a sua apresentação.',
				'',
				'Atenciosamente,',
				'Equipe LangDrift',
				'https://langdrift.md',
			].join('\n'),
		)
	})

	it('writes English for the en locale', () => {
		const message = formatConfirmationMessage(demoEmailsMessages, 'en', {
			company: 'Acme',
		})

		expect(message.subject).toBe('We received your demo request — LangDrift')
		expect(message.text).toContain(
			'Thank you for requesting a demo of LangDrift for Acme.',
		)
	})
})
