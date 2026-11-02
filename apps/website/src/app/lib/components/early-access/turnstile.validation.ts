import 'server-only'

import { z } from 'zod'

import { readTurnstile } from './read-turnstile.compute'
import { readBoundedBody } from './request.body'

const siteverifySchema = z.object({
	action: z.string(),
	hostname: z.string(),
	success: z.literal(true),
})
const testSecrets = new Set([
	'1x0000000000000000000000000000000AA',
	'2x0000000000000000000000000000000AA',
	'3x0000000000000000000000000000000AA',
])

export async function verifyEarlyAccessChallenge(
	token: string,
	parentSignal?: AbortSignal,
): Promise<'verified' | 'rejected' | 'unavailable'> {
	const invalidToken =
		typeof token !== 'string' || !token || token.length > 2048
	if (invalidToken) return 'rejected'

	const signal = parentSignal
		? AbortSignal.any([parentSignal, AbortSignal.timeout(2500)])
		: AbortSignal.timeout(2500)

	try {
		const earlyAccessServerEnv = readTurnstile()
		const response = await fetch(earlyAccessServerEnv.turnstileVerifyUrl, {
			body: JSON.stringify({
				response: token,
				secret: earlyAccessServerEnv.turnstileSecret,
			}),
			cache: 'no-store',
			headers: { 'Content-Type': 'application/json' },
			method: 'POST',
			redirect: 'error',
			signal,
		})

		if (!response.ok) {
			void response.body?.cancel().catch(() => undefined)

			return 'unavailable'
		}

		const text = await readBoundedBody(response.body, 4096, signal)

		const body: unknown = JSON.parse(text)

		if (z.object({ success: z.literal(false) }).safeParse(body).success)
			return 'rejected'

		const parsed = siteverifySchema.safeParse(body)
		if (!parsed.success) return 'unavailable'

		const result = parsed.data
		const action =
			earlyAccessServerEnv.mode !== 'production' &&
			testSecrets.has(earlyAccessServerEnv.turnstileSecret)
				? 'test'
				: 'early-access'

		const verified =
			result.success === true &&
			result.hostname === earlyAccessServerEnv.turnstileHostname &&
			result.action === action

		return verified ? 'verified' : 'rejected'
	} catch {
		return 'unavailable'
	}
}
