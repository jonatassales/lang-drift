import { z } from 'zod'

import type { EarlyAccessActionResult } from './form.types'
import { readBoundedBody } from './request.body'

export async function submitEarlyAccess(
	input: object,
	signal?: AbortSignal,
): Promise<EarlyAccessActionResult> {
	const response = await fetch('/api/early-access/', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(input),
		credentials: 'same-origin',
		cache: 'no-store',
		redirect: 'error',
		signal: signal
			? AbortSignal.any([signal, AbortSignal.timeout(10000)])
			: AbortSignal.timeout(10000),
	})

	if (response.status === 202) return { ok: true }
	if (response.status === 429) {
		const seconds = Number(response.headers.get('Retry-After'))

		return {
			ok: false,
			code: 'RATE_LIMITED',
			retryAfter:
				Number.isSafeInteger(seconds) && seconds > 0
					? Math.min(seconds, 2_678_400)
					: 30,
		}
	}
	if (response.status === 403)
		return {
			ok: false,
			code: 'CHALLENGE_FAILED',
		}
	if (response.status === 400) {
		try {
			const text = await readBoundedBody(
				response.body,
				1024,
				AbortSignal.timeout(1000),
			)

			const parsed = z
				.object({
					fieldErrors: z.object({
						email: z.enum(['emailRequired', 'emailInvalid', 'emailTooLong']),
					}),
				})
				.safeParse(JSON.parse(text))

			return {
				ok: false,
				code: 'VALIDATION_ERROR',
				...(parsed.success ? { fieldErrors: parsed.data.fieldErrors } : {}),
			}
		} catch {
			return {
				ok: false,
				code: 'VALIDATION_ERROR',
			}
		}
	}

	return {
		ok: false,
		code: 'UNAVAILABLE',
	}
}
