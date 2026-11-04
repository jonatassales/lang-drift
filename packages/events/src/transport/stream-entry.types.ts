// One entry of a stream as the transport reads it: the flat fields the
// producer wrote beside the envelope, so a consumer can route without parsing.
export interface StreamEntry {
	id: string
	type: string
	aggregateId: string
	envelope: unknown
}
