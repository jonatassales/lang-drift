import { z } from 'zod'

import { earlyAccessEmailSchema } from '@components/early-access/form.validation'
import { routing } from '@i18n/routing'

import { positionIds } from './positions'

// Every issue's message is one of these codes; the form and the route map
// them to the catalog.
export const bookDemoErrorCodes = [
	'emailRequired',
	'emailInvalid',
	'emailTooLong',
	'emailFreeMail',
	'companyRequired',
	'companyTooLong',
	'positionRequired',
] as const

export type BookDemoErrorCode = (typeof bookDemoErrorCodes)[number]

export function isBookDemoErrorCode(
	value: unknown,
): value is BookDemoErrorCode {
	return bookDemoErrorCodes.some((code) => code === value)
}

// Providers that hand out addresses to anyone; a demo request names a company.
const freeMailDomains = new Set([
	'gmail.com',
	'googlemail.com',
	'outlook.com',
	'hotmail.com',
	'live.com',
	'msn.com',
	'yahoo.com',
	'yahoo.co.uk',
	'yahoo.co.jp',
	'yahoo.com.br',
	'ymail.com',
	'icloud.com',
	'me.com',
	'mac.com',
	'proton.me',
	'protonmail.com',
	'pm.me',
	'aol.com',
	'gmx.com',
	'gmx.de',
	'mail.com',
	'zoho.com',
	'yandex.com',
	'yandex.ru',
	'fastmail.com',
	'hey.com',
	'tuta.com',
	'tutanota.com',
	'qq.com',
	'163.com',
	'126.com',
	'naver.com',
	'daum.net',
	'uol.com.br',
	'bol.com.br',
])

function domainOf(email: string) {
	return email.slice(email.indexOf('@') + 1)
}

// The early-access rules say what a valid address is and lowercase its
// domain; a demo request adds that the domain is not a free-mail provider.
const companyEmailSchema = earlyAccessEmailSchema.refine(
	(email) => !freeMailDomains.has(domainOf(email)),
	{ error: 'emailFreeMail' },
)

export const bookDemoSchema = z.object({
	email: companyEmailSchema,
	company: z
		.string({ error: 'companyRequired' })
		.trim()
		.min(1, { error: 'companyRequired' })
		.max(120, { error: 'companyTooLong' }),
	// The select's empty option is a string the enum refuses.
	position: z
		.string({ error: 'positionRequired' })
		.pipe(z.enum(positionIds, { error: 'positionRequired' })),
})

export type BookDemoInput = z.input<typeof bookDemoSchema>
export type BookDemoData = z.output<typeof bookDemoSchema>
export type BookDemoField = keyof BookDemoData

export const bookDemoSubmissionSchema = bookDemoSchema
	.extend({ locale: z.enum(routing.locales) })
	.strict()

// The first code of each field that failed.
export function fieldErrorsOf(
	error: z.ZodError,
): Partial<Record<BookDemoField, BookDemoErrorCode>> {
	const errors: Partial<Record<BookDemoField, BookDemoErrorCode>> = {}
	for (const issue of error.issues) {
		const field = issue.path[0]
		if (field !== 'email' && field !== 'company' && field !== 'position') {
			continue
		}
		if (errors[field] || !isBookDemoErrorCode(issue.message)) continue
		errors[field] = issue.message
	}

	return errors
}
