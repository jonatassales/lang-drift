import { afterEach, describe, expect, it, vi } from 'vitest'

import { createEmailSender } from './email-sender.factory.js'

describe('mailer/demo-emails/email-sender.factory', () => {
	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('throws when RESEND_API_KEY is not set', () => {
		vi.stubEnv('RESEND_API_KEY', '')

		expect(() => createEmailSender()).toThrow(/RESEND_API_KEY/)
	})

	it('builds a sender from the key without calling Resend yet', () => {
		vi.stubEnv('RESEND_API_KEY', 're_test_key')

		expect(typeof createEmailSender().send).toBe('function')
	})
})
