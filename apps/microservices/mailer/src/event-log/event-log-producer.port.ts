import type { Observable } from 'rxjs'

// What the service needs from the event log to append an event: Nest's
// client contract, narrowed to the one call it makes.
export interface EventLogProducer {
	emit(pattern: string, data: unknown): Observable<unknown>
}
