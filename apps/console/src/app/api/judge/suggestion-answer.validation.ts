import { z } from 'zod'

// The answer to a `suggestion` question, as `@langdrift/core/judge/suggestion.answer.schema.json`
// declares it: the model and the effort for the ticket, why, and at most
// three further reasons.
export const suggestionAnswerSchema = z.strictObject({
	model: z.string().min(1),
	effort: z.enum(['low', 'medium', 'high']),
	effortReason: z.string().min(1),
	modelReason: z.string().min(1),
	reasons: z.array(z.string().min(1)).max(3),
})
