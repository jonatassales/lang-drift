import type { StreamItem } from '@langdrift/events/transport'

// The aggregate id an entry carries as a flat field, read without parsing
// its envelope, so a range over the stream filters cheaply.
export function aggregateIdOf(item: StreamItem): string | null {
	const fields = item[1]
	if (!fields) return null

	for (let index = 0; index + 1 < fields.length; index += 2) {
		if (fields[index] === 'aggregateId') return fields[index + 1]
	}

	return null
}
