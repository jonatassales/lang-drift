import type { z } from 'zod'

import { episodesAnswerSchema } from './episodes-answer.validation'
import type { QuestionKind } from './judge.types'
import { reservedPartAnswerSchema } from './reserved-part-answer.validation'
import { sourcesAnswerSchema } from './sources-answer.validation'
import { suggestionAnswerSchema } from './suggestion-answer.validation'

// The schema each kind's answer text must validate against before the
// console forwards it.
export const answerSchemas: Readonly<Record<QuestionKind, z.ZodType>> = {
	sources: sourcesAnswerSchema,
	episodes: episodesAnswerSchema,
	suggestion: suggestionAnswerSchema,
	'reserved-part': reservedPartAnswerSchema,
}
