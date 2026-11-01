import { isEntityId } from './is-entity-id.compute.ts'

// Every segment of a console path is a fixed name or an entity id, and both
// keep to the alphabet; the query is the caller's to encode.
export function consolePathSafe(path: string) {
	const [pathname] = path.split('?')

	return pathname.split('/').every(isEntityId)
}
