import { z } from 'zod'

const routableEnvelopeSchema = z
	.object({
		type: z.string().min(1),
		aggregate: z
			.object({
				id: z.string().min(1),
			})
			.loose(),
	})
	.loose()

// The fields one entry carries: the type and the aggregate id flat, so a
// reader routes and filters without parsing, and the envelope as JSON.
export function streamFields(type: string, envelope: unknown): string[] {
	const routable = routableEnvelopeSchema.safeParse(envelope)
	if (!routable.success) {
		throw new Error('An event needs a type and an aggregate id to be appended')
	}

	if (routable.data.type !== type) {
		throw new Error(
			`An event is appended under its own type: ${routable.data.type}, not ${type}`,
		)
	}

	return [
		'type',
		type,
		'aggregateId',
		routable.data.aggregate.id,
		'envelope',
		JSON.stringify(envelope),
	]
}
