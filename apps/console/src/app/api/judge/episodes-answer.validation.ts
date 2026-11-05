import { z } from 'zod'

// The answer to an `episodes` question, as `@langdrift/core/judge/episodes.answer.schema.json`
// declares it: the earlier runs that bear on the ticket, each with the reason.
export const episodesAnswerSchema = z.strictObject({
	episodes: z.array(
		z.strictObject({
			run: z.string().min(1),
			reason: z.string().min(1),
		}),
	),
})
