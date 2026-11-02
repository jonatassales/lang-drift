import 'server-only'

import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { fieldErrorsOf } from '@components/book-a-demo/field-errors.compute'
import { bookDemoSubmissionSchema } from '@components/book-a-demo/form.validation'
import { isForeignRequest } from '@components/book-a-demo/is-foreign-request.validation'
import { declaresOversizedBody } from '@components/early-access/declares-oversized-body.validation'
import { isPlainJSON } from '@components/early-access/is-plain-json.validation'
import { parseSubmissionJSON } from '@components/early-access/parse-submission-json.compute'
import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'
import { sendDemoConfirmation } from '@infra/resend/send-demo-confirmation.adapter'
import { sendDemoLeadNotification } from '@infra/resend/send-demo-lead-notification.adapter'

export const runtime = 'nodejs'

const bodyLimit = 4096

// A demo request is validated against origin, headers, size and payload schema,
// then dispatches confirmation and notification emails via Resend.
export async function POST(request: Request): Promise<Response> {
	if (request.method !== 'POST') return reply(405, 'INVALID_REQUEST')
	if (isForeignRequest(request)) return reply(403, 'INVALID_ORIGIN')
	if (!isPlainJSON(request)) return reply(415, 'INVALID_REQUEST')
	if (declaresOversizedBody(request, bodyLimit))
		return reply(413, 'INVALID_REQUEST')

	let body: unknown
	try {
		const text = await readBoundedBody(request.body, bodyLimit, request.signal)

		body = parseSubmissionJSON(text)
	} catch (error) {
		const status = error instanceof BodyError ? error.status : 400

		return reply(status, 'VALIDATION_ERROR')
	}

	const result = bookDemoSubmissionSchema.safeParse(body)
	if (!result.success) {
		return reply(400, 'VALIDATION_ERROR', fieldErrorsOf(result.error))
	}

	try {
		await sendDemoConfirmation({
			company: result.data.company,
			email: result.data.email,
			locale: result.data.locale,
			position: result.data.position,
		})

		await sendDemoLeadNotification({
			company: result.data.company,
			email: result.data.email,
			locale: result.data.locale,
			position: result.data.position,
		})

		return reply(200, 'OK')
	} catch {
		return reply(500, 'DELIVERY_FAILED')
	}
}
