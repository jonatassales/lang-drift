// The aggregate an event belongs to: its kind and the id of the one instance.
export interface EventAggregate {
	type: string
	id: string
}

// One entry of the event log. The correlation id is the aggregate's id and
// the causation id names the event that caused this one, when there is one.
export interface EventEnvelope<Payload> {
	id: string
	type: string
	version: number
	aggregate: EventAggregate
	occurredAt: string
	correlationId: string
	causationId?: string
	payload: Payload
}
