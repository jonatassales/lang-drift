import { earlyAccessMode } from './early-access-mode.compute'

export interface TurnstileServerConfig {
	mode: 'production' | 'development' | 'test'
	turnstileSecret: string
	turnstileHostname: string
	turnstileVerifyUrl: string
}

export function readTurnstile(): TurnstileServerConfig {
	const mode = earlyAccessMode()
	const turnstileSecret = process.env.TURNSTILE_SECRET_KEY
	if (!turnstileSecret) {
		throw new Error('TURNSTILE_SECRET_KEY is not configured')
	}

	const turnstileSiteKey = process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
	if (!turnstileSiteKey) {
		throw new Error('NEXT_PUBLIC_TURNSTILE_SITE_KEY is not configured')
	}

	const turnstileHostname = process.env.TURNSTILE_EXPECTED_HOSTNAME
	if (
		!turnstileHostname ||
		!/^[a-zA-Z0-9](?:[a-zA-Z0-9.-]*[a-zA-Z0-9])?$/.test(turnstileHostname)
	) {
		throw new Error(
			'TURNSTILE_EXPECTED_HOSTNAME must be a hostname without a scheme, port or path',
		)
	}

	const turnstileVerifyUrl = process.env.TURNSTILE_VERIFY_URL
	if (!turnstileVerifyUrl) {
		throw new Error('TURNSTILE_VERIFY_URL is not configured')
	}

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
