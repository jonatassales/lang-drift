import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import episodesAnswer from '@langdrift/core/judge/episodes.answer.sample.json' with {
	type: 'json',
}
import reservedPartAnswer from '@langdrift/core/judge/reserved-part.answer.sample.json' with {
	type: 'json',
}
import responseSample from '@langdrift/core/judge/response.sample.json' with {
	type: 'json',
}
import sourcesAnswer from '@langdrift/core/judge/sources.answer.sample.json' with {
	type: 'json',
}
import suggestionAnswer from '@langdrift/core/judge/suggestion.answer.sample.json' with {
	type: 'json',
}

import type {
	ModelPrices,
	ProviderResult,
	QuestionKind,
} from './judge.types.ts'
import { judgeResponse } from './judge-response.compute.ts'

// The list prices of claude-sonnet-5-5, claude-opus-5-5 and claude-haiku-5-5,
// in cents per million tokens, written out so each cost below is computed by
// hand from them.
const sonnet: ModelPrices = {
	inputCentsPerMillion: 200,
	outputCentsPerMillion: 1000,
}
const opus: ModelPrices = {
	inputCentsPerMillion: 400,
	outputCentsPerMillion: 2000,
}
const haiku: ModelPrices = {
	inputCentsPerMillion: 10,
	outputCentsPerMillion: 50,
}
const samples: [QuestionKind, unknown][] = [
	['sources', sourcesAnswer],
	['episodes', episodesAnswer],
	['suggestion', suggestionAnswer],
	['reserved-part', reservedPartAnswer],
]
const finished: ProviderResult = {
	stopReason: 'end_turn',
	text: JSON.stringify(suggestionAnswer),
	usage: responseSample.usage,
}

