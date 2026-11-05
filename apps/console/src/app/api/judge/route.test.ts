import assert from 'node:assert/strict'
import { afterEach, beforeEach, describe, it, mock } from 'node:test'

import { contractValidator } from '@app/api/judge/testdata/contract-validator.mock'
import { judgeRequest } from '@app/api/judge/testdata/judge-request.mock'
import { judgeTestSecret } from '@app/api/judge/testdata/judge-secret.data'
import { messagesReply } from '@app/api/judge/testdata/messages-reply.mock'
import requestSample from '@langdrift/core/judge/request.sample.json' with {
	type: 'json',
}
import responseSchema from '@langdrift/core/judge/response.schema.json' with {
	type: 'json',
}
import suggestionAnswer from '@langdrift/core/judge/suggestion.answer.sample.json' with {
	type: 'json',
}

import { GET, POST } from './route.ts'

const environment = ['ANTHROPIC_API_KEY', 'JUDGE_MODELS', 'RUNTIME_SECRET']
const saved = new Map(environment.map((key) => [key, process.env[key]]))
const usage = {
	inputTokens: 1840,
	outputTokens: 96,
}
const contractResponse = contractValidator(responseSchema)

describe('console/api/judge/route', () => {
	beforeEach(() => {
		process.env.RUNTIME_SECRET = judgeTestSecret
		process.env.JUDGE_MODELS = 'claude-sonnet-5-5, claude-opus-5-5'
		process.env.ANTHROPIC_API_KEY = 'sk-ant-test-not-a-key'
		mock.method(console, 'info', () => undefined)
		mock.method(console, 'error', () => undefined)
	})

	afterEach(() => {
		for (const [key, value] of saved) {
			if (value === undefined) delete process.env[key]
			else process.env[key] = value
		}

		mock.restoreAll()
	})

	it('answers the connection check with 204 for the secret and 401 without it', () => {
		const connected = GET(judgeRequest({ method: 'GET' }))
		const anonymous = GET(
			judgeRequest({
				headers: { authorization: '' },
				method: 'GET',
			}),
		)
		const wrong = GET(
			judgeRequest({
				headers: { authorization: 'Bearer not-the-secret' },
				method: 'GET',
			}),
		)

		assert.equal(connected.status, 204)
		assert.equal(anonymous.status, 401)
		assert.equal(wrong.status, 401)
	})

	it('refuses a question without the bearer with 401', async () => {
		const response = await POST(
			judgeRequest({ headers: { authorization: '' } }),
		)

		assert.equal(response.status, 401)

		const reply = await response.json()

		assert.deepEqual(reply, { code: 'UNAUTHORIZED' })
	})

	it('answers no one with 503 while the deployment holds no secret', async () => {
		delete process.env.RUNTIME_SECRET

		const response = await POST(judgeRequest())

		assert.equal(response.status, 503)

		const reply = await response.json()

		assert.equal(reply.reason, 'judge-not-configured')
	})

	it('refuses another media type with 415 and a declared length over the limit with 413', async () => {
		const text = await POST(
			judgeRequest({ headers: { 'content-type': 'text/plain' } }),
		)

		const large = await POST(
			judgeRequest({ headers: { 'content-length': '1048577' } }),
		)

		assert.equal(text.status, 415)
		assert.equal(large.status, 413)
	})

	it('names the field at fault in a question outside the contract with 400', async () => {
		const withoutBound = Object.fromEntries(
			Object.entries(requestSample).filter(([key]) => key !== 'maxTokens'),
		)

		const malformed = await POST(judgeRequest({ body: '{"kind": ' }))

		const unbounded = await POST(
			judgeRequest({ body: JSON.stringify(withoutBound) }),
		)

		const extra = await POST(
			judgeRequest({
				body: JSON.stringify({
					...requestSample,
					temperature: 0,
				}),
			}),
		)

		assert.equal(malformed.status, 400)

		const malformedReply = await malformed.json()

		assert.equal(malformedReply.field, 'body')

		const unboundedReply = await unbounded.json()

		assert.equal(unboundedReply.field, 'maxTokens')

		const extraReply = await extra.json()

		assert.equal(extraReply.field, 'temperature')
	})

	it('refuses a model the console does not offer with 400 and the offered set', async () => {
		const response = await POST(
			judgeRequest({
				body: JSON.stringify({
					...requestSample,
					model: 'claude-haiku-5-5',
				}),
			}),
		)

		assert.equal(response.status, 400)

		const reply = await response.json()

		assert.deepEqual(reply, {
			code: 'MODEL_NOT_OFFERED',
			offered: ['claude-sonnet-5-5', 'claude-opus-5-5'],
		})
	})

	it('starts no call and answers 503 when the platform key is missing', async () => {
		delete process.env.ANTHROPIC_API_KEY
		const provider = mock.method(globalThis, 'fetch', async () =>
			messagesReply('{}', usage),
		)

		const response = await POST(judgeRequest())

		assert.equal(response.status, 503)

		const reply = await response.json()

		assert.equal(reply.reason, 'platform-key-missing')
		assert.equal(provider.mock.callCount(), 0)
	})

	it('answers 503 when the provider cannot be reached', async () => {
		mock.method(globalThis, 'fetch', async () => {
			throw new TypeError('fetch failed')
		})

		const response = await POST(judgeRequest())

		assert.equal(response.status, 503)

		const reply = await response.json()

		assert.equal(reply.reason, 'provider-unavailable')
	})

	it('answers a valid question with 200 and the contract response, priced at the list prices', async () => {
		const answerText = JSON.stringify(suggestionAnswer)
		const provider = mock.method(globalThis, 'fetch', async () =>
			messagesReply(answerText, usage),
		)

		const response = await POST(judgeRequest())

		assert.equal(response.status, 200)
		assert.equal(response.headers.get('cache-control'), 'no-store')

		const reply = await response.json()

		assert.equal(
			contractResponse(reply),
			true,
			JSON.stringify(contractResponse.errors),
		)
		assert.equal(reply.refused, false)
		assert.equal(reply.text, answerText)
		assert.deepEqual(reply.usage, usage)
		assert.equal(reply.cost.provenance, 'console')
		assert.equal(reply.cost.illustrative, true)
		assert.ok(reply.cost.cents > 0)
		assert.equal(provider.mock.callCount(), 1)
	})

	it('answers model text outside the kind schema as a refusal, charged and never forwarded', async () => {
		mock.method(globalThis, 'fetch', async () =>
			messagesReply('The model is claude-sonnet-5-5.', usage),
		)

		const response = await POST(judgeRequest())

		assert.equal(response.status, 200)

		const reply = await response.json()

		assert.equal(
			contractResponse(reply),
			true,
			JSON.stringify(contractResponse.errors),
		)
		assert.equal(reply.refused, true)
		assert.equal(reply.text, '')
		assert.notEqual(reply.reason, '')
		assert.deepEqual(reply.usage, usage)
		assert.ok(reply.cost.cents > 0)
	})

	it('logs one line per judgment with usage and cost and no source text', async () => {
		const logged = mock.method(console, 'info', () => undefined)
		mock.method(globalThis, 'fetch', async () =>
			messagesReply(JSON.stringify(suggestionAnswer), usage),
		)

		await POST(judgeRequest())

		assert.equal(logged.mock.callCount(), 1)

		const line = String(logged.mock.calls[0].arguments[0])
		const entry = JSON.parse(line)

		assert.equal(entry.event, 'judgment')
		assert.equal(entry.inputTokens, usage.inputTokens)
		assert.equal(entry.outputTokens, usage.outputTokens)
		assert.ok(entry.cents > 0)
		assert.ok(!line.includes(requestSample.prompt))
		assert.ok(!line.includes(requestSample.system))
	})
})
