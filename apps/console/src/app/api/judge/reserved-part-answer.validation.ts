import { z } from 'zod'

// The answer to a `reserved-part` question, as `@langdrift/core/judge/reserved-part.answer.schema.json`
// declares it: the part a person writes by hand, or the kind of ticket that
// has none. Proposed paths come with a goal, an acceptance criterion and a
// category, the schema's if/then; plannedLines against totalPlannedLines is
// left to the CLI's ParseAnswer, as the schema says.
export const reservedPartAnswerSchema = z
	.strictObject({
		paths: z.array(z.string().min(1)),
		goal: z.string(),
		acceptance: z.array(z.string().min(1)),
		category: z.enum(['core-domain', 'hard-part', 'new-stack']).nullable(),
		plannedLines: z.int().min(0),
		totalPlannedLines: z.int().min(0),
		excludedKind: z
			.enum([
				'boilerplate',
				'configuration',
				'generated',
				'rename',
				'tool-migration',
			])
			.nullable(),
		reason: z.string().min(1),
		confidence: z.number().min(0).max(1),
	})
	.superRefine((answer, context) => {
		if (answer.paths.length === 0) return

		if (answer.goal === '') {
			context.addIssue({
				code: 'custom',
				message: 'proposed paths need a goal',
				path: ['goal'],
			})
		}

		if (answer.acceptance.length === 0) {
			context.addIssue({
				code: 'custom',
				message: 'proposed paths need an acceptance criterion',
				path: ['acceptance'],
			})
		}

		if (answer.category === null) {
			context.addIssue({
				code: 'custom',
				message: 'proposed paths need a category',
				path: ['category'],
			})
		}
	})
