import { errorReason } from '@api/book-a-demo/error-reason.fmt'

// Reports a request body the reader could not cancel, through the route's
// error log, so the failure is never silently dropped.
export function reportCancelFailure(error: unknown) {
	console.error(`Request body not cancelled: ${errorReason(error)}`)
}
