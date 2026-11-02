import 'server-only'

import { randomUUID } from 'node:crypto'
import { z } from 'zod'

import { declaresOversizedBody } from '@components/early-access/declares-oversized-body.validation'
import {
	earlyAccessSubmissionSchema,
	getEmailError,
} from '@components/early-access/form.validation'
import { isPlainJSON } from '@components/early-access/is-plain-json.validation'
import { parseSubmissionJSON } from '@components/early-access/parse-submission-json.compute'
import { readRegistry } from '@components/early-access/read-registry.compute'
import { reply } from '@components/early-access/reply.response'
import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'
import { retrySeconds } from '@components/early-access/retry-seconds.compute'
import { trustedOrigin } from '@components/early-access/trusted-origin.compute'
import { verifyEarlyAccessChallenge } from '@components/early-access/turnstile.validation'

export const runtime = 'nodejs'

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
			void admission.body?.cancel().catch(() => undefined)

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
			void registered.body?.cancel().catch(() => undefined)

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

// Sign-up is paused until feature flags have a proper home. To open it,
// replace the paused response below with postEarlyAccess(request).
export async function POST(): Promise<Response> {
	return Response.json(
		{
			ok: false,
			code: 'UNAVAILABLE',
		},
		{
			status: 503,
			headers: {
				'Cache-Control': 'no-store',
				'Retry-After': '30',
				'X-Content-Type-Options': 'nosniff',
			},
		},
	)
}
