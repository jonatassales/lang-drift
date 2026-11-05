import { z } from 'zod'

// A deliverable address: printable ASCII, trimmed of spaces, at most 254
// characters, a dot-atom local part of at most 64 characters and a domain of
// valid labels with an alphabetic top level. The domain is lowercased; the
// local part keeps its case.
const dotAtomLocal =
	/^[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+)*$/
const domainLabel = /^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$/
const topLevelLabel = /^[A-Za-z]{2,63}$/

export const emailSchema = z
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

			const [local = '', domain = ''] = parts
			const localFits = local.length > 0 && local.length <= 64

			if (!localFits || !domain) return false
			if (!dotAtomLocal.test(local)) return false

			const labels = domain.split('.')
			const hasTopLevel =
				labels.length >= 2 && topLevelLabel.test(labels.at(-1) ?? '')

			return hasTopLevel && labels.every((label) => domainLabel.test(label))
		},
		{ error: 'emailInvalid' },
	)
	.transform((email) => {
		const at = email.indexOf('@')

		return email.slice(0, at + 1) + email.slice(at + 1).toLowerCase()
	})
