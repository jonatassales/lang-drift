import { describe, expect, it } from 'vitest'

import { demoRequestSubmissionSchema } from './demo-request-submission.validation'

const submission = {
	email: 'ana@acme.com',
	company: 'Acme',
	position: 'cto',
	locale: 'pt-BR',
}

describe('events/events/demo-request-submission.validation', () => {
	it('accepts a company email, a company, a position and a locale', () => {
		const result = demoRequestSubmissionSchema.safeParse(submission)

		expect(result.success).toBe(true)
	})

	it('trims the company name', () => {
		const result = demoRequestSubmissionSchema.parse({
			...submission,
			company: '  Acme  ',
		})

		expect(result.company).toBe('Acme')
	})

	it('refuses an address that is not an email', () => {
		const result = demoRequestSubmissionSchema.safeParse({
			...submission,
			email: 'ana at acme',
		})

		expect(result.success).toBe(false)
	})

	it('refuses a company name over 120 characters', () => {
		const result = demoRequestSubmissionSchema.safeParse({
			...submission,
			company: 'a'.repeat(121),
		})

		expect(result.success).toBe(false)
	})

	it('refuses a position the form does not offer', () => {
		const result = demoRequestSubmissionSchema.safeParse({
			...submission,
			position: 'intern',
		})

		expect(result.success).toBe(false)
	})

	it('refuses a locale the site does not serve', () => {
		const result = demoRequestSubmissionSchema.safeParse({
			...submission,
			locale: 'zh-Hans',
		})

		expect(result.success).toBe(false)
	})

	it('refuses a field it does not know', () => {
		const result = demoRequestSubmissionSchema.safeParse({
			...submission,
			phone: '+55 11 99999-0000',
		})

		expect(result.success).toBe(false)
	})
})
