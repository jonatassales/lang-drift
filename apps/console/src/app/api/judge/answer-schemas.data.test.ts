import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { contractValidator } from '@app/api/judge/testdata/contract-validator.mock'
import { sampleVariants } from '@app/api/judge/testdata/sample-variants.mock'
import episodesAnswer from '@langdrift/core/judge/episodes.answer.sample.json' with {
	type: 'json',
}
import episodesSchema from '@langdrift/core/judge/episodes.answer.schema.json' with {
	type: 'json',
}
import reservedPartAnswer from '@langdrift/core/judge/reserved-part.answer.sample.json' with {
	type: 'json',
}
import reservedPartSchema from '@langdrift/core/judge/reserved-part.answer.schema.json' with {
	type: 'json',
}
import sourcesAnswer from '@langdrift/core/judge/sources.answer.sample.json' with {
	type: 'json',
}
import sourcesSchema from '@langdrift/core/judge/sources.answer.schema.json' with {
	type: 'json',
}
import suggestionAnswer from '@langdrift/core/judge/suggestion.answer.sample.json' with {
	type: 'json',
}
import suggestionSchema from '@langdrift/core/judge/suggestion.answer.schema.json' with {
	type: 'json',
}

import { answerSchemas } from './answer-schemas.data.ts'
import type { QuestionKind } from './judge.types.ts'

// A reserved-part answer that proposes no part, the other side of the
// schema's if/then: empty paths need no goal, criterion or category.
const noReservedPart = {
	...reservedPartAnswer,
	acceptance: [],
	category: null,
	excludedKind: 'rename',
	goal: '',
	paths: [],
	plannedLines: 0,
}
const contracts: [QuestionKind, object, unknown][] = [
	['sources', sourcesSchema, sourcesAnswer],
	['episodes', episodesSchema, episodesAnswer],
	['suggestion', suggestionSchema, suggestionAnswer],
	['reserved-part', reservedPartSchema, reservedPartAnswer],
	['reserved-part', reservedPartSchema, noReservedPart],
]

describe('console/api/judge/answer-schemas.data', () => {
	it('accepts and refuses exactly what the contract schema of each kind does', () => {
		for (const [kind, schema, sample] of contracts) {
			const contract = contractValidator(schema)

			for (const answer of [sample, ...sampleVariants(sample)]) {
				const accepted = answerSchemas[kind].safeParse(answer).success

				assert.equal(
					accepted,
					contract(answer),
					`${kind}: ${JSON.stringify(answer)}`,
				)
			}
		}
	})
})
