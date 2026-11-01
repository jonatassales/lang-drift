// The alphabet of an entity id, as extub issues them (seed-runtime-4,
// ticket-0042): letters, digits, hyphen and underscore, 1 to 128 of them. An
// id in it is a path segment as it stands, so it is checked, never encoded:
// Next already encodes route params, and encoding one again would encode it
// twice.
const entityId = /^[A-Za-z0-9_-]{1,128}$/u

export function isEntityId(value: string) {
	return entityId.test(value)
}
