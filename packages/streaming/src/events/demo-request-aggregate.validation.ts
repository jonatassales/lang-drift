import { z } from 'zod'

import { demoRequestAggregateType } from './demo-request-stream.data.js'

// The aggregate every demo-request event belongs to.
export const demoRequestAggregateSchema = z
	.object({
		type: z.literal(demoRequestAggregateType),
		id: z.uuid(),
	})
	.strict()