describe('console/api/judge/judge-response.compute', () => {
	it('judges every sample answer of the contract and forwards its text unchanged', () => {
		for (const [kind, answer] of samples) {
			const text = JSON.stringify(answer)
			const judged = judgeResponse(
				{
					...finished,
					text,
				},
				kind,
				sonnet,
			)

			if ('reason' in judged) assert.fail(`${kind}: ${judged.reason}`)

			assert.equal(judged.text, text)
			assert.deepEqual(judged.usage, responseSample.usage)
			assert.equal(judged.judgeVersion, responseSample.judgeVersion)
			assert.deepEqual(judged.cost, responseSample.cost)
		}
	})

	it('prices the sample usage at the list prices, rounded up to the cent', () => {
		// 1,840 input tokens at 200 and 96 output tokens at 1,000 cents per
		// million: 0.368 + 0.096 = 0.464 cents, which rounds up to 1.
		const judged = judgeResponse(finished, 'suggestion', sonnet)

		if ('reason' in judged) assert.fail(judged.reason)

		assert.equal(judged.cost.cents, 1)
		assert.equal(judged.cost.illustrative, true)
		assert.equal(judged.cost.provenance, 'console')
	})

	it('prices input and output tokens apart, each at its own rate', () => {
		// 1,000,000 input at 400 plus 250,000 output at 2,000: 400 + 500 cents.
		const inputHeavy = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 1_000_000,
					outputTokens: 250_000,
				},
			},
			'suggestion',
			opus,
		)
		// 250,000 input at 400 plus 1,000,000 output at 2,000: 100 + 2,000 cents.
		const outputHeavy = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 250_000,
					outputTokens: 1_000_000,
				},
			},
			'suggestion',
			opus,
		)

		if ('reason' in inputHeavy) assert.fail(inputHeavy.reason)
		if ('reason' in outputHeavy) assert.fail(outputHeavy.reason)

		assert.equal(inputHeavy.cost.cents, 900)
		assert.equal(outputHeavy.cost.cents, 2100)
	})

	it('rounds any fraction of a cent up', () => {
		// 1,000,001 input tokens at 400 cents per million: 400.0004 cents.
		const judged = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 1_000_001,
					outputTokens: 0,
				},
			},
			'suggestion',
			opus,
		)

		if ('reason' in judged) assert.fail(judged.reason)

		assert.equal(judged.cost.cents, 401)
	})

	it('charges one cent for one token, never zero', () => {
		const oneInput = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 1,
					outputTokens: 0,
				},
			},
			'suggestion',
			haiku,
		)
		const oneOutput = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 0,
					outputTokens: 1,
				},
			},
			'suggestion',
			haiku,
		)

		if ('reason' in oneInput) assert.fail(oneInput.reason)
		if ('reason' in oneOutput) assert.fail(oneOutput.reason)

		assert.equal(oneInput.cost.cents, 1)
		assert.equal(oneOutput.cost.cents, 1)
	})

	it('refuses a non-empty answer with a usage of zero tokens', () => {
		const refused = judgeResponse(
			{
				...finished,
				usage: {
					inputTokens: 0,
					outputTokens: 0,
				},
			},
			'suggestion',
			sonnet,
		)

		assert.deepEqual(refused, {
			reason: 'the provider reported no tokens for a non-empty answer',
		})
	})

	it('refuses a usage with a missing, negative or fractional count', () => {
		const usages = [
			JSON.parse('{"inputTokens": 1840}'),
			{
				inputTokens: -1,
				outputTokens: 96,
			},
			{
				inputTokens: 1840,
				outputTokens: 0.5,
			},
		]

		for (const usage of usages) {
			const refused = judgeResponse(
				{
					...finished,
					usage,
				},
				'suggestion',
				sonnet,
			)

			assert.deepEqual(refused, {
				reason:
					'the provider reported a missing, negative or fractional token count',
			})
		}
	})

	it('refuses an answer the model did not finish, naming how it stopped', () => {
		const truncated = judgeResponse(
			{
				...finished,
				stopReason: 'max_tokens',
			},
			'suggestion',
			sonnet,
		)
		const declined = judgeResponse(
			{
				...finished,
				stopReason: 'refusal',
			},
			'suggestion',
			sonnet,
		)
		const unknown = judgeResponse(
			{
				...finished,
				stopReason: null,
			},
			'suggestion',
			sonnet,
		)

		assert.deepEqual(truncated, {
			reason: 'the model stopped with max_tokens, not end_turn',
		})
		assert.deepEqual(declined, {
			reason: 'the model stopped with refusal, not end_turn',
		})
		assert.deepEqual(unknown, {
			reason: 'the model stopped with no stop reason, not end_turn',
		})
	})

	it('refuses an empty or blank answer', () => {
		for (const text of ['', ' \n\t']) {
			const refused = judgeResponse(
				{
					...finished,
					text,
				},
				'suggestion',
				sonnet,
			)

			assert.deepEqual(refused, {
				reason: 'the model answered with no text',
			})
		}
	})

	it('refuses a text that is not JSON', () => {
		const refused = judgeResponse(
			{
				...finished,
				text: 'The model is claude-sonnet-5-5.',
			},
			'suggestion',
			sonnet,
		)

		assert.deepEqual(refused, {
			reason: 'the answer is not JSON',
		})
	})

	it('refuses a JSON answer outside the kind schema, naming the field at fault', () => {
		const { effort: _effort, ...withoutEffort } = suggestionAnswer
		const missing = judgeResponse(
			{
				...finished,
				text: JSON.stringify(withoutEffort),
			},
			'suggestion',
			sonnet,
		)
		const extra = judgeResponse(
			{
				...finished,
				text: JSON.stringify({
					...suggestionAnswer,
					temperature: 0,
				}),
			},
			'suggestion',
			sonnet,
		)
		const anotherKind = judgeResponse(
			{
				...finished,
				text: JSON.stringify(episodesAnswer),
			},
			'suggestion',
			sonnet,
		)
		const notAnObject = judgeResponse(
			{
				...finished,
				text: '[]',
			},
			'suggestion',
			sonnet,
		)

		assert.deepEqual(missing, {
			reason: 'the suggestion answer breaks its schema at effort',
		})
		assert.deepEqual(extra, {
			reason: 'the suggestion answer breaks its schema at temperature',
		})
		assert.ok('reason' in anotherKind)
		assert.deepEqual(notAnObject, {
			reason: 'the suggestion answer breaks its schema at answer',
		})
	})

	it('carries no text in any refusal', () => {
		const refusals = [
			judgeResponse(
				{
					...finished,
					stopReason: 'max_tokens',
				},
				'suggestion',
				sonnet,
			),
			judgeResponse(
				{
					...finished,
					text: '{"model": "claude-sonnet-5-5"}',
				},
				'suggestion',
				sonnet,
			),
		]

		for (const refusal of refusals) {
			assert.deepEqual(Object.keys(refusal), ['reason'])
		}
	})
})
