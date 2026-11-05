import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { contractValidator } from '@app/api/judge/testdata/contract-validator.mock'
import { sampleVariants } from '@app/api/judge/testdata/sample-variants.mock'
import requestSample from '@langdrift/core/judge/request.sample.json' with {
	type: 'json',
}
import requestSchema from '@langdrift/core/judge/request.schema.json' with {
	type: 'json',
}

import { judgeRequestSchema } from './judge-request.validation.ts'

describe('console/api/judge/judge-request.validation', () => {
	it('accepts and refuses exactly what the contract request schema does', () => {
		const contract = contractValidator(requestSchema)

		for (const question of [requestSample, ...sampleVariants(requestSample)]) {
			assert.equal(
				judgeRequestSchema.safeParse(question).success,
				contract(question),
				JSON.stringify(question),
			)
		}
	})
})
