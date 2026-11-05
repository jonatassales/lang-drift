import { z } from 'zod'

// The shape every entry of the event log has, whatever its payload. An event
// type refines the payload with its own schema.
export const eventEnvelopeSchema = z
	.object({
		id: z.uuid(),
		type: z.string().min(1),
		version: z.number().int().positive(),
		aggregate: z
			.object({
				type: z.string().min(1),
				id: z.uuid(),
			})
			.strict(),
		occurredAt: z.iso.datetime(),
		correlationId: z.uuid(),
		causationId: z.uuid().optional(),
		payload: z.unknown(),
	})
	.strict()
