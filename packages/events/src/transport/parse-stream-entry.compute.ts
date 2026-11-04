import type { StreamEntry } from './stream-entry.types.js'

// Turns the field list Redis returns for one entry into a StreamEntry, or
// throws naming what is wrong, so the caller can park the entry with a reason.
export function parseStreamEntry(
	id: string,
	fields: string[] | null,
): StreamEntry {
	if (!fields) throw new Error(`Entry ${id} has no fields`)

	const values = new Map<string, string>()
	for (let index = 0; index + 1 < fields.length; index += 2) {
		values.set(fields[index], fields[index + 1])
	}

	const type = values.get('type')
	const aggregateId = values.get('aggregateId')
	const json = values.get('envelope')

	if (!type) throw new Error(`Entry ${id} has no type field`)
	if (!aggregateId) throw new Error(`Entry ${id} has no aggregateId field`)
	if (!json) throw new Error(`Entry ${id} has no envelope field`)

	let envelope: unknown
	try {
		envelope = JSON.parse(json)
	} catch (error) {
		const reason = error instanceof Error ? error.message : String(error)

		throw new Error(`Entry ${id} has an envelope that is not JSON: ${reason}`)
	}

	return {
		id,
		type,
		aggregateId,
		envelope,
	}
}
