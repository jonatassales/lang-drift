import { z } from 'zod'

import type { EmailErrorCode } from './form.types'

export const earlyAccessEmailSchema = z
	.string({
		error: (issue) =>
			issue.input === undefined ? 'emailRequired' : 'emailInvalid',
	})
	.refine((value) => !/[^\x20-\x7e]/.test(value), { error: 'emailInvalid' })
	.transform((value) => value.replace(/^ +| +$/g, ''))
	.pipe(
		z
			.string()
			.min(1, { error: 'emailRequired' })
			.max(254, { error: 'emailTooLong' }),
	)
	.refine(
		(email) => {
			const parts = email.split('@')
			if (parts.length !== 2) return false
			const [local, domain] = parts
			if (!local || !domain || local.length > 64) return false
			if (
				!/^[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+)*$/.test(
					local,
				)
			)
				return false
			const labels = domain.split('.')

			return (
				labels.length >= 2 &&
				/^[A-Za-z]{2,63}$/.test(labels.at(-1) ?? '') &&
				labels.every((label) =>
					/^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$/.test(label),
				)
			)
		},
		{ error: 'emailInvalid' },
	)
	.transform((email) => {
		const at = email.indexOf('@')

		return email.slice(0, at + 1) + email.slice(at + 1).toLowerCase()
	})

export const earlyAccessFormSchema = z.object({
	email: earlyAccessEmailSchema,
})

export type EarlyAccessFormData = z.output<typeof earlyAccessFormSchema>

export function getEmailError(email: unknown): EmailErrorCode | null {
	const result = earlyAccessEmailSchema.safeParse(email)
	if (result.success) return null

	switch (result.error.issues[0]?.message) {
		case 'emailRequired':
			return 'emailRequired'
		case 'emailTooLong':
			return 'emailTooLong'
		default:
			return 'emailInvalid'
	}
}

export const earlyAccessSubmissionSchema = earlyAccessFormSchema
	.extend({
		locale: z.enum(['en', 'pt-BR', 'zh-Hant', 'ja', 'vi']),
		source: z.enum(['landing', 'pricing']),
		turnstileToken: z.string().min(1).max(2048),
	})
	.strict()
