import { z } from 'zod'

// The question of the judge contract, as `@langdrift/core/judge/request.schema.json`
// declares it: no field beyond these, each one present.
export const judgeRequestSchema = z.strictObject({
	kind: z.enum(['sources', 'episodes', 'suggestion', 'reserved-part']),
	model: z.string().min(1),
	system: z.string().min(1),
	prompt: z.string().min(1),
	maxTokens: z.int().min(1).max(32000),
	promptVersion: z.string().min(1),
})
