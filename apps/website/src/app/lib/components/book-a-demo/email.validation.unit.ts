import { describe, expect, it } from 'vitest'

import emails from '@components/book-a-demo/testdata/emails.json'

import { emailSchema } from './email.validation'

describe('website/book-a-demo/email.validation', () => {
	it.each(emails)('validates and normalizes $email', (fixture) => {
		const result = emailSchema.safeParse(fixture.email)

		expect(result.success).toBe(fixture.valid)

		if (result.success) expect(result.data).toBe(fixture.normalized)
	})

	it('names why an address is refused', () => {
		const empty = emailSchema.safeParse('')
		const long = emailSchema.safeParse(`${'x'.repeat(250)}@a.io`)
		const missing = emailSchema.safeParse(undefined)

		expect(empty.error?.issues[0]?.message).toBe('emailRequired')
		expect(long.error?.issues[0]?.message).toBe('emailTooLong')
		expect(missing.error?.issues[0]?.message).toBe('emailRequired')
	})
})
