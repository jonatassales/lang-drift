// Reports a request body the reader could not cancel, through the route's
// error log, so the failure is never silently dropped.
export function reportCancelFailure(error: unknown) {
	const reason = error instanceof Error ? error.message : String(error)

	console.error(`Request body not cancelled: ${reason}`)
}
