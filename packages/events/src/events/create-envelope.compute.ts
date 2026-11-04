import { randomUUID } from 'node:crypto'

import type { EventAggregate, EventEnvelope } from './event-envelope.types'

export interface CreateEnvelopeInput<Payload> {
	type: string
	version: number
	aggregate: EventAggregate
	payload: Payload
	causationId?: string
	now?: Date
}

// Stamps a new event: a fresh id, the time, and the aggregate's id as the
// correlation id, so every event of one request shares it.
export function createEnvelope<Payload>(
	input: CreateEnvelopeInput<Payload>,
): EventEnvelope<Payload> {
	const occurredAt = (input.now ?? new Date()).toISOString()
	const envelope: EventEnvelope<Payload> = {
		id: randomUUID(),
		type: input.type,
		version: input.version,
		aggregate: input.aggregate,
		occurredAt,
		correlationId: input.aggregate.id,
		payload: input.payload,
	}

	if (input.causationId) envelope.causationId = input.causationId

	return envelope
}
