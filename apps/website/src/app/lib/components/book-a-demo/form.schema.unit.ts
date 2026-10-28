import { describe, expect, it } from 'vitest'

import { bookDemoMessages } from '@i18n/messages/index'

import {
	bookDemoErrorCodes,
	bookDemoSchema,
	bookDemoSubmissionSchema,
	fieldErrorsOf,
} from './form.schema'
import { positionIds } from './positions'

const valid = {
	company: 'Acme',
	email: 'jane@acme.io',
	position: 'cto',
}

function errorsOf(input: Record<string, unknown>) {
	const result = bookDemoSchema.safeParse(input)
	if (result.success) return {}

	return fieldErrorsOf(result.error)
}

describe('a demo request', () => {
	it('accepts a company address, a company and a position, trimmed and with the domain lowercased', () => {
		const result = bookDemoSchema.safeParse({
			company: '  Acme  ',
			email: 'Jane@ACME.io',
			position: 'cto',
		})

		expect(result.success).toBe(true)
		expect(result.data).toEqual({
			company: 'Acme',
			email: 'Jane@acme.io',
			position: 'cto',
		})
	})

	it.each(positionIds)('accepts the position %s', (position) => {
		expect(
			bookDemoSchema.safeParse({
				...valid,
				position,
			}).success,
		).toBe(true)
	})

	it.each([
		['', 'emailRequired'],
		['   ', 'emailRequired'],
		['jane', 'emailInvalid'],
		['jane@', 'emailInvalid'],
		[`${'j'.repeat(64)}@${'a'.repeat(187)}.io`, 'emailTooLong'],
		['jane@gmail.com', 'emailFreeMail'],
		['Jane@GMAIL.COM', 'emailFreeMail'],
		['jane@yahoo.co.uk', 'emailFreeMail'],
		['jane@proton.me', 'emailFreeMail'],
		['jane@outlook.com', 'emailFreeMail'],
	])('refuses the address %j with %s', (email, code) => {
		expect(
			errorsOf({
				...valid,
				email,
			}),
		).toEqual({ email: code })
	})

	it.each([
		['', 'companyRequired'],
		['   ', 'companyRequired'],
		['x'.repeat(121), 'companyTooLong'],
	])('refuses the company %j with %s', (company, code) => {
		expect(
			errorsOf({
				...valid,
				company,
			}),
		).toEqual({ company: code })
	})

	it.each([
		['', 'positionRequired'],
		['monarch', 'positionRequired'],
	])('refuses the position %j with %s', (position, code) => {
		expect(
			errorsOf({
				...valid,
				position,
			}),
		).toEqual({ position: code })
	})

	it('reports one code per field when every field fails', () => {
		expect(
			errorsOf({
				company: '',
				email: 'jane@gmail.com',
				position: '',
			}),
		).toEqual({
			company: 'companyRequired',
			email: 'emailFreeMail',
			position: 'positionRequired',
		})
	})

	it('reports the missing fields of an empty submission', () => {
		expect(errorsOf({})).toEqual({
			company: 'companyRequired',
			email: 'emailRequired',
			position: 'positionRequired',
		})
	})
})

describe('a demo submission', () => {
	it('carries the locale and nothing else', () => {
		expect(
			bookDemoSubmissionSchema.safeParse({
				...valid,
				locale: 'pt-BR',
			}).success,
		).toBe(true)
		expect(
			bookDemoSubmissionSchema.safeParse({
				...valid,
				locale: 'fr',
			}).success,
		).toBe(false)
		expect(
			bookDemoSubmissionSchema.safeParse({
				...valid,
				locale: 'en',
				source: 'landing',
			}).success,
		).toBe(false)
		expect(bookDemoSubmissionSchema.safeParse(valid).success).toBe(false)
	})
})

describe('the codes and the catalog', () => {
	it('names every error code in the English errors, and nothing else', () => {
		expect([...bookDemoErrorCodes].sort()).toEqual(
			Object.keys(bookDemoMessages.en.errors).sort(),
		)
	})

	it('names every position in the English positions, and nothing else', () => {
		expect([...positionIds].sort()).toEqual(
			Object.keys(bookDemoMessages.en.positions).sort(),
		)
	})
})
