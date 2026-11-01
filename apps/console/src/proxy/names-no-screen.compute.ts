import { isEntityId } from '@data/entity-id'

export function namesNoScreen(pathname: string) {
	return !pathname.split('/').filter(Boolean).every(isEntityId)
}
