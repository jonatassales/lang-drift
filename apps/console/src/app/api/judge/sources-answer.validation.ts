import { z } from 'zod'

// The answer to a `sources` question, as `@langdrift/core/judge/sources.answer.schema.json`
// declares it: the verdict on every source read, the facts, the criteria that
// apply and the resource figures, each with where it came from.
export const sourcesAnswerSchema = z.strictObject({
	sources: z.array(
		z.strictObject({
			path: z.string().min(1),
			verdict: z.enum(['kept', 'discarded']),
			reason: z.string().min(1),
			confidence: z.number().min(0).max(1),
			role: z.string(),
			summary: z.string(),
		}),
	),
	facts: z.array(
		z.strictObject({
			topic: z.enum([
				'resource',
				'component',
				'acceptance',
				'budget',
				'limits',
			]),
			text: z.string().min(1),
			source: z.string().min(1),
		}),
	),
	criteria: z.array(
		z.strictObject({
			id: z.string().min(1),
			applies: z.string().min(1),
		}),
	),
	figures: z.array(
		z.strictObject({
			topic: z.literal('resource'),
			fact: z.string().min(1),
			source: z.string().min(1),
			cpuMillis: z.int().min(0),
			memoryMiB: z.int().min(0),
		}),
	),
})
