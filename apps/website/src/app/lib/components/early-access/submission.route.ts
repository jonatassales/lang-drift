import 'server-only'

import { randomUUID } from 'node:crypto'
import { z } from 'zod'

import { earlyAccessSubmissionSchema, getEmailError } from './form.validation'
import { BodyError, readBoundedBody } from './request.body'
import {
	declaresOversizedBody,
	isPlainJSON,
	parseSubmissionJSON,
} from './request.validation'
import { verifyEarlyAccessChallenge } from './turnstile.validation'

function reply(
	status: number,
	code: string,
	retryAfter?: number,
	fieldErrors?: { email: string },
) {
	const headers: Record<string, string> = {
		'Cache-Control': 'no-store',
		'X-Content-Type-Options': 'nosniff',
	}
	if (retryAfter) headers['Retry-After'] = String(retryAfter)

	return Response.json(
		{
			ok: status === 202,
			code,
			...(fieldErrors ? { fieldErrors } : {}),
		},
		{
			status,
			headers,
		},
	)
}

type EarlyAccessMode = 'production' | 'development' | 'test'

function earlyAccessMode(): EarlyAccessMode {
	if (process.env.NODE_ENV === 'development') return 'development'
	if (process.env.EARLY_ACCESS_MODE === 'test') return 'test'
	return 'production'
}

function assertCredentialFreeOrigin(endpoint: URL, mode: EarlyAccessMode) {
	if (
		endpoint.username ||
		endpoint.password ||
		endpoint.search ||
		endpoint.hash ||
		endpoint.pathname !== '/' ||
		!['http:', 'https:'].includes(endpoint.protocol)
	) {
		throw new Error('Early Access endpoints must be credential-free origins')
	}
	if (mode === 'production' && endpoint.protocol !== 'https:') {
		throw new Error('Production Early Access endpoints require HTTPS')
	}
	if (
		mode !== 'production' &&
		!['localhost', '127.0.0.1', 'email-store'].includes(endpoint.hostname)
	) {
		throw new Error('Development Early Access endpoints must be local')
	}
}

function readRegistry() {
	const mode = earlyAccessMode()
	const registry = new URL(process.env.EARLY_ACCESS_REGISTRY_URL ?? '')
	const website = new URL(process.env.EARLY_ACCESS_PUBLIC_ORIGIN ?? '')
	assertCredentialFreeOrigin(registry, mode)
	assertCredentialFreeOrigin(website, mode)
	const registryKey = process.env.EARLY_ACCESS_REGISTRY_KEY ?? ''
	if (registryKey.length < 32 || /[\r\n]/.test(registryKey)) {
		throw new Error(
			'EARLY_ACCESS_REGISTRY_KEY must contain at least 32 characters',
		)
	}

	return {
		mode,
		registryUrl: registry.origin,
		registryKey,
		publicOrigin: website.origin,
	}
}

function retrySeconds(response: Response): number {
	const seconds = Number(response.headers.get('Retry-After'))

	return Number.isSafeInteger(seconds) && seconds > 0
		? Math.min(seconds, 2_678_400)
		: 30
}

// Production fails closed until a trusted proxy contract for the host is
// implemented and verified; client-supplied forwarding headers never count.
function trustedOrigin(mode: string): string | null {
	if (mode !== 'production') return '127.0.0.1'

	return null
}

export async function postEarlyAccess(request: Request): Promise<Response> {
	try {
		const env = readRegistry()
		const target = new URL(request.url)
		const canonical = new URL(env.publicOrigin)
		if (request.method !== 'POST') return reply(405, 'INVALID_REQUEST')
		if (
			target.origin !== canonical.origin ||
			request.headers.get('host') !== canonical.host ||
			request.headers.get('origin') !== canonical.origin ||
			request.headers.get('sec-fetch-site') === 'cross-site'
		) {
			return reply(403, 'INVALID_ORIGIN')
		}
		if (!isPlainJSON(request)) return reply(415, 'INVALID_REQUEST')
		if (declaresOversizedBody(request, 4096))
			return reply(413, 'INVALID_REQUEST')
		const origin = trustedOrigin(env.mode)
		if (!origin) return reply(503, 'UNAVAILABLE', 30)
		const signal = AbortSignal.any([request.signal, AbortSignal.timeout(8000)])
		let input: unknown
		try {
			const body = await readBoundedBody(request.body, 4096, signal)

			input = parseSubmissionJSON(body)
		} catch (error) {
			return reply(
				error instanceof BodyError ? error.status : 400,
				'VALIDATION_ERROR',
			)
		}
		const parsed = earlyAccessSubmissionSchema.safeParse(input)
		if (!parsed.success) {
			const email =
				typeof input === 'object' && input !== null && 'email' in input
					? input.email
					: undefined
			const error = getEmailError(email)

			return reply(
				400,
				'VALIDATION_ERROR',
				undefined,
				error ? { email: error } : undefined,
			)
		}
		const { email, locale, source, turnstileToken } = parsed.data
		const requestID = randomUUID()
		const registry = async (path: string, data: object) => {
			const payload = JSON.stringify(data)
			if (Buffer.byteLength(payload) > 1024) throw new BodyError(413)

			return fetch(`${env.registryUrl}${path}`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${env.registryKey}`,
				},
				body: payload,
				cache: 'no-store',
				redirect: 'error',
				signal: AbortSignal.any([signal, AbortSignal.timeout(3500)]),
			})
		}
		const admission = await registry('/v1/admissions', {
			email,
			origin,
			request: requestID,
		})

		if (admission.status !== 200) {
			void admission.body?.cancel().catch(() => {})

			return reply(
				admission.status === 429 ? 429 : 503,
				admission.status === 429 ? 'RATE_LIMITED' : 'UNAVAILABLE',
				retrySeconds(admission),
			)
		}
		const grantText = await readBoundedBody(admission.body, 4096, signal)

		const grant = z
			.object({ grant: z.string().regex(/^[a-f0-9]{32}$/) })
			.strict()
			.safeParse(JSON.parse(grantText))
		if (!grant.success) return reply(503, 'UNAVAILABLE', 30)

		const verified = await verifyEarlyAccessChallenge(turnstileToken, signal)

		if (verified === 'unavailable') return reply(503, 'UNAVAILABLE', 30)
		if (verified === 'rejected') return reply(403, 'CHALLENGE_FAILED')
		signal.throwIfAborted()

		const registered = await registry('/v1/registrations', {
			email,
			locale,
			source,
			grant: grant.data.grant,
			request: requestID,
		})

		if (registered.status !== 202) {
			void registered.body?.cancel().catch(() => {})

			return reply(
				registered.status === 429 ? 429 : 503,
				registered.status === 429 ? 'RATE_LIMITED' : 'UNAVAILABLE',
				retrySeconds(registered),
			)
		}
		const receiptText = await readBoundedBody(registered.body, 4096, signal)

		const receipt = z
			.object({ status: z.literal('accepted') })
			.strict()
			.safeParse(JSON.parse(receiptText))

		return receipt.success
			? reply(202, 'ACCEPTED')
			: reply(503, 'UNAVAILABLE', 30)
	} catch {
		return reply(503, 'UNAVAILABLE', 30)
	}
}
