import 'server-only'

import { z } from 'zod'

import { readBoundedBody } from './request.body'

const siteverifySchema = z.object({
	success: z.literal(true),
	hostname: z.string(),
	action: z.string(),
})
const testSecrets = new Set([
	'1x0000000000000000000000000000000AA',
	'2x0000000000000000000000000000000AA',
	'3x0000000000000000000000000000000AA',
])

type EarlyAccessMode = 'production' | 'development' | 'test'

function earlyAccessMode(): EarlyAccessMode {
	const hosted =
		process.env.VERCEL === '1' ||
		['production', 'preview'].includes(process.env.VERCEL_ENV ?? '')
	if (hosted) return 'production'
	if (process.env.NODE_ENV === 'development') return 'development'
	if (process.env.EARLY_ACCESS_MODE === 'test') return 'test'
	return 'production'
}

function readTurnstile() {
	const mode = earlyAccessMode()
	const turnstileSecret = process.env.TURNSTILE_SECRET_KEY ?? ''
	const turnstileSiteKey = process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY ?? ''
	const turnstileHostname = process.env.TURNSTILE_EXPECTED_HOSTNAME ?? ''
	if (!/^[a-zA-Z0-9](?:[a-zA-Z0-9.-]*[a-zA-Z0-9])?$/.test(turnstileHostname)) {
		throw new Error(
			'TURNSTILE_EXPECTED_HOSTNAME must be a hostname without a scheme, port or path',
		)
	}
	const turnstileVerifyUrl = process.env.TURNSTILE_VERIFY_URL ?? ''
	let verifier: URL
	try {
		verifier = new URL(turnstileVerifyUrl)
	} catch {
		throw new Error('TURNSTILE_VERIFY_URL must be a valid HTTP(S) URL')
	}
	if (
		!['http:', 'https:'].includes(verifier.protocol) ||
		verifier.username ||
		verifier.password ||
		verifier.hash ||
		(mode === 'production' &&
			turnstileVerifyUrl !==
				'https://challenges.cloudflare.com/turnstile/v0/siteverify')
	) {
		throw new Error(
			'TURNSTILE_VERIFY_URL must use the official HTTPS endpoint in production and a valid HTTP(S) endpoint in development/test',
		)
	}
	if (
		mode === 'production' &&
		(/^[123]x0{20}/.test(turnstileSecret) ||
			/^[123]x0{18}/.test(turnstileSiteKey))
	) {
		throw new Error('Production requires real Turnstile credentials')
	}

	return {
		mode,
		turnstileSecret,
		turnstileHostname,
		turnstileVerifyUrl,
	}
}

export async function verifyEarlyAccessChallenge(
	token: string,
	parentSignal?: AbortSignal,
): Promise<'verified' | 'rejected' | 'unavailable'> {
	if (typeof token !== 'string' || !token || token.length > 2048)
		return 'rejected'

	const signal = parentSignal
		? AbortSignal.any([parentSignal, AbortSignal.timeout(2500)])
		: AbortSignal.timeout(2500)

	try {
		if (process.env.WEBSITE_PREVIEW_SIGNUP_DISABLED !== 'false') {
			return 'verified'
		}
		const earlyAccessServerEnv = readTurnstile()
		const response = await fetch(earlyAccessServerEnv.turnstileVerifyUrl, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				secret: earlyAccessServerEnv.turnstileSecret,
				response: token,
			}),
			cache: 'no-store',
			redirect: 'error',
			signal,
		})

		if (!response.ok) {
			void response.body?.cancel().catch(() => {})

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
