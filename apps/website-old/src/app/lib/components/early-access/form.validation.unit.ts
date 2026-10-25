import { describe, expect, it } from 'vitest'

import fixtures from '@event-testdata/registration.json'

import { earlyAccessSubmissionSchema, getEmailError } from './form.validation'

describe('early access submission', () => {
	it.each(fixtures)(
		'shares Go validation and identity for $email',
		(fixture) => {
			const result = earlyAccessSubmissionSchema.safeParse({
				email: fixture.email,
				locale: 'en',
				source: 'landing',
				turnstileToken: 'token',
			})
			expect(result.success).toBe(fixture.valid)
			if (result.success) expect(result.data.email).toBe(fixture.normalized)
		},
	)
	it('accepts known locales and sources with a bounded token', () => {
		for (const locale of ['en', 'pt-BR', 'zh-Hant', 'ja', 'vi']) {
			for (const source of ['landing', 'pricing']) {
				const result = earlyAccessSubmissionSchema.safeParse({
					email: ' a@example.com ',
					locale,
					source,
					turnstileToken: 'token',
				})
				expect(result.success).toBe(true)
				if (result.success) expect(result.data.email).toBe('a@example.com')
			}
		}
	})

	it('rejects unknown fields, invalid email and oversized tokens', () => {
		const valid = {
			email: 'a@example.com',
			locale: 'en',
			source: 'landing',
			turnstileToken: 'token',
		}
		expect(
			earlyAccessSubmissionSchema.safeParse({
				...valid,
				extra: true,
			}).success,
		).toBe(false)
		expect(
			earlyAccessSubmissionSchema.safeParse({
				...valid,
				email: 'invalid',
			}).success,
		).toBe(false)
		expect(
			earlyAccessSubmissionSchema.safeParse({
				...valid,
				turnstileToken: 'x'.repeat(2049),
			}).success,
		).toBe(false)
		expect(getEmailError('')).toBe('emailRequired')
		expect(getEmailError('x'.repeat(255))).toBe('emailTooLong')
		expect(getEmailError(undefined)).toBe('emailRequired')
		expect(getEmailError('a@example.com')).toBeNull()
	})
})
